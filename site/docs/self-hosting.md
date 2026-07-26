# Self-hosting

kilio is one Rust binary, one SQLite file, and — when you want to be reachable
from outside your own network — a one-click tunnel. There is no SaaS
backend, no cloud account, and no dependency on any hosted service. What you
run is what you own, and what you run cannot read the claims it stores.

## Three ways to run it

| Mode | Who runs it | Reachability | Identity |
|---|---|---|---|
| `desktop` | one officer, on a laptop | subprocess tunnel | local owner |
| `standalone` | an org, on a box/VPS | tunnel or reverse proxy | local admin(s), Argon2id password → session |
| `os` | behind a Vulos OS gateway | gateway | gateway-brokered, server-verified session |

`os` mode **refuses to boot without a configured auth verifier** — a
fail-closed boot gate. It never silently collapses every handler down to one
identity. Self-hosters who aren't running behind a Vulos OS gateway will use
`desktop` or `standalone`.

## Going public with no fixed infrastructure

A built-in tunnel (`cloudflared` / `ngrok` / `frp`) turns a laptop or a small
VPS into a public intake page with no fixed infrastructure, DNS, or hosting
bill. From the handler UI (owner-gated), start the tunnel and kilio surfaces
the assigned public URL. Whichever provider runs, it proxies to exactly **one**
configured loopback address, re-checked before every connection — the
inbound request's Host/URL never chooses the target.

Prefer your own reverse proxy or a Tor hidden service instead — kilio doesn't
assume anything about your network, and running behind Tor is explicitly
supported for reporters who need that property (it is just not the default).

## What a self-hoster owns

- **The branch keypair.** `kilio init` generates it locally; the private half
  never leaves your machine. It's the only thing that can ever open a claim.
- **The SQLite database.** Every claim, message, and attachment lives in one
  file, and every row in it is ciphertext plus content-free routing metadata.
- **The decision to go public.** Nothing exposes your instance until you
  start a tunnel or put it behind your own proxy.

## Multi-branch, one instance

One deployment can serve many offices, regions, or a "global" catch-all
branch — this is the supported axis of scale inside one organisation. Running
a second organisation means running a second instance, never sharing a
database with the first. See [Configuration](#configuration) for the
`KotvaDelivery` seam that optionally forwards sealed claims between separate
instances without either side sharing a server.

## Status honesty

Only the sealed-crypto core (`kilio-seal`) and the sealed store
(`kilio-core`) are landed and tested today. `kilio-server`, `kilio-cli`, the
web PWA, and the desktop app — the pieces that actually make `kilio serve`
and `kilio tunnel start` runnable — are in progress. See
[Getting started](#getting-started) for what you can build and test right
now, and the repo's `ROADMAP.md` for build order.

## Related documents

- [Getting started](#getting-started) — build/test today vs. the intended
  operator flow.
- [Security model](#security) — the privacy spine and threat model.
- [Configuration](#configuration) — deploy modes and the Delivery/Reachability
  seams in depth.
