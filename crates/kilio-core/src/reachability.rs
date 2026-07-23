//! The `Reachability` seam — make the local intake page publicly reachable.
//!
//! Mirrors wede's `Provider` interface (`start`/`stop`/`public_url`), mechanism
//! -agnostic. The default binds loopback only. The subprocess and sovereign
//! -relay providers are wired where process/network access lives (server/CLI);
//! here we define the contract and the loopback default.
//!
//! **SSRF invariant (from wede):** whatever provider runs, it proxies to
//! *exactly one* configured loopback address, re-checked before every
//! connection. The inbound request's Host/URL never chooses the target.

use std::net::SocketAddr;

use crate::CoreError;

/// A tunnel status snapshot. A token, if any, is never included.
#[derive(Clone, Debug, Default)]
pub struct TunnelStatus {
    pub running: bool,
    pub public_url: Option<String>,
    pub provider: &'static str,
}

pub trait Reachability: Send {
    /// Begin exposing `local_addr` (which must be loopback) and return the
    /// public URL once assigned.
    fn start(&mut self, local_addr: SocketAddr) -> Result<String, CoreError>;
    fn stop(&mut self) -> Result<(), CoreError>;
    fn snapshot(&self) -> TunnelStatus;
}

/// Default: no exposure. Binds loopback; use behind a reverse proxy you already
/// run, or a Tor hidden service.
#[derive(Default)]
pub struct LocalOnly {
    addr: Option<SocketAddr>,
}

impl Reachability for LocalOnly {
    fn start(&mut self, local_addr: SocketAddr) -> Result<String, CoreError> {
        if !is_loopback(&local_addr) {
            return Err(CoreError::NotLoopback);
        }
        self.addr = Some(local_addr);
        Ok(format!("http://{local_addr}"))
    }
    fn stop(&mut self) -> Result<(), CoreError> {
        self.addr = None;
        Ok(())
    }
    fn snapshot(&self) -> TunnelStatus {
        TunnelStatus {
            running: self.addr.is_some(),
            public_url: self.addr.map(|a| format!("http://{a}")),
            provider: "local-only",
        }
    }
}

/// Which subprocess tunnel to drive. Detected at runtime by binary name.
#[derive(Clone, Copy, PartialEq, Eq, Debug)]
pub enum TunnelProvider {
    Cloudflared,
    Ngrok,
    Frp,
}

impl TunnelProvider {
    pub fn binary(&self) -> &'static str {
        match self {
            TunnelProvider::Cloudflared => "cloudflared",
            TunnelProvider::Ngrok => "ngrok",
            TunnelProvider::Frp => "frpc",
        }
    }
}

/// Expose the intake page via a detected tunnel binary (cloudflared / ngrok /
/// frp), pinned to the loopback listen address. The honest "click to go public
/// with no fixed infra" default. Spawning + public-URL parsing lives in the
/// server/CLI layer; this struct records the choice and enforces the loopback
/// SSRF guard up front.
pub struct SubprocessTunnel {
    pub provider: TunnelProvider,
    pinned: Option<SocketAddr>,
}

impl SubprocessTunnel {
    pub fn new(provider: TunnelProvider) -> Self {
        Self {
            provider,
            pinned: None,
        }
    }
}

impl Reachability for SubprocessTunnel {
    fn start(&mut self, local_addr: SocketAddr) -> Result<String, CoreError> {
        if !is_loopback(&local_addr) {
            return Err(CoreError::NotLoopback); // SSRF guard: loopback target only
        }
        self.pinned = Some(local_addr);
        Err(CoreError::Unsupported(
            "SubprocessTunnel spawn/parse is wired in the server/CLI layer",
        ))
    }
    fn stop(&mut self) -> Result<(), CoreError> {
        self.pinned = None;
        Ok(())
    }
    fn snapshot(&self) -> TunnelStatus {
        TunnelStatus {
            running: false,
            public_url: None,
            provider: self.provider.binary(),
        }
    }
}

fn is_loopback(addr: &SocketAddr) -> bool {
    addr.ip().is_loopback()
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::net::{Ipv4Addr, SocketAddr};

    #[test]
    fn local_only_exposes_loopback_only() {
        let mut r = LocalOnly::default();
        let loop_addr = SocketAddr::from((Ipv4Addr::LOCALHOST, 8787));
        assert!(r.start(loop_addr).is_ok());
        assert!(r.snapshot().running);
        r.stop().unwrap();
        assert!(!r.snapshot().running);
    }

    #[test]
    fn ssrf_guard_rejects_non_loopback() {
        let mut r = LocalOnly::default();
        let public = SocketAddr::from((Ipv4Addr::new(0, 0, 0, 0), 8787));
        assert!(matches!(r.start(public), Err(CoreError::NotLoopback)));

        let mut t = SubprocessTunnel::new(TunnelProvider::Cloudflared);
        let public2 = SocketAddr::from((Ipv4Addr::new(93, 184, 216, 34), 80));
        assert!(matches!(t.start(public2), Err(CoreError::NotLoopback)));
    }
}
