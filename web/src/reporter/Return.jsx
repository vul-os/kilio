import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button, Seal } from '../ui/index.jsx'
import { IconArrowRight, IconKey } from './icons.jsx'

export default function Return() {
  const nav = useNavigate()
  const [code, setCode] = useState('')
  const words = code.trim().split(/\s+/).filter(Boolean)
  const ready = words.length === 12

  const open = (e) => {
    e.preventDefault()
    if (ready) nav('/report/thread')
  }

  return (
    <div className="rp-center">
      <div className="rp-receipt rp-in">
        <form className="card rp-return-card" onSubmit={open}>
          <div className="rp-return-head">
            <Seal size={52} />
            <h1 className="rp-return-title">Return to your report</h1>
            <p className="rp-return-sub">
              Enter the 12-word receipt code you saved. Your report is re-opened
              locally — the code is never sent to us.
            </p>
          </div>

          <textarea
            className="input rp-return-input"
            rows={3}
            value={code}
            onChange={(e) => setCode(e.target.value)}
            placeholder="harbor  willow  lantern  …"
            aria-label="12-word receipt code"
            autoFocus
          />
          <div className={`rp-return-count ${ready ? 'is-ready' : ''}`}>
            {words.length} / 12 words
          </div>

          <div style={{ marginTop: 18, display: 'flex', justifyContent: 'center' }}>
            <Button type="submit" disabled={!ready}>
              <IconKey width={16} height={16} /> Open my report <IconArrowRight width={16} height={16} />
            </Button>
          </div>

          <div className="rp-return-footer">
            Lost your code? For your safety it can't be recovered — but you can
            always <a href="/report/new">make a new report</a>.
          </div>
        </form>
      </div>
    </div>
  )
}
