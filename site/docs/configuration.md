# Configuration

kilio follows the ofisi pattern, carried over verbatim: **core defines the
interface and compiles a local default; the fancier adapter is wired only at
the composition root, and core never imports it.** Remove any adapter and the
standalone build still works. Configuration is choosing which implementation
of each seam runs — not a sprawl of flags.

## `Delivery` — where a sealed envelope ends up

```rust
trait Delivery {
    async fn deposit(&self, branch_id: &BranchId, envelope: &Envelope) -> Result<Receipt>;
    async fn collect(&self, branch_id: &BranchId, since: Cursor) -> Result<Vec<Envelope>>;
}
```

| Implementation | Default? | What it does |
|---|---|---|
| `LocalDelivery` | ✅ default | Writes the sealed envelope to the local SQLite store. Zero external dependencies — what `standalone`/`desktop` use. |
| `KotvaDelivery` | opt-in | Deposits the envelope as an opaque, content-blind blob to a kotva rendezvous mailbox. Used to forward claims to an external ombudsman or across organisations with no shared server. The envelope is already kotva-MOTE-shaped, so this is a re-wrap, never a re-encrypt. |

## `Reachability` — making the local app publicly reachable

```rust
trait Reachability {
    async fn start(&self, local_addr: SocketAddr) -> Result<PublicUrl>;
    async fn stop(&self) -> Result<()>;
    fn snapshot(&self) -> TunnelStatus;   // token always redacted
}
```

| Implementation | Default? | What it does |
|---|---|---|
| `LocalOnly` | ✅ default | Binds `127.0.0.1`, no exposure. Dev, or behind a reverse proxy the org already runs. |
| `SubprocessTunnel` | ✅ working "click to go public" path | Spawns a detected tunnel binary (`cloudflared` / `ngrok` / `frp`) pinned to the loopback listen address, parses the assigned public URL. |
| `Ephor` | ⬜ stubbed seam | A sovereign reverse-tunnel agent, wired the day an Ephor server is available to point at. |

**SSRF guard, non-negotiable:** whichever provider runs, it proxies to
exactly **one** configured loopback address, re-checked before every
connection. The inbound request's Host/URL never chooses the target.

## `Identity` / deploy mode

One typed `DEPLOY_MODE` enum, three values:

| Mode | Who runs it | Reachability | Identity |
|---|---|---|---|
| `desktop` | one officer, on a laptop | subprocess tunnel | local owner |
| `standalone` | org, on a box/VPS | tunnel or reverse proxy | local admin(s), Argon2id password → session |
| `os` | behind a Vulos OS gateway | gateway | gateway-brokered, server-verified session |

`os` mode refuses to boot without a configured auth verifier — the fail-closed
boot gate.

## Branch scoping

Every stored object — claim, message, attachment — is addressed through a
single scoped-key builder: `branch_key(branch_id, name) →
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
