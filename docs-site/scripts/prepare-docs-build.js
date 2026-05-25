const fs = require('node:fs/promises')
const path = require('node:path')

const distDir = path.resolve(__dirname, '../docs/.vitepress/dist')
const retryableCodes = new Set(['EPERM', 'EBUSY', 'ENOTEMPTY'])

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

async function removeWithRetry(target, maxAttempts = 8) {
  for (let attempt = 1; attempt <= maxAttempts; attempt += 1) {
    try {
      await fs.rm(target, { recursive: true, force: true })
      return
    } catch (error) {
      if (!retryableCodes.has(error?.code) || attempt === maxAttempts) {
        throw error
      }
      await sleep(200 * attempt)
    }
  }
}

async function renameAsFallback(target) {
  const fallback = `${target}.stale-${Date.now()}`
  await fs.rename(target, fallback)
  await removeWithRetry(fallback, 8)
}

async function main() {
  try {
    await fs.access(distDir)
  } catch {
    return
  }

  try {
    await removeWithRetry(distDir, 8)
  } catch (error) {
    if (!retryableCodes.has(error?.code)) throw error
    await renameAsFallback(distDir)
  }
}

main().catch((error) => {
  console.error('[docs:build] pre-clean failed:', error?.message || error)
  process.exit(1)
})
