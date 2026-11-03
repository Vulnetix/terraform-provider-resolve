package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/vulnetix/terraform-provider-resolve/internal/client"
)

// fakeAPI is an in-memory config API with the engine's behaviour where it
// matters to the provider: PUT merges into the stored spec, the stored
// object can carry fields nobody sent (the row's column defaults), and
// unknown tokens are refused.
type fakeAPI struct {
	t     *testing.T
	token string

	mu       sync.Mutex
	objects  map[string]client.Document // kind/name
	puts     []map[string]any           // every PUT body, in order
	deletes  []string                   // kind/name?manager
	defaults map[string]map[string]any  // kind -> spec fields set on insert
	feeds    []client.Feed
	entries  map[string][]client.FeedEntry
	locked   bool
	models   []client.AiModel
	vex      []map[string]any
	repos    []client.GitHubRepo
	// phases is each object's phase as the last apply left it.
	phases map[string]string
	// forbidden kinds answer 403, as an owner-only kind does to a token without the owner role.
	forbidden map[string]bool
}

func newFakeAPI(t *testing.T) (*fakeAPI, *httptest.Server) {
	f := &fakeAPI{
		t: t, token: "vmt_test",
		objects:   map[string]client.Document{},
		defaults:  map[string]map[string]any{},
		entries:   map[string][]client.FeedEntry{},
		phases:    map[string]string{},
		forbidden: map[string]bool{},
	}
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	return f, srv
}

func (f *fakeAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if r.Header.Get("Authorization") != "Bearer "+f.token {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "not authenticated"})
		return
	}
	p := strings.TrimPrefix(r.URL.Path, "/v1")
	switch {
	case p == "/config/kinds" && r.Method == http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]any{
			"apiVersion": "vulnetix.com/v1", "schemaVersion": "v1.0.1",
			"kinds":       []map[string]any{{"kind": "QualityGate", "schemaUrl": "https://example.test/QualityGate.schema.json", "dependsOn": []string{}, "schema": map[string]any{"title": "QualityGate"}}},
			"schemaKinds": []string{"QualityGate"},
		})
	case p == "/ai-firewall/models" && r.Method == http.MethodGet:
		f.aiModels(w, r)
	case p == "/vex/statements" && r.Method == http.MethodGet:
		f.vexStatements(w, r)
	case p == "/integrations/github/installations/repos" && r.Method == http.MethodGet:
		f.githubRepos(w, r)
	case strings.HasPrefix(p, "/config/objects/"):
		f.object(w, r, strings.TrimPrefix(p, "/config/objects/"))
	case p == "/detection-rule-feeds" && r.Method == http.MethodGet:
		if f.locked {
			writeJSON(w, http.StatusOK, map[string]any{"planLocked": true, "requiredPlan": "pro", "feeds": []any{}})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"feeds": f.feeds})
	case strings.HasPrefix(p, "/detection-rule-feeds/") && r.Method == http.MethodGet:
		id := strings.TrimPrefix(p, "/detection-rule-feeds/")
		for _, fd := range f.feeds {
			if fd.UUID == id {
				writeJSON(w, http.StatusOK, map[string]any{"feed": fd, "entries": f.entries[id]})
				return
			}
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	default:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no route " + r.Method + " " + p})
	}
}

func (f *fakeAPI) object(w http.ResponseWriter, r *http.Request, key string) {
	switch r.Method {
	case http.MethodGet:
		if f.forbidden[strings.SplitN(key, "/", 2)[0]] {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden", "message": "only an owner token may apply this kind"})
			return
		}
		d, ok := f.objects[key]
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"document": redacted(d), "status": map[string]any{"phase": f.phaseOf(key), "manager": "terraform"}})
	case http.MethodPut:
		if f.forbidden[strings.SplitN(key, "/", 2)[0]] {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden", "message": "only an owner token may apply this kind"})
			return
		}
		var body struct {
			Document client.Document `json:"document"`
			Manager  string          `json:"manager"`
		}
		raw := map[string]any{}
		dec := json.NewDecoder(r.Body)
		dec.UseNumber()
		if err := dec.Decode(&raw); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
			return
		}
		b, _ := json.Marshal(raw)
		_ = json.Unmarshal(b, &body)
		f.puts = append(f.puts, raw)
		if body.Document.Kind+"/"+body.Document.Metadata.Name != key {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "document does not match the path"})
			return
		}
		cur, exists := f.objects[key]
		if !exists {
			cur = client.Document{APIVersion: "vulnetix.com/v1", Kind: body.Document.Kind, Metadata: body.Document.Metadata, Spec: map[string]any{}}
			for k, v := range f.defaults[body.Document.Kind] {
				cur.Spec[k] = v
			}
		}
		cur.Metadata = body.Document.Metadata
		if mergeKinds[body.Document.Kind] {
			mergeSpec(cur.Spec, body.Document.Spec)
		} else {
			cur.Spec = body.Document.Spec
		}
		status := map[string]any{"phase": "Applied", "manager": body.Manager}
		if ext := f.applyStatus(body.Document, !exists); ext != nil {
			if ext["phase"] == "Failed" {
				writeJSON(w, http.StatusBadGateway, map[string]string{"error": "external_failed", "message": key + ": " + fmt.Sprint(ext["message"])})
				return
			}
			for k, v := range ext {
				status[k] = v
			}
		}
		f.objects[key] = cur
		f.phases[key] = fmt.Sprint(status["phase"])
		writeJSON(w, http.StatusOK, map[string]any{"document": redacted(cur), "action": "update", "status": status})
	case http.MethodDelete:
		f.deletes = append(f.deletes, key+"?manager="+r.URL.Query().Get("manager"))
		if _, ok := f.objects[key]; !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		delete(f.objects, key)
		w.WriteHeader(http.StatusNoContent)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method"})
	}
}

// mergeSpec is the engine's merge: a field the document leaves out keeps its
// current value.
func mergeSpec(dst, src map[string]any) {
	for k, v := range src {
		if sm, ok := v.(map[string]any); ok {
			if dm, ok := dst[k].(map[string]any); ok {
				mergeSpec(dm, sm)
				continue
			}
		}
		dst[k] = v
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (f *fakeAPI) lastPut() map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.puts) == 0 {
		return nil
	}
	return f.puts[len(f.puts)-1]
}

func (f *fakeAPI) putCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.puts)
}

// testFactories serves the generated kinds plus the test kinds.
func testFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"vulnetix": providerserver.NewProtocol6WithError(&vulnetixProvider{
			version: "test",
			kinds:   allTestKinds(),
		}),
	}
}

func providerBlock(url string) string {
	return `
provider "vulnetix" {
  endpoint = "` + url + `/v1"
  token    = "vmt_test"
}
`
}

// aiModels pages like AiFirewallModelsList: page, per_page, provider, search.
func (f *fakeAPI) aiModels(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var match []client.AiModel
	for _, m := range f.models {
		if p := q.Get("provider"); p != "" && m.ProviderUUID != p {
			continue
		}
		if s := strings.ToLower(q.Get("search")); s != "" && !strings.Contains(strings.ToLower(m.Slug+" "+m.Name), s) {
			continue
		}
		match = append(match, m)
	}
	page, per := atoiOr(q.Get("page"), 1), atoiOr(q.Get("per_page"), 25)
	from, to := min((page-1)*per, len(match)), min(page*per, len(match))
	writeJSON(w, http.StatusOK, map[string]any{"models": match[from:to], "total": len(match), "page": page, "perPage": per})
}

func atoiOr(s string, def int) int {
	if n, err := strconv.Atoi(s); err == nil && n > 0 {
		return n
	}
	return def
}

// redacted is a document as the engine answers it (config redactSpec): every
// secret replaced by a valueFrom placeholder, so a secret is never read back.
func redacted(d client.Document) client.Document {
	b, _ := json.Marshal(d)
	var out client.Document
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.UseNumber()
	_ = dec.Decode(&out)
	for _, k := range allTestKinds() {
		if k.Kind == d.Kind {
			redactFields(out.Spec, k.Fields, "VULNETIX_"+strings.ToUpper(k.TypeName))
		}
	}
	return out
}

func redactFields(spec map[string]any, fields []field, env string) {
	for _, f := range fields {
		v, ok := spec[f.JSON]
		if !ok {
			continue
		}
		switch {
		case f.WriteOnly:
			spec[f.JSON] = map[string]any{"valueFrom": map[string]any{"env": env + "_" + strings.ToUpper(f.JSON)}}
		case f.Kind == fieldObject:
			if m, ok := v.(map[string]any); ok {
				redactFields(m, f.Children, env+"_"+strings.ToUpper(f.JSON))
			}
		}
	}
}

func allTestKinds() []kindSpec {
	return append(append(append([]kindSpec{}, generatedKinds...), generatedKindsUnserved...), testKinds...)
}

// mergeKinds are the kinds whose engine adapter merges a document into the
// stored object (MergeDiff: a field the document leaves out keeps its
// value). Every other kind replaces its spec on apply, as RecordKind and the
// wholly declared adapters do.
var mergeKinds = map[string]bool{
	"AiFirewallPolicy": true, "AiProvider": true, "AwsAccount": true, "AwsConnection": true,
	"CloudSensors": true, "JailExemption": true, "JailPolicy": true, "Organization": true,
	"PackageFirewallPolicy": true, "PackageList": true, "QualityGate": true,
}

func (f *fakeAPI) phaseOf(key string) string {
	if p := f.phases[key]; p != "" {
		return p
	}
	return "Applied"
}

// siemKinds are written in Guardian, after the config database.
var siemKinds = map[string]bool{"SavedSearch": true, "DetectionRule": true, "AlertTrigger": true, "ScheduledReport": true, "LogSource": true}

// applyStatus is what the engine adds to an apply's status for the kinds that
// report more than Applied: a member's invitation is Pending, a messaging
// connection waits for a person (AwaitingAuthorization, with its URL), a
// service account's credentials and a log source's ingest token are shown
// once, and a SIEM kind named refuse-* is refused by Guardian.
func (f *fakeAPI) applyStatus(d client.Document, created bool) map[string]any {
	switch {
	case siemKinds[d.Kind] && strings.HasPrefix(d.Metadata.Name, "refuse-"):
		return map[string]any{"phase": "Failed", "message": "Guardian refused the " + d.Kind}
	case d.Kind == "Member":
		return map[string]any{"phase": "Pending", "invitationExpires": "2026-12-01T00:00:00Z"}
	case d.Kind == "MessagingConnection" && created:
		return map[string]any{"phase": "AwaitingAuthorization", "authorizationUrl": "https://auth.example.test/slack?c=" + d.Metadata.Name}
	case d.Kind == "ServiceAccount" && created:
		return map[string]any{"principal": "sa-" + d.Metadata.Name, "apiKey": "key-" + d.Metadata.Name, "secret": "shh-" + d.Metadata.Name}
	case d.Kind == "LogSource" && created:
		return map[string]any{"ingestToken": "gdn_" + d.Metadata.Name}
	}
	return nil
}

// page cuts one page like the vdb-site list routes (paginationOf).
func page[T any](q url.Values, items []T) ([]T, map[string]int) {
	p, per := atoiOr(q.Get("page"), 1), atoiOr(q.Get("perPage"), 100)
	from, to := min((p-1)*per, len(items)), min(p*per, len(items))
	return items[from:to], map[string]int{"page": p, "perPage": per, "totalRows": len(items), "totalPages": (len(items) + per - 1) / per}
}

// vexStatements filters like VexStatementsList: status, vulnerabilityId.
func (f *fakeAPI) vexStatements(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var match []map[string]any
	for _, s := range f.vex {
		if v := q.Get("status"); v != "" && s["status"] != v {
			continue
		}
		if v := q.Get("vulnerabilityId"); v != "" && !strings.EqualFold(fmt.Sprint(s["vulnerabilityId"]), v) {
			continue
		}
		match = append(match, s)
	}
	items, p := page(q, match)
	writeJSON(w, http.StatusOK, map[string]any{"statements": items, "pagination": p})
}

// githubRepos filters like GitHubInstallationReposList: installationId, q.
func (f *fakeAPI) githubRepos(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var match []client.GitHubRepo
	for _, g := range f.repos {
		if v := q.Get("installationId"); v != "" && fmt.Sprint(g.InstallationID) != v {
			continue
		}
		if v := strings.ToLower(q.Get("q")); v != "" && !strings.Contains(strings.ToLower(g.FullName), v) {
			continue
		}
		match = append(match, g)
	}
	items, p := page(q, match)
	writeJSON(w, http.StatusOK, map[string]any{"repos": items, "pagination": p})
}
