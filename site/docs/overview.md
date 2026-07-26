# Overview

**kilio** is sealed, anonymous-first intake for sensitive claims — harassment,
misconduct, whistleblowing, and safety reports. It replaces the "email HR" /
third-party ethics-hotline model, where a reporter has to trust a vendor's
cloud, hand over their identity to begin, and hope nobody in the middle is
reading along.

The difference from every incumbent: **there is no central server that can
read anything.** An organisation runs its own kilio instance — a laptop with a
tunnel, a small VPS, a Raspberry Pi, anything that runs the binary. A reporter
opens the public intake page, writes their claim, and their browser HPKE-seals
it to the destination team's public key before it ever leaves the device. The
instance — and any tunnel or relay between them — only ever stores
**ciphertext**. It is decrypted only inside a handler's app, with a key the
server never holds.

kilio is **completely standalone**. It depends on no hosted service and no
account, not even ours. It runs the same on macOS, Linux, and Windows. Optional,
off-by-default seams let it reach further when you want them — see
[Configuration](#configuration) — but nothing about the core needs them.

> **Status: 0.1.0 — early.** The sealed-crypto core (`kilio-seal`) and the
> sealed store + branch scoping (`kilio-core`) have landed and are tested,
> along with the full design record (`decisions.md`). The server, CLI, web
> surfaces, and Tauri app are in progress — see the repo's `ROADMAP.md`.

## How it works

A reporter never authenticates. A handler always does. Between them runs a
sealed, anonymous, two-way channel keyed only by a secret the reporter holds.

1. **Seal at source.** The reporter's device encrypts the claim to the
   branch's public key. The host stores ciphertext it cannot read.
2. **Receipt passphrase.** Twelve words are the reporter's only identity; they
   derive a per-claim keypair (memory-hard) for return visits and sealed
   replies. No email, no account, no recovery — by design.
3. **Anonymous two-way channel.** Handlers reply sealed to the claim key;
   reporters prove control by signing with it. Nobody learns who they are.
4. **Go public with no infra.** Expose the intake page through a one-click
   tunnel from a laptop, or run behind a reverse proxy / Tor — your choice.

## The two surfaces

kilio has exactly two audiences, and they never share a session:

| Surface | Who | What it can do | Sees plaintext? |
|---|---|---|---|
| **Reporter surface** | anonymous, no login | submit a claim, generate/re-enter a receipt passphrase, poll a claim, read/send messages on *their own* claim | Yes, but only their own claim — the client seals/opens locally |
| **Handler surface** | an authenticated handler (triager/investigator/admin) for one or more branches | open the branches they're granted, read/decrypt claims, reply | Yes, but only for branches they're granted — server never widens this |

Both surfaces are served by the same `kilio-server` binary and the same `web/`
PWA bundle in `standalone`/`os` mode; the handler surface is also available
natively in `apps/desktop` (Tauri) for a single-officer deployment that never
needs a server at all.

## Multi-branch, single trust boundary

One deployment can serve many branches — offices, regions, or a "global"
catch-all. Each branch has its own keypair; a claim seals to the branch the
reporter chose, so a handler for one branch can never open another's — denied
reads return `404`, never a forbidden that leaks existence. kilio is
**instance-per-organisation, never multi-tenant**: multi-branch is the
supported axis of scale inside one org, and multi-org is just more instances,
optionally linked over an opt-in relay.

## Related documents

- [Getting started](#getting-started) — build, test, and the intended
  operator flow.
- [Self-hosting](#self-hosting) — run your own instance, go public.
- [Security model](#security) — the privacy spine, threat model, crypto.
- [Configuration](#configuration) — deploy modes and the pluggable seams.
