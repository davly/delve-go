package firewall

import (
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/davly/delve-go/cohort/lore"
)

// repoRoot walks up from this test file to find the SDK root.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate runtime.Caller(0)")
	}
	// Walk up until we find a go.mod sibling.
	dir := filepath.Dir(file)
	for i := 0; i < 10; i++ {
		dir = filepath.Dir(dir)
		entries, err := filepath.Glob(filepath.Join(dir, "go.mod"))
		if err == nil && len(entries) > 0 {
			return dir
		}
	}
	t.Fatal("could not find repo root")
	return ""
}

// TestKAT1Parity is the cohort-canonical R151 firewall pin imported
// directly into the firewall test surface — any code path that breaks
// KAT-1 parity fails this test before the SDK ships.
//
// This is the cohort's load-bearing trust anchor: when downstream
// regulators rely on OpenSSL cold-verify, this test guarantees the
// hex literal in cohort/lore.Digest matches what the Go runtime
// computes via stdlib HMAC-SHA256.
func TestKAT1Parity(t *testing.T) {
	if !lore.AssertKAT1Parity() {
		t.Fatalf("KAT-1 cohort firewall FAILED:\n  computed: %s\n  pinned:   %s",
			lore.ComputeKAT1(), lore.Digest)
	}
}

func TestExpectedCohortPackagesSorted(t *testing.T) {
	got := ExpectedCohortPackages()
	if !sort.StringsAreSorted(got) {
		t.Errorf("ExpectedCohortPackages not sorted: %v", got)
	}
}

func TestExpectedCohortPackagesFiveOfFive(t *testing.T) {
	// R174 5-of-5: exactly five cohort packages.
	if got := len(ExpectedCohortPackages()); got != 5 {
		t.Errorf("R174 5-of-5 fail: got %d cohort packages want 5", got)
	}
}

func TestExpectedCohortPackagesHaveCanonicalSet(t *testing.T) {
	want := []string{"firewall", "honest", "lore", "manifest", "mirrormark"}
	got := ExpectedCohortPackages()
	if d := CompareSorted(got, want); d != -1 {
		t.Errorf("cohort package set drift at %d: got %v want %v", d, got, want)
	}
}

func TestCohortDirsOnDisk(t *testing.T) {
	root := repoRoot(t)
	got, err := ScanCohort(root)
	if err != nil {
		t.Fatalf("scan cohort: %v", err)
	}
	want := ExpectedCohortPackages()
	if d := CompareSorted(got, want); d != -1 {
		t.Errorf("on-disk cohort drift at %d: got %v want %v", d, got, want)
	}
}

func TestInternalDirsOnDisk(t *testing.T) {
	root := repoRoot(t)
	got, err := ScanInternal(root)
	if err != nil {
		t.Fatalf("scan internal: %v", err)
	}
	want := ExpectedInternalPackages()
	if d := CompareSorted(got, want); d != -1 {
		t.Errorf("on-disk internal drift at %d: got %v want %v", d, got, want)
	}
}

func TestCompareSortedEqualLength(t *testing.T) {
	if CompareSorted([]string{"a", "b"}, []string{"a", "b"}) != -1 {
		t.Error("equal not -1")
	}
	if CompareSorted([]string{"a", "b"}, []string{"a", "c"}) != 1 {
		t.Error("mismatch at 1 not detected")
	}
}

func TestCompareSortedDifferentLength(t *testing.T) {
	if CompareSorted([]string{"a"}, []string{"a", "b"}) != 1 {
		t.Error("shorter at right not detected")
	}
	if CompareSorted([]string{"a", "b"}, []string{"a"}) != 1 {
		t.Error("shorter at left not detected")
	}
}
