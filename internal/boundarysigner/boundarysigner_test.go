package boundarysigner

import (
	"crypto/sha256"
	"errors"
	"strings"
	"testing"
)

func TestNewSigner_EmptyKey_Rejected(t *testing.T) {
	_, err := NewSigner(nil, [sha256.Size]byte{})
	if !errors.Is(err, ErrEmptyKey) {
		t.Fatalf("expected ErrEmptyKey, got %v", err)
	}
	_, err = NewSigner([]byte{}, [sha256.Size]byte{})
	if !errors.Is(err, ErrEmptyKey) {
		t.Fatalf("expected ErrEmptyKey on zero-len slice, got %v", err)
	}
}

func TestSign_PrefixPresent(t *testing.T) {
	s, err := NewSigner([]byte("test-key"), [sha256.Size]byte{})
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	got := s.Sign([]byte("hello"))
	if !strings.HasPrefix(got, "lore@v1:") {
		t.Errorf("expected lore@v1: prefix, got %q", got)
	}
}

func TestSign_FullLength(t *testing.T) {
	s, _ := NewSigner([]byte("test-key"), [sha256.Size]byte{})
	got := s.Sign([]byte("hello"))
	// 8 prefix + 43 base64url-no-pad body = 51 chars
	if len(got) != 51 {
		t.Errorf("FullMarkLen drift: got %d want 51 (mark=%q)", len(got), got)
	}
}

func TestSign_Deterministic(t *testing.T) {
	s, _ := NewSigner([]byte("test-key"), [sha256.Size]byte{})
	a := s.Sign([]byte("hello"))
	b := s.Sign([]byte("hello"))
	if a != b {
		t.Errorf("sign non-deterministic: %s vs %s", a, b)
	}
}

func TestVerify_HappyPath(t *testing.T) {
	s, _ := NewSigner([]byte("test-key"), [sha256.Size]byte{})
	mark := s.Sign([]byte("hello"))
	if !s.Verify([]byte("hello"), mark) {
		t.Error("verify should succeed for matching payload + mark")
	}
}

func TestVerify_WrongPayload_Rejected(t *testing.T) {
	s, _ := NewSigner([]byte("test-key"), [sha256.Size]byte{})
	mark := s.Sign([]byte("hello"))
	if s.Verify([]byte("HELLO"), mark) {
		t.Error("verify should reject altered payload")
	}
}

func TestVerify_WrongKey_Rejected(t *testing.T) {
	a, _ := NewSigner([]byte("key-a"), [sha256.Size]byte{})
	b, _ := NewSigner([]byte("key-b"), [sha256.Size]byte{})
	mark := a.Sign([]byte("hello"))
	if b.Verify([]byte("hello"), mark) {
		t.Error("verify should reject across signer keys")
	}
}

func TestVerify_WrongCorpus_Rejected(t *testing.T) {
	var corpusA, corpusB [sha256.Size]byte
	corpusB[0] = 0x42
	a, _ := NewSigner([]byte("test-key"), corpusA)
	b, _ := NewSigner([]byte("test-key"), corpusB)
	mark := a.Sign([]byte("hello"))
	if b.Verify([]byte("hello"), mark) {
		t.Error("verify should reject across corpora")
	}
}

func TestVerify_BadPrefix_Rejected(t *testing.T) {
	s, _ := NewSigner([]byte("test-key"), [sha256.Size]byte{})
	if s.Verify([]byte("hello"), "lore@v2:AAAA") {
		t.Error("verify should reject non-v1 prefix")
	}
	if s.Verify([]byte("hello"), "garbage") {
		t.Error("verify should reject garbage")
	}
	if s.Verify([]byte("hello"), "") {
		t.Error("verify should reject empty")
	}
}

func TestVerify_BadBase64_Rejected(t *testing.T) {
	s, _ := NewSigner([]byte("test-key"), [sha256.Size]byte{})
	if s.Verify([]byte("hello"), "lore@v1:!!!!not-base64!!!!") {
		t.Error("verify should reject malformed base64")
	}
}

func TestVerify_WrongDigestLen_Rejected(t *testing.T) {
	s, _ := NewSigner([]byte("test-key"), [sha256.Size]byte{})
	// 16-byte (truncated) body is base64url'd, will decode but be wrong len.
	if s.Verify([]byte("hello"), "lore@v1:AAAAAAAAAAAAAAAAAAAAAA") {
		t.Error("verify should reject wrong-length digest")
	}
}

// TestKeyDefensiveCopy verifies the signer survives caller zeroing its
// key after NewSigner.
func TestKeyDefensiveCopy(t *testing.T) {
	key := []byte("test-key")
	s, _ := NewSigner(key, [sha256.Size]byte{})
	// Wipe caller's slice.
	for i := range key {
		key[i] = 0
	}
	mark := s.Sign([]byte("hello"))
	if !s.Verify([]byte("hello"), mark) {
		t.Error("defensive key copy should survive caller-zeroing")
	}
}
