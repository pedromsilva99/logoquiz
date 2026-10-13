import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    // Same-origin requests in dev, so the player cookie just works.
    proxy: { '/api': 'http://localhost:8080' },
  },
})
