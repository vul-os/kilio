// Small date-formatting helpers for the reporter surface.

export function formatDateTime(iso) {
  const d = new Date(iso)
  return new Intl.DateTimeFormat(undefined, {
    weekday: 'short', month: 'short', day: 'numeric',
    hour: 'numeric', minute: '2-digit',
  }).format(d)
}

export function formatDay(iso) {
  const d = new Date(iso)
  return new Intl.DateTimeFormat(undefined, {
    month: 'short', day: 'numeric', year: 'numeric',
  }).format(d)
}
