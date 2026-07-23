import { Link } from 'react-router-dom'
import { Logo, Button, ThemeToggle, SealBadge } from './ui/index.jsx'

// Dev landing — links to both surfaces. Not a shipped screen; the reporter and
// handler apps are the real surfaces.
export default function App() {
  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <header style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '22px 32px' }}>
        <Logo />
        <ThemeToggle />
      </header>
      <main style={{ flex: 1, display: 'grid', placeItems: 'center', padding: '32px' }}>
        <div style={{ maxWidth: 640, textAlign: 'center' }}>
          <SealBadge>Sealed, anonymous-first intake</SealBadge>
          <h1 style={{ fontSize: 'clamp(2.2rem, 6vw, 3.4rem)', margin: '18px 0 12px' }}>
            A safe place to speak.
          </h1>
          <p style={{ color: 'var(--muted)', fontSize: '1.1rem', marginBottom: 30 }}>
            kilio lets people report sensitive things without an account or an email —
            sealed so even the host can't read them.
          </p>
          <div style={{ display: 'flex', gap: 12, justifyContent: 'center', flexWrap: 'wrap' }}>
            <Button as={Link} to="/report" size="lg">Reporter surface →</Button>
            <Button as={Link} to="/handler" variant="ghost" size="lg">Handler surface →</Button>
          </div>
        </div>
      </main>
    </div>
  )
}
