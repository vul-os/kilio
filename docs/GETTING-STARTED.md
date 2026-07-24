# Getting Started with kilio

kilio is early. This guide is split cleanly into **what works today** — build
and test the crypto spine — and **the intended operator flow** once the rest
of the workspace lands. Do not skip the status markers; they are load-bearing.

---

## What works today

Only [`kilio-seal`](../crates/kilio-seal) is wired into the workspace right
now. It is the sealed-submission crypto spine: HPKE seal/open, receipt→claim
key derivation, the sealed-sender envelope, and the proof-of-work cold-contact
gate — native and `wasm32`, unit-tested.

### Prerequisites

- Rust (stable toolchain — `rustup` recommended)
- For wasm32 builds: `rustup target add wasm32-unknown-unknown`

### Build

```bash
git clone https://github.com/vul-os/kilio.git
cd kilio
cargo build --workspace
```

> **Today, `cargo build --workspace` builds exactly one crate**,
> `kilio-seal` — the workspace member list in the root `Cargo.toml` has not
> grown yet. This is expected; it will build more as `kilio-core`,
> `kilio-server`, and `kilio-cli` land (see [`../ROADMAP.md`](../ROADMAP.md)).

### Test

```bash
cargo test -p kilio-seal
```

This runs the full crypto test suite: branch key generation, the receipt→
per-claim-key derivation, seal/open roundtrips (including the full two-way
submit → reply → return-and-read flow), tampered-envelope rejection, PoW
solve/verify, and size-bucket padding. All 21 tests pass as of this writing.

### Build for wasm32

```bash
cargo build -p kilio-seal --target wasm32-unknown-unknown
```

This is the target the reporter's browser will eventually load — sealing
compiled once, run in the browser via WASM, so the browser and the server
never diverge on what "sealed" means.

### Explore the crate

Read [`crates/kilio-seal/src/lib.rs`](../crates/kilio-seal/src/lib.rs) for
the public API surface (`BranchKeys`, `Receipt`/`ClaimKeys`, `Envelope`,
`seal_to_branch`/`seal_to_claim`/`open_with_branch`/`open_with_claim`, `pow`).
Each module (`keys.rs`, `receipt.rs`, `envelope.rs`, `pow.rs`) carries a
doc-comment explaining the *why*, not just the *what* — read those before
touching the crypto. [`docs/SECURITY.md`](SECURITY.md) walks through the same
primitives end to end.

---

## The intended operator flow (aspirational — not yet built)

Everything below describes the target experience once
[`kilio-core`](../decisions.md#9-build-order), `kilio-server`, and
`kilio-cli` exist. None of these commands work yet. This section exists so
the shape of "done" is written down before the code, per decisions.md's own
stated method.

### 1. `kilio init` — generate a branch key

```bash
kilio init --name "HR — Global"
```

Generates a fresh `BranchKeys` pair (HPKE recipient keypair + Ed25519 signing
keypair, per `kilio-seal`), writes the secret half to local storage, and
prints the branch id and public key a reporter's client will seal to. Run
once per branch you intend to receive claims for.

### 2. `kilio serve` — run the intake + handler server

```bash
kilio serve --port 8080
```

Starts `kilio-server`: the intake API (reporter-facing, unauthenticated by
design), the handler API (session-gated), and the embedded PWA, bound to
`127.0.0.1` by default (the `LocalOnly` `Reachability` default — see
[ARCHITECTURE.md](ARCHITECTURE.md#2-reachability--making-the-local-app-publicly-reachable)).

### 3. Click-to-tunnel — make it publicly reachable

From the handler UI (owner-gated), start the `SubprocessTunnel` provider —
kilio detects and spawns an installed `cloudflared` / `ngrok` / `frp` binary
pinned to the loopback listen address, and surfaces the assigned public URL.
No config file editing, no reverse proxy required to get a shareable link for
"make public, hand out a URL" (decisions.md §1). A sovereign `Ephor`
provider is a stubbed seam, wired the day an Ephor server exists to
point at.

### 4. `kilio branch` — manage branches

```bash
kilio branch add --name "HR — EMEA"
kilio branch list
```

Adds and lists branches for the ofisi-style multi-branch pattern
(decisions.md §5) — one deployment, many scoped destinations, each claim
sealed to the branch it was submitted to.

---

## Next steps

- [`ARCHITECTURE.md`](ARCHITECTURE.md) — crate map, seams, data flow, branch
  scoping.
- [`SECURITY.md`](SECURITY.md) — the privacy spine, threat model, crypto
  primitives in depth.
- [`../ROADMAP.md`](../ROADMAP.md) — phased build order and open questions.
- [`../decisions.md`](../decisions.md) — the authoritative design record.
