import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
import UniUpRoot from 'uview-plus/libs/root/index.js'

// UnoCSS is ESM-only, must use dynamic import
export default defineConfig(async () => {
  const UnoCSS = (await import('unocss/vite')).default

  return {
    plugins: [
      UniUpRoot({
        rootFileName: 'App.up',
        autoCreateRootFile: false,
      }),
      uni(),
      UnoCSS(),
    ],
    css: {
      preprocessorOptions: {
        scss: {
          silenceDeprecations: ['legacy-js-api', 'import'],
        },
      },
    },
  }
})
