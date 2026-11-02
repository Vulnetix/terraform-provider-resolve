package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthorizationHeader(t *testing.T) {
	for token, want := range map[string]string{
		"vmt_abc":         "Bearer vmt_abc",
		"eyJhbGciOi.jwt":  "Bearer eyJhbGciOi.jwt",
		"":                "",
	} {
		if got := New(Options{Token: token}).authorization(); got != want {
			t.Errorf("token %q: Authorization = %q, want %q", token, got, want)
		}
	}
}

func TestDefaults(t *testing.T) {
	c := New(Options{})
	if c.endpoint != DefaultEndpoint || c.Manager() != DefaultManager {
		t.Fatalf("endpoint %q manager %q", c.endpoint, c.Manager())
	}
	if New(Options{Endpoint: "https://x.test/v1/"}).endpoint != "https://x.test/v1" {
		t.Fatal("trailing slash kept")
	}
}

func TestObjectRoundTrip(t *testing.T) {
	var gotPath, gotBody, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotQuery = r.URL.EscapedPath(), r.URL.RawQuery
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		switch r.Method {
		case http.MethodGet:
			doc := `{"apiVersion":"vulnetix.com/v1","kind":"QualityGate","metadata":{"name":"default"},"spec":{"cooldownDays":9007199254740993}}`
			io.WriteString(w, `{"document":`+doc+`,"status":{"phase":"Applied","manager":"gitops:acme/config","objectUuid":"u1"}}`)
		case http.MethodPut:
			io.WriteString(w, `{"document":{"apiVersion":"vulnetix.com/v1","kind":"QualityGate","metadata":{"name":"default"},"spec":{}}}`)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	c := New(Options{Endpoint: srv.URL, Token: "vmt_x", Manager: "terraform"})
	ctx := context.Background()

	o, err := c.GetObject(ctx, "QualityGate", "default")
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/config/objects/QualityGate/default" || o.Status.Manager != "gitops:acme/config" || o.Status.ObjectUUID != "u1" {
		t.Fatalf("path %s object %+v", gotPath, o)
	}
	// Integers survive above 2^53.
	if n, _ := o.Document.Spec["cooldownDays"].(json.Number); n.String() != "9007199254740993" {
		t.Fatalf("cooldownDays = %v", o.Document.Spec["cooldownDays"])
	}

	_, err = c.PutObject(ctx, Document{Kind: "QualityGate", Metadata: Metadata{Name: "default"}, Spec: map[string]any{"severity": "high"}})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	_ = json.Unmarshal([]byte(gotBody), &body)
	doc := body["document"].(map[string]any)
	if body["manager"] != "terraform" || doc["apiVersion"] != APIVersion || doc["spec"].(map[string]any)["severity"] != "high" {
		t.Fatalf("PUT body %s", gotBody)
	}

	if err := c.DeleteObject(ctx, "Widget", "a b"); err != nil {
		t.Fatalf("deleting an absent object: %v", err)
	}
	if gotPath != "/config/objects/Widget/a%20b" || gotQuery != "manager=terraform" {
		t.Fatalf("DELETE %s?%s", gotPath, gotQuery)
	}
}

func TestErrors(t *testing.T) {
	status, body := 0, ""
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		io.WriteString(w, body)
	}))
	defer srv.Close()
	c := New(Options{Endpoint: srv.URL})

	status, body = http.StatusNotFound, `{"error":"not found"}`
	if _, err := c.GetObject(context.Background(), "QualityGate", "default"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("404: %v", err)
	}

	status, body = http.StatusUnprocessableEntity, `{"error":"plan_blocked","message":"owned by gitops:acme/config","problems":[{"index":0,"kind":"QualityGate","name":"default","message":"spec.severity: bad"}]}`
	_, err := c.PutObject(context.Background(), Document{Kind: "QualityGate", Metadata: Metadata{Name: "default"}})
	var ae *APIError
	if !errors.As(err, &ae) || ae.Status != 422 || !strings.Contains(ae.Message, "plan_blocked: owned by gitops:acme/config") || !strings.Contains(ae.Message, "QualityGate/default spec.severity: bad") {
		t.Fatalf("422: %v", err)
	}

	status, body = http.StatusBadGateway, "upstream down"
	if _, err := c.ListKinds(context.Background()); err == nil || !strings.Contains(err.Error(), "HTTP 502: upstream down") {
		t.Fatalf("502: %v", err)
	}

	status, body = http.StatusOK, `{"planLocked":true,"requiredPlan":"pro","feeds":[]}`
	if _, err := c.ListFeeds(context.Background()); !errors.Is(err, ErrPlanLocked) {
		t.Fatalf("plan locked: %v", err)
	}
}
