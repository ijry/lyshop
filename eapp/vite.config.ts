import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'

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
    plugins: [uni(), patchUniPagesJsonImportGlob, UnoCSS()],
    css: {
      preprocessorOptions: {
        scss: {
          silenceDeprecations: ['legacy-js-api', 'import'],
        },
      },
    },
  }
})
