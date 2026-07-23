import { useState } from 'react'
import { useOutletContext } from 'react-router-dom'
import { Button, Field, Seal } from '../ui/index.jsx'
import { CopyIcon, GlobeIcon } from './icons.jsx'

const PROVIDERS = [
  { id: 'local', name: 'Local only', desc: 'Reachable only on this device / network. Most private, no public URL.' },
  { id: 'cloudflared', name: 'Cloudflared', desc: 'Cloudflare Tunnel gives a public URL without opening a port.' },
  { id: 'ngrok', name: 'ngrok', desc: 'ngrok tunnel — good for quick demos, URL changes on restart.' },
]

const MOCK_URLS = {
  cloudflared: 'https://sanctuary-4f2a.trycloudflare.com',
  ngrok: 'https://a1b2-203-0-113-9.ngrok-free.app',
}

export default function Settings() {
  const { reach, setReach, deployMode, setDeployMode } = useOutletContext()
  const [copied, setCopied] = useState(false)

  function setProvider(id) {
    setReach({ provider: id, url: id === 'local' ? '' : (MOCK_URLS[id] || reach.url) })
  }

  function copyUrl() {
    setCopied(true)
    setTimeout(() => setCopied(false), 1400)
  }

  return (
    <div className="h-page h-settings">
      <div className="h-page-head">
        <div>
          <h1>Settings</h1>
          <p className="h-page-sub">Reachability, deploy mode, and about this handler.</p>
        </div>
      </div>

      <section className="settings-section card">
        <div className="settings-section-head">
          <h2>Reachability</h2>
          <p>How reporters and this handler reach each other from outside this network.</p>
        </div>

        <div className="provider-grid" role="radiogroup" aria-label="Reachability provider">
          {PROVIDERS.map((p) => (
            <button
              key={p.id}
              type="button"
              role="radio"
              aria-checked={reach.provider === p.id}
              className={`provider-card ${reach.provider === p.id ? 'active' : ''}`}
              onClick={() => setProvider(p.id)}
            >
              <span className="provider-card-name">{p.name}</span>
              <span className="provider-card-desc">{p.desc}</span>
            </button>
          ))}
        </div>

        {reach.provider !== 'local' && (
          <Field label="Public URL" hint="Share this only with people who should be able to submit reports.">
            <div className="url-row">
              <GlobeIcon />
              <input className="input mono" readOnly value={reach.url} />
              <Button variant="ghost" type="button" onClick={copyUrl}>
                <CopyIcon /> {copied ? 'Copied' : 'Copy'}
              </Button>
            </div>
          </Field>
        )}
      </section>

      <section className="settings-section card">
        <div className="settings-section-head">
          <h2>Deploy mode</h2>
          <p>Run standalone, or hand reachability and routing to a vulos OS gateway.</p>
        </div>
        <div className="segmented" role="radiogroup" aria-label="Deploy mode">
          <button
            type="button"
            role="radio"
            aria-checked={deployMode === 'standalone'}
            className={deployMode === 'standalone' ? 'active' : ''}
            onClick={() => setDeployMode('standalone')}
          >
            Standalone
          </button>
          <button
            type="button"
            role="radio"
            aria-checked={deployMode === 'os_gateway'}
            className={deployMode === 'os_gateway' ? 'active' : ''}
            onClick={() => setDeployMode('os_gateway')}
          >
            OS gateway
          </button>
        </div>
        <p className="settings-hint">
          {deployMode === 'standalone'
            ? 'kilio manages its own reachability provider and keeps everything local to this machine.'
            : 'The vulos OS gateway handles TLS, routing, and public reachability for kilio.'}
        </p>
      </section>

      <section className="settings-section card about-section">
        <div className="about-head">
          <Seal size={40} />
          <div>
            <h2>kilio</h2>
            <p className="about-tagline">A safe place to speak — sealed, anonymous-first reporting.</p>
          </div>
        </div>
        <dl className="about-grid">
          <div><dt>Version</dt><dd className="mono">0.4.2</dd></div>
          <div><dt>Build</dt><dd className="mono">8f2a91c</dd></div>
          <div><dt>Sealing</dt><dd>kilio-seal · offline-capable</dd></div>
          <div><dt>License</dt><dd>AGPL-3.0</dd></div>
        </dl>
        <div className="about-links">
          <Button variant="ghost" size="sm">Documentation</Button>
          <Button variant="ghost" size="sm">Report an issue</Button>
          <Button variant="ghost" size="sm">Check for updates</Button>
        </div>
      </section>
    </div>
  )
}
