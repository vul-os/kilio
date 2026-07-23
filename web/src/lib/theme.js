import { useCallback, useEffect, useState } from 'react'

const KEY = 'kilio-theme'

function current() {
  if (typeof document === 'undefined') return 'light'
  return document.documentElement.getAttribute('data-theme') || 'light'
}

/** Theme hook: reads/sets data-theme on <html>, persists to localStorage. */
export function useTheme() {
  const [theme, setTheme] = useState(current)

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme)
    try { localStorage.setItem(KEY, theme) } catch (e) {}
  }, [theme])

  const toggle = useCallback(() => {
    setTheme((t) => (t === 'dark' ? 'light' : 'dark'))
  }, [])

  return { theme, setTheme, toggle }
}
