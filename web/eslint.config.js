import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'
import { defineConfig, globalIgnores } from 'eslint/config'

export default defineConfig([
  globalIgnores(['dist']),
  // web/src is TS/TSX end to end (phase 1) — parse it with the
  // typescript-eslint parser and lint it with the type-aware TS rule set.
  // kilio is a sealed anonymous-first sensitive-claims intake system; a
  // dropped promise in the crypto/seal/submission paths silently reports
  // success on a claim that never actually landed. Type-aware rules are
  // enabled (parserOptions.projectService) specifically so
  // no-floating-promises and the no-unsafe-* family actually run.
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      js.configs.recommended,
      ...tseslint.configs.recommendedTypeChecked,
      reactHooks.configs.flat.recommended,
    ],
    plugins: { 'react-refresh': reactRefresh },
    languageOptions: {
      globals: globals.browser,
      parserOptions: {
        ecmaFeatures: { jsx: true },
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
    rules: {
      'react-refresh/only-export-components': ['warn', { allowConstantExport: true }],
    },
  },
  // vite.config.js and tools/*.mjs are plain Node scripts, not part of the
  // TS-migrated app surface — syntax-only lint with eslint-recommended.
  {
    files: ['*.config.js', 'tools/**/*.mjs'],
    extends: [js.configs.recommended],
    languageOptions: {
      globals: globals.node,
    },
  },
])
