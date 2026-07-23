//! # kilio-seal
//!
//! The sealed-submission crypto spine for kilio. One implementation, compiled
//! natively (server, CLI, Tauri handler) and to `wasm32` (the reporter's
//! browser), so sealing exists exactly once and never diverges per surface.
//!
//! What it provides:
//! * [`BranchKeys`] / [`BranchPublic`] — the keyed destination a claim seals to.
//! * [`Receipt`] → [`ClaimKeys`] — the reporter's only identity, a 12-word
//!   passphrase deterministically derived (memory-hard) into per-claim keys.
//! * [`Envelope`] — a sealed-sender HPKE envelope (kotva-MOTE-shaped) with all
//!   cleartext routing fields bound into the AEAD associated data.
//! * [`pow`] — an anonymous, per-message proof-of-work cold-contact gate.
//!
//! Primitives are RFC 9180 HPKE (DHKEM-X25519 / HKDF-SHA256 / ChaCha20Poly1305)
//! via the audited `hpke` crate, Ed25519 via `ed25519-dalek`, Argon2id, and
//! BLAKE3. No primitive is hand-rolled; only their composition is ours.
//!
//! See `decisions.md` §3–4 for the privacy model this enforces.

#![forbid(unsafe_code)]

mod envelope;
mod ids;
mod keys;
mod receipt;

pub mod pow;

/// The HPKE KEM used throughout: DHKEM(X25519, HKDF-SHA256).
pub(crate) type Kem = hpke::kem::X25519HkdfSha256;

pub use envelope::{
    open_with_branch, open_with_claim, seal_to_branch, seal_to_claim, Envelope, EnvelopeKind,
    Inner, RecipientTag,
};
pub use ids::{BranchId, ClaimId};
pub use keys::{BranchKeys, BranchPublic, SecretKeyBytes};
pub use pow::PowStamp;
pub use receipt::{ClaimKeys, ClaimPublic, Receipt};

/// Errors surfaced by kilio-seal. Deliberately coarse — callers should not be
/// able to distinguish *why* a decryption failed (that would leak an oracle).
#[derive(Debug, thiserror::Error)]
pub enum SealError {
    #[error("malformed or wrong key")]
    BadKey,
    #[error("invalid signature")]
    BadSignature,
    #[error("malformed identifier")]
    BadId,
    #[error("invalid receipt phrase")]
    BadReceipt,
    #[error("key derivation failed")]
    Kdf,
    #[error("unsupported envelope version")]
    Version,
    #[error("payload too large to seal")]
    TooLarge,
    #[error("encode failed")]
    Encode,
    #[error("decode failed")]
    Decode,
    #[error("seal failed")]
    Seal,
    #[error("open failed")]
    Open,
}
