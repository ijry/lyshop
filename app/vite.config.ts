import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
import UniUpRoot from 'uview-plus/libs/root/index.js'

// UnoCSS is ESM-only, must use dynamic import
export default defineConfig(async () => {
  const UnoCSS = (await import('unocss/vite')).default
  const patchUniPagesJsonImportGlob = {
    name: 'patch-uni-pages-json-import-glob',
    enforce: 'pre' as const,
    transform(code: string, id: string) {
      if (id === 'pages-json-js' || id.endsWith('/pages-json-js') || id.endsWith('\\pages-json-js')) {
        return code.replace("import.meta.glob('./locale/*.json'", "import.meta.glob('/locale/*.json'")
      }
      return null
    },
  }

  return {
    plugins: [
      UniUpRoot({
        rootFileName: 'App.up',
        autoCreateRootFile: false,
      }),
      uni(),
      patchUniPagesJsonImportGlob,
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
