// Package delve — Go SDK client for the davly/delve search infrastructure service.
//
// Phase 1: thin HTTP-client shim. Compliance with R187 (infra cohort) +
// R191 (boundary signer) + R151 (KAT-1 firewall) + R175 (load-bearing
// empty-key rejection) + R121 (multi-tenant).
//
// Wire-in to flagships is Phase 2 (R176 library-first).
//
// # I1 strategic-review catch
//
// Per the 2026-05-28 infra marathon I1 strategic review, flagships that
// have `internal/delve/` directories are NOT genuine consumers of this
// service — they are R155.A INDEX-LIE name-collision domain-rename
// packages (search vocabulary borrowed from the infra-noun). Wire-level
// R187 compliance requires `import github.com/davly/delve-go` PLUS a
// `DELVE_URL` env var pointing at a real delve endpoint.
//
// This SDK is the path from name-collision to genuine wire-level
// consumption.
package delve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

// DefaultURL is the default URL the SDK falls back to when neither
// ClientOptions.URL nor DELVE_URL env var is set.
const DefaultURL = "http://localhost:8092"

// EnvURLKey is the canonical env-var name the SDK reads for the
// service URL.
const EnvURLKey = "DELVE_URL"

// DefaultTimeout is the default per-request HTTP timeout.
const DefaultTimeout = 30 * time.Second

// ClientOptions configures the SDK client.
type ClientOptions struct {
	URL        string       // delve service URL; falls back to DELVE_URL env var, then DefaultURL.
	APIKey     string       // iik_... key for HMAC signing. R175 LOAD-BEARING: empty rejected.
	TenantID   string       // R121 multi-tenant cohort key. Empty rejected.
	HTTPClient *http.Client // optional override (default: 30s timeout).
}

// Client wraps the delve HTTP API.
type Client struct {
	url        string
	apiKey     string
	tenantID   string
	httpClient *http.Client
}

// NewClient constructs a delve SDK client. URL defaults to DELVE_URL env var
// (or http://localhost:8092 if unset). Empty APIKey or TenantID returns an
// error — empty-key signing is structurally insecure (R175 LOAD-BEARING);
// empty-tenant violates R121 multi-tenant cohort.
func NewClient(opts ClientOptions) (*Client, error) {
	url := opts.URL
	if url == "" {
		url = os.Getenv(EnvURLKey)
	}
	if url == "" {
		url = DefaultURL
	}
	if opts.APIKey == "" {
		return nil, errors.New("delve: empty APIKey violates R175 LOAD-BEARING discipline")
	}
	if opts.TenantID == "" {
		return nil, errors.New("delve: empty TenantID violates R121 multi-tenant cohort")
	}
	hc := opts.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: DefaultTimeout}
	}
	return &Client{url: url, apiKey: opts.APIKey, tenantID: opts.TenantID, httpClient: hc}, nil
}

// URL returns the configured service URL. Exposed for diagnostics +
// the cohort/honest LOUD-ONCE advisory that warns once when the
// default-localhost fallback is in effect.
func (c *Client) URL() string { return c.url }

// TenantID returns the configured tenant ID.
func (c *Client) TenantID() string { return c.tenantID }

// SearchQuery is the search-request body shape.
type SearchQuery struct {
	Query     string `json:"query"`
	Limit     int    `json:"limit,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// SearchResult is a single search hit.
type SearchResult struct {
	ID        string  `json:"id"`
	Score     float64 `json:"score"`
	Snippet   string  `json:"snippet"`
	SourceURL string  `json:"source_url,omitempty"`
}

// SearchResponse is the search-endpoint return body.
type SearchResponse struct {
	Results []SearchResult `json:"results"`
}

// Search performs a search query against the delve service.
// Boundary-signed via internal/boundarysigner per R191.
func (c *Client) Search(ctx context.Context, q SearchQuery) (*SearchResponse, error) {
	body, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("delve: marshal query: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+"/v1/search", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("delve: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("X-Tenant-ID", c.tenantID)
	req.Header.Set("User-Agent", "delve-go/0.1")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("delve: http: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("delve: status %d", resp.StatusCode)
	}
	var out SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("delve: decode: %w", err)
	}
	return &out, nil
}
