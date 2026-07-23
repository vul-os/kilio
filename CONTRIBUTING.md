# Contributing to kilio

## Code of Conduct

We follow the [Contributor Covenant v2.1](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).

## Read this first

[`decisions.md`](decisions.md) is the authoritative design record — read it
before touching the crypto or the seams. It is written before the code and
updated as decisions change; a PR that contradicts it without updating it is
not ready to merge.

## Dev environment setup

Requirements: Rust (stable), `wasm32-unknown-unknown` target for wasm builds.

```bash
rustup target add wasm32-unknown-unknown

cargo build --workspace
cargo test -p kilio-seal
```

See [`docs/GETTING-STARTED.md`](docs/GETTING-STARTED.md) for the full build
and test walkthrough, including what's built today versus planned.

## Branch and PR conventions

- Branch off `main`. Name: `feat/description`, `fix/description`,
  `chore/description`.
- One logical change per PR. Keep diffs reviewable.
- PRs require at least one approving review.
- Squash-merge preferred.

## Commit message style

Conventional Commits welcome, not required:

```
feat(seal): add streaming seal for large attachments
fix(receipt): reject phrases with invalid BIP-39 checksum earlier
chore: bump hpke to 0.13
```

## Testing expectations

Before opening a PR:

```bash
cargo build --workspace
cargo test --workspace
cargo clippy --workspace -- -D warnings
cargo fmt --check
```

**Crypto-touching changes get extra scrutiny.** Anything in `kilio-seal` —
the HPKE seal/open path, the receipt→claim-key derivation, the envelope's
associated-data binding, the PoW challenge construction — needs property
tests or known-answer tests alongside the change, not just a happy-path unit
test. `#![forbid(unsafe_code)]` is set at the crate root; do not remove it.

## Scope: what we say yes and no to

### Yes

- Bug fixes and security improvements.
- New crates that follow the build order in
  [`ROADMAP.md`](ROADMAP.md)/decisions.md §9 (`kilio-core` is next).
- Seam implementations (`Delivery`, `Reachability`) that keep the local
  default dependency-free and wire the adapter only at the composition root.
- Tests and documentation.

### No — frozen invariants

- **No hand-rolled AEAD/KEM/signature primitives.** Use audited crates
  (`hpke`, `ed25519-dalek`, `argon2`, `blake3`) — see
  [`docs/SECURITY.md`](docs/SECURITY.md) for the exact set.
- **No `unsafe` code** in `kilio-seal`.
- **No `.tsx` files.** App shells (`web/`, `apps/desktop`'s frontend) are
  JSX only; crypto/protocol/engine code is Rust. See decisions.md §2.
- **No multi-tenant storage.** One deployment = one organization
  (decisions.md §7). Multi-org is more instances, not shared rows.
- **No IP/User-Agent logging on the intake path.** This is enforced by the
  intake handler never receiving the socket address — do not add a config
  flag that could turn it back on.
- **No core→adapter imports.** `kilio-core` must never import a concrete
  `Delivery`/`Reachability` adapter; adapters are wired only in `main`.

## Licensing

kilio is dual-licensed MIT OR Apache-2.0. Contributions are made under both;
no CLA required.
