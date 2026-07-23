//! The receipt passphrase — kilio's only reporter identity.
//!
//! At submission the client mints 128 bits of entropy as a BIP-39 12-word
//! phrase. That phrase is the *only* thing the reporter must keep. From it we
//! deterministically derive, via a memory-hard KDF, a per-claim signing key
//! (to prove control when returning) and a per-claim HPKE key (to receive
//! sealed handler replies). Nothing derived from the phrase ever reaches the
//! server except the public halves and the claim id.
//!
//! Losing the phrase means losing the claim — by design. Recovery would mean
//! someone else could impersonate the reporter.

use argon2::{Algorithm, Argon2, Params, Version};
use bip39::{Language, Mnemonic};
use ed25519_dalek::{Signature, Signer, SigningKey};
use hpke::{Deserializable, Kem as KemTrait, Serializable};
use rand::{rngs::OsRng, RngCore};
use serde::{Deserialize, Serialize};
use zeroize::Zeroizing;

use crate::ids::{BranchId, ClaimId};
use crate::{Kem, SealError};

// Argon2id cost. 64 MiB, 3 passes, 1 lane — deliberately expensive so a weak
// (but real 128-bit) phrase resists offline guessing and a stolen DB cannot be
// correlated cheaply. Do not lower without updating decisions.md §4.
const ARGON_MEM_KIB: u32 = 64 * 1024;
const ARGON_TIME: u32 = 3;
const ARGON_LANES: u32 = 1;

/// A 12-word BIP-39 receipt passphrase.
pub struct Receipt {
    phrase: Zeroizing<String>,
}

impl Receipt {
    /// Mint a fresh receipt from the OS CSPRNG (128-bit entropy → 12 words).
    pub fn generate() -> Self {
        let mut entropy = [0u8; 16];
        OsRng.fill_bytes(&mut entropy);
        let m = Mnemonic::from_entropy(&entropy).expect("16 bytes is a valid BIP-39 length");
        entropy.iter_mut().for_each(|b| *b = 0);
        Self {
            phrase: Zeroizing::new(m.to_string()),
        }
    }

    /// Parse a phrase a returning reporter typed. Validates the BIP-39 checksum,
    /// so most typos are caught before any key derivation.
    pub fn from_phrase(phrase: &str) -> Result<Self, SealError> {
        let m = Mnemonic::parse_in_normalized(Language::English, phrase)
            .map_err(|_| SealError::BadReceipt)?;
        Ok(Self {
            phrase: Zeroizing::new(m.to_string()),
        })
    }

    /// The words, to show once at submission. Never transmit these.
    pub fn phrase(&self) -> &str {
        &self.phrase
    }

    /// Derive the per-claim keys for a specific branch. The branch id is folded
    /// into the KDF salt so the same phrase yields different keys per branch —
    /// a phrase leaked for one branch reveals nothing about another.
    pub fn derive(&self, branch_id: &BranchId) -> Result<ClaimKeys, SealError> {
        let m = Mnemonic::parse_in_normalized(Language::English, &self.phrase)
            .map_err(|_| SealError::BadReceipt)?;
        let entropy = Zeroizing::new(m.to_entropy());

        let salt = blake3::derive_key("kilio/receipt-salt/v1", branch_id.as_bytes());
        let argon = Argon2::new(
            Algorithm::Argon2id,
            Version::V0x13,
            Params::new(ARGON_MEM_KIB, ARGON_TIME, ARGON_LANES, Some(32))
                .map_err(|_| SealError::Kdf)?,
        );
        let mut root = Zeroizing::new([0u8; 32]);
        argon
            .hash_password_into(&entropy, &salt[..16], root.as_mut())
            .map_err(|_| SealError::Kdf)?;

        // Two independent sub-keys from the root.
        let sign_seed = blake3::keyed_hash(&root, b"kilio/claim/sign/v1");
        let recip_ikm = blake3::keyed_hash(&root, b"kilio/claim/recip/v1");

        let signing = SigningKey::from_bytes(sign_seed.as_bytes());
        let sign_public = signing.verifying_key().to_bytes();

        let (kem_sk, kem_pk) = Kem::derive_keypair(recip_ikm.as_bytes());
        let kem_public = kem_pk.to_bytes().to_vec();
        let kem_secret = Zeroizing::new(kem_sk.to_bytes().to_vec());

        let claim_id = ClaimId::derive(&claim_id_material(&kem_public, &sign_public));

        Ok(ClaimKeys {
            claim_id,
            sign_seed: Zeroizing::new(*sign_seed.as_bytes()),
            kem_secret,
            kem_public,
            sign_public,
        })
    }
}

/// Per-claim keys derived from a receipt. Secret; lives only on the reporter's
/// device for the duration of a session.
pub struct ClaimKeys {
    pub claim_id: ClaimId,
    sign_seed: Zeroizing<[u8; 32]>,
    kem_secret: Zeroizing<Vec<u8>>,
    kem_public: Vec<u8>,
    sign_public: [u8; 32],
}

impl ClaimKeys {
    /// The public identity handed to the handler (inside the sealed payload) so
    /// they can seal replies back to this claim.
    pub fn public(&self) -> ClaimPublic {
        ClaimPublic {
            claim_id: self.claim_id,
            kem_public: self.kem_public.clone(),
            sign_public: self.sign_public,
        }
    }

    /// Sign a poll/message to prove control of this claim.
    pub fn sign(&self, msg: &[u8]) -> [u8; 64] {
        SigningKey::from_bytes(&self.sign_seed).sign(msg).to_bytes()
    }

    /// HPKE recipient private key for opening sealed handler replies.
    pub(crate) fn kem_private(&self) -> Result<<Kem as KemTrait>::PrivateKey, SealError> {
        <Kem as KemTrait>::PrivateKey::from_bytes(&self.kem_secret).map_err(|_| SealError::BadKey)
    }
}

/// Public per-claim identity. Travels sealed-sender inside a submission.
#[derive(Clone, Serialize, Deserialize, PartialEq, Eq, Debug)]
pub struct ClaimPublic {
    pub claim_id: ClaimId,
    pub kem_public: Vec<u8>,
    pub sign_public: [u8; 32],
}

impl ClaimPublic {
    pub(crate) fn kem_pk(&self) -> Result<<Kem as KemTrait>::PublicKey, SealError> {
        <Kem as KemTrait>::PublicKey::from_bytes(&self.kem_public).map_err(|_| SealError::BadKey)
    }

    /// Verify a signature made by [`ClaimKeys::sign`] (server-side poll auth).
    pub fn verify(&self, msg: &[u8], sig: &[u8; 64]) -> Result<(), SealError> {
        use ed25519_dalek::{Verifier, VerifyingKey};
        let vk = VerifyingKey::from_bytes(&self.sign_public).map_err(|_| SealError::BadKey)?;
        vk.verify(msg, &Signature::from_bytes(sig))
            .map_err(|_| SealError::BadSignature)
    }
}

fn claim_id_material(kem_public: &[u8], sign_public: &[u8; 32]) -> Vec<u8> {
    let mut m = Vec::with_capacity(kem_public.len() + 32);
    m.extend_from_slice(kem_public);
    m.extend_from_slice(sign_public);
    m
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn phrase_is_twelve_words() {
        let r = Receipt::generate();
        assert_eq!(r.phrase().split_whitespace().count(), 12);
    }

    #[test]
    fn same_phrase_same_keys() {
        let r = Receipt::generate();
        let branch = BranchId::derive(b"branch-key");
        let a = r.derive(&branch).unwrap();
        let r2 = Receipt::from_phrase(r.phrase()).unwrap();
        let b = r2.derive(&branch).unwrap();
        assert_eq!(a.claim_id, b.claim_id);
        assert_eq!(a.public(), b.public());
    }

    #[test]
    fn different_branch_different_keys() {
        let r = Receipt::generate();
        let a = r.derive(&BranchId::derive(b"branch-a")).unwrap();
        let b = r.derive(&BranchId::derive(b"branch-b")).unwrap();
        assert_ne!(a.claim_id, b.claim_id);
    }

    #[test]
    fn bad_phrase_rejected() {
        assert!(Receipt::from_phrase("not a real bip39 phrase at all nope").is_err());
    }

    #[test]
    fn claim_sign_verify() {
        let r = Receipt::generate();
        let ck = r.derive(&BranchId::derive(b"b")).unwrap();
        let sig = ck.sign(b"poll:claim");
        assert!(ck.public().verify(b"poll:claim", &sig).is_ok());
        assert!(ck.public().verify(b"other", &sig).is_err());
    }
}
