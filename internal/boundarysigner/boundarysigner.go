// Package boundarysigner implements the L43 Mirror-Mark v1 boundary
// signature for delve-go HTTP requests.
//
// R191 R-CROSS-INFRA-AUDIT-CHAIN-EMIT compliance: every outbound
// request is signed BEFORE emit. Empty key is rejected at construction
// — there is no degraded-signing fall-back (R175 LOAD-BEARING).
//
// Wire shape: "lore@v1:" || base64url( HMAC-SHA256(key, 0x01 || corpusSHA(32B) || payload) )
//
// This is the same primitive as foundation/pkg/mirrormark and as
// limitless-aiwatermark's internal/mirrormark — byte-identical.
package boundarysigner

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

const (
	prefix       = "lore@v1:"
	versionTag   = byte(0x01)
	corpusSHALen = sha256.Size
	digestLen    = sha256.Size
)

// ErrEmptyKey is returned by NewSigner when a zero-length key is passed.
// Empty-key HMAC is structurally insecure — the signature reduces to a
// keyed digest of the payload with a fixed (empty) key, which is
// cohort-public.
var ErrEmptyKey = errors.New("boundarysigner: empty signing key (R175 LOAD-BEARING)")

// Signer produces L43 Mirror-Mark v1 boundary signatures.
type Signer struct {
	key       []byte
	corpusSHA [corpusSHALen]byte
}

// NewSigner constructs a Signer from a non-empty HMAC key and a 32-byte
// corpus SHA. Returns ErrEmptyKey for a zero-length key.
//
// corpusSHA = the cohort/project corpus identifier. For Phase 1, callers
// MAY pass the zero corpus (32 zero bytes) — this matches the KAT-1
// canonical input. Production deployments override with a real corpus.
func NewSigner(key []byte, corpusSHA [corpusSHALen]byte) (*Signer, error) {
	if len(key) == 0 {
		return nil, ErrEmptyKey
	}
	// Defensive copy of key so the caller can zero theirs after.
	k := make([]byte, len(key))
	copy(k, key)
	return &Signer{key: k, corpusSHA: corpusSHA}, nil
}

// Sign returns the canonical L43 Mirror-Mark v1 wire string for payload.
// Length: PrefixLen(8) + base64url-no-pad(32B HMAC) = 8 + 43 = 51 chars.
func (s *Signer) Sign(payload []byte) string {
	mac := hmac.New(sha256.New, s.key)
	_ = writeAll(mac, []byte{versionTag})
	_ = writeAll(mac, s.corpusSHA[:])
	_ = writeAll(mac, payload)
	body := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return prefix + body
}

// Verify returns true iff mark is a syntactically valid L43 v1 mark and
// its HMAC body matches HMAC-SHA256(key, 0x01||corpus||payload).
// Constant-time comparison.
func (s *Signer) Verify(payload []byte, mark string) bool {
	if len(mark) <= len(prefix) {
		return false
	}
	if mark[:len(prefix)] != prefix {
		return false
	}
	got, err := base64.RawURLEncoding.DecodeString(mark[len(prefix):])
	if err != nil {
		return false
	}
	if len(got) != digestLen {
		return false
	}
	mac := hmac.New(sha256.New, s.key)
	_ = writeAll(mac, []byte{versionTag})
	_ = writeAll(mac, s.corpusSHA[:])
	_ = writeAll(mac, payload)
	want := mac.Sum(nil)
	return hmac.Equal(got, want)
}

// writeAll writes p to w, treating short-writes as errors. (hash.Hash
// never short-writes but this keeps the interface consistent.)
func writeAll(w interface {
	Write([]byte) (int, error)
}, p []byte) error {
	n, err := w.Write(p)
	if err != nil {
		return err
	}
	if n != len(p) {
		return errors.New("boundarysigner: short write")
	}
	return nil
}
