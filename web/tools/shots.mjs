// Screenshot the kilio surfaces (light + dark) into docs/ and site/ screenshots.
// Usage: node tools/shots.mjs  (expects the dev server on 127.0.0.1:5273)
import { chromium } from 'playwright'
import { mkdirSync, copyFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'

const __dirname = dirname(fileURLToPath(import.meta.url))
const REPO = join(__dirname, '..', '..')
const OUT = join(REPO, 'docs', 'screenshots')
const SITE = join(REPO, 'site', 'screenshots')
mkdirSync(OUT, { recursive: true })
mkdirSync(SITE, { recursive: true })

const BASE = process.env.KILIO_BASE || 'http://127.0.0.1:5273'

// [name, route, viewport]
const desktop = { width: 1440, height: 900 }
const wide = { width: 1512, height: 950 }
const phone = { width: 402, height: 860 }

const SHOTS = [
  ['reporter-landing', '/report', desktop],
  ['reporter-new', '/report/new', desktop],
  ['reporter-receipt', '/report/receipt', desktop],
  ['reporter-thread', '/report/thread', desktop],
  ['reporter-landing-mobile', '/report', phone],
  ['reporter-new-mobile', '/report/new', phone],
  ['handler-inbox', '/handler', wide],
  ['handler-case', '/handler/case/7Q4K-2M9F', wide],
  ['handler-branches', '/handler/branches', wide],
  ['handler-settings', '/handler/settings', wide],
]

const run = async () => {
  const browser = await chromium.launch()
  for (const theme of ['light', 'dark']) {
    for (const [name, route, viewport] of SHOTS) {
      const ctx = await browser.newContext({
        viewport,
        deviceScaleFactor: 2,
        colorScheme: theme,
      })
      await ctx.addInitScript((t) => {
        try { localStorage.setItem('kilio-theme', t) } catch (e) {}
      }, theme)
      const page = await ctx.newPage()
      try {
        await page.goto(BASE + route, { waitUntil: 'networkidle', timeout: 20000 })
      } catch (e) {
        await page.goto(BASE + route, { waitUntil: 'load', timeout: 20000 })
      }
      await page.waitForTimeout(600)
      const file = join(OUT, `${name}-${theme}.png`)
      await page.screenshot({ path: file })
      copyFileSync(file, join(SITE, `${name}-${theme}.png`))
      console.log('shot', name, theme)
      await ctx.close()
    }
  }
  await browser.close()
  console.log('done')
}

run().catch((e) => { console.error(e); process.exit(1) })
