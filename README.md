# delve-go

Go SDK for the Limitless `delve` infrastructure service ([infrastructure/delve](https://github.com/davly/delve)).

## Phase 1 — thin HTTP-client shim

Public API:
- `NewClient(opts ClientOptions) (*Client, error)` — constructs the client with `DELVE_URL` env-var URL + iik_… key
- `Search(ctx, request)` — wraps the HTTP boundary call
- Cohort 5-pack (`cohort/mirrormark/`, `cohort/lore/`, `cohort/honest/`, `cohort/manifest/`, `cohort/firewall/`)

## R187 R-INFRA-PROJECT-COHORT-COMPATIBLE compliance

This SDK is the WIRE-LEVEL consumer surface for `infrastructure/delve`. Flagships that import this SDK become genuine wire-protocol consumers per R192 R-INFRA-NAME-COLLISION-IS-NOT-CONSUMPTION (i.e. dir-name match does NOT count; this SDK's `import github.com/davly/delve-go` DOES count).

## R191 R-CROSS-INFRA-AUDIT-CHAIN-EMIT compliance

Every HTTP request signed via `internal/boundarysigner/` Mirror-Mark before emit. Empty key rejected.

## R151 KAT-1 firewall

KAT-1 hex `239a7d0d3f1bbe3a98aede01e2ad818c2db60b7177c02e2f015035b2b5b7dbca` pinned via `cohort/firewall/firewall_test.go`. Substrate parity test mandatory.

## License

Apache 2.0
