const assert = require('node:assert/strict')
const { readFileSync, existsSync } = require('node:fs')
const path = require('node:path')
const test = require('node:test')
const ts = require('typescript')
const axios = require('axios')

const src = path.resolve(__dirname, '../src')

function loadModules(mock) {
  const cache = new Map()
  const calls = []
  function load(filename) {
    if (filename.endsWith('.json')) return require(filename)
    if (!filename.endsWith('.ts')) filename = existsSync(filename + '.ts') ? filename + '.ts' : path.join(filename, 'index.ts')
    if (cache.has(filename)) return cache.get(filename)
    const exports = {}
    cache.set(filename, exports)
    const source = readFileSync(filename, 'utf8').replaceAll('import.meta.env.VITE_MOCK', JSON.stringify(String(mock)))
    const code = ts.transpileModule(source, {
      compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2020, esModuleInterop: true },
      fileName: filename,
    }).outputText
    new Function('require', 'exports', 'localStorage', code)((id) => {
      if (id === 'axios') return { ...axios, create(config) {
        return axios.create({ ...config, adapter: async request => {
          const url = axios.getUri(request)
          assert.ok(!mock, `Mock request unexpectedly reached network adapter: ${url}`)
          calls.push({ method: request.method, url, data: request.data })
          return { data: { code: 0, data: null }, status: 200, statusText: 'OK', headers: {}, config: request }
        } })
      } }
      if (id.startsWith('@/')) return load(path.join(src, id.slice(2)))
      if (id.startsWith('.')) return load(path.resolve(path.dirname(filename), id))
      return require(id)
    }, exports, { getItem: () => null })
    return exports
  }
  const api = load(path.join(src, 'api/distribution.ts'))
  return { api, calls }
}

test('all distribution requests use exactly one admin API prefix', async () => {
  const { api, calls } = loadModules(false)
  await api.getDistributionConfig()
  await api.updateDistributionConfig({ enabled: true })
  await api.getDistributors({ page: 2, size: 1, status: 'active' })
  await api.getDistributor(1)
  await api.updateDistributor(1, { status: 'disabled' })
  await api.getDistributionOrders()
  await api.getDistributionOrder(1)
  await api.getWithdrawals()
  await api.getWithdrawal(1)
  await api.approveWithdrawal(1)
  await api.rejectWithdrawal(1, '信息待确认')
  await api.completeWithdrawal(1)
  assert.deepEqual(calls.map(({ method, url }) => `${method} ${url}`), [
    'get /admin/api/distribution/config',
    'put /admin/api/distribution/config',
    'get /admin/api/distribution/distributors?page=2&size=1&status=active',
    'get /admin/api/distribution/distributors/1',
    'put /admin/api/distribution/distributors/1',
    'get /admin/api/distribution/orders',
    'get /admin/api/distribution/orders/1',
    'get /admin/api/distribution/withdrawals',
    'get /admin/api/distribution/withdrawals/1',
    'post /admin/api/distribution/withdrawals/1/approve',
    'post /admin/api/distribution/withdrawals/1/reject',
    'post /admin/api/distribution/withdrawals/1/complete',
  ])
  assert.deepEqual(JSON.parse(calls[10].data), { reason: '信息待确认' })
})

test('demo distributors support current fields, pagination, status filters and updates', async () => {
  const { api } = loadModules(true)
  const all = await api.getDistributors()
  assert.ok(all.total >= 3)
  for (const row of all.list) {
    assert.equal(typeof row.real_name, 'string')
    assert.ok(['pending', 'active', 'disabled'].includes(row.status))
    for (const field of ['total_earnings', 'available_amount', 'frozen_amount', 'withdrawn_amount']) {
      assert.equal(typeof row[field], 'number', field)
    }
  }
  const second = await api.getDistributors({ page: 2, size: 1 })
  assert.equal(second.total, all.total)
  assert.deepEqual(second.list, all.list.slice(1, 2))
  const pending = await api.getDistributors({ status: 'pending' })
  assert.ok(pending.list.length)
  assert.ok(pending.list.every(row => row.status === 'pending'))
  await api.updateDistributor(pending.list[0].id, { status: 'active' })
  assert.equal((await api.getDistributor(pending.list[0].id)).status, 'active')
  assert.equal((await api.getDistributors({ status: 'pending' })).total, pending.total - 1)
  assert.deepEqual((await api.getDistributors({ page: 999 })).list, [])
})

test('demo orders return numeric amounts and support status and distributor filters', async () => {
  const { api } = loadModules(true)
  const all = await api.getDistributionOrders()
  assert.ok(all.list.length)
  for (const row of all.list) {
    assert.ok(['pending', 'settled', 'cancelled'].includes(row.status))
    for (const field of ['commission', 'commission_rate', 'order_amount']) assert.equal(typeof row[field], 'number')
    assert.deepEqual(await api.getDistributionOrder(row.id), row)
  }
  const target = all.list[0]
  const filtered = await api.getDistributionOrders({ status: target.status, distributor_id: target.distributor_id })
  assert.ok(filtered.list.length)
  assert.ok(filtered.list.every(row => row.status === target.status && row.distributor_id === target.distributor_id))
})

test('demo withdrawal actions reject invalid transitions and update balances once', async () => {
  const { api } = loadModules(true)
  const pending = await api.getWithdrawals({ status: 'pending' })
  assert.ok(pending.list.length >= 2)
  const [complete, reject] = pending.list
  assert.ok(pending.list.every(row => row.status === 'pending'))
  await assert.rejects(api.completeWithdrawal(complete.id), /状态/)
  const beforeComplete = await api.getDistributor(complete.distributor_id)
  await api.approveWithdrawal(complete.id)
  assert.equal((await api.getWithdrawal(complete.id)).status, 'approved')
  await api.completeWithdrawal(complete.id)
  assert.equal((await api.getWithdrawal(complete.id)).status, 'completed')
  const afterComplete = await api.getDistributor(complete.distributor_id)
  assert.equal(afterComplete.frozen_amount, beforeComplete.frozen_amount - complete.amount)
  assert.equal(afterComplete.withdrawn_amount, beforeComplete.withdrawn_amount + complete.amount)
  await assert.rejects(api.completeWithdrawal(complete.id), /状态/)
  const beforeReject = await api.getDistributor(reject.distributor_id)
  await assert.rejects(api.rejectWithdrawal(reject.id, ''), /原因/)
  await api.rejectWithdrawal(reject.id, '收款信息待确认')
  const rejected = await api.getWithdrawal(reject.id)
  assert.equal(rejected.status, 'rejected')
  assert.equal(rejected.reject_reason, '收款信息待确认')
  const afterReject = await api.getDistributor(reject.distributor_id)
  assert.equal(afterReject.frozen_amount, beforeReject.frozen_amount - reject.amount)
  assert.equal(afterReject.available_amount, Math.round((beforeReject.available_amount + reject.amount) * 100) / 100)
  await assert.rejects(api.rejectWithdrawal(reject.id, '再次拒绝'), /状态/)
})

test('demo configuration round-trips all editable fields', async () => {
  const { api } = loadModules(true)
  const initial = await api.getDistributionConfig()
  assert.equal(initial.level1_rate, 10)
  const values = {
    enabled: false, level: 3, level1_rate: 12, level2_rate: 6, level3_rate: 3,
    min_withdraw: 200, withdraw_fee_rate: 1, auto_approve: true, require_real_name: false,
  }
  await api.updateDistributionConfig(values)
  const result = await api.getDistributionConfig()
  for (const [key, value] of Object.entries(values)) assert.equal(result[key], value, key)
})
