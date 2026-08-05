import type { ElementType, ComponentPropsWithoutRef, ReactNode } from 'react'
import './ui.css'
import { useTheme } from '../lib/theme.ts'
import type { Tone } from '../mock/data.ts'

export interface SealProps {
  size?: number
  className?: string
}

/** The kilio seal mark: a full-bleed rounded-square tile holding one leaning
 *  voice-drop and its faint echo (mirrors brand/logo.svg — see cb58ae2,
 *  which replaced this component's earlier circular medallion silhouette). */
export function Seal({ size = 34, className = '' }: SealProps) {
  return (
    <svg width={size} height={size} viewBox="0 0 64 64" className={className} aria-hidden="true">
      <defs>
        <linearGradient id="seal-g" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0" stopColor="#7B6BF0" />
          <stop offset="1" stopColor="#4B3BCF" />
        </linearGradient>
      </defs>
      <rect width="64" height="64" rx="14" fill="url(#seal-g)" />
      <g transform="rotate(-13 32 32)">
        <path d="M32 13 C38.5 20.5 41.5 26.3 41.5 33.5 A9.5 9.5 0 1 1 22.5 33.5 C22.5 26.3 25.5 20.5 32 13 Z" fill="#FBFAF6" />
        <path d="M36 24 C38.5 27 39.5 30 39 33" fill="none" stroke="#fff" strokeOpacity="0.45" strokeWidth="1.1" strokeLinecap="round" />
      </g>
      <g transform="rotate(-10 46 46)" opacity="0.5">
        <path d="M46 40 C48.2 42.6 49.2 44.6 49.2 46 A3.2 3.2 0 1 1 42.8 46 C42.8 44.6 43.8 42.6 46 40 Z" fill="#FBFAF6" />
      </g>
    </svg>
  )
}

export interface LogoProps {
  size?: number
}

export function Logo({ size = 34 }: LogoProps) {
  return (
    <span className="logo"><Seal size={size} className="tile" />kilio</span>
  )
}

type ButtonOwnProps<C extends ElementType> = {
  variant?: 'primary' | 'ghost'
  size?: 'lg' | 'sm'
  as?: C
  className?: string
  children?: ReactNode
}

export type ButtonProps<C extends ElementType> =
  ButtonOwnProps<C> & Omit<ComponentPropsWithoutRef<C>, keyof ButtonOwnProps<C>>

const DEFAULT_BUTTON_TAG = 'button'

export function Button<C extends ElementType = typeof DEFAULT_BUTTON_TAG>({
  variant = 'primary', size, as, className = '', children, ...rest
}: ButtonProps<C>) {
  const cls = `btn btn-${variant} ${size === 'lg' ? 'btn-lg' : ''} ${className}`.trim()
  const Tag = as || DEFAULT_BUTTON_TAG
  return <Tag className={cls} {...rest}>{children}</Tag>
}

export interface PillProps {
  tone?: Tone
  children?: ReactNode
}

export function Pill({ tone = 'muted', children }: PillProps) {
  return <span className={`pill pill-${tone}`}><span className="dot" />{children}</span>
}

export interface FieldProps {
  label?: ReactNode
  hint?: ReactNode
  children?: ReactNode
}

export function Field({ label, hint, children }: FieldProps) {
  return (
    <div className="field">
      {label && <label>{label}</label>}
      {children}
      {hint && <span className="hint">{hint}</span>}
    </div>
  )
}

export interface SealBadgeProps {
  children?: ReactNode
}

export function SealBadge({ children = 'Sealed end-to-end' }: SealBadgeProps) {
  return (
    <span className="seal-badge">
      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <rect x="4" y="10" width="16" height="10" rx="2.5" stroke="currentColor" strokeWidth="2" />
        <path d="M8 10V7a4 4 0 0 1 8 0v3" stroke="currentColor" strokeWidth="2" />
      </svg>
      {children}
    </span>
  )
}

export interface StepperProps {
  count: number
  active: number
}

export function Stepper({ count, active }: StepperProps) {
  return (
    <div className="stepper" aria-label={`Step ${active + 1} of ${count}`}>
      {Array.from({ length: count }).map((_, i) => (
        <span key={i} className={`step ${i < active ? 'done' : i === active ? 'active' : ''}`} />
      ))}
    </div>
  )
}

export function ThemeToggle() {
  const { theme, toggle } = useTheme()
  return (
    <button className="theme-toggle" onClick={toggle} aria-label="Toggle light or dark theme" title="Toggle theme">
      {theme === 'dark' ? (
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
          <circle cx="12" cy="12" r="4" /><path d="M12 2v2M12 20v2M4 12H2M22 12h-2M5 5l1.5 1.5M17.5 17.5 19 19M19 5l-1.5 1.5M6.5 17.5 5 19" />
        </svg>
      ) : (
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z" />
        </svg>
      )}
    </button>
  )
}
