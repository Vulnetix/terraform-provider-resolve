// Package client talks to the Vulnetix config API (vdb-site /v1).
//
// Every request and response shape the provider depends on is in this one
// file. The single-object routes (/config/objects/{kind}/{name}) were designed
// alongside the provider, so if the server's contract moves, this is the only
// file that should need to follow it.
//
//	GET    {endpoint}/config/kinds                     kinds the server applies
//	GET    {endpoint}/config/objects/{kind}/{name}     one object, 404 when absent
//	PUT    {endpoint}/config/objects/{kind}/{name}     {"document": doc, "manager": m}
//	DELETE {endpoint}/config/objects/{kind}/{name}?manager=m
//	GET    {endpoint}/detection-rule-feeds             the tenant's feeds
//	GET    {endpoint}/detection-rule-feeds/{uuid}      one feed with its entries
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

// TokenPrefix marks a config management token.
const TokenPrefix = "vmt_"

// APIVersion is the document format every kind uses today.
const APIVersion = "vulnetix.com/v1"

// ErrNotFound is returned when the object or feed does not exist.
var ErrNotFound = errors.New("not found")

// Client is a config API client. The zero value is not usable; build one
// with New.
type Client struct {
	endpoint  string
	token     string
	manager   string
	userAgent string
	http      *http.Client
}

// Options configure New.
type Options struct {
	Endpoint  string
	Token     string
	Manager   string
	UserAgent string
	HTTP      *http.Client
}

// New builds a client. Endpoint and Manager fall back to their defaults.
func New(o Options) *Client {
	c := &Client{
		endpoint:  strings.TrimRight(o.Endpoint, "/"),
		token:     o.Token,
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

// Object is what the single-object routes return.
type Object struct {
	Document Document `json:"document"`
	Manager  string   `json:"manager,omitempty"`
	UUID     string   `json:"uuid,omitempty"`
}

// APIError is a non-2xx answer. Message is the server's "error" (and
// "message", when present) field, or the raw body.
type APIError struct {
	Method, URL string
	Status      int
	Message     string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s %s: HTTP %d: %s", e.Method, e.URL, e.Status, e.Message)
}

// authorization is the Authorization header for the configured token. A
// vmt_ token is a bearer token. Anything else is taken to be the existing
// org API key, "<orgId>:<hmac>", which orgIDFromRequest already accepts, so
// the provider can be exercised before management tokens ship.
func (c *Client) authorization() string {
	if c.token == "" {
		return ""
	}
	if !strings.HasPrefix(c.token, TokenPrefix) && strings.Contains(c.token, ":") {
		return "ApiKey " + c.token
	}
	return "Bearer " + c.token
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	u := c.endpoint + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, rd)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
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
		return err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{Method: method, URL: u, Status: resp.StatusCode, Message: errorMessage(raw)}
	}
	if out == nil || len(bytes.TrimSpace(raw)) == 0 {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	// Numbers stay json.Number so an integer spec field never passes through
	// float64 on its way into Terraform state.
	dec.UseNumber()
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("%s %s: decode response: %w", method, u, err)
	}
	return nil
}

func errorMessage(raw []byte) string {
	var e struct {
		Error    string `json:"error"`
		Message  string `json:"message"`
		Problems []struct {
			Path    string `json:"path"`
			Message string `json:"message"`
		} `json:"problems"`
	}
	if json.Unmarshal(raw, &e) != nil || e.Error == "" {
		return strings.TrimSpace(string(raw))
	}
	msg := e.Error
	if e.Message != "" {
		msg += ": " + e.Message
	}
	for _, p := range e.Problems {
		msg += "\n  " + strings.TrimSpace(p.Path+" "+p.Message)
	}
	return msg
}

func objectPath(kind, name string) string {
	return "/config/objects/" + url.PathEscape(kind) + "/" + url.PathEscape(name)
}

// decodeObject accepts either {"document": doc, ...} or a bare document, so
// the provider keeps working whichever the server settles on.
func decodeObject(raw map[string]json.RawMessage) (*Object, error) {
	var o Object
	if d, ok := raw["document"]; ok {
		if err := unmarshalNumber(d, &o.Document); err != nil {
			return nil, err
		}
		if m, ok := raw["manager"]; ok {
			_ = json.Unmarshal(m, &o.Manager)
		}
		if u, ok := raw["uuid"]; ok {
			_ = json.Unmarshal(u, &o.UUID)
		}
		return &o, nil
	}
	b, _ := json.Marshal(raw)
	if err := unmarshalNumber(b, &o.Document); err != nil {
		return nil, err
	}
	return &o, nil
}

func unmarshalNumber(b []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(v)
}

// GetObject reads one object. It returns ErrNotFound when the tenant has none.
func (c *Client) GetObject(ctx context.Context, kind, name string) (*Object, error) {
	var raw map[string]json.RawMessage
	if err := c.do(ctx, http.MethodGet, objectPath(kind, name), nil, nil, &raw); err != nil {
		return nil, err
	}
	return decodeObject(raw)
}

// PutObject applies one document as this client's manager and returns the
// object as stored. An empty response body is not an error: the caller then
// keeps what it sent.
func (c *Client) PutObject(ctx context.Context, doc Document) (*Object, error) {
	if doc.APIVersion == "" {
		doc.APIVersion = APIVersion
	}
	if doc.Spec == nil {
		doc.Spec = map[string]any{}
	}
	body := map[string]any{"document": doc, "manager": c.manager}
	var raw map[string]json.RawMessage
	if err := c.do(ctx, http.MethodPut, objectPath(doc.Kind, doc.Metadata.Name), nil, body, &raw); err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, nil
	}
	return decodeObject(raw)
}

// DeleteObject deletes one object. Deleting an absent object succeeds.
func (c *Client) DeleteObject(ctx context.Context, kind, name string) error {
	err := c.do(ctx, http.MethodDelete, objectPath(kind, name), url.Values{"manager": {c.manager}}, nil, nil)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	return err
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
