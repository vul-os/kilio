//! Stable, hex-rendered 16-byte identifiers for branches and claims.
//!
//! Ids are derived from public keys (never random), so the same key always
//! yields the same id and an id is a compact, non-secret handle the server can
//! route on without learning anything a public key wouldn't already reveal.

use core::fmt;
use serde::{Deserialize, Serialize};

use crate::SealError;

/// Domain-separated BLAKE3 tag for branch ids.
const BRANCH_ID_CTX: &str = "kilio/branch-id/v1";
/// Domain-separated BLAKE3 tag for claim ids.
const CLAIM_ID_CTX: &str = "kilio/claim-id/v1";

macro_rules! id_type {
    ($name:ident, $ctx:expr, $doc:literal) => {
        #[doc = $doc]
        #[derive(Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize)]
        pub struct $name(pub [u8; 16]);

        impl $name {
            /// Derive the id from a public key (or any stable bytes).
            pub fn derive(public_bytes: &[u8]) -> Self {
                let mut h = blake3::Hasher::new_derive_key($ctx);
                h.update(public_bytes);
                let digest = h.finalize();
                let mut out = [0u8; 16];
                out.copy_from_slice(&digest.as_bytes()[..16]);
                Self(out)
            }

            /// Lowercase hex rendering (32 chars).
            pub fn to_hex(&self) -> String {
                hex::encode(self.0)
            }

            /// Parse from a 32-char hex string.
            pub fn from_hex(s: &str) -> Result<Self, SealError> {
                let v = hex::decode(s.trim()).map_err(|_| SealError::BadId)?;
                let arr: [u8; 16] = v.as_slice().try_into().map_err(|_| SealError::BadId)?;
                Ok(Self(arr))
            }

            /// Raw bytes.
            pub fn as_bytes(&self) -> &[u8; 16] {
                &self.0
            }
        }

        impl fmt::Display for $name {
            fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
                f.write_str(&self.to_hex())
            }
        }

        impl fmt::Debug for $name {
            fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
                write!(f, "{}({})", stringify!($name), self.to_hex())
            }
        }
    };
}

id_type!(BranchId, BRANCH_ID_CTX, "Identifier for a branch (a keyed destination for sealed claims).");
id_type!(ClaimId, CLAIM_ID_CTX, "Public handle for a claim; derived from the claim's signing key.");

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn hex_roundtrip() {
        let id = BranchId::derive(b"some-public-key");
        let s = id.to_hex();
        assert_eq!(s.len(), 32);
        assert_eq!(BranchId::from_hex(&s).unwrap(), id);
    }

    #[test]
    fn derivation_is_stable_and_separated() {
        assert_eq!(BranchId::derive(b"k"), BranchId::derive(b"k"));
        // Different domains must not collide for identical input.
        assert_ne!(BranchId::derive(b"k").0, ClaimId::derive(b"k").0);
    }

    #[test]
    fn rejects_bad_hex() {
        assert!(BranchId::from_hex("nothex").is_err());
        assert!(BranchId::from_hex("aa").is_err());
    }
}
