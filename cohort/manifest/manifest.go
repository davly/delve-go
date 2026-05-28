// Package manifest implements the R150 cohort-canonical schematised-
// knowledge envelope for delve-go's wire-level consumer surface.
//
// Why delve-go consumes this from inception (R174 5-of-5 strict +
// R150 R-PARALLEL-MAP-R144-REVIEW-METADATA-SIBLING):
//
// delve-go's domain content is the WIRE-LEVEL CONSUMER CONTRACT a
// flagship inherits when it imports `github.com/davly/delve-go`. The
// contract anchors to R187 (infra cohort compat), R191 (cross-infra
// audit-chain emit), R151 (KAT-1 cohort firewall), R175 (LOAD-BEARING),
// and R155.A (INDEX-LIE detection — name-collision is NOT consumption).
//
// Per R166, every entry carries ReviewedByCounsel = false honest-
// default — this SDK is a forcing-move primitive, not legal advice
// or a complete compliance solution.
package manifest

import (
	"sort"
	"time"
)

// SchemaVersion is the R150 manifest schema version pin. v1 = the
// 5-field shape (Subject / Sources / Confidence / ReviewedByCounsel /
// FreshAt with the R150.E ReviewerClass extension).
const SchemaVersion = 1

// FreshAtUnknown is the sentinel for entries whose freshness is
// not-yet-known.
var FreshAtUnknown = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

// Canonical regulatory + cohort source identifiers used by entries.
const (
	// SourceR187InfraCohortCompat — the R187 R-INFRA-PROJECT-COHORT-
	// COMPATIBLE rule (Batch 7 candidate at session close, ripe for
	// formal Batch 8 promotion).
	SourceR187InfraCohortCompat = "R187 R-INFRA-PROJECT-COHORT-COMPATIBLE (infra-project cohort compatibility; deep over-saturated ≥5/3 post-late-ships)"

	// SourceR191CrossInfraAuditChainEmit — R191 boundary-signing pin.
	SourceR191CrossInfraAuditChainEmit = "R191 R-CROSS-INFRA-AUDIT-CHAIN-EMIT (every outbound infra call emits a Mirror-Mark-signed boundary receipt)"

	// SourceR192NameCollisionNotConsumption — the I1 strategic-review
	// catch in scalar form.
	SourceR192NameCollisionNotConsumption = "R192 R-INFRA-NAME-COLLISION-IS-NOT-CONSUMPTION (dir-name match does NOT count; wire-level import DOES)"

	// SourceR151KATCohortInvariant — the R151 KAT-1 pin.
	SourceR151KATCohortInvariant = "R151 KAT-AS-COHORT-INVARIANT-CROSS-SUBSTRATE-PIN (KAT-1 hex 239a7d0d…)"

	// SourceR175MirrorMarkLoadBearingCohort — R175 posture.
	SourceR175MirrorMarkLoadBearingCohort = "R175 R-MIRROR-MARK-LOAD-BEARING-IN-PRODUCTION cohort wire posture"

	// SourceR176LibraryFirstWireLater — Phase 1 discipline anchor.
	SourceR176LibraryFirstWireLater = "R176 R-LIBRARY-FIRST-WIRE-LATER (library artefacts ship first; production wire-in to flagships is Phase 2)"

	// SourceR155AIndexLie — the INDEX-LIE class.
	SourceR155AIndexLie = "R155.A R-INDEX-LIE (declared modules / sub-packages / consumers that don't actually exist on the wire)"

	// SourceR121MultiTenant — the multi-tenant cohort key rule.
	SourceR121MultiTenant = "R121 R-MULTI-TENANT-CACHE-KEY (every request scoped by tenant-id-first)"

	// SourceR166LiabilityFooterConst — disclaimer posture.
	SourceR166LiabilityFooterConst = "R166 LIABILITY-FOOTER-CONST (founder-drafted disclaimer; not legal advice)"

	// SourceL43MirrorMarkV1Algorithm — the wire-format primitive.
	SourceL43MirrorMarkV1Algorithm = "L43 Mirror-Mark v1 (lore@v1: prefix; HMAC-SHA256 over 0x01 || corpusSHA || payload; base64url body)"

	// SourceInfraDelveCanonicalService — the upstream infra service.
	SourceInfraDelveCanonicalService = "infrastructure/delve (port 8092; canonical search service)"
)

// Confidence — R150 confidence axis.
type Confidence int

const (
	// ConfidenceHigh — multi-source corroborated, regulator-grade.
	ConfidenceHigh Confidence = 3
	// ConfidenceMedium — single-source citation, internally consistent.
	ConfidenceMedium Confidence = 2
	// ConfidenceLow — heuristic / placeholder; LOUD-ONCE on use.
	ConfidenceLow Confidence = 1
)

// ReviewerClass — R150.E REVIEWER-CLASS-EXTENSION-FIELD.
type ReviewerClass string

const (
	// ReviewerClassFounder — founder-drafted (honest default for
	// inception-stage SDKs per R166).
	ReviewerClassFounder ReviewerClass = "founder"
	// ReviewerClassRegulatoryCounsel — qualified counsel review
	// completed. Flipping to this class requires a R145.B sibling
	// branch.
	ReviewerClassRegulatoryCounsel ReviewerClass = "regulatory_counsel"
	// ReviewerClassNotifiedBody — Article 43 conformity-assessment
	// body has signed off.
	ReviewerClassNotifiedBody ReviewerClass = "notified_body"
)

// Entry is the R150 5-field manifest record.
type Entry struct {
	Subject           string
	Sources           []string
	Confidence        Confidence
	ReviewedByCounsel bool
	ReviewerClass     ReviewerClass
	FreshAt           time.Time
}

// Manifest is the SDK's canonical schematised-knowledge envelope.
type Manifest struct {
	entries []Entry
}

// New constructs an empty Manifest.
func New() *Manifest { return &Manifest{} }

// Add appends an entry to the manifest. Returns the receiver for chaining.
func (m *Manifest) Add(e Entry) *Manifest {
	m.entries = append(m.entries, e)
	return m
}

// Entries returns a sorted defensive copy of the manifest entries.
func (m *Manifest) Entries() []Entry {
	out := make([]Entry, len(m.entries))
	copy(out, m.entries)
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Subject < out[j].Subject
	})
	return out
}

// Len returns the entry count.
func (m *Manifest) Len() int { return len(m.entries) }

// Canonical builds the default delve-go manifest — the cohort-canonical
// R150 shape this SDK ships with at inception.
func Canonical() *Manifest {
	now := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC)
	m := New()
	m.Add(Entry{
		Subject: "Wire-level search call",
		Sources: []string{
			SourceR187InfraCohortCompat,
			SourceR192NameCollisionNotConsumption,
			SourceInfraDelveCanonicalService,
		},
		Confidence:        ConfidenceHigh,
		ReviewedByCounsel: false,
		ReviewerClass:     ReviewerClassFounder,
		FreshAt:           now,
	})
	m.Add(Entry{
		Subject: "Boundary-signing receipt",
		Sources: []string{
			SourceR191CrossInfraAuditChainEmit,
			SourceR175MirrorMarkLoadBearingCohort,
			SourceL43MirrorMarkV1Algorithm,
		},
		Confidence:        ConfidenceHigh,
		ReviewedByCounsel: false,
		ReviewerClass:     ReviewerClassFounder,
		FreshAt:           now,
	})
	m.Add(Entry{
		Subject: "KAT-1 cohort firewall",
		Sources: []string{
			SourceR151KATCohortInvariant,
		},
		Confidence:        ConfidenceHigh,
		ReviewedByCounsel: false,
		ReviewerClass:     ReviewerClassFounder,
		FreshAt:           now,
	})
	m.Add(Entry{
		Subject: "Phase 1 library-first discipline",
		Sources: []string{
			SourceR176LibraryFirstWireLater,
		},
		Confidence:        ConfidenceHigh,
		ReviewedByCounsel: false,
		ReviewerClass:     ReviewerClassFounder,
		FreshAt:           now,
	})
	m.Add(Entry{
		Subject: "Multi-tenant cohort key",
		Sources: []string{
			SourceR121MultiTenant,
		},
		Confidence:        ConfidenceHigh,
		ReviewedByCounsel: false,
		ReviewerClass:     ReviewerClassFounder,
		FreshAt:           now,
	})
	m.Add(Entry{
		Subject: "INDEX-LIE detection (name-collision flagships)",
		Sources: []string{
			SourceR155AIndexLie,
			SourceR192NameCollisionNotConsumption,
		},
		Confidence:        ConfidenceHigh,
		ReviewedByCounsel: false,
		ReviewerClass:     ReviewerClassFounder,
		FreshAt:           now,
	})
	m.Add(Entry{
		Subject: "Liability footer (founder-drafted)",
		Sources: []string{
			SourceR166LiabilityFooterConst,
		},
		Confidence:        ConfidenceMedium,
		ReviewedByCounsel: false,
		ReviewerClass:     ReviewerClassFounder,
		FreshAt:           now,
	})
	return m
}
