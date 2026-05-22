<script setup>
import { ref, onMounted } from 'vue'
import { useData } from 'vitepress'
const { site } = useData()
const base = site.value.base || '/'

// Animated counters
const stats = ref([
  { label: '插件模块', value: 0, target: 15, suffix: '+' },
  { label: '驱动接口', value: 0, target: 5, suffix: '种' },
  { label: '支持平台', value: 0, target: 4, suffix: '端' },
  { label: '开源协议', value: 0, target: 0, suffix: 'MIT', isText: true },
])

onMounted(() => {
  stats.value.forEach((s, i) => {
    if (s.isText) { s.value = 1; return }
    let current = 0
    const step = Math.ceil(s.target / 30)
    const timer = setInterval(() => {
      current += step
      if (current >= s.target) { current = s.target; clearInterval(timer) }
      stats.value[i].value = current
    }, 40)
  })
})

const features = [
  { icon: '🛒', title: '商品管理', desc: '多规格 SKU、分类、相册、富文本详情、AI 自动生图', color: '#ef4444' },
  { icon: '📦', title: '订单系统', desc: '购物车、状态机、发货、退款售后完整链路', color: '#f97316' },
  { icon: '🏷️', title: '营销活动', desc: '优惠券、限时秒杀、满减、积分，Redis 原子扣减', color: '#eab308' },
  { icon: '💬', title: 'IM 客服', desc: 'WebSocket 实时通信、多坐席、离线消息、自动回复', color: '#22c55e' },
  { icon: '🏭', title: '仓储管理', desc: '多仓库、入出库单、调拨盘点、库存流水追溯', color: '#06b6d4' },
  { icon: '🎨', title: 'AI 生图', desc: '通义万象 / 文心 / DALL-E 多模型聚合生成商品图', color: '#ef4444' },
  { icon: '🖌️', title: '店铺装修', desc: '可视化拖拽编辑器，9 种组件，JSON 驱动渲染', color: '#ec4899' },
  { icon: '🔌', title: '驱动抽象层', desc: '支付 / 短信 / OAuth / 存储 / AI 统一接口', color: '#dc2626' },
]

const techStack = [
  { name: 'Go', desc: '后端服务', icon: '⚡' },
  { name: 'Gin', desc: 'HTTP 框架', icon: '🌐' },
  { name: 'GORM', desc: 'ORM', icon: '🗃️' },
  { name: 'Vue 3', desc: '前端框架', icon: '💚' },
  { name: 'uni-app', desc: '多端开发', icon: '📱' },
  { name: 'TailwindCSS', desc: '原子化样式', icon: '🎨' },
  { name: 'UnoCSS', desc: '即时原子CSS', icon: '⚛️' },
  { name: 'MySQL', desc: '关系数据库', icon: '🐬' },
  { name: 'Redis', desc: '缓存 & 队列', icon: '🔴' },
  { name: 'Docker', desc: '容器化部署', icon: '🐳' },
]

const platforms = [
  { name: 'PC Web', tech: 'Vue3 + UnoCSS', desc: '桌面端商城，完整购物体验', icon: '🖥️' },
  { name: 'H5 移动端', tech: 'uni-app', desc: '手机浏览器，响应式适配', icon: '📱' },
  { name: '微信小程序', tech: 'uni-app', desc: '微信生态，无需安装', icon: '💬' },
  { name: 'iOS & Android', tech: 'uni-app', desc: '原生 App，极致体验', icon: '📲' },
  { name: '管理后台', tech: 'Vue3 + TailwindCSS', desc: '商家运营管理中心', icon: '⚙️' },
]
</script>

<template>
<div class="lyshop-home">

  <!-- Hero -->
  <section class="hero">
    <div class="hero-bg"></div>
    <div class="hero-content">
      <div class="hero-badge">
        <span class="hero-badge-dot"></span>
        <span>开源免费 · MIT 协议 · 商用无忧</span>
      </div>
      <h1 class="hero-title">
        <span class="hero-title-brand">LYShop</span>
        <br/>零云商城
      </h1>
      <p class="hero-subtitle">
        基于 <strong>Go + Vue3 + uni-app</strong> 的开源插件化商城系统<br/>
        支持 PC / H5 / 小程序 / App 四端，AI 生图，IM 客服，Docker 一键部署
      </p>
      <div class="hero-actions">
        <a :href="`${base}guide/features`" class="hero-btn hero-btn-primary">
          <span>快速开始</span>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
        </a>
        <a :href="`${base}web-demo/index.html`" target="_blank" class="hero-btn hero-btn-outline">
          <span>在线演示</span>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6M15 3h6v6M10 14L21 3"/></svg>
        </a>
        <a href="https://github.com/ijry/lyshop" target="_blank" class="hero-btn hero-btn-ghost">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
          <span>GitHub</span>
        </a>
      </div>

      <!-- Stats -->
      <div class="hero-stats">
        <div v-for="s in stats" :key="s.label" class="hero-stat">
          <div class="hero-stat-value">
            <template v-if="s.isText">{{ s.suffix }}</template>
            <template v-else>{{ s.value }}{{ s.suffix }}</template>
          </div>
          <div class="hero-stat-label">{{ s.label }}</div>
        </div>
      </div>
    </div>
  </section>

  <!-- Features -->
  <section class="section">
    <div class="container">
      <div class="section-header">
        <span class="section-tag">核心功能</span>
        <h2 class="section-title">覆盖电商全链路的插件化能力</h2>
        <p class="section-desc">每个功能模块都是独立插件，通过配置文件一行开关，按需组合</p>
      </div>
      <div class="features-grid">
        <div v-for="f in features" :key="f.title" class="feature-card">
          <div class="feature-icon" :style="{ background: f.color + '15', color: f.color }">
            {{ f.icon }}
          </div>
          <h3 class="feature-title">{{ f.title }}</h3>
          <p class="feature-desc">{{ f.desc }}</p>
        </div>
      </div>
    </div>
  </section>

  <!-- Multi-platform -->
  <section class="section section-alt">
    <div class="container">
      <div class="section-header">
        <span class="section-tag">多端覆盖</span>
        <h2 class="section-title">一套后端 API，五端前端</h2>
        <p class="section-desc">共享同一个 Go 后端和插件系统，前端按平台独立构建</p>
      </div>
      <div class="platforms-grid">
        <div v-for="p in platforms" :key="p.name" class="platform-card">
          <div class="platform-icon">{{ p.icon }}</div>
          <h3 class="platform-name">{{ p.name }}</h3>
          <div class="platform-tech">{{ p.tech }}</div>
          <p class="platform-desc">{{ p.desc }}</p>
        </div>
      </div>
    </div>
  </section>

  <!-- Tech Stack -->
  <section class="section">
    <div class="container">
      <div class="section-header">
        <span class="section-tag">技术栈</span>
        <h2 class="section-title">现代化技术选型</h2>
        <p class="section-desc">每一层都选用生态最成熟、性能最优的技术方案</p>
      </div>
      <div class="tech-grid">
        <div v-for="t in techStack" :key="t.name" class="tech-card">
          <span class="tech-icon">{{ t.icon }}</span>
          <span class="tech-name">{{ t.name }}</span>
          <span class="tech-desc">{{ t.desc }}</span>
        </div>
      </div>
    </div>
  </section>

  <!-- Architecture -->
  <section class="section section-alt">
    <div class="container">
      <div class="section-header">
        <span class="section-tag">架构设计</span>
        <h2 class="section-title">分层单体 + 插件注册表</h2>
        <p class="section-desc">单一 Go 二进制部署，插件通过 config.yaml 开关，驱动层抽象统一接口</p>
      </div>
      <div class="arch-visual">
        <div class="arch-row">
          <div class="arch-box arch-client">🖥️ PC Web</div>
          <div class="arch-box arch-client">📱 H5 / App</div>
          <div class="arch-box arch-client">⚙️ Admin</div>
        </div>
        <div class="arch-arrow">▼</div>
        <div class="arch-row">
          <div class="arch-box arch-server">
            <strong>Go Server (Gin)</strong>
            <div class="arch-inner">
              <span>product</span><span>order</span><span>marketing</span><span>im</span><span>wms</span><span>ai_image</span><span>decor</span>
            </div>
            <div class="arch-inner arch-driver">
              <span>payment</span><span>sms</span><span>oauth</span><span>storage</span><span>ai</span>
            </div>
          </div>
        </div>
        <div class="arch-arrow">▼</div>
        <div class="arch-row">
          <div class="arch-box arch-infra">🐬 MySQL</div>
          <div class="arch-box arch-infra">🔴 Redis</div>
          <div class="arch-box arch-infra">📁 OSS</div>
        </div>
      </div>
    </div>
  </section>

  <!-- CTA -->
  <section class="cta">
    <div class="container">
      <h2 class="cta-title">开始构建你的商城</h2>
      <p class="cta-desc">MIT 协议开源，免费商用，完善文档，活跃社区</p>
      <div class="cta-actions">
        <a :href="`${base}deploy/`" class="hero-btn hero-btn-primary">
          <span>部署文档</span>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
        </a>
        <a href="https://github.com/ijry/lyshop" target="_blank" class="hero-btn hero-btn-outline" style="border-color:rgba(255,255,255,0.3);color:#fff;">
          <span>Star on GitHub</span>
        </a>
      </div>
    </div>
  </section>

</div>
</template>

<style scoped>
.lyshop-home { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }

/* Hero */
.hero { position: relative; overflow: hidden; padding: 100px 24px 80px; text-align: center; }
.hero-bg { position: absolute; inset: 0; background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%); z-index: 0; }
.hero-bg::after { content: ''; position: absolute; inset: 0; background: radial-gradient(circle at 30% 40%, rgba(220,38,38,0.15) 0%, transparent 50%), radial-gradient(circle at 70% 60%, rgba(236,72,153,0.1) 0%, transparent 50%); }
.hero-content { position: relative; z-index: 1; max-width: 800px; margin: 0 auto; }
.hero-badge { display: inline-flex; align-items: center; gap: 8px; background: rgba(255,255,255,0.08); backdrop-filter: blur(8px); border: 1px solid rgba(255,255,255,0.1); border-radius: 999px; padding: 6px 18px; margin-bottom: 32px; font-size: 13px; color: rgba(255,255,255,0.7); }
.hero-badge-dot { width: 6px; height: 6px; background: #22c55e; border-radius: 50%; animation: pulse 2s infinite; }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.4; } }
.hero-title { font-size: 56px; font-weight: 800; line-height: 1.15; color: #fff; margin-bottom: 24px; letter-spacing: -1px; }
.hero-title-brand { background: linear-gradient(135deg, #f87171, #ef4444, #fca5a5); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }
.hero-subtitle { font-size: 18px; line-height: 1.7; color: rgba(255,255,255,0.6); margin-bottom: 40px; }
.hero-subtitle strong { color: rgba(255,255,255,0.9); }
.hero-actions { display: flex; gap: 12px; justify-content: center; flex-wrap: wrap; margin-bottom: 60px; }
.hero-btn { display: inline-flex; align-items: center; gap: 8px; padding: 12px 28px; border-radius: 12px; font-size: 15px; font-weight: 600; text-decoration: none; transition: all 0.2s; cursor: pointer; }
.hero-btn-primary { background: linear-gradient(135deg, #dc2626, #ef4444); color: #fff; box-shadow: 0 4px 20px rgba(220,38,38,0.4); }
.hero-btn-primary:hover { transform: translateY(-2px); box-shadow: 0 8px 30px rgba(220,38,38,0.5); }
.hero-btn-outline { background: rgba(255,255,255,0.05); color: #fff; border: 1px solid rgba(255,255,255,0.2); backdrop-filter: blur(4px); }
.hero-btn-outline:hover { background: rgba(255,255,255,0.1); border-color: rgba(255,255,255,0.3); }
.hero-btn-ghost { background: transparent; color: rgba(255,255,255,0.7); border: 1px solid rgba(255,255,255,0.12); }
.hero-btn-ghost:hover { color: #fff; border-color: rgba(255,255,255,0.25); }

/* Stats */
.hero-stats { display: flex; justify-content: center; gap: 48px; flex-wrap: wrap; }
.hero-stat { text-align: center; }
.hero-stat-value { font-size: 36px; font-weight: 800; color: #fff; }
.hero-stat-label { font-size: 13px; color: rgba(255,255,255,0.45); margin-top: 4px; }

/* Sections */
.section { padding: 80px 24px; }
.section-alt { background: #f8fafc; }
.container { max-width: 1100px; margin: 0 auto; }
.section-header { text-align: center; margin-bottom: 48px; }
.section-tag { display: inline-block; font-size: 13px; font-weight: 600; color: #dc2626; background: #fef2f2; padding: 4px 14px; border-radius: 999px; margin-bottom: 16px; }
.section-title { font-size: 32px; font-weight: 800; color: #0f172a; margin-bottom: 12px; letter-spacing: -0.5px; }
.section-desc { font-size: 16px; color: #64748b; max-width: 600px; margin: 0 auto; line-height: 1.6; }

/* Features */
.features-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; }
@media (max-width: 900px) { .features-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 500px) { .features-grid { grid-template-columns: 1fr; } }
.feature-card { background: #fff; border: 1px solid #f1f5f9; border-radius: 16px; padding: 28px; transition: all 0.25s; }
.feature-card:hover { transform: translateY(-4px); box-shadow: 0 12px 40px rgba(0,0,0,0.06); border-color: #e2e8f0; }
.feature-icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 22px; margin-bottom: 16px; }
.feature-title { font-size: 16px; font-weight: 700; color: #0f172a; margin-bottom: 8px; }
.feature-desc { font-size: 13px; color: #64748b; line-height: 1.6; }

/* Platforms */
.platforms-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 16px; }
@media (max-width: 900px) { .platforms-grid { grid-template-columns: repeat(3, 1fr); } }
@media (max-width: 500px) { .platforms-grid { grid-template-columns: repeat(2, 1fr); } }
.platform-card { background: #fff; border: 1px solid #f1f5f9; border-radius: 16px; padding: 24px 16px; text-align: center; transition: all 0.25s; }
.platform-card:hover { transform: translateY(-4px); box-shadow: 0 12px 40px rgba(0,0,0,0.06); }
.platform-icon { font-size: 36px; margin-bottom: 12px; }
.platform-name { font-size: 15px; font-weight: 700; color: #0f172a; margin-bottom: 4px; }
.platform-tech { font-size: 11px; color: #dc2626; background: #fef2f2; padding: 2px 10px; border-radius: 999px; display: inline-block; margin-bottom: 8px; }
.platform-desc { font-size: 12px; color: #94a3b8; line-height: 1.5; }

/* Tech */
.tech-grid { display: flex; flex-wrap: wrap; gap: 12px; justify-content: center; }
.tech-card { display: flex; align-items: center; gap: 10px; background: #fff; border: 1px solid #f1f5f9; border-radius: 12px; padding: 12px 20px; transition: all 0.2s; }
.tech-card:hover { border-color: #e2e8f0; box-shadow: 0 4px 16px rgba(0,0,0,0.04); }
.tech-icon { font-size: 20px; }
.tech-name { font-size: 14px; font-weight: 700; color: #0f172a; }
.tech-desc { font-size: 12px; color: #94a3b8; }

/* Architecture */
.arch-visual { max-width: 700px; margin: 0 auto; }
.arch-row { display: flex; gap: 12px; justify-content: center; }
.arch-arrow { text-align: center; font-size: 20px; color: #94a3b8; padding: 8px 0; }
.arch-box { padding: 16px 24px; border-radius: 12px; font-size: 14px; font-weight: 600; text-align: center; }
.arch-client { background: #fef2f2; color: #991b1b; border: 1px solid #fecaca; flex: 1; }
.arch-server { background: #0f172a; color: #fff; flex: 1; padding: 20px 24px; }
.arch-inner { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 10px; justify-content: center; }
.arch-inner span { background: rgba(255,255,255,0.12); padding: 3px 10px; border-radius: 6px; font-size: 12px; font-weight: 500; }
.arch-driver span { background: rgba(220,38,38,0.25); }
.arch-infra { background: #f0fdf4; color: #166534; border: 1px solid #bbf7d0; flex: 1; }

/* CTA */
.cta { background: linear-gradient(135deg, #0f172a, #1e293b); padding: 80px 24px; text-align: center; }
.cta-title { font-size: 36px; font-weight: 800; color: #fff; margin-bottom: 12px; }
.cta-desc { font-size: 16px; color: rgba(255,255,255,0.5); margin-bottom: 32px; }
.cta-actions { display: flex; gap: 12px; justify-content: center; flex-wrap: wrap; }

/* Dark mode */
.dark .section:not(.section-alt) { background: var(--vp-c-bg); }
.dark .section-alt { background: var(--vp-c-bg-soft); }
.dark .section-title { color: var(--vp-c-text-1); }
.dark .section-desc { color: var(--vp-c-text-2); }
.dark .feature-card, .dark .platform-card, .dark .tech-card { background: var(--vp-c-bg-soft); border-color: var(--vp-c-divider); }
.dark .feature-title, .dark .platform-name, .dark .tech-name { color: var(--vp-c-text-1); }
.dark .feature-desc, .dark .platform-desc, .dark .tech-desc { color: var(--vp-c-text-2); }
</style>
