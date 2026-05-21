#!/usr/bin/env node
/**
 * Build H5 demo with mock data for embedding in docs-site.
 *
 * What it does:
 *   1. Temporarily patches manifest.json to use hash router
 *   2. Builds uni-app H5 with VITE_MOCK=true
 *   3. Restores manifest.json
 *   4. Copies build output to docs-site/docs/public/demo/
 *
 * Usage:
 *   node scripts/build-demo.js [output-dir]
 */

const { execSync } = require('child_process')
const fs = require('fs')
const path = require('path')

const appDir = path.resolve(__dirname, '..')
const manifestPath = path.join(appDir, 'manifest.json')
const defaultOutput = path.resolve(appDir, '..', 'docs-site', 'docs', 'public', 'demo')
const outputDir = process.argv[2] || defaultOutput

console.log('[demo] Starting H5 demo build...')
console.log('[demo] Output dir:', outputDir)

// 1. Read and backup manifest
const originalManifest = fs.readFileSync(manifestPath, 'utf-8')
const manifest = JSON.parse(originalManifest)

// 2. Patch: hash router, relative base
if (!manifest.h5) manifest.h5 = {}
if (!manifest.h5.router) manifest.h5.router = {}
manifest.h5.router.mode = 'hash'
manifest.h5.router.base = './'

fs.writeFileSync(manifestPath, JSON.stringify(manifest, null, 2), 'utf-8')
console.log('[demo] Patched manifest.json → hash router, base="./"')

try {
  // 3. Build
  const isWin = process.platform === 'win32'
  const npx = isWin ? 'npx.cmd' : 'npx'
  execSync(`${npx} cross-env VITE_MOCK=true ${npx} uni build`, {
    cwd: appDir,
    stdio: 'inherit',
    env: { ...process.env, VITE_MOCK: 'true' },
  })
  console.log('[demo] H5 build complete')

  // 4. Copy output
  const buildDir = path.join(appDir, 'dist', 'build', 'h5')
  if (!fs.existsSync(buildDir)) {
    throw new Error(`Build output not found at ${buildDir}`)
  }

  // Clean and recreate output dir
  if (fs.existsSync(outputDir)) {
    fs.rmSync(outputDir, { recursive: true })
  }
  fs.mkdirSync(outputDir, { recursive: true })

  // Recursive copy
  copyDir(buildDir, outputDir)
  console.log(`[demo] Copied build output to ${outputDir}`)
  console.log('[demo] Done!')
} finally {
  // 5. Always restore manifest
  fs.writeFileSync(manifestPath, originalManifest, 'utf-8')
  console.log('[demo] Restored manifest.json')
}

function copyDir(src, dest) {
  const entries = fs.readdirSync(src, { withFileTypes: true })
  for (const entry of entries) {
    const srcPath = path.join(src, entry.name)
    const destPath = path.join(dest, entry.name)
    if (entry.isDirectory()) {
      fs.mkdirSync(destPath, { recursive: true })
      copyDir(srcPath, destPath)
    } else {
      fs.copyFileSync(srcPath, destPath)
    }
  }
}
