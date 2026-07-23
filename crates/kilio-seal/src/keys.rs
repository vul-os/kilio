//! Branch keys — the destination a reporter seals a claim to.
//!
//! A branch owns two keypairs:
//!  * an **HPKE (X25519) recipient keypair** claims are sealed *to*, and
//!  * an **Ed25519 signing keypair** the branch uses to authenticate its own
//!    published metadata (and, later, handler-side attestations).
//!
//! The secret half ([`BranchKeys`]) lives only with the handler / Tauri app or
//! is emitted once by `kilio init`. The public half ([`BranchPublic`]) is what
//! a reporter's browser needs, and reveals nothing sensitive.

use ed25519_dalek::{Signature, Signer, SigningKey, Verifier, VerifyingKey};
use hpke::{Deserializable, Kem as KemTrait, Serializable};
use rand::rngs::OsRng;
use serde::{Deserialize, Serialize};
use zeroize::Zeroizing;

use crate::ids::BranchId;
use crate::{Kem, SealError};

/// Secret branch material. Treat every field as sensitive.
pub struct BranchKeys {
    pub branch_id: BranchId,
    kem_secret: Zeroizing<Vec<u8>>,
    kem_public: Vec<u8>,
    sign_seed: Zeroizing<[u8; 32]>,
    sign_public: [u8; 32],
}

impl BranchKeys {
    /// Generate a fresh branch keypair from the OS CSPRNG.
    pub fn generate() -> Self {
        let (kem_sk, kem_pk) = Kem::gen_keypair(&mut OsRng);
        let signing = SigningKey::generate(&mut OsRng);
        Self::assemble(
            kem_sk.to_bytes().to_vec(),
            kem_pk.to_bytes().to_vec(),
            signing.to_bytes(),
            signing.verifying_key().to_bytes(),
        )
    }

    fn assemble(
        kem_secret: Vec<u8>,
        kem_public: Vec<u8>,
        sign_seed: [u8; 32],
        sign_public: [u8; 32],
    ) -> Self {
        let branch_id = BranchId::derive(&id_material(&kem_public, &sign_public));
        Self {
            branch_id,
            kem_secret: Zeroizing::new(kem_secret),
            kem_public,
            sign_seed: Zeroizing::new(sign_seed),
            sign_public,
        }
    }

    /// The public half a reporter needs to seal a claim to this branch.
    pub fn public(&self) -> BranchPublic {
        BranchPublic {
            branch_id: self.branch_id,
            kem_public: self.kem_public.clone(),
            sign_public: self.sign_public,
        }
    }

    /// Reconstruct the HPKE recipient private key for `open`.
    pub(crate) fn kem_private(&self) -> Result<<Kem as KemTrait>::PrivateKey, SealError> {
        <Kem as KemTrait>::PrivateKey::from_bytes(&self.kem_secret).map_err(|_| SealError::BadKey)
    }

    /// Sign branch-published bytes with the Ed25519 key.
    pub fn sign(&self, msg: &[u8]) -> [u8; 64] {
        let sk = SigningKey::from_bytes(&self.sign_seed);
        sk.sign(msg).to_bytes()
    }

    /// Serialize the secret keys to bytes for at-rest storage (caller encrypts).
    pub fn to_secret_bytes(&self) -> SecretKeyBytes {
        SecretKeyBytes {
            kem_secret: self.kem_secret.to_vec(),
            kem_public: self.kem_public.clone(),
            sign_seed: *self.sign_seed,
            sign_public: self.sign_public,
        }
    }

    /// Restore from previously serialized secret bytes.
    pub fn from_secret_bytes(b: SecretKeyBytes) -> Self {
        Self::assemble(b.kem_secret, b.kem_public, b.sign_seed, b.sign_public)
    }
}

/// On-disk / at-rest form of [`BranchKeys`]. The caller is responsible for
/// encrypting this before it touches durable storage.
#[derive(Serialize, Deserialize)]
pub struct SecretKeyBytes {
    pub kem_secret: Vec<u8>,
    pub kem_public: Vec<u8>,
    pub sign_seed: [u8; 32],
    pub sign_public: [u8; 32],
}

/// Published branch identity. Safe to hand to any reporter.
#[derive(Clone, Serialize, Deserialize)]
pub struct BranchPublic {
    pub branch_id: BranchId,
    pub kem_public: Vec<u8>,
    pub sign_public: [u8; 32],
}

impl BranchPublic {
    /// The HPKE recipient public key to seal to.
    pub(crate) fn kem_pk(&self) -> Result<<Kem as KemTrait>::PublicKey, SealError> {
        <Kem as KemTrait>::PublicKey::from_bytes(&self.kem_public).map_err(|_| SealError::BadKey)
    }

    /// Verify a signature made by [`BranchKeys::sign`].
    pub fn verify(&self, msg: &[u8], sig: &[u8; 64]) -> Result<(), SealError> {
        let vk = VerifyingKey::from_bytes(&self.sign_public).map_err(|_| SealError::BadKey)?;
        vk.verify(msg, &Signature::from_bytes(sig))
            .map_err(|_| SealError::BadSignature)
    }

    /// Recompute the branch id from the public keys; guards against a server
    /// serving a swapped key under a known id.
    pub fn expected_id(&self) -> BranchId {
        BranchId::derive(&id_material(&self.kem_public, &self.sign_public))
    }
}

fn id_material(kem_public: &[u8], sign_public: &[u8; 32]) -> Vec<u8> {
    let mut m = Vec::with_capacity(kem_public.len() + 32);
    m.extend_from_slice(kem_public);
    m.extend_from_slice(sign_public);
    m
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn branch_id_binds_public_keys() {
        let b = BranchKeys::generate();
        assert_eq!(b.branch_id, b.public().expected_id());
    }

    #[test]
    fn sign_verify_roundtrip() {
        let b = BranchKeys::generate();
        let sig = b.sign(b"published-metadata");
        assert!(b.public().verify(b"published-metadata", &sig).is_ok());
        assert!(b.public().verify(b"tampered", &sig).is_err());
    }

    #[test]
    fn secret_bytes_roundtrip() {
        let b = BranchKeys::generate();
        let restored = BranchKeys::from_secret_bytes(b.to_secret_bytes());
        assert_eq!(b.branch_id, restored.branch_id);
        // both must still open the same sealed data — checked in envelope tests
    }
}
