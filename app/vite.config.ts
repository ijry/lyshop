import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'

// UnoCSS is ESM-only, must use dynamic import
export default defineConfig(async () => {
  const UnoCSS = (await import('unocss/vite')).default

  return {
    plugins: [
      uni(),
      UnoCSS(),
    ],
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: '@import "uview-plus/theme.scss";',
          silenceDeprecations: ['legacy-js-api', 'import'],
        },
      },
    },
  }
})
