import { Routes, Route, Link, Outlet } from 'react-router-dom'
import './reporter.css'
import { Logo, ThemeToggle } from '../ui/index.jsx'
import Landing from './Landing.jsx'
import NewReport from './NewReport.jsx'
import Receipt from './Receipt.jsx'
import Return from './Return.jsx'
import Thread from './Thread.jsx'

function Layout() {
  return (
    <div className="rp-app">
      <div className="rp-grain" aria-hidden="true" />
      <header className="rp-header">
        <Link to="/report" className="rp-logo-link" aria-label="kilio home"><Logo /></Link>
        <div className="rp-header-actions"><ThemeToggle /></div>
      </header>
      <main className="rp-main"><Outlet /></main>
    </div>
  )
}

export default function ReporterApp() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route index element={<Landing />} />
        <Route path="new" element={<NewReport />} />
        <Route path="receipt" element={<Receipt />} />
        <Route path="return" element={<Return />} />
        <Route path="thread" element={<Thread />} />
      </Route>
    </Routes>
  )
}
