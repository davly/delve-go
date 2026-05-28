package mirrormark

import (
	"crypto/sha256"
	"testing"
)

func TestPrefixLiteral(t *testing.T) {
	if Prefix != "lore@v1:" {
		t.Errorf("Prefix drift: %q", Prefix)
	}
}

func TestPrefixLen(t *testing.T) {
	if PrefixLen != 8 {
		t.Errorf("PrefixLen: got %d want 8", PrefixLen)
	}
	if len(Prefix) != PrefixLen {
		t.Errorf("len(Prefix)=%d != PrefixLen=%d", len(Prefix), PrefixLen)
	}
}

func TestVersionTag(t *testing.T) {
	if VersionTag != 0x01 {
		t.Errorf("VersionTag drift: got %#x want 0x01", VersionTag)
	}
}

func TestCorpusSHALen(t *testing.T) {
	if CorpusSHALen != sha256.Size {
		t.Errorf("CorpusSHALen: got %d want %d", CorpusSHALen, sha256.Size)
	}
}

func TestDigestLen(t *testing.T) {
	if DigestLen != sha256.Size {
		t.Errorf("DigestLen: got %d want %d", DigestLen, sha256.Size)
	}
}

func TestEncodedBodyLen(t *testing.T) {
	if EncodedBodyLen != 43 {
		t.Errorf("EncodedBodyLen: got %d want 43", EncodedBodyLen)
	}
}

func TestFullMarkLen(t *testing.T) {
	want := PrefixLen + EncodedBodyLen
	if FullMarkLen != want {
		t.Errorf("FullMarkLen: got %d want %d", FullMarkLen, want)
	}
}

func TestMinInputLen(t *testing.T) {
	if MinInputLen != 1+CorpusSHALen {
		t.Errorf("MinInputLen: got %d want %d", MinInputLen, 1+CorpusSHALen)
	}
}
