//! The `Delivery` seam — where a sealed submission goes *in addition to* the
//! local store.
//!
//! The local store write always happens; `Delivery` is the optional extra hop
//! (e.g. mirroring sealed claims to an external ombudsman over a content-blind
//! kotva rendezvous mailbox). The default forwards nowhere. Adapters are wired
//! only at the composition root — core never depends on one.

use kilio_seal::{ClaimId, Envelope};

use crate::CoreError;

pub trait Delivery: Send + Sync {
    /// Forward a sealed submission onward. The envelope is already sealed at
    /// source, so a forwarder only ever handles ciphertext.
    fn forward(&self, claim_id: &ClaimId, env: &Envelope) -> Result<(), CoreError>;

    /// Human label for status/diagnostics.
    fn label(&self) -> &'static str;
}

/// Default: forwards nowhere. Standalone/desktop deployments use this — the
/// sealed claim lives only in the local store.
pub struct LocalDelivery;

impl Delivery for LocalDelivery {
    fn forward(&self, _claim_id: &ClaimId, _env: &Envelope) -> Result<(), CoreError> {
        Ok(())
    }
    fn label(&self) -> &'static str {
        "local"
    }
}

/// Forward sealed claims to a kotva rendezvous mailbox
/// (`POST {relay}/mailbox/{recipient}`, content-blind). Decentralized delivery
/// across orgs without a shared server.
///
/// The relay never sees plaintext — the envelope is sealed at source. Wiring
/// the HTTP deposit lives in the server layer (where async networking is
/// available); this struct carries the target and marks the seam.
pub struct KotvaDelivery {
    pub relay_url: String,
    pub recipient: String,
}

impl Delivery for KotvaDelivery {
    fn forward(&self, _claim_id: &ClaimId, _env: &Envelope) -> Result<(), CoreError> {
        Err(CoreError::Unsupported(
            "KotvaDelivery.forward is wired in the server layer (async mailbox deposit)",
        ))
    }
    fn label(&self) -> &'static str {
        "kotva"
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn local_forwards_nowhere() {
        let d = LocalDelivery;
        let cid = ClaimId::derive(b"c");
        // A dummy envelope is not needed — local forward is a no-op.
        let env = dummy_env();
        assert!(d.forward(&cid, &env).is_ok());
        assert_eq!(d.label(), "local");
    }

    fn dummy_env() -> Envelope {
        let b = kilio_seal::BranchKeys::generate();
        kilio_seal::seal_to_branch(
            &b.public(),
            kilio_seal::EnvelopeKind::Submission,
            &kilio_seal::Inner {
                from: None,
                created_at: 0,
                body: vec![],
            },
        )
        .unwrap()
    }
}
