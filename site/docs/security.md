# Security model

kilio exists to make one promise credible: *the person reporting a sensitive
claim does not have to trust the host, only the math.* This document states
that promise precisely — what is sealed, what is not, which primitives do the
sealing, and what residual risk remains once you've done all of that
honestly.

> **Status.** The primitives and properties described here are implemented
> and tested in `kilio-seal` today; the sealed store and the branch-scoping
> enforcement are implemented and tested in `kilio-core`. The server, CLI, and
> desktop app — the parts that would *expose* any of this over a socket — are
> specified but landing in stages. A crypto primitive being correct is not
> the same claim as a deployed system being secure; read this alongside the
> repo's status honestly.

## The privacy spine

Four properties are non-negotiable. Everything else about kilio is
negotiable; these are not.

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
submission. Contact details, if a reporter chooses to add them, are just
more sealed body content — never a required field, never a separate account.

### 3. Anonymous two-way channel

The receipt passphrase deterministically derives a per-claim keypair. The
reporter proves control of a claim by signing a poll with that key; handler
replies are sealed to that same claim key. Handler ↔ reporter messaging is
sealed in **both** directions and bound to nothing but a secret only the
reporter holds. Losing the passphrase means losing access, by design —
recovery would mean someone else could impersonate the reporter, which is a
worse failure mode than the one it would fix.

### 4. Metadata minimisation

- **No IP logging on the intake path.** The intake handlers are never given
  the client's socket address at all — enforced by omission, not by a toggle
  that could be flipped back on.
- **No third-party assets** on the reporter page: no fonts, no analytics, no
  CDNs. `Content-Security-Policy: default-src 'self'`. Everything is embedded
  in the binary.
- **Cold-contact proof-of-work, not an account** — stops bulk abuse without
  demanding identity.
- **Padded submission sizes** — ciphertext is padded to size buckets so wire
  length cannot fingerprint the claim.

## Threat model

| Adversary | Capability | kilio's answer |
|---|---|---|
| Network observer / tunnel operator | sees all traffic | TLS + sealed-at-source; only ciphertext + size-bucket transit |
| Malicious/curious host admin | full DB + disk | claims sealed to branch key held only by handlers; DB is ciphertext |
| Key-substituting host | serves a branch key it controls, so new claims seal to *it* | branch key **pinning**: the key is accepted once, and thereafter only a rekey countersigned by the pinned key can change it. Host-side, stored branch keys are immutable and a branch id must derive from its keys. **First contact is still TOFU** — compare the `branch_id` out of band if you need more |
| Compelled host (subpoena) | can be forced to hand over data | can only hand over ciphertext + content-free metadata; no keys, no IPs |
| Retaliatory insider (a handler) | valid handler creds for branch A | branch scoping + 404-on-deny; content-free audit log records every open |
| Spammer / DoS | floods intake | per-branch PoW cold-contact stamp; size caps; no unauth injection |
| Reporter deanonymisation | correlate metadata | no IP/UA logging, size padding, sealed sender, PoW is unlinkable |
| Passphrase thief | steals the receipt phrase | full access to that one claim — accepted; mitigated by Argon2id + user guidance |

## Crypto primitives

No primitive is hand-rolled; only their composition is kilio's.

| Primitive | Used for | Crate |
|---|---|---|
| **HPKE (RFC 9180)**, mode Base, DHKEM(X25519, HKDF-SHA256), HKDF-SHA256, ChaCha20Poly1305 | Sealing a claim/message to a branch or claim public key | the audited `hpke` crate |
| **Ed25519** | Branch identity signatures; claim-control signatures (proving control of a claim without transmitting the passphrase) | `ed25519-dalek` |
| **Argon2id** | Deriving a per-claim key root from the receipt passphrase (memory-hard, resists offline guessing) | `argon2` |
| **BLAKE3** | Domain-separated KDF sub-keys, branch/claim id derivation, PoW challenge derivation | `blake3` |
| **BIP-39** | Rendering 128 bits of entropy as a 12-word receipt phrase | `bip39` |
| **CBOR (ciborium)** | Deterministic envelope/AAD encoding | `ciborium` |

### The sealed envelope

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

### The proof-of-work cold-contact gate

```
challenge = BLAKE3::derive_key("kilio/pow-challenge/v1", env.v || env.enc || env.ciphertext)
stamp     = { nonce, bits }   where BLAKE3(challenge || nonce) has ≥ bits leading zero bits
```

A human reporter pays a fraction of a second of CPU, once, per message. A
spammer pays it for **every** message, and the cost cannot be amortised
across submissions because the challenge is derived from the envelope's own
ciphertext — a stamp cannot be pre-computed and replayed onto a different
message. Difficulty is tunable per branch and can be raised under active
abuse.

### Branch key pinning

Sealing to the right key is what makes every property above matter, and the
reporter gets that key from the host. A self-consistency check on the key's own
id proves nothing: a host that substitutes a key it controls substitutes the
matching id too. So kilio pins, following aql's pairing discipline rather than
a looser first-key-wins:

```
descriptor  = { v, branch_id, kem_public, sign_public, name, pow_bits, epoch, issued_at }
signature   = Ed25519(branch_sign_key, "kilio/branch-descriptor/v1\0" || CBOR(descriptor))
rekey cert  = Ed25519(OLD branch_sign_key, "kilio/branch-rekey/v1\0" || CBOR(prev_id, next_descriptor))
```

- The key is accepted **once**; every later fetch is checked against the pin,
  and any key change or epoch rollback is refused.
- A pinned key rotates only via a rekey certificate countersigned by the
  currently pinned key, with a strictly greater epoch. The two signing domains
  are separated, so a descriptor signature can never be replayed as a rekey.
- The only other escape is an explicit unpin — the factory-reset analogue,
  never reachable from a fetch path.
- Host-side, stored branch keys are immutable and a branch id must derive from
  the keys it publishes.

## Explicit residual risks

Writing these down rather than pretending they don't exist:

- **A compromised or hostile client build could exfiltrate before sealing.**
  If the JS/WASM the reporter's browser actually loaded has been tampered
  with, sealing happens *after* the tamper, and no server-side control can
  detect that. Mitigated by CSP, embedded assets (no third-party script
  injection surface), and subresource integrity — reproducible builds are
  roadmap, not shipped.
- **Traffic analysis of *when* someone submits is not defeated without Tor.**
  kilio does not hide submission timing from a network observer positioned
  to watch the reporter's own connection. Running kilio behind a Tor hidden
  service is explicitly supported for reporters who need that property; it
  is not the default.
- **A passphrase thief gets that one claim, fully.** This is an accepted
  design tradeoff, not a bug — the alternative, a recovery mechanism, would
  mean someone other than the reporter could regain access, which defeats the
  property recovery would supposedly restore.
- **Argon2id parameters are a cost tradeoff, not a guarantee.** 64 MiB / 3
  passes raises the bar for offline guessing of a 128-bit phrase; it does not
  make guessing impossible against an adversary with enough resources.
- **Pinning is trust-on-first-use.** A host hostile from the very first fetch
  is pinned to the attacker's key, and every check afterwards faithfully
  enforces the wrong key. Pinning defeats *later* substitution, not a bad first
  contact — compare the `branch_id` against a value published out of band.
- **Nothing offline-verifiable authorizes a handler yet.** A handler's branch
  grant is resolved in-process; there is no signed, independently verifiable
  capability for "this handler may open branch X". That must land with
  `kilio-server`.
- **The `Os` deploy mode's security depends on the Vulos OS gateway's
  session verification being correctly configured.** The fail-closed boot gate
  that refuses to start without a configured verifier is specified, not built —
  it lands with `kilio-server`. A misconfigured *verifier* is not kilio's
  problem to solve.

## Reporting a vulnerability

See the repo-level `SECURITY.md` for the disclosure process and supported
versions. Report cryptographic issues — anything touching `kilio-seal`, the
envelope's AAD binding, the receipt derivation chain, or the PoW challenge
construction — with priority; that crate is the entire trust boundary this
document describes.
