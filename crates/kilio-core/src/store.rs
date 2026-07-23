//! The sealed store: SQLite that holds **ciphertext and routing metadata only**.
//!
//! No column in this store contains anything a human wrote — claim bodies,
//! attachments, and the sender identity all live inside the sealed
//! [`kilio_seal::Envelope`] bytes in `messages.envelope`, which this store can
//! move but cannot read. Reads that cross a branch boundary are refused through
//! the [`Requester`] choke point and return `Ok(None)` (never an existence
//! leak).

use rusqlite::{params, Connection, OptionalExtension};

use kilio_seal::{BranchId, ClaimId, Envelope, RecipientTag};

use crate::domain::*;
use crate::scoping::Requester;
use crate::CoreError;

/// A sealed SQLite store. Single-connection; wrap in a mutex / run on a
/// blocking thread from an async server.
pub struct SealedStore {
    conn: Connection,
}

impl SealedStore {
    /// Open (creating if needed) a store at `path` and apply the schema.
    pub fn open(path: &str) -> Result<Self, CoreError> {
        let conn = Connection::open(path)?;
        Self::from_conn(conn)
    }

    /// An ephemeral in-memory store (tests).
    pub fn open_in_memory() -> Result<Self, CoreError> {
        Self::from_conn(Connection::open_in_memory()?)
    }

    fn from_conn(conn: Connection) -> Result<Self, CoreError> {
        conn.pragma_update(None, "journal_mode", "WAL")?;
        conn.pragma_update(None, "foreign_keys", "ON")?;
        conn.execute_batch(SCHEMA)?;
        Ok(Self { conn })
    }

    // ---- branches -------------------------------------------------------

    pub fn put_branch(&self, b: &BranchRecord) -> Result<(), CoreError> {
        self.conn.execute(
            "INSERT INTO branches (id,name,kem_public,sign_public,pow_bits,created_at,active)
             VALUES (?1,?2,?3,?4,?5,?6,?7)
             ON CONFLICT(id) DO UPDATE SET name=?2, pow_bits=?5, active=?7",
            params![
                b.id.to_hex(),
                b.name,
                b.kem_public,
                b.sign_public.to_vec(),
                b.pow_bits as i64,
                b.created_at as i64,
                b.active as i64,
            ],
        )?;
        Ok(())
    }

    pub fn get_branch(&self, id: &BranchId) -> Result<Option<BranchRecord>, CoreError> {
        self.conn
            .query_row(
                "SELECT id,name,kem_public,sign_public,pow_bits,created_at,active
                 FROM branches WHERE id=?1",
                params![id.to_hex()],
                row_to_branch,
            )
            .optional()
            .map_err(Into::into)
    }

    pub fn list_branches(&self, only_active: bool) -> Result<Vec<BranchRecord>, CoreError> {
        let mut stmt = self.conn.prepare(
            "SELECT id,name,kem_public,sign_public,pow_bits,created_at,active
             FROM branches WHERE (?1=0 OR active=1) ORDER BY name",
        )?;
        let rows = stmt.query_map(params![only_active as i64], row_to_branch)?;
        Ok(rows.collect::<Result<_, _>>()?)
    }

    // ---- ingest & messages ---------------------------------------------

    /// Ingest a first-contact submission. `claim_id` is the reporter's public,
    /// cleartext claim handle; the sealed body stays sealed. The caller
    /// (server) must have already verified the proof-of-work stamp.
    pub fn ingest_submission(
        &mut self,
        claim_id: ClaimId,
        env: &Envelope,
    ) -> Result<ClaimId, CoreError> {
        let branch_id = match &env.recipient {
            RecipientTag::Branch(b) => *b,
            RecipientTag::Claim(_) => return Err(CoreError::WrongRecipient),
        };
        let branch = self.get_branch(&branch_id)?.ok_or(CoreError::UnknownBranch)?;
        if !branch.active {
            return Err(CoreError::BranchInactive);
        }
        if self.claim_exists(&claim_id)? {
            return Err(CoreError::ClaimExists);
        }

        let now = now_ms();
        let bytes = encode_envelope(env)?;
        let tx = self.conn.transaction()?;
        tx.execute(
            "INSERT INTO claims (claim_id,branch_id,status,size_bucket,created_at,updated_at,message_count)
             VALUES (?1,?2,?3,?4,?5,?5,1)",
            params![
                claim_id.to_hex(),
                branch_id.to_hex(),
                ClaimStatus::New.as_str(),
                env.size_bucket as i64,
                now as i64,
            ],
        )?;
        insert_message(&tx, &claim_id, Direction::Reporter, now, &bytes)?;
        tx.commit()?;
        Ok(claim_id)
    }

    /// A returning reporter appends a follow-up. Server must have verified the
    /// reporter's signature over the claim first.
    pub fn append_reporter_message(
        &mut self,
        claim_id: &ClaimId,
        env: &Envelope,
    ) -> Result<(), CoreError> {
        if !self.claim_exists(claim_id)? {
            return Err(CoreError::UnknownClaim);
        }
        let now = now_ms();
        let bytes = encode_envelope(env)?;
        let tx = self.conn.transaction()?;
        insert_message(&tx, claim_id, Direction::Reporter, now, &bytes)?;
        bump_claim(&tx, claim_id, now)?;
        tx.commit()?;
        Ok(())
    }

    /// A handler replies, sealed to the claim key. Authorized through the
    /// choke point; writes a content-free audit entry.
    pub fn append_handler_reply(
        &mut self,
        requester: &Requester,
        claim_id: &ClaimId,
        env: &Envelope,
    ) -> Result<(), CoreError> {
        let claim = self.claim_row(claim_id)?.ok_or(CoreError::NotFound)?;
        if !requester.may_access_claim(&claim.branch_id, claim_id) {
            return Err(CoreError::NotFound); // no existence leak
        }
        let now = now_ms();
        let bytes = encode_envelope(env)?;
        let actor = requester.actor_label();
        let tx = self.conn.transaction()?;
        insert_message(&tx, claim_id, Direction::Handler, now, &bytes)?;
        bump_claim(&tx, claim_id, now)?;
        insert_audit(&tx, &actor, "handler_reply", Some(claim_id), now)?;
        tx.commit()?;
        Ok(())
    }

    // ---- scoped reads ---------------------------------------------------

    /// List claims visible to `requester`. Handlers see their branches (admins
    /// all); reporters see only their own claim; the public sees nothing.
    pub fn list_claims(&self, requester: &Requester) -> Result<Vec<ClaimRecord>, CoreError> {
        let all = self.all_claims()?;
        Ok(all
            .into_iter()
            .filter(|c| requester.may_access_claim(&c.branch_id, &c.claim_id))
            .collect())
    }

    /// Fetch one claim, or `None` if it doesn't exist *or* the requester may not
    /// see it — the two are indistinguishable to the caller by design.
    pub fn get_claim(
        &self,
        requester: &Requester,
        claim_id: &ClaimId,
    ) -> Result<Option<ClaimRecord>, CoreError> {
        match self.claim_row(claim_id)? {
            Some(c) if requester.may_access_claim(&c.branch_id, claim_id) => Ok(Some(c)),
            _ => Ok(None),
        }
    }

    /// Sealed messages for a claim, oldest first, authorized through the choke
    /// point. Used by a returning reporter (their own claim) and by handlers.
    pub fn messages_for_claim(
        &self,
        requester: &Requester,
        claim_id: &ClaimId,
    ) -> Result<Vec<MessageRecord>, CoreError> {
        let claim = match self.claim_row(claim_id)? {
            Some(c) => c,
            None => return Ok(vec![]),
        };
        if !requester.may_access_claim(&claim.branch_id, claim_id) {
            return Ok(vec![]);
        }
        let mut stmt = self.conn.prepare(
            "SELECT id,claim_id,direction,created_at,envelope FROM messages
             WHERE claim_id=?1 ORDER BY created_at ASC, id ASC",
        )?;
        let rows = stmt.query_map(params![claim_id.to_hex()], row_to_message)?;
        Ok(rows.collect::<Result<_, _>>()?)
    }

    /// Update a claim's status (handler action, audited).
    pub fn set_status(
        &mut self,
        requester: &Requester,
        claim_id: &ClaimId,
        status: ClaimStatus,
    ) -> Result<(), CoreError> {
        let claim = self.claim_row(claim_id)?.ok_or(CoreError::NotFound)?;
        if !requester.may_access_claim(&claim.branch_id, claim_id) {
            return Err(CoreError::NotFound);
        }
        let now = now_ms();
        let actor = requester.actor_label();
        let tx = self.conn.transaction()?;
        tx.execute(
            "UPDATE claims SET status=?2, updated_at=?3 WHERE claim_id=?1",
            params![claim_id.to_hex(), status.as_str(), now as i64],
        )?;
        insert_audit(&tx, &actor, &format!("status:{}", status.as_str()), Some(claim_id), now)?;
        tx.commit()?;
        Ok(())
    }

    /// The content-free audit trail for a claim.
    pub fn audit_for_claim(
        &self,
        requester: &Requester,
        claim_id: &ClaimId,
    ) -> Result<Vec<AuditEvent>, CoreError> {
        let claim = match self.claim_row(claim_id)? {
            Some(c) => c,
            None => return Ok(vec![]),
        };
        if !requester.may_access_claim(&claim.branch_id, claim_id) {
            return Ok(vec![]);
        }
        let mut stmt = self.conn.prepare(
            "SELECT id,actor,action,claim_id,at FROM audit WHERE claim_id=?1 ORDER BY at ASC",
        )?;
        let rows = stmt.query_map(params![claim_id.to_hex()], row_to_audit)?;
        Ok(rows.collect::<Result<_, _>>()?)
    }

    // ---- internals ------------------------------------------------------

    fn claim_exists(&self, claim_id: &ClaimId) -> Result<bool, CoreError> {
        Ok(self
            .conn
            .query_row(
                "SELECT 1 FROM claims WHERE claim_id=?1",
                params![claim_id.to_hex()],
                |_| Ok(()),
            )
            .optional()?
            .is_some())
    }

    fn claim_row(&self, claim_id: &ClaimId) -> Result<Option<ClaimRecord>, CoreError> {
        self.conn
            .query_row(
                "SELECT claim_id,branch_id,status,size_bucket,created_at,updated_at,message_count
                 FROM claims WHERE claim_id=?1",
                params![claim_id.to_hex()],
                row_to_claim,
            )
            .optional()
            .map_err(Into::into)
    }

    fn all_claims(&self) -> Result<Vec<ClaimRecord>, CoreError> {
        let mut stmt = self.conn.prepare(
            "SELECT claim_id,branch_id,status,size_bucket,created_at,updated_at,message_count
             FROM claims ORDER BY updated_at DESC",
        )?;
        let rows = stmt.query_map([], row_to_claim)?;
        Ok(rows.collect::<Result<_, _>>()?)
    }
}

const SCHEMA: &str = r#"
CREATE TABLE IF NOT EXISTS branches (
  id          TEXT PRIMARY KEY,
  name        TEXT NOT NULL,
  kem_public  BLOB NOT NULL,
  sign_public BLOB NOT NULL,
  pow_bits    INTEGER NOT NULL,
  created_at  INTEGER NOT NULL,
  active      INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS claims (
  claim_id      TEXT PRIMARY KEY,
  branch_id     TEXT NOT NULL REFERENCES branches(id),
  status        TEXT NOT NULL,
  size_bucket   INTEGER NOT NULL,
  created_at    INTEGER NOT NULL,
  updated_at    INTEGER NOT NULL,
  message_count INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_claims_branch ON claims(branch_id);
CREATE TABLE IF NOT EXISTS messages (
  id         TEXT PRIMARY KEY,
  claim_id   TEXT NOT NULL REFERENCES claims(claim_id),
  direction  TEXT NOT NULL,
  created_at INTEGER NOT NULL,
  envelope   BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_messages_claim ON messages(claim_id);
CREATE TABLE IF NOT EXISTS audit (
  id       TEXT PRIMARY KEY,
  actor    TEXT NOT NULL,
  action   TEXT NOT NULL,
  claim_id TEXT,
  at       INTEGER NOT NULL
);
"#;

fn insert_message(
    tx: &rusqlite::Transaction,
    claim_id: &ClaimId,
    dir: Direction,
    now: u64,
    envelope: &[u8],
) -> Result<(), CoreError> {
    tx.execute(
        "INSERT INTO messages (id,claim_id,direction,created_at,envelope) VALUES (?1,?2,?3,?4,?5)",
        params![new_id(), claim_id.to_hex(), dir.as_str(), now as i64, envelope],
    )?;
    Ok(())
}

fn bump_claim(tx: &rusqlite::Transaction, claim_id: &ClaimId, now: u64) -> Result<(), CoreError> {
    tx.execute(
        "UPDATE claims SET updated_at=?2, message_count=message_count+1 WHERE claim_id=?1",
        params![claim_id.to_hex(), now as i64],
    )?;
    Ok(())
}

fn insert_audit(
    tx: &rusqlite::Transaction,
    actor: &str,
    action: &str,
    claim_id: Option<&ClaimId>,
    now: u64,
) -> Result<(), CoreError> {
    tx.execute(
        "INSERT INTO audit (id,actor,action,claim_id,at) VALUES (?1,?2,?3,?4,?5)",
        params![
            new_id(),
            actor,
            action,
            claim_id.map(|c| c.to_hex()),
            now as i64
        ],
    )?;
    Ok(())
}

fn row_to_branch(r: &rusqlite::Row) -> rusqlite::Result<BranchRecord> {
    let sign: Vec<u8> = r.get(3)?;
    let mut sign_public = [0u8; 32];
    if sign.len() == 32 {
        sign_public.copy_from_slice(&sign);
    }
    Ok(BranchRecord {
        id: BranchId::from_hex(&r.get::<_, String>(0)?).unwrap_or(BranchId([0; 16])),
        name: r.get(1)?,
        kem_public: r.get(2)?,
        sign_public,
        pow_bits: r.get::<_, i64>(4)? as u8,
        created_at: r.get::<_, i64>(5)? as u64,
        active: r.get::<_, i64>(6)? != 0,
    })
}

fn row_to_claim(r: &rusqlite::Row) -> rusqlite::Result<ClaimRecord> {
    Ok(ClaimRecord {
        claim_id: ClaimId::from_hex(&r.get::<_, String>(0)?).unwrap_or(ClaimId([0; 16])),
        branch_id: BranchId::from_hex(&r.get::<_, String>(1)?).unwrap_or(BranchId([0; 16])),
        status: ClaimStatus::from_str(&r.get::<_, String>(2)?).unwrap_or(ClaimStatus::New),
        size_bucket: r.get::<_, i64>(3)? as u32,
        created_at: r.get::<_, i64>(4)? as u64,
        updated_at: r.get::<_, i64>(5)? as u64,
        message_count: r.get::<_, i64>(6)? as u32,
    })
}

fn row_to_message(r: &rusqlite::Row) -> rusqlite::Result<MessageRecord> {
    Ok(MessageRecord {
        id: r.get(0)?,
        claim_id: ClaimId::from_hex(&r.get::<_, String>(1)?).unwrap_or(ClaimId([0; 16])),
        direction: Direction::from_str(&r.get::<_, String>(2)?).unwrap_or(Direction::Reporter),
        created_at: r.get::<_, i64>(3)? as u64,
        envelope: r.get(4)?,
    })
}

fn row_to_audit(r: &rusqlite::Row) -> rusqlite::Result<AuditEvent> {
    Ok(AuditEvent {
        id: r.get(0)?,
        actor: r.get(1)?,
        action: r.get(2)?,
        claim_id: r
            .get::<_, Option<String>>(3)?
            .and_then(|s| ClaimId::from_hex(&s).ok()),
        at: r.get::<_, i64>(4)? as u64,
    })
}

fn encode_envelope(env: &Envelope) -> Result<Vec<u8>, CoreError> {
    let mut v = Vec::new();
    ciborium::into_writer(env, &mut v).map_err(|_| CoreError::Encode)?;
    Ok(v)
}

fn now_ms() -> u64 {
    std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .map(|d| d.as_millis() as u64)
        .unwrap_or(0)
}

fn new_id() -> String {
    use rand::RngCore;
    let mut b = [0u8; 16];
    rand::rngs::OsRng.fill_bytes(&mut b);
    hex::encode(b)
}
