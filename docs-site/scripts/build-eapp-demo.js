#!/usr/bin/env node
/**
 * Build eapp H5 demo for docs-site.
 *
 * What it does:
 *   1. Temporarily patches eapp manifest.json to use hash router + relative base
 *   2. Builds eapp H5
 *   3. Restores manifest.json
 *   4. Copies build output to docs-site/docs/public/eapp-demo/
 *
 * Usage:
 *   node scripts/build-eapp-demo.js [output-dir]
 */

const { execSync } = require('child_process')
const fs = require('fs')
const path = require('path')

const rootDir = path.resolve(__dirname, '..', '..')
const eappDir = path.join(rootDir, 'eapp')
const manifestPath = path.join(eappDir, 'manifest.json')
const defaultOutput = path.resolve(__dirname, '..', 'docs', 'public', 'eapp-demo')
const outputDir = process.argv[2] || defaultOutput

console.log('[eapp-demo] Starting H5 demo build...')
console.log('[eapp-demo] Output dir:', outputDir)

const originalManifest = fs.readFileSync(manifestPath, 'utf-8')
const manifest = JSON.parse(originalManifest)

if (!manifest.h5) manifest.h5 = {}
if (!manifest.h5.router) manifest.h5.router = {}
manifest.h5.router.mode = 'hash'
manifest.h5.router.base = './'
manifest.h5.publicPath = './'

fs.writeFileSync(manifestPath, JSON.stringify(manifest, null, 2), 'utf-8')
console.log('[eapp-demo] Patched manifest.json → hash router, base="./"')

try {
  const isWin = process.platform === 'win32'
  const uniCmd = path.join(eappDir, 'node_modules', '.bin', isWin ? 'uni.cmd' : 'uni')
  execSync(`${uniCmd} build`, {
    cwd: eappDir,
    stdio: 'inherit',
    env: { ...process.env, UNI_INPUT_DIR: eappDir },
  })
  console.log('[eapp-demo] H5 build complete')

  const buildDir = path.join(eappDir, 'dist', 'build', 'h5')
  if (!fs.existsSync(buildDir)) {
    throw new Error(`Build output not found at ${buildDir}`)
  }

  if (fs.existsSync(outputDir)) {
    fs.rmSync(outputDir, { recursive: true, force: true })
  }
  fs.mkdirSync(outputDir, { recursive: true })
  copyDir(buildDir, outputDir)

  console.log(`[eapp-demo] Copied build output to ${outputDir}`)
  console.log('[eapp-demo] Done!')
} finally {
  fs.writeFileSync(manifestPath, originalManifest, 'utf-8')
  console.log('[eapp-demo] Restored manifest.json')
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

