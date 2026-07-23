import { useMemo, useRef, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { Button, Pill, SealBadge } from '../ui/index.jsx'
import { CLAIMS, STATUS, branchName } from '../mock/data.js'
import { fullTime, maskReceipt, timeAgo } from './utils.js'
import { ArrowLeftIcon, ChevronDownIcon, LockIcon, SendIcon } from './icons.jsx'

const STATUS_ORDER = ['new', 'triaged', 'in_progress', 'resolved', 'closed']

export default function CaseDetail() {
  const { id } = useParams()
  const base = useMemo(() => CLAIMS.find((c) => c.id === id) || CLAIMS[0], [id])

  const [status, setStatus] = useState(base.status)
  const [thread, setThread] = useState(base.thread)
  const [audit, setAudit] = useState(base.audit)
  const [draft, setDraft] = useState('')
  const listRef = useRef(null)

  // Reset local session state when navigating to a different case.
  const caseKey = useRef(base.id)
  if (caseKey.current !== base.id) {
    caseKey.current = base.id
    setStatus(base.status)
    setThread(base.thread)
    setAudit(base.audit)
    setDraft('')
  }

  function changeStatus(next) {
    if (next === status) return
    setStatus(next)
    setAudit((a) => [...a, { at: new Date().toISOString(), text: `Status → ${STATUS[next].label}` }])
  }

  function sendReply(e) {
    e.preventDefault()
    const body = draft.trim()
    if (!body) return
    const now = new Date().toISOString()
    setThread((t) => [...t, { dir: 'handler', at: now, body }])
    setAudit((a) => [...a, { at: now, text: 'You replied to the reporter' }])
    setDraft('')
  }

  return (
    <div className="h-page case-page">
      <div className="case-main">
        <header className="case-header">
          <Link to="/handler" className="case-back">
            <ArrowLeftIcon /> Inbox
          </Link>

          <div className="case-header-row">
            <div className="case-header-id">
              <LockIcon className="h-lock opened" />
              <h1 className="mono">{base.id}</h1>
            </div>
            <StatusMenu value={status} onChange={changeStatus} />
          </div>

          <div className="case-header-tags">
            <span className="case-tag">{branchName(base.branchId)}</span>
            <span className="case-tag-sep">·</span>
            <span className="case-tag">{base.category}</span>
            <span className="case-tag-sep">·</span>
            <span className="case-tag muted">Opened {timeAgo(base.createdAt)}</span>
          </div>
          <p className="case-title-line">{base.title}</p>
        </header>

        <div className="case-thread" ref={listRef}>
          <div className="case-sysnote">
            <LockIcon />
            Report opened and decrypted locally on this device · {fullTime(base.createdAt)}
          </div>

          {thread.map((m, i) => (
            <div key={i} className={`bubble-row ${m.dir}`}>
              <div className="bubble">
                <div className="bubble-meta">
                  {m.dir === 'reporter' ? `Anonymous reporter · receipt ${maskReceipt(base.id)}` : 'You · Handler'}
                  <span className="bubble-time">{fullTime(m.at)}</span>
                </div>
                <p>{m.body}</p>
              </div>
            </div>
          ))}
        </div>

        <form className="case-composer" onSubmit={sendReply}>
          <SealBadge>Sent sealed, only the reporter can read it</SealBadge>
          <div className="case-composer-row">
            <textarea
              className="textarea"
              placeholder="Write a reply to the reporter…"
              value={draft}
              onChange={(e) => setDraft(e.target.value)}
              onKeyDown={(e) => {
                if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') sendReply(e)
              }}
              aria-label="Reply to reporter"
              rows={3}
            />
            <Button type="submit" disabled={!draft.trim()}>
              <SendIcon /> Send
            </Button>
          </div>
          <span className="case-composer-hint">⌘ / Ctrl + Enter to send</span>
        </form>
      </div>

      <aside className="case-rail" aria-label="Audit trail">
        <h2>Audit trail</h2>
        <p className="case-rail-note">
          kilio logs that actions happened — never their content. Nothing here reveals what was
          read or written.
        </p>
        <ol className="audit-list">
          {audit.slice().reverse().map((entry, i) => (
            <li key={i}>
              <span className="audit-dot" aria-hidden="true" />
              <div>
                <p className="audit-text">{entry.text}</p>
                <time className="audit-time">{fullTime(entry.at)}</time>
              </div>
            </li>
          ))}
        </ol>
      </aside>
    </div>
  )
}

function StatusMenu({ value, onChange }) {
  const [open, setOpen] = useState(false)
  const tone = STATUS[value].tone

  return (
    <div className="status-menu" data-tone={tone}>
      <button
        type="button"
        className="status-menu-btn"
        aria-haspopup="listbox"
        aria-expanded={open}
        onClick={() => setOpen((o) => !o)}
        onBlur={() => setTimeout(() => setOpen(false), 120)}
      >
        <Pill tone={tone}>{STATUS[value].label}</Pill>
        <ChevronDownIcon />
      </button>
      {open && (
        <ul className="status-menu-list" role="listbox" aria-label="Change status">
          {STATUS_ORDER.map((key) => (
            <li key={key}>
              <button
                type="button"
                role="option"
                aria-selected={key === value}
                className={key === value ? 'active' : ''}
                onMouseDown={(e) => { e.preventDefault(); onChange(key); setOpen(false) }}
              >
                <span className={`status-swatch tone-${STATUS[key].tone}`} aria-hidden="true" />
                {STATUS[key].label}
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
