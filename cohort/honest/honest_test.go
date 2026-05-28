package honest

import (
	"bytes"
	"strings"
	"sync"
	"testing"
)

func TestAdvisoryCount_Six(t *testing.T) {
	if got := AdvisoryCount(); got != 6 {
		t.Errorf("R143 cohort shape: AdvisoryCount = %d, want 6", got)
	}
}

func TestSeverityMix_ThreeErrorThreeWarn(t *testing.T) {
	if got := ErrorAdvisoryCount(); got != 3 {
		t.Errorf("ErrorAdvisoryCount = %d, want 3", got)
	}
	if got := WarnAdvisoryCount(); got != 3 {
		t.Errorf("WarnAdvisoryCount = %d, want 3", got)
	}
}

func TestAdvisoryCodesNonEmpty(t *testing.T) {
	for i, a := range Advisories {
		if a.Code == "" {
			t.Errorf("Advisories[%d].Code empty", i)
		}
		if a.Message == "" {
			t.Errorf("Advisories[%d].Message empty", i)
		}
	}
}

func TestAllAdvisoryCodesUnique(t *testing.T) {
	seen := make(map[string]bool)
	for _, a := range Advisories {
		if seen[a.Code] {
			t.Errorf("duplicate advisory code: %s", a.Code)
		}
		seen[a.Code] = true
	}
}

func TestAllAdvisoryCodesDelvePrefixed(t *testing.T) {
	for _, a := range Advisories {
		if !strings.HasPrefix(a.Code, "DELVE_") {
			t.Errorf("advisory code %q missing DELVE_ prefix", a.Code)
		}
	}
}

func TestFireOnce_FirstFireReturnsTrue(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	if !r.FireOnce("DELVE_EMPTY_KEY_REJECTED_AT_CONSTRUCTION") {
		t.Error("first fire should return true")
	}
	if !strings.Contains(buf.String(), LoudOncePrefix) {
		t.Errorf("expected prefix %q in output, got: %s", LoudOncePrefix, buf.String())
	}
}

func TestFireOnce_SecondFireReturnsFalse(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.FireOnce("DELVE_EMPTY_KEY_REJECTED_AT_CONSTRUCTION")
	if r.FireOnce("DELVE_EMPTY_KEY_REJECTED_AT_CONSTRUCTION") {
		t.Error("second fire should return false")
	}
}

func TestFireOnce_UnknownCode_NoFire(t *testing.T) {
	r := NewReporter(nil)
	if r.FireOnce("UNKNOWN_CODE_DOES_NOT_EXIST") {
		t.Error("unknown code should not fire")
	}
}

func TestFireOnce_NilWriter_Silent(t *testing.T) {
	r := NewReporter(nil)
	if !r.FireOnce("DELVE_EMPTY_KEY_REJECTED_AT_CONSTRUCTION") {
		t.Error("first fire should return true even with nil writer")
	}
}

func TestReset_ClearsFiredSet(t *testing.T) {
	r := NewReporter(nil)
	r.FireOnce("DELVE_EMPTY_KEY_REJECTED_AT_CONSTRUCTION")
	if !r.HasFired("DELVE_EMPTY_KEY_REJECTED_AT_CONSTRUCTION") {
		t.Error("should report fired after FireOnce")
	}
	r.Reset()
	if r.HasFired("DELVE_EMPTY_KEY_REJECTED_AT_CONSTRUCTION") {
		t.Error("should report not-fired after Reset")
	}
}

func TestFireOnce_GoroutineSafe(t *testing.T) {
	r := NewReporter(nil)
	var wg sync.WaitGroup
	fires := make([]bool, 100)
	for i := range fires {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fires[i] = r.FireOnce("DELVE_EMPTY_KEY_REJECTED_AT_CONSTRUCTION")
		}(i)
	}
	wg.Wait()
	count := 0
	for _, f := range fires {
		if f {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 fire across 100 goroutines, got %d", count)
	}
}
