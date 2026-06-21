package boundarysigner

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

// TestSign_GoldenWireBytes pins the EXACT wire-byte output of Sign by independently re-deriving the
// L43 Mirror-Mark v1 formula from the stdlib:
//
//	"lore@v1:" + base64url-no-pad( HMAC-SHA256(key, 0x01 || corpusSHA || payload) )
//
// Any drift in the on-wire format (prefix, version tag, HMAC input order, base64 alphabet/padding) fails
// this test. It is the cross-language parity anchor: every -go SDK port — and the future server-side
// verifier — must produce byte-identical marks for the same (key, corpus, payload).
func TestSign_GoldenWireBytes(t *testing.T) {
	key := []byte("kat-1-golden-key")
	var corpus [corpusSHALen]byte // zero corpus = the KAT-1 canonical input
	payload := []byte("hello")

	s, err := NewSigner(key, corpus)
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	got := s.Sign(payload)

	// Independent re-derivation via the stdlib (not the package's own helpers).
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte{0x01}) // versionTag
	mac.Write(corpus[:])
	mac.Write(payload)
	want := "lore@v1:" + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if got != want {
		t.Fatalf("wire-byte golden drift:\n got  %q\n want %q", got, want)
	}
	if len(got) != 51 {
		t.Fatalf("expected a 51-char mark, got %d (%q)", len(got), got)
	}
	// Sanity: the independently-derived mark must verify through the package's own Verify.
	if !s.Verify(payload, got) {
		t.Fatal("golden mark does not verify through Signer.Verify")
	}
}
