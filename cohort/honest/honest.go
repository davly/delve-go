// Package honest implements the cohort R143 LOUD-ONCE-WARNING-FLAG
// discipline for delve-go, with R157 substrate-native idiom (Go's
// `sync.Once` via a fired-map under a mutex — same shape as
// limitless-aiwatermark/cohort/honest).
//
// delve-go ships SIX canonical advisories — three Error (R153 strict
// surfaces) + three Warn (cadence / posture / wire-format) — that name
// the boundary between what the SDK proves and what the deploying
// tenant remains responsible for.
//
//  1. DELVE_EMPTY_KEY_REJECTED_AT_CONSTRUCTION — Error. NewClient
//     rejects empty APIKey per R175 LOAD-BEARING. There is no
//     "degraded signing" fallback — that posture would silently
//     emit cohort-public-key signatures.
//  2. DELVE_MULTITENANT_TENANT_ID_REQUIRED — Error. Empty TenantID is
//     rejected per R121 multi-tenant cohort. A rogue-empty-tenant
//     request leaks data across tenant boundaries.
//  3. DELVE_DEFAULT_LOCALHOST_NOT_PRODUCTION — Error. If neither
//     ClientOptions.URL nor DELVE_URL is set, the SDK falls back to
//     http://localhost:8092. This is fine for dev/test but MUST NOT
//     reach production — fire this advisory on first request when
//     the localhost fallback is in effect.
//  4. DELVE_NAME_COLLISION_IS_NOT_CONSUMPTION — Warn. Flagships with
//     `internal/delve/` are R155.A INDEX-LIE name-collision packages,
//     NOT wire-protocol consumers. The path to genuine R187
//     compliance is `import github.com/davly/delve-go` + a real
//     DELVE_URL.
//  5. DELVE_PHASE_1_THIN_SHIM_NOT_FULL_API — Warn. This SDK is the
//     Phase 1 thin HTTP-client shim (R176 library-first). The full
//     delve API surface (ingest / dig / forge-archaeology / scheduler
//     / lost-ideas) is NOT yet wrapped. Phase 2 expands surface area.
//  6. DELVE_KAT1_HEX_LITERAL_IS_LOAD_BEARING — Warn. The KAT-1 hex
//     literal in cohort/lore is the cohort firewall. Any drift = a
//     cohort-wide outage signal.
//
// R157 substrate-native: delve-go uses Go's standard once-fire idiom
// (sync.Mutex + fired-map) — not a foreign-ported once primitive.
package honest

import (
	"fmt"
	"io"
	"sync"
)

// LoudOncePrefix is the canonical prefix every advisory output uses
// so log scrapers can grep for advisories.
const LoudOncePrefix = "[LOUD-ONCE-WARNING]"

// Severity classifies an advisory's impact.
type Severity string

const (
	SeverityInfo  Severity = "INFO"
	SeverityWarn  Severity = "WARN"
	SeverityError Severity = "ERROR"
)

// Advisory describes a single LOUD-ONCE advisory.
type Advisory struct {
	Code     string
	Severity Severity
	Message  string
	DocLink  string
}

// Advisories is the canonical list this SDK ships with.
var Advisories = []Advisory{
	{
		Code:     "DELVE_EMPTY_KEY_REJECTED_AT_CONSTRUCTION",
		Severity: SeverityError,
		Message:  "NewClient rejects empty APIKey per R175 LOAD-BEARING. No degraded-signing fallback.",
		DocLink:  "client.go NewClient",
	},
	{
		Code:     "DELVE_MULTITENANT_TENANT_ID_REQUIRED",
		Severity: SeverityError,
		Message:  "NewClient rejects empty TenantID per R121 multi-tenant cohort.",
		DocLink:  "client.go NewClient",
	},
	{
		Code:     "DELVE_DEFAULT_LOCALHOST_NOT_PRODUCTION",
		Severity: SeverityError,
		Message:  "Default URL http://localhost:8092 fires when neither ClientOptions.URL nor DELVE_URL is set. MUST NOT reach production.",
		DocLink:  "client.go DefaultURL",
	},
	{
		Code:     "DELVE_NAME_COLLISION_IS_NOT_CONSUMPTION",
		Severity: SeverityWarn,
		Message:  "Flagships with internal/delve/ are R155.A INDEX-LIE name-collision packages, NOT wire-protocol consumers. Path to R187: import github.com/davly/delve-go + DELVE_URL.",
		DocLink:  "README.md R187 section",
	},
	{
		Code:     "DELVE_PHASE_1_THIN_SHIM_NOT_FULL_API",
		Severity: SeverityWarn,
		Message:  "Phase 1 thin HTTP-client shim (R176 library-first). Full delve surface (ingest/dig/forge-archaeology/scheduler/lost-ideas) not yet wrapped.",
		DocLink:  "README.md Phase 1 section",
	},
	{
		Code:     "DELVE_KAT1_HEX_LITERAL_IS_LOAD_BEARING",
		Severity: SeverityWarn,
		Message:  "KAT-1 hex literal in cohort/lore is the cohort firewall. Any drift = cohort-wide outage signal.",
		DocLink:  "cohort/lore/lore.go Digest",
	},
}

// Reporter is the R143 LOUD-ONCE emitter. Each advisory code fires at
// most once per Reporter lifetime; subsequent FireOnce calls with the
// same code are silent.
type Reporter struct {
	mu     sync.Mutex
	fired  map[string]bool
	output io.Writer
}

// NewReporter constructs a Reporter writing to w. If w is nil the
// reporter is silent (used for tests).
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{fired: make(map[string]bool), output: w}
}

// FireOnce emits the advisory if it has not been fired yet.
// Returns true iff the message was emitted (i.e. first fire).
//
// Goroutine-safe per R143 canonical shape — concurrent calls
// serialise on the reporter's mutex.
func (r *Reporter) FireOnce(code string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.fired[code] {
		return false
	}
	r.fired[code] = true
	adv, ok := byCode(code)
	if !ok {
		return false
	}
	if r.output != nil {
		fmt.Fprintf(r.output, "%s [%s] %s — %s (see %s)\n",
			LoudOncePrefix, adv.Severity, adv.Code, adv.Message, adv.DocLink)
	}
	return true
}

// HasFired reports whether code has been fired during this reporter's
// lifetime.
func (r *Reporter) HasFired(code string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.fired[code]
}

// Reset clears the fired set — used by tests + by hot-reload scenarios.
func (r *Reporter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fired = make(map[string]bool)
}

// AdvisoryCount returns the count of advisories shipped — exposed for
// the firewall package's R145.C tests to enforce the canonical 6.
func AdvisoryCount() int { return len(Advisories) }

// ErrorAdvisoryCount returns the count of SeverityError advisories.
func ErrorAdvisoryCount() int {
	n := 0
	for _, a := range Advisories {
		if a.Severity == SeverityError {
			n++
		}
	}
	return n
}

// WarnAdvisoryCount returns the count of SeverityWarn advisories.
func WarnAdvisoryCount() int {
	n := 0
	for _, a := range Advisories {
		if a.Severity == SeverityWarn {
			n++
		}
	}
	return n
}

// byCode returns the advisory for code (and whether it exists).
func byCode(code string) (Advisory, bool) {
	for _, a := range Advisories {
		if a.Code == code {
			return a, true
		}
	}
	return Advisory{}, false
}
