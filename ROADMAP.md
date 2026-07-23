# kilio — Roadmap

kilio is a self-hostable, sealed, anonymous-first intake app for sensitive
claims — harassment, misconduct, whistleblowing, safety concerns. One small
binary an organization runs itself, reachable over a tunnel with no fixed
infra, usable across branches, standalone by default, decentralized when it
wants to be. Full design reasoning lives in [`decisions.md`](decisions.md);
this document tracks *build order* and *status* against it.

> **Stack invariants (FROZEN):** one Rust workspace for everything that
> touches crypto or domain logic; JSX (never `.tsx`) for app shells; sealing
> exists in exactly one implementation (`kilio-seal`, compiled native +
> `wasm32`) and is never reimplemented per surface; instance-per-org, never
> multi-tenant (decisions.md §7).

The build order below is decisions.md §9, followed verbatim — it is a
dependency chain, not a wishlist: each stage is built on the guarantees the
one before it already proved.

---

## Status — Now / Next / Later

### Now — the spine, shipped

- ✅ **`kilio-seal`** — real HPKE seal/open (RFC 9180, DHKEM-X25519/HKDF-SHA256/
  ChaCha20Poly1305 via the audited `hpke` crate), receipt-passphrase → per-claim
  key derivation (Argon2id → Ed25519 + X25519), the sealed-sender envelope
  (kotva-MOTE-shaped, AAD-bound routing fields, size-bucketed padding), and the
  anonymous proof-of-work cold-contact gate. Compiles native + `wasm32`.
  21 unit tests green (`cargo test -p kilio-seal`), covering the full two-way
  submit → reply → return-and-read flow, tamper rejection, and wrong-key
  rejection. See [`docs/SECURITY.md`](docs/SECURITY.md) for the primitives in
  depth.

### Next — actively planned, in build order

- ⬜ **`kilio-core`** — domain model (Branch/Claim/Message/Attachment/Handler/
  AuditEvent per decisions.md §5), SQLite sealed store, the `Delivery` /
  `Reachability` / `Identity` seam traits with their local defaults compiled
  in, the ofisi-style `branch_key()` scoped-storage builder, and the single
  `requesterID()` choke point that resolves a handler's branch grants
  server-side and returns `404` (never a distinguishable `403`) on denial.
- ⬜ **`kilio-server` + `kilio-cli`** — axum intake API (reporter-facing, no
  auth, never given the client socket addr) and handler API (session-gated),
  the embedded PWA, owner-gated tunnel start/stop, and the `kilio init |
  serve | tunnel | branch` commands.

### Later — completes the two surfaces

- ⬜ **`web/`** — the JSX PWA: reporter surface (submit, receipt passphrase
  display, return-and-poll, read replies) and handler surface (inbox, open,
  reply), sealing entirely client-side via `kilio-seal` compiled to WASM, a
  service worker, and a TWA manifest for installability.
- ⬜ **`apps/desktop`** — the Tauri v2 handler app. Embeds `kilio-core`
  natively (not WASM) so a single officer can run kilio, decrypt, and reply
  from a laptop with the subprocess tunnel as the default public-reachability
  path — no server infrastructure required at all.
- ⬜ **Integrate → end-to-end vertical slice → security-verify the seams.**
  The slice that must work first, unchanged from decisions.md §9: anonymous
  sealed submit → receipt passphrase → handler decrypts in inbox → sealed
  reply → reporter returns with passphrase and reads it. Everything else
  decorates that spine; nothing ships ahead of it working end to end.

---

## Open questions (decide as we build — decisions.md §10)

These are explicitly undecided. They are not silently deferred; each will get
a decisions.md update the moment it's resolved.

- **Attachment size ceiling & streaming-seal for large files** (e.g. video
  evidence) — `kilio-seal`'s current envelope seals a payload in one shot;
  large attachments likely need a streaming or chunked seal design before
  attachments ship for real.
- **Whether to ship a Tor-hidden-service helper in-box, or just document it.**
  kilio already doesn't preclude running behind Tor (decisions.md §8); the
  question is whether kilio should make that easier itself.
- **Branch key rotation & re-sealing of open claims on rotation.** No design
  yet for what happens to claims already sealed to a branch key that then
  rotates.
- **`KotvaDelivery` cursor/ack semantics against the relay's 48h TTL.** The
  kotva rendezvous mailbox path is a real, working opt-in seam
  (decisions.md §6.1), but `Delivery::collect`'s cursor semantics against a
  relay that only retains messages for 48 hours are not yet specified.

---

## Explicit non-goals

Carried from decisions.md §1 — restated here because a roadmap that doesn't
say what it *won't* build is incomplete:

- **Not a case-management ERP.** kilio manages claims and the conversation
  about them, nothing more (the ofisi lesson: stay narrow).
- **Not a Tor hidden service** by default, though it must never preclude
  running behind one.
- **Not a social network.** No directory, no discovery, no profiles.
- **Not multi-tenant SaaS.** One deployment = one organization. Multi-branch
  (within one org) is the supported axis of scale; multi-org is more
  instances, optionally federated over `KotvaDelivery` — never a shared
  database across organizations.

---

## Related documents

- [`decisions.md`](decisions.md) — the authoritative design record; read it
  before proposing a change to build order or scope.
- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — the crate/component map
  and the three seams.
- [`docs/SECURITY.md`](docs/SECURITY.md) — the privacy spine and threat
  model.
- [`docs/GETTING-STARTED.md`](docs/GETTING-STARTED.md) — build/test today vs.
  the intended operator flow.
