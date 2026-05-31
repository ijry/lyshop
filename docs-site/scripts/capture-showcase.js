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
  },
  {
    name: 'eapp-dashboard',
    path: '/eapp-demo/index.html#/pages/index/index',
    viewport: { width: 390, height: 844 },
    waitFor: 'uni-page-body',
  },
  {
    name: 'eapp-products',
    path: '/eapp-demo/index.html#/pages/product/list',
    viewport: { width: 390, height: 844 },
    waitFor: 'uni-page-body',
  },
  {
    name: 'eapp-marketing-index',
    path: '/eapp-demo/index.html#/pages/marketing/index',
    viewport: { width: 390, height: 844 },
    waitFor: 'uni-page-body',
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

  let browser
  try {
    browser = await chromium.launch()
    for (const item of captures) {
      await capture(browser, item)
    }
  } finally {
    if (browser) await browser.close()
    await closeServer(server)
  }

  console.log(`[showcase] Screenshots saved to ${outputDir}`)
}

async function capture(browser, item) {
  const page = await browser.newPage({
    viewport: item.viewport,
    deviceScaleFactor: item.viewport.width <= 430 ? 2 : 1,
  })
  const target = `${origin}${item.path}`
  const filePath = path.join(outputDir, `${item.name}.png`)

  console.log(`[showcase] Capturing ${item.name}: ${target}`)
  await page.goto(target, { waitUntil: 'networkidle' })
  await page.waitForSelector(item.waitFor, { timeout: 15000 })
  await page.waitForTimeout(1200)
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
