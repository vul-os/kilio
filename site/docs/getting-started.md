# Getting started

kilio is early. This guide is split cleanly into **what works today** — build
and test the crypto spine — and **the intended operator flow** once the rest
of the workspace lands. The status markers are load-bearing; don't skip them.

## What works today

`kilio-seal` and `kilio-core` are wired into the workspace. `kilio-seal` is
the sealed-submission crypto spine: HPKE seal/open, receipt→claim key
derivation, the sealed-sender envelope, and the proof-of-work cold-contact
gate — native and `wasm32`, unit-tested. `kilio-core` layers the domain model
and sealed SQLite store on top of it.

### Prerequisites

- Rust (stable toolchain — `rustup` recommended)
- For wasm32 builds: `rustup target add wasm32-unknown-unknown`

### Build

```bash
git clone https://github.com/vul-os/kilio.git
cd kilio
cargo build --workspace
```

### Test

```bash
./scripts/test-gate.sh        # runs the suite AND asserts it actually ran
```

The gate is what CI runs: it executes each crate's tests and then checks the
harness's own summary line, so a suite that silently stops running cannot
report green. Run the crates directly if you prefer:

```bash
cargo test -p kilio-seal      # sealed-submission crypto + branch key pinning
cargo test -p kilio-core      # sealed store, branch scoping, the seams
```

`kilio-seal` covers branch key generation, the receipt→per-claim-key
derivation, seal/open roundtrips (including the full two-way submit → reply →
return-and-read flow), tampered-envelope rejection, PoW solve/verify,
size-bucket padding, and branch key pinning. `kilio-core` covers the sealed
store's end-to-end flow, branch isolation, and the fail-closed branch-key
checks.

### Build for wasm32

```bash
cargo build -p kilio-seal --target wasm32-unknown-unknown
```

This is the target the reporter's browser will eventually load — sealing
compiled once, run in the browser via WASM, so the browser and the server
never diverge on what "sealed" means.

### Explore the crate

Read `crates/kilio-seal/src/lib.rs` for the public API surface
(`BranchKeys`, `Receipt`/`ClaimKeys`, `Envelope`,
`seal_to_branch`/`seal_to_claim`/`open_with_branch`/`open_with_claim`, `pow`).
Each module (`keys.rs`, `receipt.rs`, `envelope.rs`, `pow.rs`) carries a
doc-comment explaining the *why*, not just the *what* — read those before
touching the crypto. See [Security model](#security) for the same primitives
walked through end to end.

---

## The intended operator flow

Everything below describes the target experience once `kilio-server` and
`kilio-cli` land on top of the already-built `kilio-core`. This section
exists so the shape of "done" is written down before the CLI does it.

### 1. `kilio init` — generate a branch key

```bash
kilio init --name "HR — Global"
```

Generates a fresh `BranchKeys` pair (HPKE recipient keypair + Ed25519 signing
keypair), writes the secret half to local storage, and prints the branch id
and public key a reporter's client will seal to. Run once per branch you
intend to receive claims for.

### 2. `kilio serve` — run the intake + handler server

```bash
kilio serve --port 8080
```

Starts `kilio-server`: the intake API (reporter-facing, unauthenticated by
design), the handler API (session-gated), and the embedded PWA, bound to
`127.0.0.1` by default.

### 3. Click-to-tunnel — make it publicly reachable

From the handler UI (owner-gated), start the built-in tunnel provider — kilio
detects and spawns an installed `cloudflared` / `ngrok` / `frp` binary pinned
to the loopback listen address, and surfaces the assigned public URL. No
config file editing, no reverse proxy required to get a shareable link.

### 4. `kilio branch` — manage branches

```bash
kilio branch add --name "HR — EMEA"
kilio branch list
```

Adds and lists branches — one deployment, many scoped destinations, each
claim sealed to the branch it was submitted to.

## Next steps

- [Self-hosting](#self-hosting) — deploy modes and going public.
- [Security model](#security) — the privacy spine and threat model in depth.
- [Configuration](#configuration) — the pluggable Delivery/Reachability/Identity seams.
