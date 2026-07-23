# Changelog

All notable changes to kilio are documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
kilio uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

Nothing yet — `kilio-core` is the next crate up. See
[ROADMAP.md](ROADMAP.md).

---

## [0.1.0] - 2026-07-23

### Added — the sealed-crypto core and the spec it implements

- **`decisions.md`**, the authoritative design record: the privacy spine
  (sealed at source, no mandatory identity, anonymous two-way channel,
  metadata minimization), the receipt-passphrase identity primitive, the
  data model, the ofisi-style branch-scoping pattern, the three seams
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
