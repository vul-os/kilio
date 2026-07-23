import { Link } from 'react-router-dom'
import { Button, Seal } from '../ui/index.jsx'
import { IconLock, IconShield, IconUser, IconArrowRight } from './icons.jsx'

export default function Landing() {
  return (
    <section className="rp-hero">
      <Seal size={210} className="rp-hero-seal" />
      <div className="rp-hero-inner">
        <h1 className="rp-hero-title rp-in" style={{ '--i': 0 }}>
          Speak safely.<br />You stay <em>anonymous</em>.
        </h1>
        <p className="rp-hero-lede rp-in" style={{ '--i': 1 }}>
          Report something that matters — without an account, without an email.
          Your words are sealed the moment you send them, and only the team you
          choose can ever open them.
        </p>
        <div className="rp-hero-actions rp-in" style={{ '--i': 2 }}>
          <Button as={Link} to="/report/new" size="lg">
            Make a report <IconArrowRight width={17} height={17} />
          </Button>
          <Button as={Link} to="/report/return" variant="ghost" size="lg">
            I have a receipt code
          </Button>
        </div>
        <div className="rp-trust rp-in" style={{ '--i': 3 }}>
          <span className="rp-trust-item"><IconLock width={15} height={15} /> No account or email</span>
          <span className="rp-trust-sep">·</span>
          <span className="rp-trust-item"><IconShield width={15} height={15} /> Sealed end-to-end</span>
          <span className="rp-trust-sep">·</span>
          <span className="rp-trust-item"><IconUser width={15} height={15} /> You choose if you say who you are</span>
        </div>
      </div>
    </section>
  )
}
