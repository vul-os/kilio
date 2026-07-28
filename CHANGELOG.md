# Changelog

All notable changes to kilio are documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
kilio uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added

- **`kilio-core`** — the domain model, the sealed SQLite store (ciphertext and
  routing metadata only), diwan-style branch scoping through a single
  `Requester` choke point and `branch_scoped_key()` builder, and the
  `Delivery` / `Reachability` seams with their local defaults.
- **Branch key pinning** (`kilio-seal::pin`) — signed, versioned branch
  descriptors; `BranchPin` trust-on-first-use where the key is accepted exactly
  once; rekey certificates countersigned by the currently pinned key as the
  only rotation path; epoch monotonicity as a rollback guard. Closes the gap
  where a host could serve a self-consistent substituted branch key and have it
  accepted.
- **CI** (`.github/workflows/ci.yml`) — `cargo fmt --check`, `clippy -D
  warnings`, native and `wasm32` builds, `scripts/test-gate.sh`, and the web
  build plus `scripts/web-assets-gate.sh`. Both gates fail closed: they assert
  the checks actually ran (per-crate passing-test floors; a non-empty bundle to
  scan) rather than trusting an exit code.

### Changed

- `SealedStore::put_branch()` now fails closed on branch key substitution: a
  stored branch's keys are immutable (`BranchKeyChangeRefused`) and a branch id
  must derive from the keys it publishes (`BranchIdMismatch`). It previously
  accepted such a write and silently kept the old keys.
- Store row decoding fails closed. Corrupt id / status / direction columns
  return a conversion error instead of defaulting to a zero id or `New` —
  the old behaviour made authorization decisions on garbage.
- `ClaimStatus::from_str` / `Direction::from_str` renamed to `from_db` (storage
  codecs, not `std::str::FromStr`).

### Removed

- Unused dependencies: `x25519-dalek`, `subtle`, `rand_core` (declared by
  `kilio-seal`, referenced nowhere), and `chacha20poly1305`, `time` (declared in
  `[workspace.dependencies]`, used by no member).

### Fixed

- Docs reconciled against the code: `kilio-core` is no longer described as
  planned; the `Delivery` / `Reachability` trait signatures match the
  implemented ones; `SubprocessTunnel` and `KotvaDelivery` are marked as the
  partial seams they are; the non-existent `Ephor` reachability variant is
  removed from every doc (Ephor is its own product,
  `github.com/vul-os/ephor`, not a wede component); `DeployMode` is documented
  with its two real values; `branch_key()`/`requesterID()` renamed to the
  implemented `branch_scoped_key()`/`Requester`; the envelope field list and
  the `hpke` crate name corrected in `decisions.md`.

---

## [0.1.0] - 2026-07-23

### Added — the sealed-crypto core and the spec it implements

- **`decisions.md`**, the authoritative design record: the privacy spine
  (sealed at source, no mandatory identity, anonymous two-way channel,
  metadata minimization), the receipt-passphrase identity primitive, the
  data model, the diwan-style branch-scoping pattern, the three seams
  (`Delivery`, `Reachability`, `Identity`/deploy-mode), the one-org-per-
  deployment stance, the abbreviated threat model, and the build order.
- **`kilio-seal`** — the sealed-submission crypto spine, native + `wasm32`:
  - HPKE (RFC 9180) seal/open to a branch or a claim, mode Base,
    DHKEM(X25519, HKDF-SHA256), HKDF-SHA256, ChaCha20Poly1305, via the
    audited `hpke` crate.
  - `Receipt` → `ClaimKeys`: a BIP-39 12-word passphrase deterministically
    derived (Argon2id, branch-id-salted) into a per-claim Ed25519 signing
    key and X25519 recipient key.
  - `Envelope`/`Inner`, a kotva-MOTE-shaped sealed-sender envelope: cleartext
    outer routing metadata (kind, recipient tag, size bucket) bound into the
    AEAD's associated data, sealed inner payload carrying the sender's
    per-claim identity and body.
  - Size-bucketed padding (4 KiB … 4 MiB, then 4 MiB steps) so ciphertext
    length cannot fingerprint a claim.
  - `pow` — an anonymous, per-message proof-of-work cold-contact gate bound
    to the envelope's own ciphertext, so a stamp cannot be precomputed or
    replayed onto a different message.
  - 21 unit tests: branch-id binding, sign/verify roundtrips, receipt→
    claim-key determinism and per-branch separation, the full two-way
    submit → reply → return-and-read flow, tampered-envelope and
    wrong-key rejection, and PoW solve/verify/reject-insufficient-work.
- Workspace scaffolding: `Cargo.toml` workspace with shared crypto/encoding
  dependency pins, dual MIT/Apache-2.0 licensing.

[Unreleased]: https://github.com/vul-os/kilio/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/vul-os/kilio/releases/tag/v0.1.0
