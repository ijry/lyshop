import { defineConfig, type Plugin } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
import { resolve } from 'path'
import { cpSync, existsSync } from 'fs'

/**
 * H5 构建时把 static/ 同步到产物根目录。
 * pages.json 的 tabBar iconPath 在 H5 端会被解析为 /static/... 绝对路径，
 * uni 的 Vite 构建默认不复制 static 目录，需要在此显式带上，避免图标 404。
 */
function copyStaticAssets(): Plugin {
  let outDir = ''
  return {
    name: 'eapp-copy-static',
    apply: 'build',
    enforce: 'post',
    configResolved(config) {
      outDir = resolve(config.root, config.build.outDir)
    },
    closeBundle() {
      if (process.env.UNI_PLATFORM && process.env.UNI_PLATFORM !== 'h5') return
      const src = resolve(__dirname, 'static')
      if (!outDir || !existsSync(src)) return
      cpSync(src, resolve(outDir, 'static'), { recursive: true })
    },
  }
}

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
    plugins: [uni(), patchUniPagesJsonImportGlob, UnoCSS(), copyStaticAssets()],
    resolve: {
      alias: {
        '@': resolve(__dirname, '.'),
      },
    },
    css: {
      preprocessorOptions: {
        scss: {
          silenceDeprecations: ['legacy-js-api', 'import'],
        },
      },
    },
  }
})
