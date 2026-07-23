import { useState } from 'react'
import { Button, SealBadge } from '../ui/index.jsx'
import { CLAIMS, STATUS, branchName } from '../mock/data.js'
import { formatDateTime } from './format.js'
import { IconLock, IconSend } from './icons.jsx'

export default function Thread() {
  const claim = CLAIMS[0]
  const status = STATUS[claim.status]
  const [draft, setDraft] = useState('')

  return (
    <div className="rp-thread">
      <div className="rp-thread-top rp-in">
        <div>
          <h1 className="rp-thread-title">{claim.title}</h1>
          <div className="rp-thread-meta">
            <span className="rp-thread-id">Receipt {claim.id}</span>
            <span className="rp-trust-sep">·</span>
            <span className="rp-thread-id">{branchName(claim.branchId)}</span>
          </div>
        </div>
        <SealBadge>Sealed to {branchName(claim.branchId)}</SealBadge>
      </div>

      <div className={`rp-banner rp-banner-tone-${status.tone} rp-in`} style={{ '--i': 1 }}>
        <span className="rp-banner-dot" />
        <span className="rp-banner-text">
          <strong>{status.label}.</strong> The team has your report and can reply
          here. You'll see their messages the next time you return with your code.
        </span>
      </div>

      <div className="rp-messages">
        {claim.thread.map((m, i) => (
          <div key={i} className={`rp-msg rp-msg-${m.dir} rp-in`} style={{ '--i': i + 2 }}>
            <div className="rp-msg-bubble">{m.body}</div>
            <div className="rp-msg-meta">
              <IconLock width={12} height={12} />
              {m.dir === 'handler' ? branchName(claim.branchId) : 'You'} · {formatDateTime(m.at)}
            </div>
          </div>
        ))}
      </div>

      <div className="card rp-composer rp-in" style={{ '--i': 8 }}>
        <textarea className="textarea" value={draft} onChange={(e) => setDraft(e.target.value)}
          placeholder="Add anything else you'd like the team to know…" />
        <div className="rp-composer-row">
          <span className="rp-composer-hint"><IconLock width={14} height={14} /> Sealed before it leaves your device</span>
          <Button disabled={!draft.trim()}>
            <IconSend width={16} height={16} /> Send sealed reply
          </Button>
        </div>
      </div>
    </div>
  )
}
