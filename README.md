# delve-go

Go SDK for the Limitless `delve` infrastructure service ([infrastructure/delve](https://github.com/davly/delve)).

## Phase 1 — thin HTTP-client shim

Public API:
- `NewClient(opts ClientOptions) (*Client, error)` — constructs the client with `DELVE_URL` env-var URL + iik_… key
- `Search(ctx, request)` — wraps the HTTP boundary call
- Cohort 5-pack (`cohort/mirrormark/`, `cohort/lore/`, `cohort/honest/`, `cohort/manifest/`, `cohort/firewall/`)

## R187 R-INFRA-PROJECT-COHORT-COMPATIBLE compliance

This SDK is the WIRE-LEVEL consumer surface for `infrastructure/delve`. Flagships that import this SDK become genuine wire-protocol consumers per R192 R-INFRA-NAME-COLLISION-IS-NOT-CONSUMPTION (i.e. dir-name match does NOT count; this SDK's `import github.com/davly/delve-go` DOES count).

## R191 R-CROSS-INFRA-AUDIT-CHAIN-EMIT — Phase 1 status (signer present, not yet wired)

> **Honesty correction (2026-06-04).** An earlier version of this section
> read *"Every HTTP request signed via `internal/boundarysigner/` Mirror-Mark
> before emit."* That was an **over-claim** against the Phase 1 code on disk
> and has been corrected below.

The L43 Mirror-Mark v1 signer **exists and is unit-tested in isolation**
(`internal/boundarysigner/`: `NewSigner` / `Sign` / `Verify`, empty-key
rejected via `ErrEmptyKey`, KAT-1 parity pinned). **However, the Phase 1
client does NOT yet call it:** `client.go`'s `Search` method imports only
the standard library and sets `Authorization: Bearer …` + `X-Tenant-ID`
headers — it does **not** import `internal/boundarysigner`, does not invoke
`Sign`, and attaches **no** `X-Mirror-Mark` header. Requests therefore
currently reach the wire **unsigned**.

Concretely, on disk today:
- `Search` (request dispatch) is **not** boundary-signed; the only auth on
  the wire is the bearer API key + tenant header.
- The signer is a correct, self-contained library awaiting Phase 2
  wire-in once a `LORE_SIGNING_KEY` + corpus-SHA pair is plumbed through
  `ClientOptions`.

So R191 boundary-signing is **deferred to Phase 2**, not satisfied in
Phase 1. The `cohort/firewall` `ExpectedInternalPackages()` /
`TestInternalDirsOnDisk` checks only assert the signer package is
**present on disk** — they do not (and cannot) assert it is wired into the
request path. The `cohort/honest` advisory
`DELVE_PHASE_1_THIN_SHIM_NOT_FULL_API` already flags this Phase-1 boundary.

What **is** enforced at the client boundary in Phase 1:
- `R175` empty-`APIKey` rejection at `NewClient` (no degraded-signing fallback).
- `R121` empty-`TenantID` rejection at `NewClient` (multi-tenant isolation).
- `R151` KAT-1 firewall parity (`cohort/firewall/firewall_test.go`).

## R151 KAT-1 firewall

KAT-1 hex `239a7d0d3f1bbe3a98aede01e2ad818c2db60b7177c02e2f015035b2b5b7dbca` pinned via `cohort/firewall/firewall_test.go`. Substrate parity test mandatory.

## License

Apache 2.0
