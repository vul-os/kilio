//! Branch scoping — the ofisi pattern, ported.
//!
//! Two invariants make multi-branch isolation trustworthy:
//!
//! 1. **One `requester` choke point.** A caller's branch access is resolved
//!    from a server-verified session into a [`Requester`], never from a value
//!    the caller supplied. The store's read paths take a `&Requester` and go
//!    through [`Requester::may_access`] — the single place authorization is
//!    decided.
//! 2. **One scoped-key builder.** [`branch_scoped_key`] is the only way to
//!    address per-branch storage, with segment sanitization so a name can never
//!    escape its namespace.
//!
//! Denied reads return "not found", never "forbidden" — existence is not leaked.

use kilio_seal::BranchId;

/// How the deployment resolves identity (the ofisi `DEPLOY_MODE` enum).
#[derive(Clone, Copy, PartialEq, Eq, Debug)]
pub enum DeployMode {
    /// One sovereign instance; handler accounts are local.
    Standalone,
    /// Behind a Vulos OS gateway; identity is gateway-brokered. Fails closed
    /// without a configured verifier (enforced at the server's boot).
    Os,
}

/// A server-verified caller. Constructed only from an authenticated session or
/// a proven claim-control signature — never from raw request input.
#[derive(Clone, Debug)]
pub enum Requester {
    /// The public intake path: may submit, may not read anyone's claims.
    Public,
    /// A returning reporter who proved control of exactly one claim.
    Reporter { claim: kilio_seal::ClaimId },
    /// An authenticated handler, scoped to the branches they may open.
    Handler {
        id: String,
        branches: Vec<BranchId>,
        admin: bool,
    },
}

impl Requester {
    /// The single authorization decision for handler branch access.
    pub fn may_access_branch(&self, branch: &BranchId) -> bool {
        match self {
            Requester::Handler { branches, admin, .. } => *admin || branches.contains(branch),
            _ => false,
        }
    }

    /// Whether this requester may read a specific claim in a branch.
    pub fn may_access_claim(&self, branch: &BranchId, claim: &kilio_seal::ClaimId) -> bool {
        match self {
            Requester::Handler { .. } => self.may_access_branch(branch),
            Requester::Reporter { claim: owned } => owned == claim,
            Requester::Public => false,
        }
    }

    /// A stable string identifying the actor for the content-free audit log.
    pub fn actor_label(&self) -> String {
        match self {
            Requester::Public => "public".into(),
            Requester::Reporter { .. } => "reporter".into(),
            Requester::Handler { id, .. } => format!("handler:{id}"),
        }
    }
}

/// Build the one authoritative per-branch storage key: `<branch_id>/<name>`,
/// with segments sanitized so `name` can never contain a path separator or
/// traversal sequence.
pub fn branch_scoped_key(branch: &BranchId, name: &str) -> String {
    format!("{}/{}", branch.to_hex(), sanitize_segment(name))
}

fn sanitize_segment(seg: &str) -> String {
    seg.chars()
        .map(|c| match c {
            '/' | '\\' => '_',
            c if c.is_control() => '_',
            c => c,
        })
        .collect::<String>()
        .replace("..", "__")
}

#[cfg(test)]
mod tests {
    use super::*;

    fn bid(seed: &[u8]) -> BranchId {
        BranchId::derive(seed)
    }

    #[test]
    fn handler_scoped_to_its_branches() {
        let a = bid(b"a");
        let b = bid(b"b");
        let h = Requester::Handler {
            id: "h1".into(),
            branches: vec![a],
            admin: false,
        };
        assert!(h.may_access_branch(&a));
        assert!(!h.may_access_branch(&b));
    }

    #[test]
    fn admin_sees_all_branches() {
        let h = Requester::Handler {
            id: "root".into(),
            branches: vec![],
            admin: true,
        };
        assert!(h.may_access_branch(&bid(b"anything")));
    }

    #[test]
    fn reporter_only_its_own_claim() {
        let c1 = kilio_seal::ClaimId::derive(b"c1");
        let c2 = kilio_seal::ClaimId::derive(b"c2");
        let r = Requester::Reporter { claim: c1 };
        assert!(r.may_access_claim(&bid(b"a"), &c1));
        assert!(!r.may_access_claim(&bid(b"a"), &c2));
    }

    #[test]
    fn public_reads_nothing() {
        let p = Requester::Public;
        assert!(!p.may_access_branch(&bid(b"a")));
        assert!(!p.may_access_claim(&bid(b"a"), &kilio_seal::ClaimId::derive(b"c")));
    }

    #[test]
    fn scoped_key_sanitizes_traversal() {
        let a = bid(b"a");
        let k = branch_scoped_key(&a, "../../etc/passwd");
        assert!(!k.contains(".."));
        assert!(!k.contains('/') || k.matches('/').count() == 1); // only the branch/name separator
        assert!(k.starts_with(&a.to_hex()));
    }
}
