package delve

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/davly/delve-go/internal/boundarysigner"
)

// TestSearch_NoBoundarySignatureByDefault is the characterization test that PROVES the R191 wiring is
// byte-neutral by default: a client built without a BoundarySignerKey must send exactly the same headers as
// before — no X-Boundary-Signature. This freezes the existing wire contract so the opt-in wiring cannot
// silently change bytes-on-wire for existing callers.
func TestSearch_NoBoundarySignatureByDefault(t *testing.T) {
	var sig string
	var present bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sig = r.Header.Get("X-Boundary-Signature")
		_, present = r.Header["X-Boundary-Signature"]
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(SearchResponse{Results: []SearchResult{}})
	}))
	defer srv.Close()

	c, err := NewClient(ClientOptions{URL: srv.URL, APIKey: "iik_test", TenantID: "t1"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Search(context.Background(), SearchQuery{Query: "hello"}); err != nil {
		t.Fatal(err)
	}
	if present || sig != "" {
		t.Fatalf("default client must NOT send X-Boundary-Signature (got %q, present=%v) — opt-in wiring must be byte-neutral", sig, present)
	}
}

// TestSearch_SignsAndVerifiesWhenKeySet proves that when a key IS configured the request carries a valid
// L43 mark that round-trips through an independent verifier over the exact body the server received.
func TestSearch_SignsAndVerifiesWhenKeySet(t *testing.T) {
	key := []byte("test-boundary-key")
	var sig string
	var body []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sig = r.Header.Get("X-Boundary-Signature")
		body, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(SearchResponse{Results: []SearchResult{}})
	}))
	defer srv.Close()

	c, err := NewClient(ClientOptions{URL: srv.URL, APIKey: "iik_test", TenantID: "t1", BoundarySignerKey: key})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Search(context.Background(), SearchQuery{Query: "hello"}); err != nil {
		t.Fatal(err)
	}
	if len(sig) != 51 || !strings.HasPrefix(sig, "lore@v1:") {
		t.Fatalf("expected a 51-char lore@v1: mark, got %q (len %d)", sig, len(sig))
	}
	v, _ := boundarysigner.NewSigner(key, [32]byte{})
	if !v.Verify(body, sig) {
		t.Fatal("the X-Boundary-Signature the server received does not verify against the request body")
	}
}
