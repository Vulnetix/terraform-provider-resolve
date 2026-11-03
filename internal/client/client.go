// Package client talks to the Vulnetix config API (vdb-site /v1).
//
// Every request and response shape the provider depends on is in this one
// file, so if the server's contract moves, this is the only file that should
// need to follow it. The contract is the OpenAPI document the server
// publishes at {endpoint}/config/openapi.yaml (public, no auth):
//
//	GET    {endpoint}/config/kinds                     kinds the server applies
//	GET    {endpoint}/config/objects/{kind}/{name}     ObjectResult, 404 when absent
//	PUT    {endpoint}/config/objects/{kind}/{name}     ObjectPutRequest -> ObjectResult
//	DELETE {endpoint}/config/objects/{kind}/{name}?manager=m
//	GET    {endpoint}/config/ownership                 {"objects": [Ownership]}
//	GET    {endpoint}/config/export?format=yaml        the tenant's documents, secrets redacted
//	GET    {endpoint}/detection-rule-feeds             the tenant's feeds
//	GET    {endpoint}/detection-rule-feeds/{uuid}      one feed with its entries
//	GET    {endpoint}/ai-firewall/models               the AI model catalogue
//	GET    {endpoint}/vex/statements                   the tenant's VEX statements
//	GET    {endpoint}/integrations/github/installations/repos  GitHub App repositories
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultEndpoint is the console API every signed-in surface calls. The
// config routes are mounted flat under its /v1.
const DefaultEndpoint = "https://www.vulnetix.com/api/site/v1"

// DefaultManager is how objects this provider applies are recorded in the
// server's ownership table, so a GitOps apply with --prune never deletes them.
const DefaultManager = "terraform"

// APIVersion is the document format every kind uses today.
const APIVersion = "vulnetix.com/v1"

// ErrNotFound is returned when the object or feed does not exist. A 404 is
// also an *APIError that wraps it, so the server's message is not lost: the
// same status means an unknown kind, which a PUT reports as an error.
var ErrNotFound = errors.New("not found")

// Client is a config API client. The zero value is not usable; build one
// with New.
type Client struct {
	endpoint  string
	token     string
	apiKey    string
	orgID     string
	manager   string
	userAgent string
	http      *http.Client
}

// Options configure New.
type Options struct {
	Endpoint string
	// Token is a bearer token: an Authentik OAuth2 access token or a
	// session JWT. It is sent as "Authorization: Bearer <token>" and wins
	// over APIKey when both are set.
	Token string
	// APIKey and OrgID are an organisation ApiKey, sent as
	// "Authorization: ApiKey <org>:<key>". APIKey may already carry the
	// "<org>:" prefix.
	APIKey    string
	OrgID     string
	Manager   string
	UserAgent string
	HTTP      *http.Client
}

// New builds a client. Endpoint and Manager fall back to their defaults.
func New(o Options) *Client {
	c := &Client{
		endpoint:  strings.TrimRight(o.Endpoint, "/"),
		token:     o.Token,
		apiKey:    o.APIKey,
		orgID:     o.OrgID,
		manager:   o.Manager,
		userAgent: o.UserAgent,
		http:      o.HTTP,
	}
	if c.endpoint == "" {
		c.endpoint = DefaultEndpoint
	}
	if c.manager == "" {
		c.manager = DefaultManager
	}
	if c.http == nil {
		c.http = &http.Client{Timeout: 60 * time.Second}
	}
	return c
}

// Manager is the ownership manager this client applies as.
func (c *Client) Manager() string { return c.manager }

// Document is one vulnetix.com/v1 config document, the same shape the CLI
// reads from .vulnetix/config and the engine exports.
type Document struct {
	APIVersion string         `json:"apiVersion"`
	Kind       string         `json:"kind"`
	Metadata   Metadata       `json:"metadata"`
	Spec       map[string]any `json:"spec"`
}

// Metadata is a document's metadata block.
type Metadata struct {
	Name        string            `json:"name"`
	Org         string            `json:"org,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Object is the OpenAPI ObjectResult: what every single-object route
// returns.
type Object struct {
	Document Document     `json:"document"`
	Status   ObjectStatus `json:"status"`
	// Action is what a PUT or DELETE did: create, update, noop, delete.
	Action   string      `json:"action,omitempty"`
	Warnings []Condition `json:"warnings,omitempty"`
}

// ObjectStatus is the OpenAPI ObjectStatus. A kind may report more fields
// beside these (Extra): the authorization URL a messaging connection waits
// on, an invitation's expiry, a service account's credentials, shown once.
type ObjectStatus struct {
	// Phase is Applied, Unchanged, Deleted, Absent or Unmanaged; or, for a
	// kind that waits on someone, Pending or AwaitingAuthorization; or, for
	// a write in Guardian that failed, Failed (Message says why).
	Phase       string `json:"phase"`
	Manager     string `json:"manager,omitempty"`
	ObjectUUID  string `json:"objectUuid,omitempty"`
	AppliedAt   int64  `json:"appliedAt,omitempty"`
	AppliedHash string `json:"appliedHash,omitempty"`
	SourceRef   string `json:"sourceRef,omitempty"`
	Drift       bool   `json:"drift,omitempty"`
	// Message says why a write failed.
	Message string `json:"message,omitempty"`
	// Extra holds every other field of the status, as the kind wrote it.
	Extra map[string]any `json:"-"`
}

// Status phases beyond the engine's own Applied and Unchanged.
const (
	PhaseApplied               = "Applied"
	PhasePending               = "Pending"
	PhaseAwaitingAuthorization = "AwaitingAuthorization"
	PhaseFailed                = "Failed"
)

// UnmarshalJSON reads the known fields and keeps the rest in Extra.
func (s *ObjectStatus) UnmarshalJSON(b []byte) error {
	type plain ObjectStatus
	var p plain
	if err := unmarshalNumber(b, &p); err != nil {
		return err
	}
	var all map[string]any
	if err := unmarshalNumber(b, &all); err != nil {
		return err
	}
	for _, k := range []string{"phase", "manager", "objectUuid", "appliedAt", "appliedHash", "sourceRef", "drift", "message"} {
		delete(all, k)
	}
	if len(all) > 0 {
		p.Extra = all
	}
	*s = ObjectStatus(p)
	return nil
}

// Condition is a plan warning or blocking condition (OpenAPI Condition).
type Condition struct {
	Code    string `json:"code"`
	Kind    string `json:"kind,omitempty"`
	Name    string `json:"name,omitempty"`
	Message string `json:"message"`
}

// APIError is a non-2xx answer. Code is the server's "error" field
// (plan_blocked, managed_by_other, not_managed and so on); Message carries it
// with the details, or the raw body.
type APIError struct {
	Method, URL string
	Status      int
	Code        string
	Message     string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s %s: HTTP %d: %s", e.Method, e.URL, e.Status, e.Message)
}

// Unwrap makes a 404 match ErrNotFound.
func (e *APIError) Unwrap() error {
	if e.Status == http.StatusNotFound {
		return ErrNotFound
	}
	return nil
}

// authorization is the Authorization header for the configured credential:
// "Bearer <token>" for an Authentik access token or a session JWT, or
// "ApiKey <org>:<key>" for an organisation ApiKey, the format the Vulnetix CLI
// sends. A key that already starts with "<org>:" is not prefixed twice.
func (c *Client) authorization() string {
	switch {
	case c.token != "":
		return "Bearer " + c.token
	case c.apiKey != "" && c.orgID != "":
		return "ApiKey " + c.orgID + ":" + strings.TrimPrefix(c.apiKey, c.orgID+":")
	default:
		return ""
	}
}

// request sends one request and returns the body of a 2xx answer.
func (c *Client) request(ctx context.Context, method, path string, query url.Values, body any, accept string) ([]byte, error) {
	u := c.endpoint + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, rd)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if a := c.authorization(); a != "" {
		req.Header.Set("Authorization", a)
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		code, msg := errorMessage(raw)
		return nil, &APIError{Method: method, URL: u, Status: resp.StatusCode, Code: code, Message: msg}
	}
	return raw, nil
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	raw, err := c.request(ctx, method, path, query, body, "application/json")
	if err != nil {
		return err
	}
	if out == nil || len(bytes.TrimSpace(raw)) == 0 {
		return nil
	}
	// Numbers stay json.Number so an integer spec field never passes through
	// float64 on its way into Terraform state.
	if err := unmarshalNumber(raw, out); err != nil {
		return fmt.Errorf("%s %s: decode response: %w", method, path, err)
	}
	return nil
}

// errorMessage reads the server's error bodies: Error and Conflict
// ({error, message}), PlanError (with plan.blocking conditions) and the
// validation answer (with problems).
func errorMessage(raw []byte) (code, msg string) {
	var e struct {
		Error    string `json:"error"`
		Message  string `json:"message"`
		Problems []struct {
			Kind    string `json:"kind"`
			Name    string `json:"name"`
			Message string `json:"message"`
		} `json:"problems"`
		Plan *struct {
			Blocking []Condition `json:"blocking"`
		} `json:"plan"`
	}
	if json.Unmarshal(raw, &e) != nil || e.Error == "" {
		return "", strings.TrimSpace(string(raw))
	}
	msg = e.Error
	if e.Message != "" {
		msg += ": " + e.Message
	}
	for _, p := range e.Problems {
		msg += "\n  " + strings.TrimSpace(objectRef(p.Kind, p.Name)+" "+p.Message)
	}
	if e.Plan != nil {
		for _, b := range e.Plan.Blocking {
			msg += "\n  " + strings.TrimSpace(b.Code+" "+objectRef(b.Kind, b.Name)+": "+b.Message)
		}
	}
	return e.Error, msg
}

func objectRef(kind, name string) string {
	if kind == "" {
		return name
	}
	return kind + "/" + name
}

func objectPath(kind, name string) string {
	return "/config/objects/" + url.PathEscape(kind) + "/" + url.PathEscape(name)
}

func unmarshalNumber(b []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(v)
}

// GetObject reads one object. It returns ErrNotFound when the tenant has none.
func (c *Client) GetObject(ctx context.Context, kind, name string) (*Object, error) {
	var o Object
	if err := c.do(ctx, http.MethodGet, objectPath(kind, name), nil, nil, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// PutObject applies one document as this client's manager (server-side
// apply: plan and apply in one transaction) and returns the result. A plan
// with blocking conditions is refused with HTTP 422, whose APIError lists
// them.
func (c *Client) PutObject(ctx context.Context, doc Document) (*Object, error) {
	if doc.APIVersion == "" {
		doc.APIVersion = APIVersion
	}
	if doc.Spec == nil {
		doc.Spec = map[string]any{}
	}
	body := map[string]any{"document": doc, "manager": c.manager}
	var o Object
	if err := c.do(ctx, http.MethodPut, objectPath(doc.Kind, doc.Metadata.Name), nil, body, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// DeleteObject deletes one object as this client's manager. Deleting an
// absent object succeeds. An object another manager owns, or that no
// manager owns (made in the console), is refused with HTTP 409.
func (c *Client) DeleteObject(ctx context.Context, kind, name string) error {
	err := c.do(ctx, http.MethodDelete, objectPath(kind, name), url.Values{"manager": {c.manager}}, nil, nil)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	return err
}

// Ownership is one record of GET /config/ownership: an object managed as
// code, by which manager and from where.
type Ownership struct {
	TenantUUID  string            `json:"tenantUuid"`
	Kind        string            `json:"kind"`
	Name        string            `json:"name"`
	ObjectUUID  string            `json:"objectUuid"`
	Manager     string            `json:"manager"`
	AppliedHash string            `json:"appliedHash"`
	AppliedAt   int64             `json:"appliedAt"`
	AppliedBy   string            `json:"appliedBy"`
	SourceRef   string            `json:"sourceRef,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// ListOwnership lists ownership records, sorted by kind and name. kind and
// name filter when set.
func (c *Client) ListOwnership(ctx context.Context, kind, name string) ([]Ownership, error) {
	q := url.Values{}
	if kind != "" {
		q.Set("kind", kind)
	}
	if name != "" {
		q.Set("name", name)
	}
	var out struct {
		Objects []Ownership `json:"objects"`
	}
	if err := c.do(ctx, http.MethodGet, "/config/ownership", q, nil, &out); err != nil {
		return nil, err
	}
	return out.Objects, nil
}

// ExportFilter narrows GET /config/export. Empty fields do not filter.
type ExportFilter struct {
	Kind, Name, Manager string
}

// Export returns the tenant's current state as a YAML stream of documents,
// in dependency order, with secrets replaced by valueFrom.env placeholders.
func (c *Client) Export(ctx context.Context, f ExportFilter) (string, error) {
	q := url.Values{"format": {"yaml"}}
	if f.Kind != "" {
		q.Set("kind", f.Kind)
	}
	if f.Name != "" {
		q.Set("name", f.Name)
	}
	if f.Manager != "" {
		q.Set("manager", f.Manager)
	}
	raw, err := c.request(ctx, http.MethodGet, "/config/export", q, nil, "application/yaml")
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// Kinds is GET /config/kinds.
type Kinds struct {
	APIVersion    string `json:"apiVersion"`
	SchemaVersion string `json:"schemaVersion"`
	Kinds         []Kind `json:"kinds"`
	// SchemaKinds have a released schema the server does not apply yet.
	SchemaKinds []string `json:"schemaKinds"`
}

// Kind is one entry of Kinds.
type Kind struct {
	Kind      string          `json:"kind"`
	SchemaURL string          `json:"schemaUrl"`
	DependsOn []string        `json:"dependsOn"`
	Schema    json.RawMessage `json:"schema"`
}

// ListKinds lists the kinds the server applies.
func (c *Client) ListKinds(ctx context.Context) (*Kinds, error) {
	var k Kinds
	if err := c.do(ctx, http.MethodGet, "/config/kinds", nil, nil, &k); err != nil {
		return nil, err
	}
	return &k, nil
}

// Feed is a detection rule feed (handler/detection_rule_feeds.go
// detectionRuleFeedRow).
type Feed struct {
	UUID          string  `json:"uuid"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Status        string  `json:"status"`
	FeedType      string  `json:"feedType"`
	SchemaVersion int64   `json:"schemaVersion"`
	RuleCount     int64   `json:"ruleCount"`
	PublicBaseURL *string `json:"publicBaseUrl"`
	PublishedAt   *int64  `json:"publishedAt"`
	RevokedAt     *int64  `json:"revokedAt"`
	CreatedAt     int64   `json:"createdAt"`
	UpdatedAt     *int64  `json:"updatedAt"`
}

// FeedEntry is one rule in a feed (detectionRuleFeedEntryRow).
type FeedEntry struct {
	ID         int64  `json:"id"`
	RuleTable  string `json:"ruleTable"`
	RuleType   string `json:"ruleType"`
	RuleUUID   string `json:"ruleUuid"`
	RuleName   string `json:"ruleName"`
	RuleSource string `json:"ruleSource"`
	MatchedCve string `json:"matchedCve"`
	Disabled   bool   `json:"disabled"`
	AddedAt    int64  `json:"addedAt"`
}

// RuleTableVirtualPatch is the ruleTable of an entry that is a Vulnetix
// generated virtual patch (handler/defence_formats.go TableVirtualPatch).
const RuleTableVirtualPatch = "virtual_patch"

// ErrPlanLocked is returned when the tenant's plan does not include feeds.
// The feed routes answer 200 with planLocked rather than 403.
var ErrPlanLocked = errors.New("detection rule feeds need a Pro plan or higher")

type planLock struct {
	PlanLocked   bool   `json:"planLocked"`
	RequiredPlan string `json:"requiredPlan"`
}

func (p planLock) err() error {
	if !p.PlanLocked {
		return nil
	}
	if p.RequiredPlan != "" {
		return fmt.Errorf("%w (requires %s)", ErrPlanLocked, p.RequiredPlan)
	}
	return ErrPlanLocked
}

// ListFeeds lists the tenant's detection rule feeds, newest first.
func (c *Client) ListFeeds(ctx context.Context) ([]Feed, error) {
	var out struct {
		planLock
		Feeds []Feed `json:"feeds"`
	}
	if err := c.do(ctx, http.MethodGet, "/detection-rule-feeds", nil, nil, &out); err != nil {
		return nil, err
	}
	if err := out.err(); err != nil {
		return nil, err
	}
	return out.Feeds, nil
}

// GetFeed reads one feed and its entries.
func (c *Client) GetFeed(ctx context.Context, uuid string) (*Feed, []FeedEntry, error) {
	var out struct {
		planLock
		Feed    *Feed       `json:"feed"`
		Entries []FeedEntry `json:"entries"`
	}
	if err := c.do(ctx, http.MethodGet, "/detection-rule-feeds/"+url.PathEscape(uuid), nil, nil, &out); err != nil {
		return nil, nil, err
	}
	if err := out.err(); err != nil {
		return nil, nil, err
	}
	if out.Feed == nil {
		return nil, nil, ErrNotFound
	}
	return out.Feed, out.Entries, nil
}

// AiModel is one entry of the AI firewall model catalogue
// (handler/ai_firewall_models.go AiModelCatalogEntry): a built-in
// provider's model or one of the tenant's custom provider's.
type AiModel struct {
	UUID         string `json:"uuid"`
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	ProviderUUID string `json:"providerUuid"`
	ProviderSlug string `json:"providerSlug"`
	IsActive     bool   `json:"isActive"`
	ReleasedAt   int64  `json:"releasedAt"`
}

// AiModelFilter narrows ListAiModels. Empty fields do not filter.
type AiModelFilter struct {
	// ProviderUUID is a provider's uuid.
	ProviderUUID string
	// Search matches the model's slug or name, case-insensitively.
	Search string
}

// ListAiModels reads every page of GET /ai-firewall/models.
func (c *Client) ListAiModels(ctx context.Context, f AiModelFilter) ([]AiModel, error) {
	var all []AiModel
	for page := 1; ; page++ {
		q := url.Values{"page": {fmt.Sprint(page)}, "per_page": {"100"}}
		if f.ProviderUUID != "" {
			q.Set("provider", f.ProviderUUID)
		}
		if f.Search != "" {
			q.Set("search", f.Search)
		}
		var out struct {
			Models []AiModel `json:"models"`
			Total  int64     `json:"total"`
		}
		if err := c.do(ctx, http.MethodGet, "/ai-firewall/models", q, nil, &out); err != nil {
			return nil, err
		}
		all = append(all, out.Models...)
		if len(out.Models) == 0 || int64(len(all)) >= out.Total {
			return all, nil
		}
	}
}

// pagination is the paging block of the list routes (handler paginationOf).
type pagination struct {
	TotalPages int `json:"totalPages"`
}

// listPages reads a paged route at perPage until it runs out or limit items
// have been read (0: no limit). page decodes one answer and returns its items
// and paging block.
func listPages[T any](ctx context.Context, c *Client, path string, q url.Values, perPage, limit int,
	page func(raw json.RawMessage) ([]T, pagination, error)) ([]T, error) {
	var all []T
	for n := 1; ; n++ {
		qq := url.Values{}
		for k, v := range q {
			qq[k] = v
		}
		qq.Set("page", fmt.Sprint(n))
		qq.Set("perPage", fmt.Sprint(perPage))
		var raw json.RawMessage
		if err := c.do(ctx, http.MethodGet, path, qq, nil, &raw); err != nil {
			return nil, err
		}
		items, p, err := page(raw)
		if err != nil {
			return nil, fmt.Errorf("GET %s: decode response: %w", path, err)
		}
		all = append(all, items...)
		if limit > 0 && len(all) >= limit {
			return all[:limit], nil
		}
		if len(items) == 0 || n >= p.TotalPages {
			return all, nil
		}
	}
}

// VexStatement is one VEX statement the tenant holds
// (handler/vex_statements.go vexStatementView).
type VexStatement struct {
	UUID              string          `json:"uuid"`
	DocumentUUID      string          `json:"documentUuid"`
	DocumentID        string          `json:"documentId"`
	Author            string          `json:"author"`
	VulnerabilityID   string          `json:"vulnerabilityId"`
	VulnerabilityName string          `json:"vulnerabilityName"`
	Aliases           []string        `json:"aliases"`
	Status            string          `json:"status"`
	Justification     string          `json:"justification"`
	StatusNotes       string          `json:"statusNotes"`
	ImpactStatement   string          `json:"impactStatement"`
	ActionStatement   string          `json:"actionStatement"`
	Supplier          string          `json:"supplier"`
	Products          json.RawMessage `json:"products"`
	Purls             []string        `json:"purls"`
	// Timestamp and LastUpdated are Unix milliseconds.
	Timestamp          int64  `json:"timestamp"`
	LastUpdated        int64  `json:"lastUpdated"`
	Disposition        string `json:"disposition"`
	SuppressionUUID    string `json:"suppressionUuid"`
	RepositoryFullName string `json:"repositoryFullName"`
	Grounded           bool   `json:"grounded"`
	SourceFormat       string `json:"sourceFormat"`
	MintedBy           string `json:"mintedBySuppressionUuid"`
}

// VexStatementFilter narrows ListVexStatements. Empty fields do not filter.
type VexStatementFilter struct {
	// Status is an OpenVEX status: not_affected, affected, fixed or
	// under_investigation.
	Status string
	// VulnerabilityID matches exactly, ignoring case.
	VulnerabilityID string
	// Purl is a substring of the statement's products.
	Purl string
	// Repository is a repository's owner/name, exactly.
	Repository string
	// Limit stops after this many statements; 0 reads them all.
	Limit int
}

// ListVexStatements reads the tenant's VEX statements, newest first, from
// GET /vex/statements.
func (c *Client) ListVexStatements(ctx context.Context, f VexStatementFilter) ([]VexStatement, error) {
	q := url.Values{}
	for k, v := range map[string]string{"status": f.Status, "vulnerabilityId": f.VulnerabilityID,
		"purl": f.Purl, "repository": f.Repository} {
		if v != "" {
			q.Set(k, v)
		}
	}
	return listPages(ctx, c, "/vex/statements", q, 500, f.Limit, func(raw json.RawMessage) ([]VexStatement, pagination, error) {
		var out struct {
			Statements []VexStatement `json:"statements"`
			Pagination pagination     `json:"pagination"`
		}
		err := unmarshalNumber(raw, &out)
		return out.Statements, out.Pagination, err
	})
}

// GitHubRepo is one repository a GitHub App installation of the tenant
// reaches (handler/github_installation_repos.go githubInstallationRepo).
type GitHubRepo struct {
	InstallationID int64  `json:"installationId"`
	AccountLogin   string `json:"accountLogin"`
	RepoID         int64  `json:"repoId"`
	OwnerLogin     string `json:"ownerLogin"`
	Name           string `json:"name"`
	FullName       string `json:"fullName"`
	Private        bool   `json:"private"`
	Archived       bool   `json:"archived"`
	Host           string `json:"host"`
	// UpdatedAt is Unix milliseconds.
	UpdatedAt int64 `json:"updatedAt"`
}

// GitHubRepoFilter narrows ListGitHubInstallationRepos. Zero fields do not
// filter.
type GitHubRepoFilter struct {
	InstallationID int64
	// Query is a substring of owner/name, ignoring case.
	Query string
}

// ListGitHubInstallationRepos reads every repository the tenant's active
// GitHub App installations reach, ordered by full name, from GET
// /integrations/github/installations/repos.
func (c *Client) ListGitHubInstallationRepos(ctx context.Context, f GitHubRepoFilter) ([]GitHubRepo, error) {
	q := url.Values{}
	if f.InstallationID > 0 {
		q.Set("installationId", fmt.Sprint(f.InstallationID))
	}
	if f.Query != "" {
		q.Set("q", f.Query)
	}
	return listPages(ctx, c, "/integrations/github/installations/repos", q, 500, 0, func(raw json.RawMessage) ([]GitHubRepo, pagination, error) {
		var out struct {
			Repos      []GitHubRepo `json:"repos"`
			Pagination pagination   `json:"pagination"`
		}
		err := unmarshalNumber(raw, &out)
		return out.Repos, out.Pagination, err
	})
}
