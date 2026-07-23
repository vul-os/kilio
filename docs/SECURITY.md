# kilio Security Model

kilio exists to make one promise credible: *the person reporting a sensitive
claim does not have to trust the host, only the math.* This document states
that promise precisely — what is sealed, what is not, which primitives do the
sealing, and what residual risk remains once you've done all of that
honestly. For the vulnerability-disclosure process, see the repo-level
[`SECURITY.md`](../SECURITY.md); this document is the model, that one is the
contact form.

> **Status.** The primitives and properties described here are implemented
> and tested in [`kilio-seal`](../crates/kilio-seal) today. Everything that
> *uses* them — the server, the store, the branch-scoping enforcement, the
> desktop app — is specified but not yet built (see
> [ARCHITECTURE.md](ARCHITECTURE.md) and [`../ROADMAP.md`](../ROADMAP.md)).
> A crypto primitive being correct is not the same claim as a deployed system
> being secure; read this alongside that status honestly.

---

## The privacy spine

Four properties are non-negotiable. Everything else about kilio is
negotiable; these are not (decisions.md §3).

### 1. Sealed at source

A claim is HPKE-sealed **in the reporter's browser or app, to the destination
branch's public key, before it leaves the device.** The server, any tunnel,
any relay, and the database only ever hold ciphertext. Decryption happens
only in the handler's Tauri app or an authenticated handler session holding
the branch private key. This is *honest* privacy in the specific sense that
matters for a compelled-disclosure threat: **the host cannot read claims even
if compelled, because it never has the key.**

### 2. No mandatory identity

The reporter is never required to supply a name, email, or contact detail.
The only identity is a **receipt passphrase**, minted client-side at
submission (below). Contact details, if a reporter chooses to add them, are
just more sealed body content — never a required field, never a separate
account.

### 3. Anonymous two-way channel

The receipt passphrase deterministically derives a per-claim keypair. The
reporter proves control of a claim by signing a poll with that key; handler
replies are sealed to that same claim key. Handler ↔ reporter messaging is
sealed in **both** directions and bound to nothing but a secret only the
reporter holds. Losing the passphrase means losing access, by design —
recovery would mean someone else could impersonate the reporter, which is a
worse failure mode than the one it would fix.

### 4. Metadata minimization

- **No IP logging on the intake path.** The intake handlers are never given
  the client's socket address at all — this is enforced by omission, not by a
  toggle that could be flipped back on.
- **No third-party assets** on the reporter page: no fonts, no analytics, no
  CDNs. `Content-Security-Policy: default-src 'self'`. Everything is embedded
  in the binary.
- **Cold-contact proof-of-work, not an account** (below) — stops bulk abuse
  without demanding identity.
- **Padded submission sizes** — ciphertext is padded to size buckets so wire
  length cannot fingerprint the claim.

---

## Threat model

From decisions.md §8, reproduced as the canonical table:

| Adversary | Capability | kilio's answer |
|---|---|---|
| Network observer / tunnel operator | sees all traffic | TLS + sealed-at-source; only ciphertext + size-bucket transit |
| Malicious/curious host admin | full DB + disk | claims sealed to branch key held only by handlers; DB is ciphertext |
| Compelled host (subpoena) | can be forced to hand over data | can only hand over ciphertext + content-free metadata; no keys, no IPs |
| Retaliatory insider (a handler) | valid handler creds for branch A | branch scoping + 404-on-deny; content-free audit log records every open |
| Spammer / DoS | floods intake | per-branch PoW cold-contact stamp; size caps; no unauth injection |
| Reporter deanonymization | correlate metadata | no IP/UA logging, size padding, sealed sender, PoW is unlinkable |
| Passphrase thief | steals the receipt phrase | full access to that one claim — accepted; mitigated by Argon2id + user guidance |

---

## Crypto primitives

No primitive is hand-rolled; only their composition is kilio's. Every choice
below is implemented in [`kilio-seal`](../crates/kilio-seal/src/lib.rs) today.

| Primitive | Used for | Crate |
|---|---|---|
| **HPKE (RFC 9180)**, mode Base, DHKEM(X25519, HKDF-SHA256), HKDF-SHA256, ChaCha20Poly1305 | Sealing a claim/message to a branch or claim public key | the audited `hpke` crate |
| **Ed25519** | Branch identity signatures; claim-control signatures (proving control of a claim without transmitting the passphrase) | `ed25519-dalek` |
| **Argon2id** | Deriving a per-claim key root from the receipt passphrase (memory-hard, resists offline guessing) | `argon2` |
| **BLAKE3** | Domain-separated KDF sub-keys, branch/claim id derivation, PoW challenge derivation | `blake3` |
| **BIP-39** | Rendering 128 bits of entropy as a 12-word receipt phrase | `bip39` |
| **CBOR (ciborium)** | Deterministic envelope/AAD encoding | `ciborium` |

We do not hand-roll AEAD or KEM constructions. The one place kilio composes
primitives itself is the receipt→claim-key derivation chain below, and the
envelope's associated-data binding — both are small, tested, and documented
here precisely so they can be scrutinized.

### The sealed envelope

The envelope mirrors kotva's MOTE `Envelope`/`Payload` split
([`envelope.rs`](../crates/kilio-seal/src/envelope.rs)):

- **Outer, cleartext:** `v` (version), `kind` (`Submission` /
  `ReporterMessage` / `HandlerReply`), `recipient` (a `Branch` or `Claim`
  tag), `enc` (the HPKE encapsulated key), `ciphertext`, `size_bucket`.
- **Inner, sealed:** `from` (the sender's public claim identity — present
  only on a `Submission`, absent otherwise), `created_at`, `body`.

**Sealed sender.** The reporter's identity lives *inside* the sealed inner
payload, never in the cleartext outer envelope. Intermediaries — server,
tunnel, relay — see only an ephemeral, unlinkable recipient tag, never who
sent a message.

**Every cleartext routing field is bound into the AEAD's associated data**
(`v`, `kind`, `recipient` tag, `size_bucket`). None of it can be altered in
transit without decryption failing — a network observer or a malicious host
cannot re-tag a `Submission` as a `HandlerReply`, or reroute an envelope to a
different branch, without the open failing.

**Size-bucketed padding.** Ciphertext is zero-padded up to the smallest of a
fixed set of buckets (4 KiB … 4 MiB, then 4 MiB steps above that) before
sealing, so wire length cannot fingerprint the claim's true size.

### Receipt → per-claim-key derivation

The reporter's *entire* identity is a 12-word BIP-39 phrase minted from 128
bits of OS-CSPRNG entropy. From it (decisions.md §4,
[`receipt.rs`](../crates/kilio-seal/src/receipt.rs)):

```
seed        = Argon2id(passphrase, salt = "kilio/receipt-salt/v1" || branch_id,
                        m=64MiB, t=3, p=1) -> 32 bytes
sign_seed   = BLAKE3::keyed_hash(seed, "kilio/claim/sign/v1")   # Ed25519 — proves control
recip_ikm   = BLAKE3::keyed_hash(seed, "kilio/claim/recip/v1")  # X25519 — receives sealed replies
claim_id    = BLAKE3("kilio/claim-id/v1" || claim_pk)[0..16]     # public handle
```

- **The branch id is folded into the Argon2id salt.** The same phrase
  produces *different* keys per branch, so a phrase compromised for one
  branch reveals nothing about that reporter's claims to any other branch.
- **The server never sees the passphrase or the seed** — only the derived
  public keys and the claim id. To return, the reporter re-enters the
  phrase; keys are re-derived **locally**, and a `poll` is signed with the
  re-derived key. No password is ever transmitted.
- **Argon2id's cost (64 MiB / 3 passes / 1 lane)** is deliberately expensive:
  it makes offline guessing of a weak-but-real 128-bit phrase impractical and
  blunts a stolen-DB correlation attack. Do not lower this without updating
  decisions.md §4 and this document together.
- **Why derive instead of storing a random browser key:** a reporter may
  submit from a locked-down or shared machine and needs to walk away leaving
  nothing behind, then return from a different device later. The passphrase
  is portable state carried in the reporter's head or on paper — it is the
  only state that has to survive.

### The proof-of-work cold-contact gate

kilio refuses fully-unauthenticated injection but also refuses to demand an
identity from a reporter. The reconciliation (kotva §9.2's model,
[`pow.rs`](../crates/kilio-seal/src/pow.rs)): first contact carries a small
proof-of-work stamp bound to *that exact sealed message* —

```
challenge = BLAKE3::derive_key("kilio/pow-challenge/v1", env.v || env.enc || env.ciphertext)
stamp     = { nonce, bits }   where BLAKE3(challenge || nonce) has ≥ bits leading zero bits
```

A human reporter pays a fraction of a second of CPU, once, per message. A
spammer pays it for **every** message, and the cost cannot be amortized
across submissions because the challenge is derived from the envelope's own
ciphertext — a stamp cannot be pre-computed and replayed onto a different
message. Difficulty is tunable per branch (18–22 bits ≈ well under a second
to multiple seconds on a laptop) and can be raised under active abuse.
Verification checks the stamp's claimed `bits` against the actual leading-zero
count, so a stamp cannot lie about how much work it represents.

---

## Explicit residual risks

Writing these down rather than pretending they don't exist (decisions.md §8):

- **A compromised or hostile client build could exfiltrate before sealing.**
  If the JS/WASM the reporter's browser actually loaded has been tampered
  with, sealing happens *after* the tamper, and no server-side control can
  detect that. Mitigated by CSP, embedded assets (no third-party script
  injection surface), subresource integrity, and reproducible builds — the
  last of these is roadmap, not shipped.
- **Traffic analysis of *when* someone submits is not defeated without Tor.**
  kilio does not hide submission timing from a network observer positioned
  to watch the reporter's own connection. Running kilio behind a Tor hidden
  service is explicitly supported for reporters who need that property; it
  is not the default.
- **A passphrase thief gets that one claim, fully.** This is an accepted
  design tradeoff (§3.3 above), not a bug — the alternative, a recovery
  mechanism, would mean someone other than the reporter could regain access,
  which defeats the property recovery would supposedly restore.
- **Argon2id parameters are a cost tradeoff, not a guarantee.** 64 MiB / 3
  passes raises the bar for offline guessing of a 128-bit phrase; it does not
  make guessing impossible against an adversary with enough resources. The
  phrase's 128 bits of entropy is doing the real work — Argon2id makes
  *cheap* correlation of a stolen DB expensive, not impossible.
- **The `os` deploy mode's security depends on the Vulos OS gateway's
  session verification being correctly configured.** kilio's own boot gate
  refuses to start in `os` mode without a configured verifier, but a
  misconfigured *verifier* (not kilio's problem to solve) would weaken
  identity brokering it depends on.

---

## Reporting a vulnerability

See the repo-level [`SECURITY.md`](../SECURITY.md) for the disclosure process
and supported versions. Report cryptographic issues — anything touching
`kilio-seal`, the envelope's AAD binding, the receipt derivation chain, or
the PoW challenge construction — with priority; that crate is the entire
trust boundary this document describes.
