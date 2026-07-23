import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button, Field, Stepper } from '../ui/index.jsx'
import { BRANCHES, CATEGORIES, branchName } from '../mock/data.js'
import {
  IconArrowLeft, IconArrowRight, IconPaperclip, IconLock, IconCheck, IconShield,
} from './icons.jsx'

const STEPS = ['Team', 'Category', 'Your account', 'Attachments', 'Contact', 'Review']

const REASSURE = [
  'Choose who receives this. It stays sealed until it reaches them.',
  'A rough category helps the right people respond — nothing here is binding.',
  'Write as much or as little as you like. This is sealed on your device first.',
  'Optional. Anything you attach is sealed too, before it leaves your device.',
  'Completely optional. Leave this blank to stay fully anonymous.',
  'Nothing has been sent yet. Review, then seal and send.',
]

export default function NewReport() {
  const nav = useNavigate()
  const [step, setStep] = useState(0)
  const [branch, setBranch] = useState('')
  const [category, setCategory] = useState('')
  const [title, setTitle] = useState('')
  const [body, setBody] = useState('')
  const [files, setFiles] = useState([])
  const [contact, setContact] = useState('')

  const canNext =
    (step === 0 && branch) ||
    (step === 1 && category) ||
    (step === 2 && title.trim() && body.trim()) ||
    step === 3 || step === 4 || step === 5

  const next = () => (step < STEPS.length - 1 ? setStep(step + 1) : nav('/report/receipt'))
  const back = () => (step > 0 ? setStep(step - 1) : nav('/report'))

  const addMockFile = () =>
    setFiles((f) => [...f, { name: `evidence-${f.length + 1}.pdf`, size: '248 KB' }])

  return (
    <div className="rp-flow">
      <div className="rp-flow-top">
        <span className="rp-flow-eyebrow">New report · {STEPS[step]}</span>
        <Stepper count={STEPS.length} active={step} />
      </div>

      <div className="card rp-flow-card rp-in">
        <h2 className="rp-step-title">{stepTitle(step)}</h2>
        <p className="rp-step-reassure"><IconLock width={15} height={15} /> {REASSURE[step]}</p>

        {step === 0 && (
          <div className="rp-branches">
            {BRANCHES.map((b) => (
              <div className="rp-branch" key={b.id}>
                <input type="radio" id={b.id} name="branch" checked={branch === b.id}
                  onChange={() => setBranch(b.id)} />
                <label htmlFor={b.id}>
                  <div className="rp-branch-row">
                    <span className="rp-branch-name">{b.name}</span>
                    <span className="rp-branch-check"><IconCheck width={13} height={13} /></span>
                  </div>
                  <span className="rp-branch-blurb">{b.blurb}</span>
                </label>
              </div>
            ))}
          </div>
        )}

        {step === 1 && (
          <div className="rp-chips">
            {CATEGORIES.map((c) => (
              <div className="rp-chip" key={c}>
                <input type="radio" id={`cat-${c}`} name="cat" checked={category === c}
                  onChange={() => setCategory(c)} />
                <label htmlFor={`cat-${c}`}>{c}</label>
              </div>
            ))}
          </div>
        )}

        {step === 2 && (
          <div className="rp-fields">
            <Field label="A short title">
              <input className="input" value={title} maxLength={120}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="e.g. Repeated comments from a team lead" />
            </Field>
            <Field label="What happened?" hint="Include dates, places, and anyone involved if you can — only if you're comfortable.">
              <textarea className="textarea rp-textarea-lg" value={body}
                onChange={(e) => setBody(e.target.value)}
                placeholder="Tell it in your own words. There's no wrong way to do this." />
              <span className="rp-charcount">{body.length} characters</span>
            </Field>
          </div>
        )}

        {step === 3 && (
          <div>
            <div className="rp-dropzone" role="button" tabIndex={0} onClick={addMockFile}
              onKeyDown={(e) => e.key === 'Enter' && addMockFile()}>
              <IconPaperclip width={22} height={22} />
              <div><strong>Add files</strong> — drag them here or click to browse</div>
              <div className="rp-dropzone-hint">Documents, photos, screenshots. Sealed before they leave your device.</div>
            </div>
            {files.length > 0 && (
              <div className="rp-files">
                {files.map((f, i) => (
                  <div className="rp-file" key={i}>
                    <IconPaperclip width={15} height={15} />
                    <span className="rp-file-name">{f.name}</span>
                    <span className="rp-file-size">{f.size}</span>
                    <button className="rp-file-remove" aria-label="Remove"
                      onClick={() => setFiles(files.filter((_, j) => j !== i))}>×</button>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {step === 4 && (
          <div className="rp-fields">
            <Field label="A way to reach you (optional)"
              hint="A phone, email, or handle — only if you want us to be able to follow up directly. You can always reply here with your receipt code instead.">
              <input className="input" value={contact} onChange={(e) => setContact(e.target.value)}
                placeholder="Leave blank to stay fully anonymous" />
            </Field>
          </div>
        )}

        {step === 5 && (
          <>
            <div className="rp-review">
              <Row label="Team" value={branch ? branchName(branch) : ''} />
              <Row label="Category" value={category} />
              <Row label="Title" value={title} />
              <Row label="Account" value={body} body />
              <Row label="Attachments" value={files.length ? `${files.length} file(s)` : ''} />
              <Row label="Contact" value={contact || ''} />
            </div>
            <div className="rp-submit-note">
              <IconShield width={18} height={18} />
              <span>When you send, your report is <strong>sealed on this device</strong> to
                {branch ? ` ${branchName(branch)}` : ' the chosen team'}. We can't read it in
                transit or at rest — and you'll get a receipt code that's the only way back.</span>
            </div>
          </>
        )}

        <div className="rp-flow-actions">
          <Button variant="ghost" onClick={back}>
            <IconArrowLeft width={16} height={16} /> {step === 0 ? 'Cancel' : 'Back'}
          </Button>
          <div className="rp-flow-actions-right">
            <Button onClick={next} disabled={!canNext}>
              {step === STEPS.length - 1 ? <>Seal &amp; send <IconLock width={16} height={16} /></>
                : <>Continue <IconArrowRight width={16} height={16} /></>}
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}

function stepTitle(step) {
  return [
    'Who should receive this?',
    'What kind of concern is it?',
    'Tell us what happened',
    'Add anything that helps',
    'Do you want us to reach you?',
    'Ready to send',
  ][step]
}

function Row({ label, value, body }) {
  return (
    <div className="rp-review-row">
      <div className="rp-review-label">{label}</div>
      <div className={`rp-review-value ${body ? 'rp-review-body' : ''}`}>
        {value ? value : <span className="rp-review-empty">Not provided</span>}
      </div>
    </div>
  )
}
