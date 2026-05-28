package delve

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestNewClient_EmptyAPIKey_R175 asserts the cohort R175 LOAD-BEARING
// firewall: empty key MUST be rejected at construction time. Empty-key
// HMAC signing is structurally insecure; deferring the failure to first
// request is too late.
func TestNewClient_EmptyAPIKey_R175(t *testing.T) {
	_, err := NewClient(ClientOptions{
		URL:      "http://localhost:8092",
		APIKey:   "",
		TenantID: "tenant-1",
	})
	if err == nil {
		t.Fatal("R175 LOAD-BEARING: empty APIKey must be rejected at NewClient")
	}
}

// TestNewClient_EmptyTenantID_R121 asserts the cohort R121 multi-tenant
// rule: empty tenant ID MUST be rejected at construction time. A
// rogue-empty-tenant request leaks data across tenant boundaries.
func TestNewClient_EmptyTenantID_R121(t *testing.T) {
	_, err := NewClient(ClientOptions{
		URL:      "http://localhost:8092",
		APIKey:   "iik_test_abc123",
		TenantID: "",
	})
	if err == nil {
		t.Fatal("R121 multi-tenant: empty TenantID must be rejected at NewClient")
	}
}

// TestNewClient_DefaultURL_FromEnv verifies the SDK reads DELVE_URL.
func TestNewClient_DefaultURL_FromEnv(t *testing.T) {
	t.Setenv(EnvURLKey, "http://example.test:8092")
	c, err := NewClient(ClientOptions{
		APIKey:   "iik_test_abc123",
		TenantID: "tenant-1",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if c.URL() != "http://example.test:8092" {
		t.Errorf("URL: got %q want %q", c.URL(), "http://example.test:8092")
	}
}

// TestNewClient_DefaultURL_Localhost verifies the SDK falls back to
// DefaultURL (localhost:8092) when neither URL nor env is set.
func TestNewClient_DefaultURL_Localhost(t *testing.T) {
	// Clear env in case the test harness inherited it.
	if err := os.Unsetenv(EnvURLKey); err != nil {
		t.Fatalf("unset env: %v", err)
	}
	c, err := NewClient(ClientOptions{
		APIKey:   "iik_test_abc123",
		TenantID: "tenant-1",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if c.URL() != DefaultURL {
		t.Errorf("URL: got %q want %q", c.URL(), DefaultURL)
	}
}

// TestNewClient_ExplicitURLOverridesEnv verifies an explicit
// ClientOptions.URL beats the env var.
func TestNewClient_ExplicitURLOverridesEnv(t *testing.T) {
	t.Setenv(EnvURLKey, "http://envval.test:8092")
	c, err := NewClient(ClientOptions{
		URL:      "http://explicit.test:8092",
		APIKey:   "iik_test_abc123",
		TenantID: "tenant-1",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if c.URL() != "http://explicit.test:8092" {
		t.Errorf("URL: got %q want %q", c.URL(), "http://explicit.test:8092")
	}
}

// TestNewClient_TenantIDExposed verifies TenantID() returns the
// configured tenant for diagnostics.
func TestNewClient_TenantIDExposed(t *testing.T) {
	c, err := NewClient(ClientOptions{
		URL:      "http://localhost:8092",
		APIKey:   "iik_test_abc123",
		TenantID: "tenant-7",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if c.TenantID() != "tenant-7" {
		t.Errorf("TenantID: got %q want %q", c.TenantID(), "tenant-7")
	}
}

// TestSearch_HappyPath drives a successful Search call against a
// httptest.NewServer. Asserts headers (Authorization, X-Tenant-ID,
// Content-Type) + decoded body.
func TestSearch_HappyPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method: got %s want POST", r.Method)
		}
		if r.URL.Path != "/v1/search" {
			t.Errorf("path: got %s want /v1/search", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer iik_test_abc123" {
			t.Errorf("Authorization header missing or wrong: %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-Tenant-ID") != "tenant-1" {
			t.Errorf("X-Tenant-ID header missing or wrong: %q", r.Header.Get("X-Tenant-ID"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type header: %q", r.Header.Get("Content-Type"))
		}
		var got SearchQuery
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if got.Query != "hello" {
			t.Errorf("query: got %q want %q", got.Query, "hello")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(SearchResponse{
			Results: []SearchResult{
				{ID: "r1", Score: 0.95, Snippet: "hello world"},
			},
		})
	}))
	defer srv.Close()
	c, err := NewClient(ClientOptions{
		URL:      srv.URL,
		APIKey:   "iik_test_abc123",
		TenantID: "tenant-1",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	resp, err := c.Search(context.Background(), SearchQuery{Query: "hello", Limit: 10})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("Results: got %d want 1", len(resp.Results))
	}
	if resp.Results[0].ID != "r1" {
		t.Errorf("Result.ID: got %q want %q", resp.Results[0].ID, "r1")
	}
}

// TestSearch_ServerError verifies non-200 statuses propagate as errors.
func TestSearch_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()
	c, _ := NewClient(ClientOptions{
		URL:      srv.URL,
		APIKey:   "iik_test_abc123",
		TenantID: "tenant-1",
	})
	_, err := c.Search(context.Background(), SearchQuery{Query: "x"})
	if err == nil {
		t.Fatal("Search should error on 500")
	}
}

// TestSearch_NetworkError verifies network failures surface as errors.
func TestSearch_NetworkError(t *testing.T) {
	c, _ := NewClient(ClientOptions{
		URL:      "http://127.0.0.1:1", // unbound port
		APIKey:   "iik_test_abc123",
		TenantID: "tenant-1",
	})
	_, err := c.Search(context.Background(), SearchQuery{Query: "x"})
	if err == nil {
		t.Fatal("Search should error on unbound port")
	}
}

// TestDefaultURLConst asserts the canonical port mapping.
func TestDefaultURLConst(t *testing.T) {
	if DefaultURL != "http://localhost:8092" {
		t.Errorf("DefaultURL drift: got %q want %q", DefaultURL, "http://localhost:8092")
	}
}

// TestEnvURLKeyConst asserts the canonical env var name.
func TestEnvURLKeyConst(t *testing.T) {
	if EnvURLKey != "DELVE_URL" {
		t.Errorf("EnvURLKey drift: got %q want %q", EnvURLKey, "DELVE_URL")
	}
}
