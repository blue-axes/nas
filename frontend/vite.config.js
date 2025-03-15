import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  base: "/static/",
  server: {
    host: '0.0.0.0',
    proxy: {
      '/simple_upload': {
        target: 'http://127.0.0.1:8088',
        changeOrigin: true
      }
    }
  }
})
