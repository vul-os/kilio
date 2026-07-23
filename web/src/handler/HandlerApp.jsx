import { Route, Routes } from 'react-router-dom'
import './handler.css'
import Shell from './Shell.jsx'
import Inbox from './Inbox.jsx'
import CaseDetail from './CaseDetail.jsx'
import Branches from './Branches.jsx'
import Settings from './Settings.jsx'

// Handler surface — the case-worker desktop app. Nested under /handler.
export default function HandlerApp() {
  return (
    <Routes>
      <Route element={<Shell />}>
        <Route index element={<Inbox />} />
        <Route path="case/:id" element={<CaseDetail />} />
        <Route path="branches" element={<Branches />} />
        <Route path="settings" element={<Settings />} />
      </Route>
    </Routes>
  )
}
