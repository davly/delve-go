package lore

import "testing"

// TestAssertKAT1Parity is the cohort-canonical R151 firewall pin.
// Drift here = cohort-wide outage signal.
func TestAssertKAT1Parity(t *testing.T) {
	if !AssertKAT1Parity() {
		t.Fatalf("KAT-1 parity FAILED:\n  got:  %s\n  want: %s",
			ComputeKAT1(), Digest)
	}
}

func TestDigestLiteral(t *testing.T) {
	if Digest != "239a7d0d3f1bbe3a98aede01e2ad818c2db60b7177c02e2f015035b2b5b7dbca" {
		t.Fatalf("Digest literal drift: %q", Digest)
	}
}

func TestInputLen(t *testing.T) {
	if InputLen != 33 {
		t.Fatalf("InputLen drift: got %d want 33", InputLen)
	}
}

func TestVersionTag(t *testing.T) {
	if VersionTag != 0x01 {
		t.Fatalf("VersionTag drift: got %#x want 0x01", VersionTag)
	}
}

func TestComputeKAT1Deterministic(t *testing.T) {
	a := ComputeKAT1()
	b := ComputeKAT1()
	if a != b {
		t.Fatalf("ComputeKAT1 non-deterministic: %s vs %s", a, b)
	}
}

func TestComputeKAT1MatchesPinned(t *testing.T) {
	if got := ComputeKAT1(); got != Digest {
		t.Errorf("got %s\nwant %s", got, Digest)
	}
}
