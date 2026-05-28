// Package lore pins the ecosystem-canonical KAT-1 HMAC-SHA256 invariant
// for the R151 cross-substrate cohort pin within delve-go.
//
// delve-go is the WIRE-LEVEL Go SDK for the davly/delve infrastructure
// service (port 8092). Per the I1 strategic-review catch in the
// 2026-05-28 infra marathon, flagships with `internal/delve/` are NOT
// genuine consumers; this SDK is the path to wire-level consumption.
//
// Cohort-port from inception per R174 5-of-5 strict.
//
// Cold-verify recipe (OpenSSL one-liner — no Go toolchain involved):
//
//	printf '\x01' > /tmp/kat1.bin
//	printf '\x00%.0s' {1..32} >> /tmp/kat1.bin
//	openssl dgst -sha256 -mac hmac -macopt key: /tmp/kat1.bin
//	# → 239a7d0d3f1bbe3a98aede01e2ad818c2db60b7177c02e2f015035b2b5b7dbca
package lore

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Digest is the cohort-canonical KAT-1 HMAC-SHA256 digest, hex-encoded.
// Pinned byte-identical to foundation/pkg/mirrormark.KAT1Digest and to
// every cohort port across the ecosystem.
const Digest = "239a7d0d3f1bbe3a98aede01e2ad818c2db60b7177c02e2f015035b2b5b7dbca"

// InputLen is the canonical KAT-1 input length: 1 byte version tag +
// 32 bytes zero corpus = 33 bytes.
const InputLen = 33

// VersionTag is the v1 1-byte tag prefix. Bumping this byte to v2
// invalidates every mark in flight — necessary if the canonicalization
// rule ever changes.
const VersionTag byte = 0x01

// ComputeKAT1 deterministically reproduces the KAT-1 digest using
// only stdlib crypto. Returns the hex-encoded digest.
//
// The implementation is pure FIPS PUB 180-4 SHA-256 + RFC 2104 HMAC
// + RFC 4648 hex-encoding — no Limitless types or invariants. A
// regulator's OpenSSL one-liner produces byte-identical output.
func ComputeKAT1() string {
	input := make([]byte, InputLen)
	input[0] = VersionTag
	// bytes 1..32 are zero corpus.
	mac := hmac.New(sha256.New, nil) // empty key
	_, _ = mac.Write(input)
	return hex.EncodeToString(mac.Sum(nil))
}

// AssertKAT1Parity returns true iff ComputeKAT1() == Digest.
// The R151 firewall pin: drift here breaks the cohort-wide trust
// chain on every CI run.
func AssertKAT1Parity() bool {
	return ComputeKAT1() == Digest
}
