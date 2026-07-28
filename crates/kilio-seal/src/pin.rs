//! Branch key pinning — trust-on-first-use for the key a reporter seals to.
//!
//! ## Why this exists
//!
//! A reporter's browser fetches the branch public key from the very host the
//! privacy model says it must not have to trust. [`BranchPublic::expected_id`]
//! only proves a served key is *self-consistent* with its own id: an attacker
//! who substitutes a key it controls also substitutes the matching id, and the
//! check passes. Self-consistency is not authenticity.
//!
//! Pinning is the missing half, and it follows aql's pairing discipline
//! (`aql/controller/internal/pairing`) rather than a looser first-key-wins:
//!
//! * [`BranchPin::pin`] is the **only** moment a branch key is accepted.
//! * Every later fetch goes through [`BranchPin::check`], which refuses any key
//!   change ([`SealError::KeyChangeRefused`]) and any epoch rollback
//!   ([`SealError::Rollback`]).
//! * The single way to rotate a pinned key is a [`RekeyCert`] signed by the
//!   **currently pinned** signing key — aql's `repair` command, in kilio's
//!   shape. A host that loses or replaces a branch key cannot rotate the pin.
//! * The only other escape is [`BranchPin::unpin`], which consumes the pin: the
//!   analogue of a physical factory reset, deliberately explicit and never a
//!   side effect of a fetch.
//!
//! ## What a signature here does and does not prove
//!
//! [`SignedDescriptor::verify_self`] proves the publisher held the branch
//! signing key and that `branch_id` binds both published keys. It establishes
//! **no trust at all** on its own — an attacker signs its own descriptor just
//! as easily. Trust comes only from the pin, or from a `branch_id` the reporter
//! obtained out of band (a poster, a card, a printed URL).

use ed25519_dalek::{Signature, Verifier, VerifyingKey};
use serde::{Deserialize, Serialize};

use crate::ids::BranchId;
use crate::keys::{BranchKeys, BranchPublic};
use crate::SealError;

/// Descriptor format version. Anything else fails closed.
pub const DESCRIPTOR_VERSION: u8 = 1;

const DESC_DOMAIN: &[u8] = b"kilio/branch-descriptor/v1\0";
const REKEY_DOMAIN: &[u8] = b"kilio/branch-rekey/v1\0";

/// A branch's published identity: the keys plus the cleartext settings a
/// reporter's client needs before it can seal anything.
///
/// `epoch` is a monotonic counter the branch owner increments on every
/// republish. It is the rollback guard: a pinned client never accepts a
/// descriptor older than the newest one it has already seen, so a host cannot
/// replay an old descriptor to force a weaker `pow_bits` or a retired key.
#[derive(Clone, PartialEq, Eq, Debug, Serialize, Deserialize)]
pub struct BranchDescriptor {
    pub v: u8,
    pub branch_id: BranchId,
    pub kem_public: Vec<u8>,
    pub sign_public: [u8; 32],
    pub name: String,
    pub pow_bits: u8,
    pub epoch: u64,
    /// Millis since epoch, from the publisher's clock. Advisory only — never a
    /// security input; `epoch` is what orders descriptors.
    pub issued_at: u64,
}

impl BranchDescriptor {
    /// The public half a reporter seals to.
    pub fn branch_public(&self) -> BranchPublic {
        BranchPublic {
            branch_id: self.branch_id,
            kem_public: self.kem_public.clone(),
            sign_public: self.sign_public,
        }
    }

    /// Deterministic signing input. Verification re-encodes the *decoded*
    /// struct, so a re-encoding of the same descriptor with different CBOR
    /// framing does not verify — there is one signable form, not many.
    fn signing_bytes(&self, domain: &[u8]) -> Result<Vec<u8>, SealError> {
        let mut buf = domain.to_vec();
        ciborium::into_writer(self, &mut buf).map_err(|_| SealError::Encode)?;
        Ok(buf)
    }
}

/// A [`BranchDescriptor`] signed by its own branch signing key.
#[derive(Clone, Serialize, Deserialize)]
pub struct SignedDescriptor {
    pub desc: BranchDescriptor,
    /// Raw 64-byte Ed25519 signature over `signing_bytes`.
    pub sig: Vec<u8>,
}

impl SignedDescriptor {
    /// Structural validation: correct version, `branch_id` binds both published
    /// keys, and the self-signature verifies.
    ///
    /// This is a prerequisite for trust, never a source of it — see the module
    /// docs. Use [`BranchPin`] to decide whether the key may be *believed*.
    pub fn verify_self(&self) -> Result<&BranchDescriptor, SealError> {
        if self.desc.v != DESCRIPTOR_VERSION {
            return Err(SealError::Version);
        }
        if self.desc.branch_id != self.desc.branch_public().expected_id() {
            return Err(SealError::IdMismatch);
        }
        let sig: [u8; 64] = self
            .sig
            .as_slice()
            .try_into()
            .map_err(|_| SealError::BadSignature)?;
        let vk = VerifyingKey::from_bytes(&self.desc.sign_public).map_err(|_| SealError::BadKey)?;
        vk.verify(
            &self.desc.signing_bytes(DESC_DOMAIN)?,
            &Signature::from_bytes(&sig),
        )
        .map_err(|_| SealError::BadSignature)?;
        Ok(&self.desc)
    }
}

/// Authorization to move a pin from one branch key to another: the next
/// descriptor, countersigned by the key currently pinned.
///
/// This is aql's `repair` command. The countersignature is what a hostile host
/// cannot forge, so re-publishing under a fresh key can never rotate a pin.
#[derive(Clone, Serialize, Deserialize)]
pub struct RekeyCert {
    /// The branch id the countersigning (old) key belongs to.
    pub prev_branch_id: BranchId,
    pub next: SignedDescriptor,
    /// Raw 64-byte Ed25519 signature by the OLD branch signing key.
    pub sig: Vec<u8>,
}

/// The exact bytes a rekey countersignature covers.
#[derive(Serialize)]
struct RekeyBody<'a> {
    prev_branch_id: &'a BranchId,
    next: &'a BranchDescriptor,
}

fn rekey_signing_bytes(
    prev_branch_id: &BranchId,
    next: &BranchDescriptor,
) -> Result<Vec<u8>, SealError> {
    let mut buf = REKEY_DOMAIN.to_vec();
    ciborium::into_writer(
        &RekeyBody {
            prev_branch_id,
            next,
        },
        &mut buf,
    )
    .map_err(|_| SealError::Encode)?;
    Ok(buf)
}

impl BranchKeys {
    /// Publish this branch's signed descriptor at `epoch`.
    pub fn publish(
        &self,
        name: &str,
        pow_bits: u8,
        epoch: u64,
        issued_at: u64,
    ) -> Result<SignedDescriptor, SealError> {
        let p = self.public();
        let desc = BranchDescriptor {
            v: DESCRIPTOR_VERSION,
            branch_id: p.branch_id,
            kem_public: p.kem_public,
            sign_public: p.sign_public,
            name: name.to_string(),
            pow_bits,
            epoch,
            issued_at,
        };
        let sig = self.sign(&desc.signing_bytes(DESC_DOMAIN)?).to_vec();
        Ok(SignedDescriptor { desc, sig })
    }

    /// Countersign `next` with *these* (the outgoing, currently pinned) keys so
    /// pinned clients will follow the rotation.
    pub fn sign_rekey(&self, next: &SignedDescriptor) -> Result<RekeyCert, SealError> {
        next.verify_self()?;
        if next.desc.branch_id == self.branch_id {
            // Rotating to the same keys is not a rotation; refuse rather than
            // mint a certificate that authorizes nothing.
            return Err(SealError::KeyChangeRefused);
        }
        let sig = self
            .sign(&rekey_signing_bytes(&self.branch_id, &next.desc)?)
            .to_vec();
        Ok(RekeyCert {
            prev_branch_id: self.branch_id,
            next: next.clone(),
            sig,
        })
    }
}

/// What a reporter's device remembers about a branch after first contact.
///
/// Persist this next to the receipt phrase. It is not secret — it is a
/// commitment.
#[derive(Clone, PartialEq, Eq, Debug, Serialize, Deserialize)]
pub struct BranchPin {
    pub branch_id: BranchId,
    pub kem_public: Vec<u8>,
    pub sign_public: [u8; 32],
    /// Highest descriptor epoch accepted so far.
    pub epoch: u64,
}

impl BranchPin {
    /// Trust on first use. **The only moment a branch key is accepted.**
    ///
    /// Callers that know the `branch_id` out of band (a poster, a printed card)
    /// should compare it *before* calling this; then the pin is not TOFU at all
    /// but an out-of-band verification.
    pub fn pin(sd: &SignedDescriptor) -> Result<Self, SealError> {
        let d = sd.verify_self()?;
        Ok(Self {
            branch_id: d.branch_id,
            kem_public: d.kem_public.clone(),
            sign_public: d.sign_public,
            epoch: d.epoch,
        })
    }

    /// Validate a freshly fetched descriptor against the pin.
    ///
    /// Fails closed on: a bad self-signature, an id that does not bind its keys,
    /// *any* key change ([`SealError::KeyChangeRefused`]), and any epoch older
    /// than the pinned one ([`SealError::Rollback`]).
    pub fn check<'a>(&self, sd: &'a SignedDescriptor) -> Result<&'a BranchDescriptor, SealError> {
        let d = sd.verify_self()?;
        if d.branch_id != self.branch_id
            || d.sign_public != self.sign_public
            || d.kem_public != self.kem_public
        {
            return Err(SealError::KeyChangeRefused);
        }
        if d.epoch < self.epoch {
            return Err(SealError::Rollback);
        }
        Ok(d)
    }

    /// [`BranchPin::check`], then advance the pinned epoch. Use this on every
    /// fetch so a descriptor can never be rolled back after being seen.
    pub fn observe<'a>(
        &mut self,
        sd: &'a SignedDescriptor,
    ) -> Result<&'a BranchDescriptor, SealError> {
        let d = self.check(sd)?;
        self.epoch = d.epoch;
        Ok(d)
    }

    /// Rotate the pin — the *only* way a pinned key changes without a reset.
    ///
    /// Requires a certificate countersigned by the currently pinned signing key
    /// and a strictly greater epoch.
    pub fn accept_rekey(&mut self, cert: &RekeyCert) -> Result<(), SealError> {
        if cert.prev_branch_id != self.branch_id {
            return Err(SealError::KeyChangeRefused);
        }
        let next = cert.next.verify_self()?;
        if next.epoch <= self.epoch {
            return Err(SealError::Rollback);
        }
        let sig: [u8; 64] = cert
            .sig
            .as_slice()
            .try_into()
            .map_err(|_| SealError::BadSignature)?;
        let vk = VerifyingKey::from_bytes(&self.sign_public).map_err(|_| SealError::BadKey)?;
        vk.verify(
            &rekey_signing_bytes(&self.branch_id, next)?,
            &Signature::from_bytes(&sig),
        )
        .map_err(|_| SealError::KeyChangeRefused)?;

        self.branch_id = next.branch_id;
        self.kem_public = next.kem_public.clone();
        self.sign_public = next.sign_public;
        self.epoch = next.epoch;
        Ok(())
    }

    /// Forget this pin. The factory-reset analogue: consuming, explicit, and
    /// never reachable from a fetch path.
    pub fn unpin(self) {}
}

/// Signing-key-only handle used by tests to forge a foreign countersignature.
#[cfg(test)]
fn foreign_sign(seed: [u8; 32], msg: &[u8]) -> Vec<u8> {
    use ed25519_dalek::{Signer, SigningKey};
    SigningKey::from_bytes(&seed).sign(msg).to_vec()
}

#[cfg(test)]
mod tests {
    use super::*;

    fn published(b: &BranchKeys, epoch: u64) -> SignedDescriptor {
        b.publish("corporate", 20, epoch, 1_600_000_000_000)
            .unwrap()
    }

    #[test]
    fn self_signed_descriptor_verifies() {
        let b = BranchKeys::generate();
        let sd = published(&b, 1);
        let d = sd.verify_self().unwrap();
        assert_eq!(d.branch_id, b.branch_id);
        assert_eq!(d.pow_bits, 20);
    }

    #[test]
    fn tampered_field_breaks_self_signature() {
        let b = BranchKeys::generate();
        let mut sd = published(&b, 1);
        sd.desc.pow_bits = 1; // downgrade the cold-contact gate
        assert!(matches!(sd.verify_self(), Err(SealError::BadSignature)));
    }

    #[test]
    fn id_must_bind_the_published_keys() {
        let b = BranchKeys::generate();
        let other = BranchKeys::generate();
        let mut sd = published(&b, 1);
        sd.desc.branch_id = other.branch_id;
        assert!(matches!(sd.verify_self(), Err(SealError::IdMismatch)));
    }

    #[test]
    fn wrong_version_fails_closed() {
        let b = BranchKeys::generate();
        let mut sd = published(&b, 1);
        sd.desc.v = 2;
        assert!(matches!(sd.verify_self(), Err(SealError::Version)));
    }

    #[test]
    fn short_signature_rejected() {
        let b = BranchKeys::generate();
        let mut sd = published(&b, 1);
        sd.sig.truncate(63);
        assert!(matches!(sd.verify_self(), Err(SealError::BadSignature)));
    }

    #[test]
    fn pin_accepts_the_same_branch_again() {
        let b = BranchKeys::generate();
        let pin = BranchPin::pin(&published(&b, 1)).unwrap();
        assert!(pin.check(&published(&b, 2)).is_ok());
    }

    /// The attack `expected_id()` alone cannot stop: a hostile host serves a
    /// *fully self-consistent* descriptor for a key it controls.
    #[test]
    fn pin_refuses_a_substituted_key_that_is_self_consistent() {
        let real = BranchKeys::generate();
        let attacker = BranchKeys::generate();
        let evil = attacker.publish("corporate", 20, 99, 1).unwrap();

        // The substitute passes every self-consistency check...
        assert!(evil.verify_self().is_ok());
        assert_eq!(evil.desc.branch_id, evil.desc.branch_public().expected_id());

        // ...and is still refused by a pinned client.
        let pin = BranchPin::pin(&published(&real, 1)).unwrap();
        assert!(matches!(pin.check(&evil), Err(SealError::KeyChangeRefused)));
    }

    #[test]
    fn pin_refuses_epoch_rollback() {
        let b = BranchKeys::generate();
        let mut pin = BranchPin::pin(&published(&b, 5)).unwrap();
        assert!(matches!(
            pin.check(&published(&b, 4)),
            Err(SealError::Rollback)
        ));
        // observe() raises the floor, so a previously valid epoch is now stale.
        pin.observe(&published(&b, 7)).unwrap();
        assert_eq!(pin.epoch, 7);
        assert!(matches!(
            pin.check(&published(&b, 6)),
            Err(SealError::Rollback)
        ));
    }

    #[test]
    fn rekey_signed_by_the_pinned_key_rotates_the_pin() {
        let old = BranchKeys::generate();
        let new = BranchKeys::generate();
        let mut pin = BranchPin::pin(&published(&old, 1)).unwrap();

        let next = published(&new, 2);
        let cert = old.sign_rekey(&next).unwrap();
        pin.accept_rekey(&cert).unwrap();

        assert_eq!(pin.branch_id, new.branch_id);
        assert!(pin.check(&published(&new, 3)).is_ok());
        // The retired key is no longer believed.
        assert!(matches!(
            pin.check(&published(&old, 4)),
            Err(SealError::KeyChangeRefused)
        ));
    }

    #[test]
    fn rekey_not_signed_by_the_pinned_key_is_refused() {
        let old = BranchKeys::generate();
        let attacker = BranchKeys::generate();
        let new = BranchKeys::generate();
        let mut pin = BranchPin::pin(&published(&old, 1)).unwrap();

        // A certificate minted by a key that is not the pinned one.
        let mut cert = attacker.sign_rekey(&published(&new, 2)).unwrap();
        cert.prev_branch_id = old.branch_id; // claim to be the outgoing branch
        assert!(matches!(
            pin.accept_rekey(&cert),
            Err(SealError::KeyChangeRefused)
        ));
        assert_eq!(pin.branch_id, old.branch_id);
    }

    #[test]
    fn rekey_countersignature_is_bound_to_the_outgoing_branch() {
        let old = BranchKeys::generate();
        let new = BranchKeys::generate();
        let mut pin = BranchPin::pin(&published(&old, 1)).unwrap();
        let cert = old.sign_rekey(&published(&new, 2)).unwrap();

        // Same signature, replayed against a different outgoing branch id.
        let mut replayed = cert.clone();
        replayed.prev_branch_id = BranchId::derive(b"someone else");
        assert!(matches!(
            pin.accept_rekey(&replayed),
            Err(SealError::KeyChangeRefused)
        ));
    }

    #[test]
    fn rekey_must_advance_the_epoch() {
        let old = BranchKeys::generate();
        let new = BranchKeys::generate();
        let mut pin = BranchPin::pin(&published(&old, 5)).unwrap();
        let cert = old.sign_rekey(&published(&new, 5)).unwrap();
        assert!(matches!(pin.accept_rekey(&cert), Err(SealError::Rollback)));
    }

    #[test]
    fn rekey_body_is_domain_separated_from_a_descriptor() {
        // A descriptor self-signature must never be reusable as a rekey
        // countersignature (or the reverse).
        let old = BranchKeys::generate();
        let new = BranchKeys::generate();
        let mut pin = BranchPin::pin(&published(&old, 1)).unwrap();
        let next = published(&new, 2);

        let mut forged = old.sign_rekey(&next).unwrap();
        // Swap in the *descriptor* signature made by the same old key.
        forged.sig = published(&old, 1).sig;
        assert!(matches!(
            pin.accept_rekey(&forged),
            Err(SealError::KeyChangeRefused)
        ));
    }

    #[test]
    fn rekey_to_the_same_keys_is_refused_at_issue_time() {
        let b = BranchKeys::generate();
        assert!(matches!(
            b.sign_rekey(&published(&b, 2)),
            Err(SealError::KeyChangeRefused)
        ));
    }

    #[test]
    fn forged_countersignature_by_a_random_seed_is_refused() {
        let old = BranchKeys::generate();
        let new = BranchKeys::generate();
        let mut pin = BranchPin::pin(&published(&old, 1)).unwrap();
        let next = published(&new, 2);
        let body = rekey_signing_bytes(&old.branch_id, &next.desc).unwrap();
        let cert = RekeyCert {
            prev_branch_id: old.branch_id,
            next,
            sig: foreign_sign([9u8; 32], &body),
        };
        assert!(matches!(
            pin.accept_rekey(&cert),
            Err(SealError::KeyChangeRefused)
        ));
    }

    #[test]
    fn pin_roundtrips_through_cbor() {
        let b = BranchKeys::generate();
        let pin = BranchPin::pin(&published(&b, 3)).unwrap();
        let mut buf = Vec::new();
        ciborium::into_writer(&pin, &mut buf).unwrap();
        let back: BranchPin = ciborium::from_reader(&buf[..]).unwrap();
        assert_eq!(pin, back);
    }
}
