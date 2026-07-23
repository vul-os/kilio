// Small inline icon set for the reporter surface. Stroke-based, currentColor,
// sized to sit inline with text/buttons. Kept local to this surface.

const base = {
  width: 18, height: 18, viewBox: '0 0 24 24', fill: 'none',
  stroke: 'currentColor', strokeWidth: 1.8, strokeLinecap: 'round', strokeLinejoin: 'round',
  'aria-hidden': 'true',
}

export function IconLock(props) {
  return (
    <svg {...base} {...props}>
      <rect x="4.5" y="10.5" width="15" height="10" rx="2.5" />
      <path d="M8 10.5V7.8a4 4 0 0 1 8 0v2.7" />
    </svg>
  )
}

export function IconArrowRight(props) {
  return (
    <svg {...base} {...props}>
      <path d="M5 12h14M13 6l6 6-6 6" />
    </svg>
  )
}

export function IconArrowLeft(props) {
  return (
    <svg {...base} {...props}>
      <path d="M19 12H5M11 6l-6 6 6 6" />
    </svg>
  )
}

export function IconCopy(props) {
  return (
    <svg {...base} {...props}>
      <rect x="9" y="9" width="11" height="11" rx="2" />
      <path d="M5 15V6a2 2 0 0 1 2-2h9" />
    </svg>
  )
}

export function IconDownload(props) {
  return (
    <svg {...base} {...props}>
      <path d="M12 4v11m0 0-4-4m4 4 4-4" />
      <path d="M5 17.5V19a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-1.5" />
    </svg>
  )
}

export function IconCheck(props) {
  return (
    <svg {...base} {...props}>
      <path d="M5 12.5 10 17l9-10" />
    </svg>
  )
}

export function IconPaperclip(props) {
  return (
    <svg {...base} {...props}>
      <path d="M20 11.5 12.4 19a4.5 4.5 0 0 1-6.4-6.4l8-8a3 3 0 0 1 4.3 4.3l-7.9 7.9a1.5 1.5 0 0 1-2.1-2.1l7.1-7.1" />
    </svg>
  )
}

export function IconSend(props) {
  return (
    <svg {...base} {...props}>
      <path d="M4.5 12 20 4.5 15 19l-3.5-6L4.5 12Z" />
    </svg>
  )
}

export function IconShield(props) {
  return (
    <svg {...base} {...props}>
      <path d="M12 3.5c2.6 0 5 .8 5.8 1.1.4.15.6.5.6.9v5.4c0 4.3-2.7 6.9-6.4 8.6-3.7-1.7-6.4-4.3-6.4-8.6V5.5c0-.4.2-.75.6-.9.8-.3 3.2-1.1 5.8-1.1Z" />
    </svg>
  )
}

export function IconUser(props) {
  return (
    <svg {...base} {...props}>
      <circle cx="12" cy="8.5" r="3.2" />
      <path d="M5 20c1-3.4 4-5.2 7-5.2s6 1.8 7 5.2" />
    </svg>
  )
}

export function IconKey(props) {
  return (
    <svg {...base} {...props}>
      <circle cx="8" cy="15" r="3.5" />
      <path d="M10.6 12.4 18 5m0 0v3.5M18 5h-3.5" />
    </svg>
  )
}
