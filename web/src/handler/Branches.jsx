import { useState } from 'react'
import { Button, Field } from '../ui/index.jsx'
import { BRANCHES } from '../mock/data.js'
import { keyFingerprint } from './utils.js'
import { KeyIcon, PlusIcon } from './icons.jsx'

export default function Branches() {
  const [branches, setBranches] = useState(BRANCHES.map((b) => ({ ...b, active: true })))
  const [adding, setAdding] = useState(false)
  const [draft, setDraft] = useState({ name: '', blurb: '', powBits: 20 })

  function toggleActive(id) {
    setBranches((list) => list.map((b) => (b.id === id ? { ...b, active: !b.active } : b)))
  }

  function addBranch(e) {
    e.preventDefault()
    const name = draft.name.trim()
    if (!name) return
    const id = 'b_' + name.toLowerCase().replace(/[^a-z0-9]+/g, '_').slice(0, 24) + '_' + branches.length
    setBranches((list) => [...list, { id, name, blurb: draft.blurb.trim() || 'No description yet.', powBits: Number(draft.powBits) || 20, active: true }])
    setDraft({ name: '', blurb: '', powBits: 20 })
    setAdding(false)
  }

  return (
    <div className="h-page h-branches">
      <div className="h-page-head">
        <div>
          <h1>Branches &amp; keys</h1>
          <p className="h-page-sub">Each branch has its own sealing key and proof-of-work floor for intake.</p>
        </div>
        <Button variant="ghost" onClick={() => setAdding((a) => !a)}>
          <PlusIcon /> Add branch
        </Button>
      </div>

      {adding && (
        <form className="card branch-add-form" onSubmit={addBranch}>
          <div className="branch-add-grid">
            <Field label="Name">
              <input
                className="input"
                value={draft.name}
                onChange={(e) => setDraft((d) => ({ ...d, name: e.target.value }))}
                placeholder="e.g. Data & Privacy"
                autoFocus
              />
            </Field>
            <Field label="PoW difficulty" hint="Higher slows spam, but also slows honest senders.">
              <select
                className="select"
                value={draft.powBits}
                onChange={(e) => setDraft((d) => ({ ...d, powBits: e.target.value }))}
              >
                {[16, 18, 20, 22, 24].map((n) => <option key={n} value={n}>{n} bits</option>)}
              </select>
            </Field>
          </div>
          <Field label="Blurb" hint="Shown to reporters when choosing where to send a report.">
            <input
              className="input"
              value={draft.blurb}
              onChange={(e) => setDraft((d) => ({ ...d, blurb: e.target.value }))}
              placeholder="What belongs in this branch?"
            />
          </Field>
          <div className="branch-add-actions">
            <Button variant="ghost" type="button" onClick={() => setAdding(false)}>Cancel</Button>
            <Button type="submit">Create branch</Button>
          </div>
        </form>
      )}

      <div className="branch-grid">
        {branches.map((b) => (
          <article key={b.id} className={`card branch-card ${b.active ? '' : 'inactive'}`}>
            <div className="branch-card-top">
              <h3>{b.name}</h3>
              <label className="switch" title={b.active ? 'Active — accepting reports' : 'Inactive — intake paused'}>
                <input type="checkbox" checked={b.active} onChange={() => toggleActive(b.id)} />
                <span className="switch-track"><span className="switch-thumb" /></span>
                <span className="switch-label">{b.active ? 'Active' : 'Paused'}</span>
              </label>
            </div>
            <p className="branch-blurb">{b.blurb}</p>
            <div className="branch-stats">
              <span className="branch-stat">
                <span className="branch-stat-label">PoW floor</span>
                <span className="branch-stat-value">{b.powBits} bits</span>
              </span>
            </div>
            <div className="branch-fingerprint">
              <KeyIcon />
              <div>
                <span className="branch-fp-label">Sealing key fingerprint</span>
                <span className="mono branch-fp-value">{keyFingerprint(b.id)}</span>
              </div>
            </div>
          </article>
        ))}
      </div>
    </div>
  )
}
