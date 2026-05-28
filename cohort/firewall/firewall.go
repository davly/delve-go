// Package firewall implements the R145.C FIREWALL-TEST-DISCIPLINE pin
// for delve-go — structural firewall against root + cohort/ + internal/
// drift, plus the R151 KAT-1 directly-imported parity gate.
//
// Cohort-port from inception per R174 5-of-5 strict — delve-go ships
// dedicated cohort/firewall/ from day one.
//
// delve-go's layout is the FLAT shape:
//   - root: client.go + client_test.go + go.mod + README.md + LICENSE
//   - cohort/{lore, mirrormark, honest, manifest, firewall}
//   - internal/boundarysigner — the Mirror-Mark signer impl
package firewall

import (
	"os"
	"path/filepath"
	"sort"
)

// ExpectedCohortPackages returns the canonical R174 5-of-5 cohort
// packages this SDK ships under `cohort/`.
//
// The cohort 5-pack is:
//
//	firewall   — R145.C structural firewall (this package)
//	honest     — R143 LOUD-ONCE advisories
//	lore       — R151 KAT-1 pin
//	manifest   — R150 schematised knowledge
//	mirrormark — L43 v1 wire-format constants
func ExpectedCohortPackages() []string {
	return []string{
		"firewall",
		"honest",
		"lore",
		"manifest",
		"mirrormark",
	}
}

// ExpectedInternalPackages returns the internal/ package names.
func ExpectedInternalPackages() []string {
	return []string{
		"boundarysigner",
	}
}

// ScanCohort returns the directory names under cohort/ that contain
// at least one .go file.
func ScanCohort(repoRoot string) ([]string, error) {
	return scanGoSubtree(filepath.Join(repoRoot, "cohort"))
}

// ScanInternal returns the directory names under internal/ that
// contain at least one .go file.
func ScanInternal(repoRoot string) ([]string, error) {
	return scanGoSubtree(filepath.Join(repoRoot, "internal"))
}

// scanGoSubtree returns the immediate sub-directory names of root
// that contain at least one .go file (recursive within each
// sub-dir). Returns sorted result for deterministic comparisons.
func scanGoSubtree(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		has, err := hasGoFile(filepath.Join(root, e.Name()))
		if err != nil {
			return nil, err
		}
		if has {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}

// hasGoFile returns true iff dir contains at least one .go file
// (recursive).
func hasGoFile(dir string) (bool, error) {
	found := false
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".go" {
			found = true
		}
		return nil
	})
	return found, err
}

// CompareSorted reports whether a and b are equal as sorted string
// slices. Returns the first mismatch index (or -1 on equality).
func CompareSorted(a, b []string) int {
	if len(a) != len(b) {
		if len(a) < len(b) {
			return len(a)
		}
		return len(b)
	}
	for i := range a {
		if a[i] != b[i] {
			return i
		}
	}
	return -1
}
