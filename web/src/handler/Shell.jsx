import { useMemo, useState } from 'react'
import { NavLink, Outlet, useNavigate, useSearchParams, useLocation } from 'react-router-dom'
import { Logo, Pill, ThemeToggle } from '../ui/index.jsx'
import { BRANCHES, CLAIMS } from '../mock/data.js'
import { InboxIcon, BranchIcon, GearIcon, SearchIcon } from './icons.jsx'

const REACH_LABEL = {
  local: { label: 'Local only', tone: 'muted' },
  cloudflared: { label: 'Public · tunnel on', tone: 'ok' },
  ngrok: { label: 'Public · ngrok on', tone: 'ok' },
}

/** App shell: left nav + branch filters, top bar with search and reachability.
 * Holds the small bits of cross-page state (search text, reachability,
 * deploy mode) and hands them to child routes via <Outlet context>. */
export default function Shell() {
  const navigate = useNavigate()
  const location = useLocation()
  const [params] = useSearchParams()
  const [search, setSearch] = useState('')
  const [reach, setReach] = useState({ provider: 'cloudflared', url: 'https://sanctuary-4f2a.trycloudflare.com' })
  const [deployMode, setDeployMode] = useState('standalone')

  const activeBranch = params.get('branch') || ''
  const onInbox = location.pathname === '/handler' || location.pathname === '/handler/'

  const branchCounts = useMemo(() => {
    const m = new Map()
    for (const c of CLAIMS) {
      if (c.status === 'resolved' || c.status === 'closed') continue
      m.set(c.branchId, (m.get(c.branchId) || 0) + 1)
    }
    return m
  }, [])

  const totalOpen = useMemo(
    () => CLAIMS.filter((c) => c.status !== 'resolved' && c.status !== 'closed').length,
    []
  )

  function goToBranch(id) {
    navigate(id ? `/handler?branch=${id}` : '/handler')
  }

  const reachInfo = REACH_LABEL[reach.provider]

  return (
    <div className="h-shell">
      <aside className="h-nav" aria-label="Handler navigation">
        <div className="h-nav-top">
          <Logo size={28} />
          <span className="h-nav-kicker">Handler</span>
        </div>

        <nav className="h-nav-list" aria-label="Sections">
          <NavLink to="/handler" end className={({ isActive }) => `h-nav-link ${isActive ? 'active' : ''}`}>
            <InboxIcon /> Inbox
            {totalOpen > 0 && <span className="h-nav-count">{totalOpen}</span>}
          </NavLink>
          <NavLink to="/handler/branches" className={({ isActive }) => `h-nav-link ${isActive ? 'active' : ''}`}>
            <BranchIcon /> Branches
          </NavLink>
          <NavLink to="/handler/settings" className={({ isActive }) => `h-nav-link ${isActive ? 'active' : ''}`}>
            <GearIcon /> Settings
          </NavLink>
        </nav>

        <div className="h-nav-section">
          <span className="h-nav-heading">Filter by branch</span>
          <ul className="h-branch-filter" role="list">
            <li>
              <button
                type="button"
                className={`h-branch-item ${onInbox && !activeBranch ? 'active' : ''}`}
                onClick={() => goToBranch('')}
              >
                <span className="h-branch-dot all" aria-hidden="true" />
                All branches
              </button>
            </li>
            {BRANCHES.map((b) => (
              <li key={b.id}>
                <button
                  type="button"
                  className={`h-branch-item ${onInbox && activeBranch === b.id ? 'active' : ''}`}
                  onClick={() => goToBranch(b.id)}
                  title={b.blurb}
                >
                  <span className="h-branch-dot" aria-hidden="true" />
                  <span className="h-branch-name">{b.name}</span>
                  {branchCounts.get(b.id) > 0 && <span className="h-nav-count">{branchCounts.get(b.id)}</span>}
                </button>
              </li>
            ))}
          </ul>
        </div>
      </aside>

      <header className="h-topbar">
        <label className="h-search">
          <SearchIcon />
          <input
            type="search"
            placeholder="Search receipt id, title, or category…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            aria-label="Search cases"
          />
        </label>
        <div className="h-topbar-right">
          <button type="button" className="h-reach-chip" onClick={() => navigate('/handler/settings')} title="Reachability — open in Settings">
            <Pill tone={reachInfo.tone}>{reachInfo.label}</Pill>
          </button>
          <ThemeToggle />
        </div>
      </header>

      <main className="h-main">
        <Outlet context={{ search, reach, setReach, deployMode, setDeployMode }} />
      </main>
    </div>
  )
}
