import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig(async () => {
  const UnoCSS = (await import('unocss/vite')).default

  return {
    plugins: [vue(), UnoCSS()],
    resolve: {
      alias: { '@': resolve(__dirname, 'src') }
    },
    server: {
      port: 9527,
      fs: { allow: ['..'] },
      proxy: {
        '/admin/api': { target: 'http://localhost:8080', changeOrigin: true },
        '/api':       { target: 'http://localhost:8080', changeOrigin: true }
      }
    }
  }
})
