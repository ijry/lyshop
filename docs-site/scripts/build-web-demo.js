#!/usr/bin/env node

const { execSync } = require('child_process')
const fs = require('fs')
const path = require('path')

const rootDir = path.resolve(__dirname, '..', '..')
const webDir = path.join(rootDir, 'web')
const outputDir = path.join(rootDir, 'docs-site', 'docs', 'public', 'web-demo')

console.log('[web-demo] Building PC demo...')

execSync('npm run build:demo -- --base=./', {
  cwd: webDir,
  stdio: 'inherit',
})

const buildDir = path.join(webDir, 'dist')
if (!fs.existsSync(buildDir)) {
  throw new Error(`Build output not found at ${buildDir}`)
}

fs.rmSync(outputDir, { recursive: true, force: true })
fs.mkdirSync(outputDir, { recursive: true })
copyDir(buildDir, outputDir)

console.log(`[web-demo] Copied build output to ${outputDir}`)

function copyDir(src, dest) {
  for (const entry of fs.readdirSync(src, { withFileTypes: true })) {
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
