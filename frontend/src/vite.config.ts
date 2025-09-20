import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  preview: {
    allowedHosts: ['luincpong.com'],
    port: 5173
  },
  server: {
    allowedHosts: true,
    port: 5173
  },
})