import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// kilio web — reporter + handler surfaces. Static SPA, no backend needed in dev
// (screens run on mock data from src/mock).
export default defineConfig({
  plugins: [react()],
  server: { port: 5273, host: '127.0.0.1' },
  build: { outDir: 'dist', sourcemap: false },
})
