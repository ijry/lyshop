#!/usr/bin/env node
/**
 * Build H5 for embedded standalone runtime.
 *
 * What it does:
 *   1. Temporarily patches manifest.json to hash router + relative base/publicPath
 *   2. Builds uni-app H5 with real API mode (no VITE_MOCK)
 *   3. Restores manifest.json
 */

const { execSync } = require('child_process')
const fs = require('fs')
const path = require('path')

const appDir = path.resolve(__dirname, '..')
const manifestPath = path.join(appDir, 'manifest.json')

console.log('[embed] Starting H5 embed build...')

const originalManifest = fs.readFileSync(manifestPath, 'utf-8')
const manifest = JSON.parse(originalManifest)

if (!manifest.h5) manifest.h5 = {}
if (!manifest.h5.router) manifest.h5.router = {}
manifest.h5.router.mode = 'hash'
manifest.h5.router.base = './'
manifest.h5.publicPath = './'

fs.writeFileSync(manifestPath, JSON.stringify(manifest, null, 2), 'utf-8')
console.log('[embed] Patched manifest.json for embedded runtime')

try {
  const isWin = process.platform === 'win32'
  const uniCmd = path.join(appDir, 'node_modules', '.bin', isWin ? 'uni.cmd' : 'uni')
  execSync(`${uniCmd} build`, {
    cwd: appDir,
    stdio: 'inherit',
    env: { ...process.env, UNI_INPUT_DIR: appDir },
  })
  console.log('[embed] H5 embed build complete')
} finally {
  fs.writeFileSync(manifestPath, originalManifest, 'utf-8')
  console.log('[embed] Restored manifest.json')
}
