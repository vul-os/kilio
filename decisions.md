# kilio — design decisions

> **kilio** (Swahili: *an outcry, a cry for help*) — a self-hostable, sealed,
> anonymous-first intake app for sensitive claims (harassment, misconduct,
> whistleblowing, safety concerns). Built the Vulos way: one small binary an
> organization runs itself, reachable over a tunnel with no fixed infra,
> usable across branches, standalone by default, decentralized when it wants
> to be.

This document is the authoritative record of *why* kilio is shaped the way it
is. It is written before the code and updated as decisions change. Read it
before touching the crypto or the seams.

---

## 1. The problem

An organization needs to receive reports of sensitive events. The people who
report them are, at the moment they report, often frightened, junior, or
implicated. The single biggest determinant of whether they report at all is
**whether they believe the channel is safe** — specifically:

1. that they do **not** have to identify themselves to begin, and
2. that nobody in the middle (IT, the host, a cloud vendor, a network
   observer) can read what they wrote or tie it back to them, and
3. that they can still have a **two-way conversation** — answer follow-up
   questions, receive an outcome — without ever giving up (1) or (2).

Existing tools (GlobaLeaks, SecureDrop) get this right but assume Tor and a
dedicated operator. kilio targets the ordinary enterprise: an HR or ethics
office that can run one binary, click "make public," and hand out a URL —
while keeping the same source-protecting guarantees.

### Non-goals

- Not a case-management ERP. It manages *claims and the conversation about
  them*, nothing more (the ofisi lesson: stay narrow).
- Not a Tor hidden service (though it must not preclude running behind one).
- Not a social network. There is no directory, no discovery, no profiles.
- Not multi-tenant SaaS. One deployment = one organization (§7).

---

## 2. Shape and stack

One Rust workspace, one shared web frontend, three ways to run it.

```
kilio/
  crates/
    kilio-seal     # sealed-submission crypto. native + wasm32. THE spine.
    kilio-core     # domain model, sealed store, seams. no I/O framework.
    kilio-server   # axum: intake API + handler API + embedded PWA + tunnel
    kilio-cli      # `kilio init | serve | tunnel | branch`
  apps/
    desktop/       # Tauri v2 handler app (embeds kilio-core, decrypts natively)
  web/             # React/JSX PWA — reporter + handler surfaces; seal via WASM
  docs/
```

**Why Rust.** The crypto substrate we most want to align with — kotva-core
(sealed-sender HPKE, the anonymous-but-accountable cold-contact model) — is
Rust, and Tauri (the desktop handler surface the user asked for) is Rust. One
language spans the browser (via WASM), the server, the CLI, and the desktop
app, so **the sealing code exists exactly once** and is never reimplemented
per surface. That single-implementation property is the whole reason to pick
Rust here; it is the opposite of the "five hand-rolled sync engines" mistake.

**Why JSX for the UI, TS/Rust for crypto.** Vulos language policy: app shells
are JSX; protocol/crypto/engine code is typed. kilio honors this — all
sealing is Rust-in-WASM with a thin typed wrapper; the surfaces are plain JSX.

**Three run modes** (the ofisi `DEPLOY_MODE` pattern, one typed enum):

| Mode | Who runs it | Reachability | Identity |
|------|-------------|--------------|----------|
| `desktop` | one officer, on a laptop | subprocess tunnel | local owner |
| `standalone` | org, on a box/VPS | tunnel or reverse proxy | local admin(s) |
| `os` | behind a Vulos OS gateway | gateway | gateway-brokered |

`os` mode **refuses to boot without an auth posture** (ofisi's fail-closed
boot gate — never silently collapse all handlers to one identity).

---

## 3. The privacy spine (this is the point)

Everything else is negotiable; these four properties are not.

### 3.1 Sealed at source

The reporter's claim is HPKE-sealed **in their browser/app, to the destination
branch's public key, before it leaves the device.** The server, the tunnel,
any relay, and the DB only ever hold ciphertext. Decryption happens only in
the handler's Tauri app (or an authenticated handler session holding the
branch private key). This is *honest* privacy — the host cannot read claims
even if compelled, because it never has the key.

- Sealing: **HPKE (RFC 9180)**, mode Base, DHKEM(X25519, HKDF-SHA256),
  HKDF-SHA256, ChaCha20Poly1305. Chosen crate: `hpke-rs` (audited,
  RustCrypto-adjacent, wasm-friendly). We do **not** hand-roll AEAD/KEM.
- The envelope mirrors kotva's MOTE `Envelope`/`Payload` split: a cleartext
  outer (`ciphertext`, ephemeral `enc`, `kind`, `challenge`) and a sealed
  inner (`from`, `body`, `attachments`, `created_at`). The sender identity
  lives *inside* the sealed inner — **sealed sender** — so intermediaries
  never see who sent it, only an ephemeral, unlinkable key.

### 3.2 No mandatory identity

The reporter is never required to supply a name, email, or any contact detail.
The **only** identity is a **receipt passphrase** minted at submission (§4).
Contact details, if the reporter chooses to add them later, are just more
sealed body content — never a required field, never a separate account.

### 3.3 Anonymous two-way channel

The receipt passphrase deterministically derives a per-claim keypair
(Argon2id → Ed25519/X25519 seed). The reporter proves control of a claim by
signing a poll with that key; handler replies are sealed to that same claim
key. So handler ↔ reporter messaging is sealed **both directions** and bound
to nothing but a secret the reporter alone holds. Losing the passphrase =
losing access, by design; there is no recovery, because recovery would mean
someone else could impersonate the reporter.

### 3.4 Metadata minimization

- **No IP logging on the intake path.** The intake handlers never write client
  IP, User-Agent, or timing to any store. (`decisions`: this is enforced by
  the handler *not being given* the socket addr, not by a config flag.)
- **No third-party assets** on the reporter page — no fonts, analytics, CDNs.
  CSP is `default-src 'self'`. Everything is embedded in the binary.
- **Cold-contact gate, not an account.** To stop bulk abuse without demanding
  identity, first contact carries a **proof-of-work** stamp (kotva's PoW
  challenge model). A human reporter pays a second of CPU; a spammer pays for
  every message. Tunable per branch; can be raised under attack.
- **Padded submission sizes.** Ciphertext is padded to size buckets so the
  wire length does not fingerprint the claim.

---

## 4. The receipt passphrase (the identity primitive)

At submission the client generates 128 bits of entropy and renders it as a
**BIP-39 12-word passphrase** ("the only thing you need to keep"). From it:

```
seed        = Argon2id(passphrase, salt = "kilio/receipt/v1" || branch_id,
                        m=64MiB, t=3, p=1) -> 32 bytes
claim_sk    = Ed25519 from seed[0..32]        # signs polls / new messages
claim_x_sk  = X25519 derived from same seed   # receives sealed replies
claim_id    = BLAKE3("kilio/claim-id/v1" || claim_pk)[0..16]  # public handle
```

- The server stores `claim_pk` / `claim_x_pk` / `claim_id` only. It never sees
  the passphrase or the seed.
- To return, the reporter re-enters the passphrase; the client re-derives the
  keys locally and signs a `poll` for `claim_id`. No password is transmitted.
- Argon2id makes offline guessing of a weak-but-real 128-bit phrase
  impractical and blunts a stolen-DB correlation attack.

**Why derive rather than store a random key in the browser?** Because the
reporter may submit from a locked-down or shared machine and must be able to
walk away leaving nothing behind, then return from a different device. The
passphrase is portable state they carry in their head or on paper.

---

## 5. Data model

All fields below marked 🔒 are *inside* the sealed inner payload and never
readable by the server.

```
Branch      { id, name, x_pub, ed_pub, pow_bits, created_at, active }
              # one keypair per branch; claims seal to the branch that owns them.
Claim       { claim_id, branch_id, size_bucket, pow_stamp, created_at,
              status, 🔒category, 🔒title, 🔒body, 🔒attachments[],
              🔒contact_optin }
Message     { id, claim_id, direction(reporter|handler), created_at,
              sealed_blob }   # each turn of the conversation, sealed to the
                              # other side's per-claim / branch key.
Attachment  { id, claim_id, size, sealed_blob }   # content-sealed, size-bucketed
Handler     { id, display_name, ed_pub, role(triager|investigator|admin),
              branch_ids[] }  # who may open which branch's claims
AuditEvent  { id, actor, action, claim_id?, at }  # append-only, content-free
```

- **Server-side row = envelope + routing metadata only.** `status`,
  `created_at`, `size_bucket`, `pow_stamp`, `branch_id` are cleartext because
  the server must route and rate-limit on them. Everything a human wrote is
  sealed.
- **Audit log is content-free**: it records *that* handler X opened claim Y at
  time T, never what Y said. It exists for the org's own governance and for
  the reporter's trust ("who has looked at my report").

### Branch scoping (the ofisi multi-branch pattern)

One deployment, many branches (offices, regions, "corporate/global"). Copied
from ofisi verbatim in spirit:

- **One scoped-key builder.** `branch_key(branch_id, name)` →
  `<branch_id>/<name>` with segment sanitization (no `/`, `\`, `..`). Every
  stored object is addressed through it. This is the single isolation
  primitive.
- **One `requesterID()` choke point.** A handler's branch access is resolved
  server-side from their authenticated session, **never** from a client
  header. A handler for branch A cannot read branch B's claims — denied reads
  return `404`, never "403 exists-but-forbidden" (no existence leak).
- Reporters choose a branch at submission (or submit to "global"); the claim
  seals to *that* branch's key, so even a misrouted claim is unreadable by the
  wrong team.

---

## 6. Seams (thin interface, local default, adapter at the composition root)

The ofisi rule: **core defines the interface and compiles a local default;
the fancy adapter is wired only in `main`, and core never imports it.** Remove
any adapter and the standalone build still works.

### 6.1 `Delivery` — where sealed claims go

```rust
trait Delivery {
    async fn deposit(&self, branch_id: &BranchId, envelope: &Envelope) -> Result<Receipt>;
    async fn collect(&self, branch_id: &BranchId, since: Cursor) -> Result<Vec<Envelope>>;
}
```

- **`LocalDelivery` (default):** write the sealed envelope to the local
  SQLite store. Zero dependencies. This is what `standalone`/`desktop` use.
- **`KotvaDelivery` (opt-in, decentralized):** deposit the sealed envelope as
  an opaque blob to a **kotva rendezvous mailbox** (`POST {relay}/mailbox/{to}`,
  content-blind, key-addressed — the working ephor Go relay + `@vulos/relay-client`).
  Used to forward claims to an **external ombudsman** or across orgs without a
  shared server. The relay never sees plaintext (already sealed at source).
  This is the "decentralized, using kotva/ephor" path — real, but optional.

kilio's envelope is deliberately kotva-MOTE-shaped so `KotvaDelivery` is a
re-wrap, not a re-encrypt.

### 6.2 `Reachability` — make the local app publicly reachable

Mirrors wede's `Provider` interface (`start/stop/public_url/snapshot`),
mechanism-agnostic:

```rust
trait Reachability {
    async fn start(&self, local_addr: SocketAddr) -> Result<PublicUrl>;
    async fn stop(&self) -> Result<()>;
    fn snapshot(&self) -> TunnelStatus;   // token always redacted
}
```

- **`LocalOnly` (default):** bind `127.0.0.1`, no exposure. For dev / behind a
  reverse proxy the org already runs.
- **`SubprocessTunnel` (working default for "click to go public"):** spawn a
  detected tunnel binary — `cloudflared` / `ngrok` / `frp` — pinned to the
  loopback listen addr, parse the assigned public URL. This is the honest,
  runnable-today ngrok-like path (wede's built-in relay needs a relay *server*
  that isn't in-tree yet).
- **`Ephor` (seam, stubbed):** the wede sovereign reverse-tunnel agent,
  wired the day an Ephor server is available.

**SSRF guard (non-negotiable, from wede):** whichever provider runs, it
proxies to **exactly one** configured loopback address, re-checked before
every connection. The inbound request's Host/URL never chooses the target.

### 6.3 `Identity` / deploy mode

`standalone`/`desktop`: local handler accounts (Argon2id password → session).
`os`: identity brokered by the Vulos OS gateway via server-verified session,
never a client header. Boot gate refuses `os` without a configured verifier.

---

## 7. Why one-org-per-deployment (not multi-tenant)

A tool that promises "the host cannot read your claim" must not then run every
org's claims through one shared operator's database. Instance-per-org keeps the
trust boundary honest: **the org that receives the claims is the only party
with the keys, and it runs its own box.** Multi-branch (within one org) is the
supported axis of scale; multi-*org* is achieved by running more instances,
optionally federated over kotva delivery. This is exactly ofisi's stance and
it is a *security* decision, not just an ops one.

---

## 8. Threat model (abbreviated)

| Adversary | Capability | kilio's answer |
|-----------|-----------|----------------|
| Network observer / tunnel operator | sees all traffic | TLS + sealed-at-source; only ciphertext + size-bucket transit |
| Malicious/curious host admin | full DB + disk | claims sealed to branch key held only by handlers; DB is ciphertext |
| Compelled host (subpoena) | can be forced to hand over data | can only hand over ciphertext + content-free metadata; no keys, no IPs |
| Retaliatory insider (a handler) | valid handler creds for branch A | branch scoping + `404`-on-deny; content-free audit log records every open |
| Spammer / DoS | floods intake | per-branch PoW cold-contact stamp; size caps; no unauth injection |
| Reporter deanonymization | correlate metadata | no IP/UA logging, size padding, sealed sender, PoW is unlinkable |
| Passphrase thief | steals the receipt phrase | full access to that one claim — accepted; mitigated by Argon2id + user guidance |

**Explicit residual risks** (write them down, don't pretend): a
compromised/hostile *client build* could exfiltrate before sealing → mitigated
by CSP, embedded assets, subresource integrity, reproducible builds (roadmap).
Traffic-analysis of *when* someone submits is not defeated without Tor;
documented, and running kilio behind a Tor hidden service is supported.

---

## 9. Build order

1. `kilio-seal` — real HPKE seal/open, receipt derivation, PoW, sealed-sender
   envelope; native + wasm32; property tests + KATs. **(Opus-authored.)**
2. `kilio-core` — model, SQLite sealed store, seams with local defaults,
   branch scoping + `requesterID` choke point.
3. `kilio-server` + `kilio-cli` — axum intake/handler APIs, embedded PWA,
   owner-gated tunnel control.
4. `web/` — JSX PWA (reporter + handler), WASM sealing, SW + TWA manifest.
5. `apps/desktop` — Tauri v2 handler app.
6. Integrate → end-to-end vertical slice → security-verify the seams.

The vertical slice that must work first: **anonymous sealed submit → receipt
passphrase → handler decrypts in inbox → sealed reply → reporter returns with
passphrase and reads it.** Everything else decorates that spine.

---

## 10. Open questions (decide as we build)

- Attachment size ceiling & streaming-seal for large files (video evidence).
- Whether to ship a Tor-hidden-service helper in-box or document it.
- Branch key rotation & re-sealing of open claims on rotation.
- `KotvaDelivery` cursor/ack semantics against the relay's 48h TTL.
