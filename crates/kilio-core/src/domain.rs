//! Domain records. Everything a human wrote lives *sealed* inside an
//! [`kilio_seal::Envelope`]; these records hold only the routing and lifecycle
//! metadata the server legitimately needs to see in the clear.

use serde::{Deserialize, Serialize};

use kilio_seal::{BranchId, ClaimId};

/// Where a claim sits in its lifecycle. Cleartext — the server routes and
/// filters on it, but it says nothing about the claim's content.
#[derive(Clone, Copy, PartialEq, Eq, Debug, Serialize, Deserialize)]
pub enum ClaimStatus {
    New,
    Triaged,
    InProgress,
    Resolved,
    Closed,
}

impl ClaimStatus {
    pub fn as_str(&self) -> &'static str {
        match self {
            ClaimStatus::New => "new",
            ClaimStatus::Triaged => "triaged",
            ClaimStatus::InProgress => "in_progress",
            ClaimStatus::Resolved => "resolved",
            ClaimStatus::Closed => "closed",
        }
    }
    pub fn from_str(s: &str) -> Option<Self> {
        Some(match s {
            "new" => ClaimStatus::New,
            "triaged" => ClaimStatus::Triaged,
            "in_progress" => ClaimStatus::InProgress,
            "resolved" => ClaimStatus::Resolved,
            "closed" => ClaimStatus::Closed,
            _ => return None,
        })
    }
}

/// Who a message came from. The reporter is never named — this is a role, not
/// an identity.
#[derive(Clone, Copy, PartialEq, Eq, Debug, Serialize, Deserialize)]
pub enum Direction {
    Reporter,
    Handler,
}

impl Direction {
    pub fn as_str(&self) -> &'static str {
        match self {
            Direction::Reporter => "reporter",
            Direction::Handler => "handler",
        }
    }
    pub fn from_str(s: &str) -> Option<Self> {
        match s {
            "reporter" => Some(Direction::Reporter),
            "handler" => Some(Direction::Handler),
            _ => None,
        }
    }
}

/// A branch: a keyed destination claims seal to. The public keys are safe to
/// serve to any reporter; the private half never lives here (see the keystore).
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct BranchRecord {
    pub id: BranchId,
    pub name: String,
    pub kem_public: Vec<u8>,
    pub sign_public: [u8; 32],
    /// Proof-of-work difficulty (leading zero bits) required for cold contact.
    pub pow_bits: u8,
    pub created_at: u64,
    pub active: bool,
}

/// A claim: the routing/lifecycle shell around a sealed conversation.
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct ClaimRecord {
    pub claim_id: ClaimId,
    pub branch_id: BranchId,
    pub status: ClaimStatus,
    pub size_bucket: u32,
    pub created_at: u64,
    pub updated_at: u64,
    pub message_count: u32,
}

/// One turn of the conversation, sealed. `envelope` is the CBOR-encoded
/// [`kilio_seal::Envelope`] exactly as it will be opened by the other side.
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct MessageRecord {
    pub id: String,
    pub claim_id: ClaimId,
    pub direction: Direction,
    pub created_at: u64,
    pub envelope: Vec<u8>,
}

/// A content-free audit entry: records *that* something happened, never what.
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct AuditEvent {
    pub id: String,
    pub actor: String,
    pub action: String,
    pub claim_id: Option<ClaimId>,
    pub at: u64,
}
