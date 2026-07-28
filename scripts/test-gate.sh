#!/usr/bin/env bash
#
# kilio test gate — runs the suite and asserts it actually ran.
#
# The failure mode this exists to prevent is a green gate that verified
# nothing: a filter that matched no tests, a crate silently dropped from the
# workspace, a harness that compiled and exited 0 without executing a case.
# So this script does not trust "cargo test exited 0". It parses the harness's
# own summary line per crate and asserts:
#
#   * a `test result:` line exists at all      (else: FAIL, naming the crate)
#   * failed == 0
#   * ignored == 0                             (a skipped test is not coverage)
#   * passed >= the floor recorded below
#
# Any unparsable output is a FAILURE, never a pass — there is no "could not
# determine, assume fine" path here. Raise a floor when you add tests; never
# lower one to make a red gate green.

set -euo pipefail

cd "$(dirname "$0")/.."

# crate:minimum-passing-tests
FLOORS=(
  "kilio-seal:37"
  "kilio-core:13"
)

status=0
total=0

for entry in "${FLOORS[@]}"; do
  crate="${entry%%:*}"
  floor="${entry##*:}"

  echo "== ${crate} (floor: ${floor} passing) =="
  out=""
  if ! out="$(cargo test -p "${crate}" --lib 2>&1)"; then
    echo "${out}"
    echo "GATE FAIL: ${crate}: cargo test exited non-zero."
    status=1
    continue
  fi

  summary="$(printf '%s\n' "${out}" | grep -E '^test result:' | tail -1 || true)"
  if [ -z "${summary}" ]; then
    echo "${out}"
    echo "GATE FAIL: ${crate}: no 'test result:' line in the harness output."
    echo "           NOT VERIFIED: whether a single ${crate} test executed."
    status=1
    continue
  fi

  passed="$(printf '%s' "${summary}" | sed -n 's/.*ok\. \([0-9]*\) passed.*/\1/p')"
  failed="$(printf '%s' "${summary}" | sed -n 's/.*; \([0-9]*\) failed.*/\1/p')"
  ignored="$(printf '%s' "${summary}" | sed -n 's/.*; \([0-9]*\) ignored.*/\1/p')"

  if [ -z "${passed}" ] || [ -z "${failed}" ] || [ -z "${ignored}" ]; then
    echo "GATE FAIL: ${crate}: could not parse counts from: ${summary}"
    echo "           NOT VERIFIED: pass/fail/ignore counts for ${crate}."
    status=1
    continue
  fi

  echo "   ${summary}"

  if [ "${failed}" -ne 0 ]; then
    echo "GATE FAIL: ${crate}: ${failed} test(s) failed."
    status=1
  fi
  if [ "${ignored}" -ne 0 ]; then
    echo "GATE FAIL: ${crate}: ${ignored} test(s) ignored — an ignored test is not coverage."
    status=1
  fi
  if [ "${passed}" -lt "${floor}" ]; then
    echo "GATE FAIL: ${crate}: ${passed} passing < floor ${floor}. Coverage went backwards."
    status=1
  fi

  total=$(( total + passed ))
done

echo
if [ "${status}" -eq 0 ]; then
  echo "GATE OK: ${total} tests executed and passed across ${#FLOORS[@]} crates."
else
  echo "GATE FAILED."
fi
exit "${status}"
