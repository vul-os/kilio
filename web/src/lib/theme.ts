import { useCallback, useEffect, useState, type Dispatch, type SetStateAction } from 'react'

const KEY = 'kilio-theme'

export type Theme = 'light' | 'dark'

function current(): Theme {
  if (typeof document === 'undefined') return 'light'
  return (document.documentElement.getAttribute('data-theme') as Theme | null) || 'light'
}

export interface UseThemeResult {
  theme: Theme
  setTheme: Dispatch<SetStateAction<Theme>>
  toggle: () => void
}

/** Theme hook: reads/sets data-theme on <html>, persists to localStorage. */
export function useTheme(): UseThemeResult {
  const [theme, setTheme] = useState<Theme>(current)

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme)
    try { localStorage.setItem(KEY, theme) } catch { /* localStorage may be unavailable (private mode / disabled) */ }
  }, [theme])

  const toggle = useCallback(() => {
    setTheme((t) => (t === 'dark' ? 'light' : 'dark'))
  }, [])

  return { theme, setTheme, toggle }
}
