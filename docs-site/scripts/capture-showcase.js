#!/usr/bin/env node

const fs = require('fs')
const http = require('http')
const path = require('path')
const { chromium } = require('playwright')

const publicDir = path.resolve(__dirname, '..', 'docs', 'public')
const outputDir = path.join(publicDir, 'showcase')
const port = Number(process.env.LYSHOP_SHOWCASE_PORT || 4187)
const origin = `http://127.0.0.1:${port}`

const captures = [
  {
    name: 'h5-home',
    path: '/demo/index.html#/pages/index/index',
    viewport: { width: 390, height: 844 },
    waitFor: 'uni-page-body',
  },
  {
    name: 'h5-products',
    path: '/demo/index.html#/pages/product/list',
    viewport: { width: 390, height: 844 },
    waitFor: 'uni-page-body',
  },
  {
    name: 'h5-orders',
    path: '/demo/index.html#/pages/order/list',
    viewport: { width: 390, height: 844 },
    waitFor: 'uni-page-body',
  },
  {
    name: 'web-home',
    path: '/web-demo/index.html#/',
    viewport: { width: 1440, height: 960 },
    waitFor: '#app',
  },
  // eapp (merchant mini-program H5)
  {
    name: 'eapp-orders',
    path: '/eapp-demo/index.html#/pages/order/list',
    viewport: { width: 390, height: 844 },
    waitFor: 'uni-page-body',
    expect: '订单',
  },
  {
    name: 'eapp-dashboard',
    path: '/eapp-demo/index.html#/pages/dashboard/index',
    viewport: { width: 390, height: 844 },
    waitFor: 'uni-page-body',
    expect: '工作台',
  },
  {
    name: 'eapp-products',
    path: '/eapp-demo/index.html#/pages/product/list',
    viewport: { width: 390, height: 844 },
    waitFor: 'uni-page-body',
    expect: '商品',
  },
  {
    name: 'eapp-marketing-index',
    path: '/eapp-demo/index.html#/pages/marketing/index',
    viewport: { width: 390, height: 844 },
    waitFor: 'uni-page-body',
    expect: '营销中心',
  },
  {
    name: 'eapp-marketing-seckill',
    path: '/eapp-demo/index.html#/pages/marketing/seckill',
    viewport: { width: 390, height: 844 },
    waitFor: 'uni-page-body',
    expect: '秒杀',
  },
]

main().catch((error) => {
  console.error('[showcase] Failed:', error)
  process.exitCode = 1
})

async function main() {
  ensureDemoExists('demo/index.html')
  ensureDemoExists('web-demo/index.html')
  ensureDemoExists('eapp-demo/index.html')
  fs.mkdirSync(outputDir, { recursive: true })

  const server = await createServer(publicDir, port)
  console.log(`[showcase] Serving ${publicDir} at ${origin}`)

  // 支持只重截某一端，避免更新商家端时覆盖未改动的买家端 / PC 截图
  const onlyArg = process.argv.find((arg) => arg.startsWith('--only='))
  const only = process.env.LYSHOP_SHOWCASE_ONLY || (onlyArg ? onlyArg.slice('--only='.length) : '')
  const selected = only ? captures.filter((item) => item.name.startsWith(only)) : captures
  console.log(`[showcase] ${selected.length}/${captures.length} captures${only ? ` (--only=${only})` : ''}`)

  let browser
  try {
    browser = await launchBrowser()
    for (const item of selected) {
      await capture(browser, item)
    }
  } finally {
    if (browser) await browser.close()
    await closeServer(server)
  }

  console.log(`[showcase] Screenshots saved to ${outputDir}`)
}

// 优先使用 Playwright 自带 Chromium；未执行 playwright install 的机器上
// 回退到系统 Chrome / Edge，避免仅因浏览器版本不匹配就无法更新官网截图。
async function launchBrowser() {
  const candidates = [
    { label: 'Playwright Chromium', options: {} },
    { label: '系统 Chrome', options: { channel: 'chrome' } },
    { label: '系统 Edge', options: { channel: 'msedge' } },
  ]
  let lastError
  for (const candidate of candidates) {
    try {
      const browser = await chromium.launch(candidate.options)
      console.log(`[showcase] Browser: ${candidate.label}`)
      return browser
    } catch (error) {
      lastError = error
      console.warn(`[showcase] ${candidate.label} 不可用：${String(error.message).split('\n')[0]}`)
    }
  }
  throw lastError
}

async function capture(browser, item) {
  const page = await browser.newPage({
    viewport: item.viewport,
    deviceScaleFactor: item.viewport.width <= 430 ? 2 : 1,
  })
  // 商家端 demo 需要登录态；mock 模式使用固定演示 token，避免首帧闪登录页
  if (item.name.startsWith('eapp-')) {
    await page.addInitScript(() => {
      try {
        localStorage.setItem('eapp_admin_token', 'demo_admin_token')
        localStorage.setItem('eapp_admin_username', 'admin')
      } catch (_) { /* 存储不可用时忽略 */ }
    })
  }
  const target = `${origin}${item.path}`
  const filePath = path.join(outputDir, `${item.name}.png`)

  const pageErrors = []
  page.on('pageerror', (error) => pageErrors.push(String(error.message)))

  console.log(`[showcase] Capturing ${item.name}: ${target}`)
  // uni-app H5 demo 存在轮询请求，networkidle 可能永远不达成，按 DOM + 页面文本等待即可
  await page.goto(target, { waitUntil: 'domcontentloaded', timeout: 30000 })
  await page.waitForSelector(item.waitFor, { timeout: 15000 })
  if (item.expect) {
    try {
      await page.waitForFunction(
        (text) => (document.body?.innerText || '').includes(text),
        item.expect,
        { timeout: 12000 },
      )
    } catch (_) {
      const body = (await page.evaluate(() => (document.body?.innerText || '').slice(0, 200))) || ''
      throw new Error(
        `${item.name} 未渲染出预期内容「${item.expect}」，实际内容：${body || '(空)'}` +
          (pageErrors.length ? `；页面错误：${pageErrors.slice(0, 2).join(' | ')}` : ''),
      )
    }
  }
  await page.waitForTimeout(1500)
  if (pageErrors.length) {
    console.warn(`[showcase] ${item.name} 页面错误（不影响截图）：${pageErrors.slice(0, 2).join(' | ')}`)
  }
  await page.screenshot({ path: filePath, fullPage: false })
  await page.close()
}

function ensureDemoExists(relativeFile) {
  const filePath = path.join(publicDir, relativeFile)
  if (!fs.existsSync(filePath)) {
    throw new Error(`${relativeFile} not found. Run npm run showcase:build to rebuild demos before capturing.`)
  }
}

function createServer(root, listenPort) {
  const server = http.createServer((req, res) => {
    const url = new URL(req.url || '/', origin)
    const decoded = decodeURIComponent(url.pathname)
    const safePath = path.normalize(decoded).replace(/^(\.\.[/\\])+/, '')
    let filePath = path.join(root, safePath)

    if (!filePath.startsWith(root)) {
      res.writeHead(403)
      res.end('Forbidden')
      return
    }

    if (fs.existsSync(filePath) && fs.statSync(filePath).isDirectory()) {
      filePath = path.join(filePath, 'index.html')
    }

    if (!fs.existsSync(filePath)) {
      res.writeHead(404)
      res.end('Not found')
      return
    }

    res.writeHead(200, { 'Content-Type': contentType(filePath) })
    fs.createReadStream(filePath).pipe(res)
  })

  return new Promise((resolve, reject) => {
    server.once('error', reject)
    server.listen(listenPort, '127.0.0.1', () => resolve(server))
  })
}

function closeServer(server) {
  return new Promise((resolve, reject) => {
    server.close((error) => (error ? reject(error) : resolve()))
  })
}

function contentType(filePath) {
  const ext = path.extname(filePath)
  return {
    '.css': 'text/css; charset=utf-8',
    '.html': 'text/html; charset=utf-8',
    '.js': 'text/javascript; charset=utf-8',
    '.json': 'application/json; charset=utf-8',
    '.png': 'image/png',
    '.svg': 'image/svg+xml',
    '.webp': 'image/webp',
    '.woff2': 'font/woff2',
  }[ext] || 'application/octet-stream'
}
