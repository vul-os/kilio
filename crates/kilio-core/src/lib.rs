//! # kilio-core
//!
//! The domain heart of kilio: the sealed store, the seams, and branch scoping.
//! It holds the shape of a claim's lifecycle and the rules for who may read
//! what — but it never sees claim content. Everything a human wrote stays
//! sealed inside [`kilio_seal::Envelope`] bytes this crate moves but cannot
//! open.
//!
//! - [`SealedStore`] — SQLite that stores ciphertext + routing metadata only.
//! - [`Requester`] / [`branch_scoped_key`] — the ofisi scoping invariants: one
//!   authorization choke point, one scoped-key builder.
//! - [`Delivery`] / [`Reachability`] — the seams, with local defaults compiled
//!   in and adapters wired only at the composition root.

#![forbid(unsafe_code)]

mod domain;
mod scoping;
mod store;

pub mod delivery;
pub mod reachability;

pub use domain::{
    AuditEvent, BranchRecord, ClaimRecord, ClaimStatus, Direction, MessageRecord,
};
pub use scoping::{branch_scoped_key, DeployMode, Requester};
pub use store::SealedStore;

pub use delivery::{Delivery, KotvaDelivery, LocalDelivery};
pub use reachability::{LocalOnly, Reachability, SubprocessTunnel, TunnelProvider, TunnelStatus};

/// Errors from kilio-core.
#[derive(Debug, thiserror::Error)]
pub enum CoreError {
    #[error("storage error: {0}")]
    Sqlite(#[from] rusqlite::Error),
    #[error("seal error: {0}")]
    Seal(#[from] kilio_seal::SealError),
    #[error("could not encode envelope")]
    Encode,
    #[error("envelope is not addressed to a branch")]
    WrongRecipient,
    #[error("unknown branch")]
    UnknownBranch,
    #[error("branch is inactive")]
    BranchInactive,
    #[error("a claim with that id already exists")]
    ClaimExists,
    #[error("unknown claim")]
    UnknownClaim,
    /// Returned for both "does not exist" and "not authorized" — the two are
    /// deliberately indistinguishable so existence is never leaked.
    #[error("not found")]
    NotFound,
    #[error("address is not loopback")]
    NotLoopback,
    #[error("unsupported: {0}")]
    Unsupported(&'static str),
}

#[cfg(test)]
mod integration_tests {
    use super::*;
    use kilio_seal::{
        open_with_branch, open_with_claim, seal_to_branch, seal_to_claim, BranchKeys,
        EnvelopeKind, Inner, Receipt,
    };

    fn branch_record(b: &BranchKeys, name: &str) -> BranchRecord {
        let p = b.public();
        BranchRecord {
            id: p.branch_id,
            name: name.into(),
            kem_public: p.kem_public,
            sign_public: p.sign_public,
            pow_bits: 8,
            created_at: 1,
            active: true,
        }
    }

    fn submission(b: &BranchKeys, claim_pub: kilio_seal::ClaimPublic, body: &[u8]) -> kilio_seal::Envelope {
        seal_to_branch(
            &b.public(),
            EnvelopeKind::Submission,
            &Inner {
                from: Some(claim_pub),
                created_at: 1,
                body: body.to_vec(),
            },
        )
        .unwrap()
    }

    #[test]
    fn end_to_end_sealed_two_way_flow() {
        let mut store = SealedStore::open_in_memory().unwrap();
        let branch = BranchKeys::generate();
        store.put_branch(&branch_record(&branch, "corporate")).unwrap();

        // Reporter mints a receipt, derives a claim identity, seals a submission.
        let receipt = Receipt::generate();
        let claim = receipt.derive(&branch.branch_id).unwrap();
        let env = submission(&branch, claim.public(), b"I have a report");
        let claim_id = store.ingest_submission(claim.claim_id, &env).unwrap();

        // Handler is scoped to this branch and can list + open the claim.
        let handler = Requester::Handler {
            id: "h1".into(),
            branches: vec![branch.branch_id],
            admin: false,
        };
        let claims = store.list_claims(&handler).unwrap();
        assert_eq!(claims.len(), 1);
        assert_eq!(claims[0].claim_id, claim_id);

        // Handler reads the sealed message and decrypts it with the branch key.
        let msgs = store.messages_for_claim(&handler, &claim_id).unwrap();
        assert_eq!(msgs.len(), 1);
        let env0: kilio_seal::Envelope = ciborium::from_reader(&msgs[0].envelope[..]).unwrap();
        let inner = open_with_branch(&branch, &env0).unwrap();
        assert_eq!(inner.body, b"I have a report");
        let reporter_pub = inner.from.expect("sealed sender present");

        // Handler seals a reply to the claim key and stores it.
        let reply = seal_to_claim(
            &reporter_pub,
            EnvelopeKind::HandlerReply,
            &Inner { from: None, created_at: 2, body: b"received".to_vec() },
        )
        .unwrap();
        store.append_handler_reply(&handler, &claim_id, &reply).unwrap();

        // Reporter returns with the same passphrase and reads the sealed reply.
        let reporter = Requester::Reporter { claim: claim_id };
        let msgs2 = store.messages_for_claim(&reporter, &claim_id).unwrap();
        assert_eq!(msgs2.len(), 2);
        let reply_env: kilio_seal::Envelope =
            ciborium::from_reader(&msgs2[1].envelope[..]).unwrap();
        let claim2 = Receipt::from_phrase(receipt.phrase())
            .unwrap()
            .derive(&branch.branch_id)
            .unwrap();
        let got = open_with_claim(&claim2, &reply_env).unwrap();
        assert_eq!(got.body, b"received");
    }

    #[test]
    fn branch_isolation_returns_not_found() {
        let mut store = SealedStore::open_in_memory().unwrap();
        let branch_a = BranchKeys::generate();
        let branch_b = BranchKeys::generate();
        store.put_branch(&branch_record(&branch_a, "a")).unwrap();
        store.put_branch(&branch_record(&branch_b, "b")).unwrap();

        let receipt = Receipt::generate();
        let claim = receipt.derive(&branch_a.branch_id).unwrap();
        let env = submission(&branch_a, claim.public(), b"secret to A");
        let claim_id = store.ingest_submission(claim.claim_id, &env).unwrap();

        // A handler scoped only to branch B cannot see A's claim — get returns
        // None and list is empty (no existence leak).
        let handler_b = Requester::Handler {
            id: "hb".into(),
            branches: vec![branch_b.branch_id],
            admin: false,
        };
        assert!(store.get_claim(&handler_b, &claim_id).unwrap().is_none());
        assert!(store.list_claims(&handler_b).unwrap().is_empty());
        assert!(store.messages_for_claim(&handler_b, &claim_id).unwrap().is_empty());

        // A reply attempt across the boundary is refused as NotFound.
        let reply = seal_to_claim(
            &claim.public(),
            EnvelopeKind::HandlerReply,
            &Inner { from: None, created_at: 2, body: b"x".to_vec() },
        )
        .unwrap();
        assert!(matches!(
            store.append_handler_reply(&handler_b, &claim_id, &reply),
            Err(CoreError::NotFound)
        ));
    }

    #[test]
    fn duplicate_submission_rejected() {
        let mut store = SealedStore::open_in_memory().unwrap();
        let branch = BranchKeys::generate();
        store.put_branch(&branch_record(&branch, "b")).unwrap();
        let claim = Receipt::generate().derive(&branch.branch_id).unwrap();
        let env = submission(&branch, claim.public(), b"first");
        store.ingest_submission(claim.claim_id, &env).unwrap();
        assert!(matches!(
            store.ingest_submission(claim.claim_id, &env),
            Err(CoreError::ClaimExists)
        ));
    }
}
