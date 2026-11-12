/**
 * 生成商家端 tabBar 图标（常态 / 选中态）。
 *
 * 小程序 tabBar 的 iconPath 只接受本地 PNG，且不支持 SVG，
 * 因此这里用纯 Node（zlib + 手写 PNG 编码）把与 AppIcon 同源的线性图形
 * 光栅化为 81x81 的透明 PNG，保证三端图标风格一致。
 *
 * 用法：node scripts/gen-tabbar-icons.cjs
 */
const fs = require('fs')
const path = require('path')
const zlib = require('zlib')

const SIZE = 81
const PAD = 7
const SCALE = (SIZE - PAD * 2) / 24
const STROKE = 1.7 * SCALE
const OUT_DIR = path.join(__dirname, '..', 'static', 'tabbar')

/* ---------- 24x24 坐标系下的几何原语 ---------- */
function rectSegs(x, y, w, h) {
  return [
    [x, y, x + w, y],
    [x + w, y, x + w, y + h],
    [x + w, y + h, x, y + h],
    [x, y + h, x, y],
  ]
}

function circleSegs(cx, cy, r, steps = 40) {
  const out = []
  for (let i = 0; i < steps; i++) {
    const a0 = (i / steps) * Math.PI * 2
    const a1 = ((i + 1) / steps) * Math.PI * 2
    out.push([
      cx + r * Math.cos(a0), cy + r * Math.sin(a0),
      cx + r * Math.cos(a1), cy + r * Math.sin(a1),
    ])
  }
  return out
}

function arcSegs(cx, cy, r, a0, a1, steps = 24) {
  const out = []
  for (let i = 0; i < steps; i++) {
    const t0 = a0 + ((a1 - a0) * i) / steps
    const t1 = a0 + ((a1 - a0) * (i + 1)) / steps
    out.push([
      cx + r * Math.cos(t0), cy + r * Math.sin(t0),
      cx + r * Math.cos(t1), cy + r * Math.sin(t1),
    ])
  }
  return out
}

const ICONS = {
  dashboard: [
    ...rectSegs(3, 3, 7.5, 7.5),
    ...rectSegs(13.5, 3, 7.5, 7.5),
    ...rectSegs(3, 13.5, 7.5, 7.5),
    ...rectSegs(13.5, 13.5, 7.5, 7.5),
  ],
  order: [
    ...rectSegs(5, 3, 14, 18),
    [9, 8, 15, 8],
    [9, 12, 15, 12],
    [9, 16, 13, 16],
  ],
  chat: [
    ...rectSegs(3, 5, 18, 12),
    [7, 17, 7, 20],
    [7, 20, 12, 17],
    [8, 10, 16, 10],
    [8, 13, 13, 13],
  ],
  me: [
    ...circleSegs(12, 8, 3.5),
    ...arcSegs(12, 20, 7, Math.PI, Math.PI * 2),
  ],
}

/* ---------- 光栅化 ---------- */
function distToSeg(px, py, s) {
  const dx = s[2] - s[0]
  const dy = s[3] - s[1]
  const len2 = dx * dx + dy * dy
  let t = len2 ? ((px - s[0]) * dx + (py - s[1]) * dy) / len2 : 0
  t = Math.max(0, Math.min(1, t))
  const qx = s[0] + t * dx
  const qy = s[1] + t * dy
  return Math.hypot(px - qx, py - qy)
}

function renderIdat(segs, rgb) {
  const pxSegs = segs.map((s) => s.map((v) => PAD + v * SCALE))
  const raw = Buffer.alloc(SIZE * (1 + SIZE * 4))
  const half = STROKE / 2
  for (let y = 0; y < SIZE; y++) {
    const rowStart = y * (1 + SIZE * 4)
    raw[rowStart] = 0 // filter: none
    for (let x = 0; x < SIZE; x++) {
      const px = x + 0.5
      const py = y + 0.5
      let d = Infinity
      for (const s of pxSegs) {
        const dd = distToSeg(px, py, s)
        if (dd < d) d = dd
      }
      // 1.2px 过渡带做简单抗锯齿
      const alpha = Math.max(0, Math.min(1, (half + 0.6 - d) / 1.2))
      const o = rowStart + 1 + x * 4
      raw[o] = rgb[0]
      raw[o + 1] = rgb[1]
      raw[o + 2] = rgb[2]
      raw[o + 3] = Math.round(alpha * 255)
    }
  }
  return zlib.deflateSync(raw, { level: 9 })
}

/* ---------- 最小 PNG 编码器 ---------- */
const CRC_TABLE = (() => {
  const table = new Uint32Array(256)
  for (let n = 0; n < 256; n++) {
    let c = n
    for (let k = 0; k < 8; k++) c = c & 1 ? 0xedb88320 ^ (c >>> 1) : c >>> 1
    table[n] = c >>> 0
  }
  return table
})()

function crc32(buf) {
  let c = 0xffffffff
  for (const b of buf) c = CRC_TABLE[(c ^ b) & 0xff] ^ (c >>> 8)
  return (c ^ 0xffffffff) >>> 0
}

function chunk(type, data) {
  const len = Buffer.alloc(4)
  len.writeUInt32BE(data.length, 0)
  const t = Buffer.from(type, 'ascii')
  const crc = Buffer.alloc(4)
  crc.writeUInt32BE(crc32(Buffer.concat([t, data])), 0)
  return Buffer.concat([len, t, data, crc])
}

function buildPng(idat) {
  const sig = Buffer.from([137, 80, 78, 71, 13, 10, 26, 10])
  const ihdr = Buffer.alloc(13)
  ihdr.writeUInt32BE(SIZE, 0)
  ihdr.writeUInt32BE(SIZE, 4)
  ihdr[8] = 8 // bit depth
  ihdr[9] = 6 // RGBA
  return Buffer.concat([
    sig,
    chunk('IHDR', ihdr),
    chunk('IDAT', idat),
    chunk('IEND', Buffer.alloc(0)),
  ])
}

/* ---------- 输出 ---------- */
const COLORS = {
  normal: [0x98, 0xa2, 0xb3],
  active: [0x1d, 0x4e, 0xd8],
}

fs.mkdirSync(OUT_DIR, { recursive: true })
let count = 0
for (const [name, segs] of Object.entries(ICONS)) {
  for (const [state, rgb] of Object.entries(COLORS)) {
    const file = path.join(OUT_DIR, `${name}${state === 'active' ? '-active' : ''}.png`)
    fs.writeFileSync(file, buildPng(renderIdat(segs, rgb)))
    count++
  }
}
console.log(`generated ${count} tabbar icons -> ${OUT_DIR}`)
