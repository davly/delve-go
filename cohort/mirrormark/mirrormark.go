// Package mirrormark exposes the L43 Mirror-Mark v1 wire-format
// constants the delve-go SDK uses for boundary signing.
//
// L43 Mirror-Mark v1 wire shape:
//
//	"lore@v1:" || base64url( HMAC-SHA256(key, 0x01 || corpusSHA(32B) || payload) )
//
// The HMAC is computed over (VersionTag || CorpusSHA || Payload). The
// 8-character prefix "lore@v1:" plus the base64url-encoded 32-byte
// digest body is the on-the-wire artefact.
//
// This package is wire-format-only — actual sign/verify lives in
// internal/boundarysigner. Externalising the constants here lets
// consumers cold-verify wire-format invariants without importing the
// signer.
//
// Cohort-port from inception per R174 5-of-5 strict.
package mirrormark

import "crypto/sha256"

// Prefix is the canonical L43 wire prefix. 8 bytes.
const Prefix = "lore@v1:"

// PrefixLen is the byte length of Prefix.
const PrefixLen = 8

// VersionTag is the v1 HMAC input prefix byte. 0x01.
const VersionTag byte = 0x01

// CorpusSHALen is the byte length of the corpus-SHA component (32B).
const CorpusSHALen = sha256.Size

// DigestLen is the HMAC-SHA256 output byte length (32B).
const DigestLen = sha256.Size

// EncodedBodyLen is the base64url-no-pad encoded length of a 32-byte
// digest. ceil(32 * 4 / 3) = 44; base64url-no-pad strips trailing '='.
// HMAC-SHA256 → 32 bytes → 43 base64url-no-pad chars.
const EncodedBodyLen = 43

// FullMarkLen is the canonical full-mark length (Prefix + encoded body).
const FullMarkLen = PrefixLen + EncodedBodyLen

// MinInputLen is the minimum byte length of a Mirror-Mark HMAC input
// (VersionTag + CorpusSHA + 0-byte payload).
const MinInputLen = 1 + CorpusSHALen
