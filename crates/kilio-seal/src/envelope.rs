//! The sealed envelope — kilio's wire and at-rest unit.
//!
//! Shape mirrors kotva's MOTE `Envelope`/`Payload` split so a `KotvaDelivery`
//! seam can re-wrap rather than re-encrypt:
//!
//! * the **outer** [`Envelope`] is cleartext routing metadata + HPKE ciphertext;
//! * the **inner** [`Inner`] is sealed and carries the sender identity
//!   (*sealed sender* — intermediaries never learn who sent it) plus the body.
//!
//! Everything the server may route/rate-limit on (kind, recipient tag, size
//! bucket, version) is bound into the AEAD's associated data, so none of it can
//! be altered in transit without decryption failing.

use hpke::{
    aead::ChaCha20Poly1305, single_shot_open, single_shot_seal, Deserializable, OpModeR, OpModeS,
    Serializable,
};
use rand::rngs::OsRng;
use serde::{Deserialize, Serialize};

use crate::keys::{BranchKeys, BranchPublic};
use crate::receipt::{ClaimKeys, ClaimPublic};
use crate::{Kem, SealError};

type Aead = ChaCha20Poly1305;
type Kdf = hpke::kdf::HkdfSha256;

const SEAL_VERSION: u8 = 1;
const HPKE_INFO: &[u8] = b"kilio/hpke/v1";

/// Size buckets (bytes). Ciphertext is padded up to the smallest bucket that
/// fits, so wire length cannot fingerprint a claim. Above the top bucket we
/// round up to whole 4 MiB steps.
const BUCKETS: [usize; 6] = [
    4 * 1024,
    16 * 1024,
    64 * 1024,
    256 * 1024,
    1024 * 1024,
    4 * 1024 * 1024,
];
const BUCKET_STEP: usize = 4 * 1024 * 1024;

/// What a message is, in the clear (bound into AAD).
#[derive(Clone, Copy, PartialEq, Eq, Debug, Serialize, Deserialize)]
pub enum EnvelopeKind {
    /// First contact: a new anonymous claim sealed to a branch.
    Submission,
    /// A follow-up from the reporter on an existing claim.
    ReporterMessage,
    /// A reply from a handler, sealed to the claim's key.
    HandlerReply,
}

/// Which key the envelope is sealed to (bound into AAD, used for routing).
#[derive(Clone, PartialEq, Eq, Debug, Serialize, Deserialize)]
pub enum RecipientTag {
    Branch(crate::ids::BranchId),
    Claim(crate::ids::ClaimId),
}

/// The sealed inner payload. Only a holder of the recipient private key sees it.
#[derive(Clone, Serialize, Deserialize)]
pub struct Inner {
    /// Sealed sender. On a `Submission` this carries the reporter's per-claim
    /// public identity so the handler can seal replies back. Absent otherwise.
    pub from: Option<ClaimPublic>,
    /// Millis since epoch, set by the *sealer's* clock (unverified; advisory).
    pub created_at: u64,
    /// Opaque application content. kilio-core defines the schema; seal is
    /// content-agnostic and only moves bytes.
    pub body: Vec<u8>,
}

/// The cleartext outer envelope: what is stored and transmitted.
#[derive(Clone, Serialize, Deserialize)]
pub struct Envelope {
    pub v: u8,
    pub kind: EnvelopeKind,
    pub recipient: RecipientTag,
    /// HPKE encapsulated key.
    pub enc: Vec<u8>,
    /// AEAD ciphertext of the padded, CBOR-encoded [`Inner`].
    pub ciphertext: Vec<u8>,
    /// The padded plaintext length (a bucket size), for wire-length uniformity.
    pub size_bucket: u32,
}

/// Seal a payload to a branch (a reporter submitting or messaging).
pub fn seal_to_branch(
    branch: &BranchPublic,
    kind: EnvelopeKind,
    inner: &Inner,
) -> Result<Envelope, SealError> {
    let tag = RecipientTag::Branch(branch.branch_id);
    seal(&branch.kem_pk()?, tag, kind, inner)
}

/// Seal a reply to a claim (a handler answering a reporter).
pub fn seal_to_claim(
    claim: &ClaimPublic,
    kind: EnvelopeKind,
    inner: &Inner,
) -> Result<Envelope, SealError> {
    let tag = RecipientTag::Claim(claim.claim_id);
    seal(&claim.kem_pk()?, tag, kind, inner)
}

/// Open an envelope sealed to a branch (handler side).
pub fn open_with_branch(keys: &BranchKeys, env: &Envelope) -> Result<Inner, SealError> {
    open(&keys.kem_private()?, env)
}

/// Open an envelope sealed to a claim (reporter side).
pub fn open_with_claim(keys: &ClaimKeys, env: &Envelope) -> Result<Inner, SealError> {
    open(&keys.kem_private()?, env)
}

fn seal(
    recipient_pk: &<Kem as hpke::Kem>::PublicKey,
    tag: RecipientTag,
    kind: EnvelopeKind,
    inner: &Inner,
) -> Result<Envelope, SealError> {
    let mut plaintext = Vec::new();
    ciborium::into_writer(inner, &mut plaintext).map_err(|_| SealError::Encode)?;

    let bucket = bucketize(plaintext.len());
    plaintext.resize(bucket, 0); // zero-pad; ciborium ignores trailing bytes on read
    let size_bucket = u32::try_from(bucket).map_err(|_| SealError::TooLarge)?;

    let aad = associated_data(SEAL_VERSION, kind, &tag, size_bucket);
    let (encapped, ciphertext) = single_shot_seal::<Aead, Kdf, Kem, _>(
        &OpModeS::Base,
        recipient_pk,
        HPKE_INFO,
        &plaintext,
        &aad,
        &mut OsRng,
    )
    .map_err(|_| SealError::Seal)?;

    Ok(Envelope {
        v: SEAL_VERSION,
        kind,
        recipient: tag,
        enc: encapped.to_bytes().to_vec(),
        ciphertext,
        size_bucket,
    })
}

fn open(recipient_sk: &<Kem as hpke::Kem>::PrivateKey, env: &Envelope) -> Result<Inner, SealError> {
    if env.v != SEAL_VERSION {
        return Err(SealError::Version);
    }
    let encapped = <Kem as hpke::Kem>::EncappedKey::from_bytes(&env.enc)
        .map_err(|_| SealError::BadKey)?;
    let aad = associated_data(env.v, env.kind, &env.recipient, env.size_bucket);
    let padded = single_shot_open::<Aead, Kdf, Kem>(
        &OpModeR::Base,
        recipient_sk,
        &encapped,
        HPKE_INFO,
        &env.ciphertext,
        &aad,
    )
    .map_err(|_| SealError::Open)?;

    // ciborium reads exactly one value and stops, ignoring the zero padding.
    let inner: Inner = ciborium::from_reader(&padded[..]).map_err(|_| SealError::Decode)?;
    Ok(inner)
}

/// Deterministic associated data binding every cleartext routing field.
fn associated_data(v: u8, kind: EnvelopeKind, tag: &RecipientTag, size_bucket: u32) -> Vec<u8> {
    #[derive(Serialize)]
    struct Aad<'a> {
        v: u8,
        kind: EnvelopeKind,
        tag: &'a RecipientTag,
        size_bucket: u32,
    }
    let mut buf = Vec::new();
    ciborium::into_writer(
        &Aad {
            v,
            kind,
            tag,
            size_bucket,
        },
        &mut buf,
    )
    .expect("AAD serialization is infallible");
    buf
}

fn bucketize(len: usize) -> usize {
    for &b in &BUCKETS {
        if len <= b {
            return b;
        }
    }
    len.div_ceil(BUCKET_STEP) * BUCKET_STEP
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::ids::BranchId;
    use crate::receipt::Receipt;

    fn inner(body: &[u8], from: Option<ClaimPublic>) -> Inner {
        Inner {
            from,
            created_at: 1_600_000_000_000,
            body: body.to_vec(),
        }
    }

    #[test]
    fn branch_seal_open_roundtrip() {
        let branch = BranchKeys::generate();
        let env = seal_to_branch(&branch.public(), EnvelopeKind::Submission, &inner(b"hello", None))
            .unwrap();
        let got = open_with_branch(&branch, &env).unwrap();
        assert_eq!(got.body, b"hello");
    }

    #[test]
    fn full_two_way_flow() {
        // Reporter derives a claim identity and seals a submission to the branch.
        let branch = BranchKeys::generate();
        let receipt = Receipt::generate();
        let claim = receipt.derive(&branch.branch_id).unwrap();
        let submission = seal_to_branch(
            &branch.public(),
            EnvelopeKind::Submission,
            &inner(b"I was harassed", Some(claim.public())),
        )
        .unwrap();

        // Handler opens it and recovers the sealed-sender claim identity.
        let opened = open_with_branch(&branch, &submission).unwrap();
        assert_eq!(opened.body, b"I was harassed");
        let reporter_pub = opened.from.expect("submission carries sealed sender");

        // Handler seals a reply to the claim key.
        let reply =
            seal_to_claim(&reporter_pub, EnvelopeKind::HandlerReply, &inner(b"we are looking", None))
                .unwrap();

        // Reporter returns with the same phrase and opens the reply.
        let claim2 = Receipt::from_phrase(receipt.phrase())
            .unwrap()
            .derive(&branch.branch_id)
            .unwrap();
        let got = open_with_claim(&claim2, &reply).unwrap();
        assert_eq!(got.body, b"we are looking");
    }

    #[test]
    fn wrong_key_cannot_open() {
        let branch = BranchKeys::generate();
        let other = BranchKeys::generate();
        let env =
            seal_to_branch(&branch.public(), EnvelopeKind::Submission, &inner(b"secret", None))
                .unwrap();
        assert!(open_with_branch(&other, &env).is_err());
    }

    #[test]
    fn tampered_kind_fails_to_open() {
        let branch = BranchKeys::generate();
        let mut env =
            seal_to_branch(&branch.public(), EnvelopeKind::Submission, &inner(b"x", None)).unwrap();
        env.kind = EnvelopeKind::HandlerReply; // AAD no longer matches
        assert!(open_with_branch(&branch, &env).is_err());
    }

    #[test]
    fn size_is_bucketed() {
        let branch = BranchKeys::generate();
        let env =
            seal_to_branch(&branch.public(), EnvelopeKind::Submission, &inner(b"tiny", None))
                .unwrap();
        assert_eq!(env.size_bucket, 4 * 1024);
    }

    #[test]
    fn bucketize_steps() {
        assert_eq!(bucketize(1), 4 * 1024);
        assert_eq!(bucketize(4 * 1024), 4 * 1024);
        assert_eq!(bucketize(4 * 1024 + 1), 16 * 1024);
        assert_eq!(bucketize(5 * 1024 * 1024), 8 * 1024 * 1024);
    }
}
