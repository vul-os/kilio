import { Route, Routes } from 'react-router-dom'
import './handler.css'
import Shell from './Shell.tsx'
import Inbox from './Inbox.tsx'
import CaseDetail from './CaseDetail.tsx'
import Branches from './Branches.tsx'
import Settings from './Settings.tsx'

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
