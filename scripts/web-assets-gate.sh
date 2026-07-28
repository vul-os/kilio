#!/usr/bin/env bash
#
# kilio web asset gate — "no third-party assets on the reporter page" is a
# stated privacy property (docs/SECURITY.md §4, metadata minimization). This
# checks the BUILT bundle, because that is what a reporter's browser loads.
#
# It asserts, and fails closed on:
#   * dist/index.html exists, and at least one CSS and one JS asset were
#     emitted — otherwise the scan below would "pass" by having nothing to scan
#   * no <script src>/<link href>/<img src> in index.html points at an
#     absolute origin
#   * no `url(http…)` in any emitted CSS (fonts must be bundled, not fetched)
#
# Deliberately NOT flagged: XML namespace URIs (http://www.w3.org/2000/svg and
# friends) and URLs that appear only as display strings — neither causes a
# fetch. Grepping for "any absolute URL" would flag those and get switched off.

set -euo pipefail

cd "$(dirname "$0")/../web"

DIST="dist"
status=0

if [ ! -f "${DIST}/index.html" ]; then
  echo "GATE FAIL: ${DIST}/index.html missing — run 'npm run build' first."
  echo "           NOT VERIFIED: any third-party asset reference in the bundle."
  exit 1
fi

css_count="$(find "${DIST}" -name '*.css' | wc -l | tr -d ' ')"
js_count="$(find "${DIST}" -name '*.js' | wc -l | tr -d ' ')"
if [ "${css_count}" -lt 1 ] || [ "${js_count}" -lt 1 ]; then
  echo "GATE FAIL: expected >=1 CSS and >=1 JS asset, found ${css_count} CSS / ${js_count} JS."
  echo "           NOT VERIFIED: the bundle's asset references — there was nothing to scan."
  exit 1
fi

echo "scanning ${DIST}: 1 index.html, ${css_count} css, ${js_count} js"

html_hits="$(grep -oE '(src|href)="https?://[^"]+"' "${DIST}/index.html" || true)"
if [ -n "${html_hits}" ]; then
  echo "GATE FAIL: index.html loads an external origin:"
  printf '%s\n' "${html_hits}"
  status=1
fi

css_hits="$(find "${DIST}" -name '*.css' -exec grep -oE "url\(['\"]?https?://[^)]+\)" {} + || true)"
if [ -n "${css_hits}" ]; then
  echo "GATE FAIL: a stylesheet fetches from an external origin:"
  printf '%s\n' "${css_hits}"
  status=1
fi

echo
if [ "${status}" -eq 0 ]; then
  echo "GATE OK: no third-party asset reference in the built bundle."
else
  echo "GATE FAILED."
fi
exit "${status}"
