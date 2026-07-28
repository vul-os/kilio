# Configuration

kilio follows the diwan pattern, carried over verbatim: **core defines the
interface and compiles a local default; the fancier adapter is wired only at
the composition root, and core never imports it.** Remove any adapter and the
standalone build still works. Configuration is choosing which implementation
of each seam runs — not a sprawl of flags.

## `Delivery` — where a sealed envelope ends up

The local store write always happens; `Delivery` is the optional *extra* hop.

```rust
// crates/kilio-core/src/delivery.rs — as implemented
pub trait Delivery: Send + Sync {
    fn forward(&self, claim_id: &ClaimId, env: &Envelope) -> Result<(), CoreError>;
    fn label(&self) -> &'static str;
}
```

| Implementation | Default? | What it does | Status |
|---|---|---|---|
| `LocalDelivery` | ✅ default | Forwards nowhere — the sealed claim lives only in the local store. Zero external dependencies. | ✅ built |
| `KotvaDelivery` | opt-in | Carries the target relay + recipient for a content-blind kotva rendezvous mailbox, to forward claims to an external ombudsman or across organisations with no shared server. The envelope is already kotva-MOTE-shaped, so this is a re-wrap, never a re-encrypt. | 🚧 seam only — `forward()` returns `Unsupported`; the async deposit lands with `kilio-server` |

## `Reachability` — making the local app publicly reachable

```rust
// crates/kilio-core/src/reachability.rs — as implemented
pub trait Reachability: Send {
    fn start(&mut self, local_addr: SocketAddr) -> Result<String, CoreError>;
    fn stop(&mut self) -> Result<(), CoreError>;
    fn snapshot(&self) -> TunnelStatus;   // never carries a token
}
```

| Implementation | Default? | What it does | Status |
|---|---|---|---|
| `LocalOnly` | ✅ default | Binds loopback, no exposure. Dev, or behind a reverse proxy the org already runs. | ✅ built |
| `SubprocessTunnel` | opt-in | Records the chosen tunnel binary (`cloudflared` / `ngrok` / `frp`) and enforces the loopback SSRF guard before anything is pinned. | 🚧 guard + choice only — `start()` returns `Unsupported`; spawning the binary and parsing the assigned URL lands with the server/CLI |

There is **no `Ephor` variant in the tree.** Pointing kilio at an
[Ephor](https://github.com/vul-os/ephor) broker is recorded design intent, not
a shipped seam — do not configure it.

**SSRF guard, non-negotiable:** whichever provider runs, it proxies to
exactly **one** configured loopback address, re-checked before every
connection. The inbound request's Host/URL never chooses the target.

## `Identity` / deploy mode

One typed `DeployMode` enum. It has **two** values in `kilio-core`; the
single-officer laptop is a *shape* of `Standalone`, not a third variant:

| Mode | Who runs it | Reachability | Identity |
|---|---|---|---|
| `Standalone` | one officer on a laptop, or an org on a box/VPS | subprocess tunnel, or a reverse proxy the org already runs | local admin(s), Argon2id password → session |
| `Os` | behind a Vulos OS gateway | gateway | gateway-brokered, server-verified session |

`Os` mode is specified to refuse to boot without a configured auth verifier —
the fail-closed boot gate. That gate lands with `kilio-server`; today
`DeployMode` is a typed enum and nothing boots.

## Branch scoping

Every stored object — claim, message, attachment — is addressed through a
single scoped-key builder: `branch_scoped_key(branch_id, name) →
"<branch_id>/<name>"`, with segment sanitisation (no `/`, `\`, `..`). This is
the *only* isolation primitive; there is no second path that reaches storage.
A handler's branch access is resolved server-side from their authenticated
session, never from a client header — a handler for branch A cannot read
branch B's claims, and denied reads return `404`, never a distinguishable
`403`.

Reporters choose their branch at submission (or "global"), and the claim is
HPKE-sealed to *that* branch's public key — so branch scoping is enforced
twice: once by the server's access check, and once unconditionally by the
math, even if the server check were ever bypassed.

## Related documents

- [Self-hosting](#self-hosting) — choosing a deploy mode in practice.
- [Security model](#security) — why these seams exist, in threat-model terms.
