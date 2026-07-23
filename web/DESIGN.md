# kilio web — design system ("Sanctuary")

The aesthetic for kilio is **Sanctuary**: calm, editorial, trustworthy, humane.
A person reaching this interface may be frightened. Every choice should lower
anxiety and signal safety — never a cold corporate SaaS, never techy neon,
never generic purple-on-white.

## Principles
- **Calm over clever.** Generous whitespace, clear hierarchy, unhurried motion.
- **Trust made visible.** The "seal" motif (a shield holding a voice) recurs;
  sealed content shows a wax-seal/lock; copy states plainly what can and cannot
  be read.
- **Editorial, not app-y.** A warm serif display voice (Fraunces) over a
  humanist sans (Hanken Grotesk). Real typographic hierarchy.
- **Restraint with the accent.** Mostly ink-on-paper (or paper-on-ink); indigo
  appears for key actions and the seal, not everywhere.
- **Accessible by default.** Semantic HTML, labelled controls, visible focus,
  AA contrast in both themes, full keyboard operation, `prefers-reduced-motion`.

## Tokens
Defined in `src/styles/tokens.css` as CSS variables, themed by
`:root` (light) and `:root[data-theme="dark"]`. **Use the variables — never
hard-code colors.** Key ones:

- Surfaces: `--paper`, `--paper-2`, `--panel`, `--ink`, `--ink-2`, `--muted`,
  `--faint`, `--line`, `--line-2`
- Accent: `--indigo`, `--indigo-strong`, `--indigo-soft` (tinted surface),
  `--indigo-ink` (text on indigo)
- Signal: `--ok`, `--warn`, `--danger` (used sparingly, muted)
- Radii: `--r-sm` 10px, `--r` 14px, `--r-lg` 20px
- Type: `--serif` (Fraunces), `--sans` (Hanken Grotesk), `--mono`
- Shadow: `--shadow-1`, `--shadow-2` (soft, low-opacity)

## Fonts
Bundled via `@fontsource-variable/*` (no network at runtime, matching kilio's
ethos). Fraunces for display; Hanken Grotesk for UI/body.

## Themes
Light = warm paper + ink. Dark = deep ink + warm off-white. A manual toggle
(`src/lib/theme.js`) sets `data-theme` on `<html>`; default follows
`prefers-color-scheme`. **Every screen must look intentional in both.**

## Surfaces (two apps, one system)
- **Reporter** (`src/reporter/`, route `/report`) — the public, PWA surface.
  Anonymous, no account. Screens: landing, new-report (stepped), receipt
  (12-word passphrase), return-with-code, sealed thread.
- **Handler** (`src/handler/`, route `/handler`) — the case-worker surface
  (desktop/Tauri). Screens: inbox, case detail (decrypted thread + reply +
  status + content-free audit), branches, settings.

Both consume mock data from `src/mock/` so the app runs and screenshots without
a backend. Real crypto lives in the `kilio-seal` crate; the UI treats sealing as
a given and focuses on the human experience.

## Shared primitives (`src/ui/`)
`Logo`/`Seal`, `Button`, `Pill`, `Field`, `TextArea`, `Card`, `ThemeToggle`,
`Stepper`. Build screens from these; add surface-specific components within each
surface's folder.
