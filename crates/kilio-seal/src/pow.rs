//! Anonymous cold-contact proof-of-work.
//!
//! kilio refuses fully-unauthenticated injection but also refuses to demand an
//! identity from a reporter. The reconciliation (kotva §9.2's model): first
//! contact carries a small proof-of-work stamp bound to *this exact sealed
//! message*. A human reporter pays a fraction of a second of CPU once; a
//! spammer pays it for every message, and the cost is unlinkable across
//! submissions. The difficulty is per-branch and can be raised under attack.

use serde::{Deserialize, Serialize};

use crate::envelope::Envelope;

/// A solved proof-of-work stamp attached to a submission.
#[derive(Clone, Copy, PartialEq, Eq, Debug, Serialize, Deserialize)]
pub struct PowStamp {
    pub nonce: u64,
    pub bits: u8,
}

/// The challenge a stamp must be solved against: bound to the envelope's sealed
/// bytes so a stamp cannot be replayed onto a different message.
pub fn envelope_challenge(env: &Envelope) -> [u8; 32] {
    let mut h = blake3::Hasher::new_derive_key("kilio/pow-challenge/v1");
    h.update(&[env.v]);
    h.update(&env.enc);
    h.update(&env.ciphertext);
    *h.finalize().as_bytes()
}

/// Solve a proof-of-work of `bits` leading zero bits over `challenge`.
///
/// Cost is ~`2^bits` hashes. 20 bits ≈ a second on a laptop; keep reporter-side
/// difficulty modest (18–22) and raise only under abuse.
pub fn solve(challenge: &[u8; 32], bits: u8) -> PowStamp {
    let mut nonce: u64 = 0;
    loop {
        if meets(challenge, nonce, bits) {
            return PowStamp { nonce, bits };
        }
        nonce = nonce.wrapping_add(1);
    }
}

/// Verify a stamp against a challenge, requiring at least `min_bits` of work.
pub fn verify(challenge: &[u8; 32], stamp: &PowStamp, min_bits: u8) -> bool {
    stamp.bits >= min_bits && meets(challenge, stamp.nonce, stamp.bits)
}

fn meets(challenge: &[u8; 32], nonce: u64, bits: u8) -> bool {
    let mut h = blake3::Hasher::new();
    h.update(challenge);
    h.update(&nonce.to_le_bytes());
    leading_zero_bits(h.finalize().as_bytes()) >= u32::from(bits)
}

fn leading_zero_bits(hash: &[u8]) -> u32 {
    let mut n = 0u32;
    for &b in hash {
        if b == 0 {
            n += 8;
        } else {
            n += b.leading_zeros(); // u8::leading_zeros() is already 0..=8
            break;
        }
    }
    n
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn solve_then_verify() {
        let challenge = [7u8; 32];
        let stamp = solve(&challenge, 12);
        assert!(verify(&challenge, &stamp, 12));
        // Also satisfies any lower requirement.
        assert!(verify(&challenge, &stamp, 8));
    }

    #[test]
    fn rejects_insufficient_work() {
        let challenge = [7u8; 32];
        let stamp = solve(&challenge, 8);
        // Claiming more bits than actually solved must fail.
        let lying = PowStamp {
            nonce: stamp.nonce,
            bits: 30,
        };
        assert!(!verify(&challenge, &lying, 30));
    }

    #[test]
    fn rejects_wrong_challenge() {
        let stamp = solve(&[1u8; 32], 12);
        assert!(!verify(&[2u8; 32], &stamp, 12));
    }

    #[test]
    fn leading_zero_count() {
        assert_eq!(leading_zero_bits(&[0x00, 0xff]), 8);
        assert_eq!(leading_zero_bits(&[0x0f, 0xff]), 4);
        assert_eq!(leading_zero_bits(&[0xff]), 0);
        assert_eq!(leading_zero_bits(&[0x00, 0x00, 0x80]), 16);
    }
}
