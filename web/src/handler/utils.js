// Small helpers for the handler surface. No backend — everything here formats
// or derives from the mock data already loaded.

/** "3h ago" / "2d ago" style relative time, falling back to a date for
 * anything older than a week (keeps distant mock dates legible). */
export function timeAgo(iso) {
  const then = new Date(iso).getTime()
  const now = Date.now()
  const diff = Math.max(0, now - then)
  const min = Math.round(diff / 60000)
  if (min < 1) return 'just now'
  if (min < 60) return `${min}m ago`
  const hr = Math.round(min / 60)
  if (hr < 24) return `${hr}h ago`
  const day = Math.round(hr / 24)
  if (day < 7) return `${day}d ago`
  return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

/** Full readable timestamp for thread bubbles / audit entries. */
export function fullTime(iso) {
  return new Date(iso).toLocaleString(undefined, {
    weekday: 'short', month: 'short', day: 'numeric',
    hour: 'numeric', minute: '2-digit',
  })
}

/** Mask a receipt id down to its last group, e.g. "7Q4K-2M9F" -> "••••-2M9F". */
export function maskReceipt(id) {
  const parts = String(id).split('-')
  if (parts.length < 2) return '•'.repeat(Math.max(4, id.length - 2)) + id.slice(-2)
  return parts.slice(0, -1).map((p) => '•'.repeat(p.length)).join('-') + '-' + parts[parts.length - 1]
}

/** Deterministic pseudo-fingerprint for a branch id, styled like a key
 * fingerprint. Not real crypto — display only, mock demo data. */
export function keyFingerprint(seed) {
  let h1 = 0x811c9dc5, h2 = 0x1000193
  for (let i = 0; i < seed.length; i++) {
    const c = seed.charCodeAt(i)
    h1 = (h1 ^ c) * 16777619 >>> 0
    h2 = (h2 + c * 2654435761) >>> 0
  }
  const hex = (n) => n.toString(16).padStart(8, '0')
  const raw = (hex(h1) + hex(h2) + hex(h1 ^ h2) + hex((h1 + h2) >>> 0)).slice(0, 32)
  return raw.toUpperCase().match(/.{1,4}/g).join(' ')
}

export function initialsFrom(text) {
  return text.split(/\s+/).filter(Boolean).slice(0, 2).map((w) => w[0]).join('').toUpperCase()
}
