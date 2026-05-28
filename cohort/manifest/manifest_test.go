package manifest

import (
	"testing"
	"time"
)

func TestSchemaVersion(t *testing.T) {
	if SchemaVersion != 1 {
		t.Errorf("SchemaVersion drift: got %d want 1", SchemaVersion)
	}
}

func TestNew_Empty(t *testing.T) {
	m := New()
	if m.Len() != 0 {
		t.Errorf("New().Len() = %d, want 0", m.Len())
	}
	if len(m.Entries()) != 0 {
		t.Errorf("New().Entries() len = %d, want 0", len(m.Entries()))
	}
}

func TestAdd_ChainsReceiver(t *testing.T) {
	now := time.Now()
	m := New().Add(Entry{Subject: "a", FreshAt: now}).Add(Entry{Subject: "b", FreshAt: now})
	if m.Len() != 2 {
		t.Errorf("Len after 2 Adds = %d, want 2", m.Len())
	}
}

func TestEntries_SortedBySubject(t *testing.T) {
	m := New().Add(Entry{Subject: "z"}).Add(Entry{Subject: "a"}).Add(Entry{Subject: "m"})
	es := m.Entries()
	if es[0].Subject != "a" || es[1].Subject != "m" || es[2].Subject != "z" {
		t.Errorf("Entries not sorted: %v %v %v", es[0].Subject, es[1].Subject, es[2].Subject)
	}
}

func TestCanonical_HasAtLeastSixEntries(t *testing.T) {
	m := Canonical()
	if m.Len() < 6 {
		t.Errorf("Canonical().Len() = %d, want at least 6", m.Len())
	}
}

func TestCanonical_AllFreshAtSet(t *testing.T) {
	m := Canonical()
	for _, e := range m.Entries() {
		if e.FreshAt.IsZero() || e.FreshAt.Equal(FreshAtUnknown) {
			t.Errorf("Canonical entry %q has unset/unknown FreshAt", e.Subject)
		}
	}
}

func TestCanonical_AllReviewerClassFounder(t *testing.T) {
	// R166 honest default — Phase 1 SDK is founder-drafted.
	m := Canonical()
	for _, e := range m.Entries() {
		if e.ReviewerClass != ReviewerClassFounder {
			t.Errorf("entry %q ReviewerClass = %q, want founder (R166)", e.Subject, e.ReviewerClass)
		}
	}
}

func TestCanonical_AllReviewedByCounselFalse(t *testing.T) {
	// R166 honest default — flipping to true requires R145.B sibling branch.
	m := Canonical()
	for _, e := range m.Entries() {
		if e.ReviewedByCounsel {
			t.Errorf("entry %q ReviewedByCounsel = true; R166 violation (Phase 1 must be founder-drafted)", e.Subject)
		}
	}
}

func TestCanonical_AllSourcesNonEmpty(t *testing.T) {
	m := Canonical()
	for _, e := range m.Entries() {
		if len(e.Sources) == 0 {
			t.Errorf("entry %q has zero Sources", e.Subject)
		}
		for i, s := range e.Sources {
			if s == "" {
				t.Errorf("entry %q Sources[%d] empty", e.Subject, i)
			}
		}
	}
}

func TestEntries_DefensiveCopy(t *testing.T) {
	m := New().Add(Entry{Subject: "a"})
	es := m.Entries()
	es[0].Subject = "MUTATED"
	if m.Entries()[0].Subject != "a" {
		t.Error("Entries should return a defensive copy")
	}
}
