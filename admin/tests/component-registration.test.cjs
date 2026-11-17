const assert = require('node:assert/strict')
const { readFileSync, readdirSync } = require('node:fs')
const path = require('node:path')
const test = require('node:test')
const ts = require('typescript')
const vue = require('vue')
const { parse, compileTemplate } = require('@vue/compiler-sfc')
const { renderToString } = require('vue/server-renderer')

const src = path.resolve(__dirname, '../src')

function evaluate(source, filename, resolveImport) {
  const code = ts.transpileModule(source, {
    compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2020, esModuleInterop: true },
    fileName: filename,
  }).outputText
  const exports = {}
  new Function('require', 'exports', code)(resolveImport, exports)
  return exports
}

function distributorPage() {
  const filename = path.join(src, 'views/distribution/DistributorList.vue')
  const { descriptor } = parse(readFileSync(filename, 'utf8'))
  const compiled = compileTemplate({ source: descriptor.template.content, filename, id: 'distributor-regression' })
  assert.deepEqual(compiled.errors, [])
  const { render } = evaluate(compiled.code, filename, require)
  return {
    render,
    setup: () => ({
      list: [], loading: false, page: 1, size: 20, total: 0,
      statusFilter: '', statusOptions: [{ label: '全部', value: '' }], statusMap: {},
      loadData() {}, handlePageChange() {}, handleStatusChange() {},
    }),
  }
}

// Execute the real entry point, replacing only routing, side effects and DOM mount.
// Element Plus imports and app.use calls stay real so missing registration fails.
function bootstrap(root) {
  let app
  evaluate(readFileSync(path.join(src, 'main.ts'), 'utf8'), 'main.ts', (id) => {
    if (id === 'vue') return { ...vue, createApp: () => {
      app = vue.createSSRApp(root)
      app.mount = () => {}
      return app
    } }
    if (id === './App.vue') return { default: root, __esModule: true }
    if (id === './router') return { default: { install() {} }, __esModule: true }
    if (id === './locales') return { i18n: { install() {} } }
    if (id === './utils/notify') return { setNotifyHandler() {} }
    if (id === './utils/toast') return { pushToast() {} }
    if (id === 'element-plus/es/locale/lang/zh-cn') return require('element-plus/lib/locale/lang/zh-cn')
    if (id.endsWith('.css')) return {}
    return require(id)
  })
  // SSR-only providers; the production app is client rendered.
  const { ID_INJECTION_KEY, ZINDEX_INJECTION_KEY, zIndexContextKey } = require('element-plus')
  app.provide(ID_INJECTION_KEY, { prefix: 1, current: 0 })
  app.provide(ZINDEX_INJECTION_KEY, { current: 0 })
  app.provide(zIndexContextKey, { current: 0 })
  return app
}

function vueFiles(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap(entry => {
    const file = path.join(dir, entry.name)
    return entry.isDirectory() ? vueFiles(file) : file.endsWith('.vue') ? [file] : []
  })
}

function tags(node, result = new Set()) {
  if (node.type === 1 && /^el-/.test(node.tag)) result.add(node.tag)
  for (const child of node.children || []) tags(child, result)
  return result
}

test('unregistered table reproduces the reported row destructuring exception', async () => {
  const app = vue.createSSRApp(distributorPage())
  app.config.warnHandler = () => {}
  await assert.rejects(renderToString(app), /Cannot destructure property 'row'/)
})

test('real admin bootstrap renders an empty distributor table without component warnings', async () => {
  const app = bootstrap(distributorPage())
  const warnings = []
  app.config.warnHandler = message => warnings.push(message)
  const html = await renderToString(app)
  assert.match(html, /el-table/)
  assert.match(html, /el-pagination/)
  assert.doesNotMatch(html, /<el-table[\s>]/)
  assert.deepEqual(warnings, [])
})

test('all Element Plus template tags and loading directive are registered by bootstrap', () => {
  const app = bootstrap({ render: () => null })
  const names = new Set()
  for (const filename of vueFiles(src)) {
    const { descriptor } = parse(readFileSync(filename, 'utf8'))
    if (!descriptor.template) continue
    for (const name of tags(descriptor.template.ast)) names.add(name)
  }
  assert.ok(names.size >= 17, 'Expected all currently used Element Plus components')
  for (const name of names) {
    const componentName = name.split('-').map(part => part[0].toUpperCase() + part.slice(1)).join('')
    assert.ok(app.component(componentName), `${name} is not registered`)
  }
  assert.ok(app.directive('loading'), 'v-loading is not registered')
})
