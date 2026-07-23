// Mock data so both surfaces run and screenshot without a backend.
// Real sealing/decryption lives in the kilio-seal crate; here the handler
// side shows already-"opened" content and the reporter side simulates submit.

export const BRANCHES = [
  { id: 'b_people', name: 'People & Culture', blurb: 'Harassment, discrimination, workplace conduct.', powBits: 20 },
  { id: 'b_ethics', name: 'Ethics & Compliance', blurb: 'Fraud, conflicts of interest, policy breaches.', powBits: 20 },
  { id: 'b_safety', name: 'Health & Safety', blurb: 'Unsafe conditions, incidents, near-misses.', powBits: 18 },
]

export const CATEGORIES = [
  'Harassment', 'Discrimination', 'Bullying', 'Safety concern',
  'Fraud or theft', 'Conflict of interest', 'Retaliation', 'Something else',
]

export const STATUS = {
  new:        { label: 'Received',     tone: 'indigo' },
  triaged:    { label: 'Under review', tone: 'warn' },
  in_progress:{ label: 'In progress',  tone: 'warn' },
  resolved:   { label: 'Resolved',     tone: 'ok' },
  closed:     { label: 'Closed',       tone: 'muted' },
}

// Short, non-graphic sample threads. Reporter is always anonymous.
export const CLAIMS = [
  {
    id: '7Q4K-2M9F', branchId: 'b_people', category: 'Harassment', status: 'triaged',
    title: 'Repeated comments from a team lead', createdAt: '2026-07-20T09:12:00Z',
    updatedAt: '2026-07-22T14:03:00Z', unread: true,
    thread: [
      { dir: 'reporter', at: '2026-07-20T09:12:00Z', body: 'A senior team lead keeps making comments about my appearance in front of others. I have asked them to stop twice. I do not want to give my name yet.' },
      { dir: 'handler',  at: '2026-07-20T15:40:00Z', body: 'Thank you for telling us — that took courage. You can stay anonymous for as long as you like. Could you share roughly when the most recent comment happened, and whether anyone else was present?' },
      { dir: 'reporter', at: '2026-07-22T14:03:00Z', body: 'It was this Monday in the afternoon stand-up. Two other people on my team were there.' },
    ],
    audit: [
      { at: '2026-07-20T09:12:00Z', text: 'Report received · sealed' },
      { at: '2026-07-20T15:39:00Z', text: 'Opened by you' },
      { at: '2026-07-20T15:40:00Z', text: 'Status → Under review' },
    ],
  },
  {
    id: 'B8XR-J1P0', branchId: 'b_ethics', category: 'Conflict of interest', status: 'in_progress',
    title: 'Vendor selection may be biased', createdAt: '2026-07-18T11:00:00Z',
    updatedAt: '2026-07-21T10:22:00Z', unread: false,
    thread: [
      { dir: 'reporter', at: '2026-07-18T11:00:00Z', body: 'A manager pushed hard for a vendor that a close family member works for. It was not disclosed in the review meeting.' },
      { dir: 'handler',  at: '2026-07-19T08:15:00Z', body: 'Understood, and thank you. We are reviewing the procurement records. Do you know the approximate date of the review meeting?' },
    ],
    audit: [
      { at: '2026-07-18T11:00:00Z', text: 'Report received · sealed' },
      { at: '2026-07-19T08:14:00Z', text: 'Opened by you' },
      { at: '2026-07-19T08:16:00Z', text: 'Status → In progress' },
    ],
  },
  {
    id: 'L3TT-9WQ2', branchId: 'b_safety', category: 'Safety concern', status: 'resolved',
    title: 'Blocked fire exit on level 2', createdAt: '2026-07-10T07:45:00Z',
    updatedAt: '2026-07-12T16:30:00Z', unread: false,
    thread: [
      { dir: 'reporter', at: '2026-07-10T07:45:00Z', body: 'The fire exit near the level 2 kitchen has been blocked by stacked boxes for over a week.' },
      { dir: 'handler',  at: '2026-07-10T09:02:00Z', body: 'Thank you — we are sending facilities to clear it today and will confirm here.' },
      { dir: 'handler',  at: '2026-07-12T16:30:00Z', body: 'The exit has been cleared and facilities added a monthly check. Marking this resolved. Thank you again for flagging it.' },
    ],
    audit: [
      { at: '2026-07-10T07:45:00Z', text: 'Report received · sealed' },
      { at: '2026-07-10T09:01:00Z', text: 'Opened by you' },
      { at: '2026-07-12T16:30:00Z', text: 'Status → Resolved' },
    ],
  },
  {
    id: 'K0M5-4RD7', branchId: 'b_people', category: 'Retaliation', status: 'new',
    title: 'Shift changes after raising a concern', createdAt: '2026-07-23T06:20:00Z',
    updatedAt: '2026-07-23T06:20:00Z', unread: true,
    thread: [
      { dir: 'reporter', at: '2026-07-23T06:20:00Z', body: 'After I raised a concern about overtime, my shifts were suddenly cut. I think it is connected.' },
    ],
    audit: [ { at: '2026-07-23T06:20:00Z', text: 'Report received · sealed' } ],
  },
]

// A small, friendly word pool for the demo receipt passphrase (BIP-39-like
// look; the real receipt is a 12-word BIP-39 phrase minted client-side).
const WORDS = ['harbor','willow','lantern','pebble','cedar','meadow','copper','anchor',
  'saffron','ripple','marble','thistle','ember','quartz','hollow','bramble',
  'orchard','velvet','pewter','cove','fable','indigo','maple','wren']

export function makeReceipt() {
  const out = []
  const pool = [...WORDS]
  for (let i = 0; i < 12; i++) {
    const n = (i * 7 + 3) % pool.length // deterministic for stable screenshots
    out.push(pool.splice(n % pool.length, 1)[0] || WORDS[i % WORDS.length])
  }
  return out
}

export function branchName(id) {
  return (BRANCHES.find((b) => b.id === id) || {}).name || 'General'
}
