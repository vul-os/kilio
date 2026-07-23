# kilio brand mark

## Concept: "Sealed Cry"

Three directions were sketched before committing:

1. **Sheltering arc** — an open, hand-like arc curving protectively over a
   small voice-dot. Rejected: reads too close to generic "protection /
   umbrella" iconography, and doesn't scale down to a distinctive glyph.
2. **Twin quote-mark** — two small comma-drops nested like an opening
   quotation mark, literally "someone is speaking." Elegant, but as a pair
   of small shapes it lost clarity at 16px and felt more like punctuation
   than a mark.
3. **Sealed Cry (chosen)** — a circular medallion (a wax seal, not a
   rounded-square app-icon blob) with a fine pressed inner rim, holding one
   leaning, ink-drop-shaped voice mark near its center, with a soft engraved
   highlight and a faint smaller echo-drop nearby.

**Sealed Cry** was chosen because every element is doing thematic work:
- The **circle** reads as a seal, a coin, a medallion — permanence, formal
  protection, dignity — and immediately differentiates kilio from generic
  rounded-square app tiles.
- The **pressed inner rim** suggests an impression stamped into wax: this
  mark has been formally sealed, not just decorated.
- The **single leaning drop** is the outcry itself — one voice, captured and
  held safely, tilted as if written by hand rather than stamped flat and
  static. It doubles as a drop of ink (a claim being set down in writing)
  and a comma (a pause, a breath, being heard).
- The **faint echo-drop** is a second, quieter voice nearby — solidarity and
  witness, without literalism.

It reads clearly as a simple circle-with-drop silhouette at 16px (the rim
and echo fade to near-invisible detail, which is intentional — they reward
larger sizes without competing at small ones), and scales up to a refined,
own-able medallion in headers and marketing contexts.

## Colors

- Indigo accent (light UI): `#5B4BD6`
- Indigo accent (dark UI): `#8E7EFF`
- Mark gradient (light contexts): `#7B6BF0` → `#4B3BCF`
- Mark gradient (dark contexts, brightened): `#8E7EFF` → `#5B4BE6`
- Drop / seal impression: warm paper white `#FBFAF6`
- Wordmark ink (light backgrounds): `#18181F`
- Wordmark paper (dark backgrounds): `#F3F1EC`

The medallion itself is a self-toned emblem (opaque gradient fill), so
`mark.svg` and `favicon.svg` are used as-is on both light and dark
backgrounds without needing separate variants. Only the wordmark text color
changes between `lockup-light.svg` and `lockup-dark.svg`.

## Files

- `mark.svg` — full glyph (medallion + rim + drop + echo), viewBox 0 0 64 64.
- `favicon.svg` (in `web/public/`) — simplified glyph (medallion + drop
  only, hairline details removed) for crispness at 16px.
- `lockup-light.svg` / `lockup-dark.svg` — mark + "kilio" wordmark, viewBox
  0 0 290 84, serif system stack.

The wordmark sets "kilio" in a serif system stack
(`Georgia, 'Iowan Old Style', 'Palatino Linotype', 'Times New Roman', serif`)
to keep the editorial, calm register of the Sanctuary aesthetic.
