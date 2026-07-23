import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Button, Seal } from '../ui/index.jsx'
import { makeReceipt } from '../mock/data.js'
import { IconCopy, IconDownload, IconCheck, IconArrowRight } from './icons.jsx'

export default function Receipt() {
  const words = useMemo(() => makeReceipt(), [])
  const phrase = words.join(' ')
  const [copied, setCopied] = useState(false)

  const copy = async () => {
    try { await navigator.clipboard.writeText(phrase) } catch (e) {}
    setCopied(true)
    setTimeout(() => setCopied(false), 2200)
  }

  const download = () => {
    const blob = new Blob([`kilio receipt code\n\n${phrase}\n\nKeep this safe — it is the only way back to your report.`], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url; a.download = 'kilio-receipt.txt'; a.click()
    URL.revokeObjectURL(url)
  }

  return (
    <div className="rp-center">
      <div className="rp-receipt rp-in">
        <div className="card rp-receipt-card">
          <div className="rp-receipt-head">
            <Seal size={56} />
            <h1 className="rp-receipt-title">Your report is sealed and sent</h1>
            <p className="rp-receipt-sub">
              Here is your receipt code. It's the <strong>only</strong> way to return,
              read replies, and add more later — while staying anonymous.
            </p>
          </div>

          <ol className="rp-receipt-words">
            {words.map((w, i) => <li key={i}>{w}</li>)}
          </ol>

          <div className="rp-receipt-warning">
            <IconCheck width={18} height={18} style={{ color: 'var(--warn)' }} />
            <span><strong>Save these 12 words now.</strong> We never store them and
              cannot recover them. Anyone with them can open your report — and without
              them, no one (including us) can link it back to you.</span>
          </div>

          <div className="rp-receipt-actions">
            <Button variant="ghost" onClick={copy}>
              <IconCopy width={16} height={16} /> Copy code
            </Button>
            <Button variant="ghost" onClick={download}>
              <IconDownload width={16} height={16} /> Download
            </Button>
          </div>
          {copied && <div className="rp-copied"><IconCheck width={14} height={14} /> Copied to clipboard</div>}

          <div className="rp-receipt-continue">
            <Button as={Link} to="/report/thread">
              Go to my report <IconArrowRight width={16} height={16} />
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
