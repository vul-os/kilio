import { useMemo, useState } from 'react'
import { Link, useOutletContext, useSearchParams } from 'react-router-dom'
import { Pill } from '../ui/index.jsx'
import { CLAIMS, STATUS, branchName } from '../mock/data.js'
import { timeAgo } from './utils.js'
import { LockIcon, LockOpenIcon } from './icons.jsx'

const STATUS_ORDER = ['new', 'triaged', 'in_progress', 'resolved', 'closed']

export default function Inbox() {
  const { search } = useOutletContext()
  const [params] = useSearchParams()
  const branchId = params.get('branch') || ''
  const [statusFilter, setStatusFilter] = useState('all')

  const branchScoped = useMemo(
    () => (branchId ? CLAIMS.filter((c) => c.branchId === branchId) : CLAIMS),
    [branchId]
  )

  const counts = useMemo(() => {
    const m = { all: branchScoped.length }
    for (const key of STATUS_ORDER) m[key] = branchScoped.filter((c) => c.status === key).length
    return m
  }, [branchScoped])

  const rows = useMemo(() => {
    const q = search.trim().toLowerCase()
    return branchScoped
      .filter((c) => statusFilter === 'all' || c.status === statusFilter)
      .filter((c) => {
        if (!q) return true
        return (
          c.id.toLowerCase().includes(q) ||
          c.title.toLowerCase().includes(q) ||
          c.category.toLowerCase().includes(q)
        )
      })
      .slice()
      .sort((a, b) => new Date(b.updatedAt) - new Date(a.updatedAt))
  }, [branchScoped, statusFilter, search])

  return (
    <div className="h-page h-inbox">
      <div className="h-page-head">
        <div>
          <h1>Inbox</h1>
          <p className="h-page-sub">
            {branchId ? `Filtered to ${branchName(branchId)}` : 'All branches'} · {rows.length} of {branchScoped.length} cases
          </p>
        </div>
      </div>

      <div className="h-summary" role="tablist" aria-label="Filter by status">
        <button
          type="button"
          role="tab"
          aria-selected={statusFilter === 'all'}
          className={`h-summary-card ${statusFilter === 'all' ? 'active' : ''}`}
          onClick={() => setStatusFilter('all')}
        >
          <span className="h-summary-n">{counts.all}</span>
          <span className="h-summary-label">All cases</span>
        </button>
        {STATUS_ORDER.map((key) => (
          <button
            key={key}
            type="button"
            role="tab"
            aria-selected={statusFilter === key}
            className={`h-summary-card tone-${STATUS[key].tone} ${statusFilter === key ? 'active' : ''}`}
            onClick={() => setStatusFilter(key)}
          >
            <span className="h-summary-n">{counts[key]}</span>
            <span className="h-summary-label">{STATUS[key].label}</span>
          </button>
        ))}
      </div>

      <div className="h-case-list card" role="table" aria-label="Cases">
        <div className="h-case-row h-case-row-head" role="row">
          <span role="columnheader">Receipt</span>
          <span role="columnheader">Branch &amp; category</span>
          <span role="columnheader">Status</span>
          <span role="columnheader">Last activity</span>
        </div>

        {rows.length === 0 && (
          <div className="h-empty">No cases match this view.</div>
        )}

        {rows.map((c) => (
          <Link key={c.id} to={`/handler/case/${c.id}`} className={`h-case-row ${c.unread ? 'unread' : ''}`} role="row">
            <span className="h-case-receipt" role="cell">
              {c.unread ? <span className="h-unread-dot" aria-label="Unread" /> : <span className="h-unread-dot spacer" aria-hidden="true" />}
              {c.unread ? <LockIcon className="h-lock" /> : <LockOpenIcon className="h-lock opened" />}
              <span className="mono">{c.id}</span>
            </span>
            <span className="h-case-meta" role="cell">
              <span className="h-case-title">{c.title}</span>
              <span className="h-case-sub">{branchName(c.branchId)} · {c.category}</span>
            </span>
            <span role="cell">
              <Pill tone={STATUS[c.status].tone}>{STATUS[c.status].label}</Pill>
            </span>
            <span className="h-case-time" role="cell">{timeAgo(c.updatedAt)}</span>
          </Link>
        ))}
      </div>
    </div>
  )
}
