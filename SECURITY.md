# Security policy

kilio exists to protect people who report sensitive things. We take security
reports seriously and will work with you in good faith.

## Supported versions

kilio is pre-1.0 (`0.1.x`). Only the latest `main` is supported. Until a 1.0
release, treat every interface as subject to change.

## Reporting a vulnerability

**Please report privately — do not open a public issue for a security bug.**

- Use GitHub's **[private vulnerability reporting](https://github.com/vul-os/kilio/security/advisories/new)**
  ("Report a vulnerability" under the Security tab), or
- open a minimal public issue that says only "security report — please open a
  private channel" with no details, and we will follow up.

Please include: affected component (crate/surface), version or commit, a
description, and a proof-of-concept if you have one. We aim to acknowledge
within a few days.

### In scope
The sealing/crypto (`kilio-seal`), the intake and handler APIs, branch
key handling and at-rest storage, the tunnel/reachability path (SSRF, exposure),
and anything that could **deanonymize a reporter** or **let the host read a
sealed claim**. Deanonymization and sealed-content-disclosure are our highest
severity classes.

### Out of scope
Findings that require a compromised handler device or the reporter's receipt
passphrase (both are trusted by design — see the threat model), and traffic
*timing* analysis without Tor (documented as a residual risk).

## The security model

The privacy spine, threat model, cryptographic primitives, and explicit residual
risks are documented in **[docs/SECURITY.md](docs/SECURITY.md)** and the design
rationale in **[decisions.md](decisions.md)**. Read those before filing a report
— they state what kilio does and does not defend against, so we can focus on
real gaps.

## Cryptography

kilio uses RFC 9180 HPKE (DHKEM-X25519 / HKDF-SHA256 / ChaCha20Poly1305) via the
audited `hpke` crate, Ed25519 (`ed25519-dalek`), Argon2id, and BLAKE3. No
primitive is hand-rolled; only their composition is ours. Review of that
composition is especially welcome.
