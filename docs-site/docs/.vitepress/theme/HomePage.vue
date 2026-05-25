<script setup>
import { computed } from 'vue'
import { useData } from 'vitepress'

const { site } = useData()
const base = computed(() => site.value.base || '/')

const metrics = [
  { value: '4+', label: '前端交付形态' },
  { value: '15+', label: '可插拔业务模块' },
  { value: '5类', label: '支付/短信/存储等驱动' },
  { value: 'MIT', label: '开源协议' },
]

const showcase = [
  { title: 'H5 商城首页', image: 'h5-home.png', desc: '装修组件、营销入口、推荐商品聚合展示' },
  { title: '商品列表', image: 'h5-products.png', desc: '分类筛选、商品瀑布流与移动端购买路径' },
  { title: '订单中心', image: 'h5-orders.png', desc: '订单状态、售后入口与履约进度聚合' },
]

const capabilities = [
  { title: '商品与交易', desc: '商品、SKU、购物车、订单、支付、发货与售后串成完整电商闭环。' },
  { title: '营销增长', desc: '优惠券、秒杀、拼团、砍价、积分与签到按插件组织，便于按业务阶段启用。' },
  { title: '运营后台', desc: '后台覆盖商品、订单、仓储、消息、装修、评价和系统配置，适合真实运营场景。' },
  { title: 'AI 与内容', desc: 'AI 生图驱动适配商品素材生产，店铺装修通过 JSON 数据渲染到多端。' },
  { title: '客服与消息', desc: 'IM 客服、站内信和订单通知为用户服务与复购触达提供基础设施。' },
  { title: '开放驱动层', desc: '支付、短信、OAuth、存储、物流、AI 驱动统一抽象，方便替换云服务厂商。' },
]

const channels = [
  { name: 'PC 商城', desc: '面向桌面端用户的完整商城体验', href: 'web-demo/index.html' },
  { name: 'H5 / App / 小程序', desc: '基于 uni-app 复用移动端购物链路', href: 'demo/index.html' },
  { name: '管理后台', desc: '面向商家运营与系统管理', href: 'guide/features' },
]
</script>

<template>
  <main class="lyshop-home">
    <section class="hero-shell">
      <div class="hero-inner">
        <div class="hero-copy">
          <div class="eyebrow">开源插件化多端商城系统</div>
          <h1>LYShop 零云商城</h1>
          <p class="hero-lead">
            基于 Go、Vue 3 与 uni-app 构建，覆盖 PC 商城、H5、小程序、App 与管理后台，帮助团队更快交付可二次开发、可部署、可运营的电商系统。
          </p>
          <div class="hero-actions">
            <a class="btn primary" :href="`${base}guide/features`">查看功能</a>
            <a class="btn secondary" :href="`${base}web-demo/index.html`" target="_blank" rel="noreferrer">打开 PC 演示</a>
          </div>
          <div class="metric-row" aria-label="项目能力指标">
            <div v-for="item in metrics" :key="item.label" class="metric-item">
              <strong>{{ item.value }}</strong>
              <span>{{ item.label }}</span>
            </div>
          </div>
        </div>

        <div class="hero-visual" aria-label="LYShop 演示截图">
          <div class="browser-frame">
            <div class="browser-bar">
              <span></span><span></span><span></span>
              <em>PC Web Storefront</em>
            </div>
            <img :src="`${base}showcase/web-home.png`" alt="LYShop PC 商城首页截图" />
          </div>
          <div class="phone-frame">
            <img :src="`${base}showcase/h5-home.png`" alt="LYShop H5 商城首页截图" />
          </div>
        </div>
      </div>
    </section>

    <section class="section intro-band">
      <div class="container split">
        <div>
          <span class="section-kicker">真实演示截图</span>
          <h2>首页直接展示可运行的商城能力</h2>
        </div>
        <p>
          截图由 docs-site 脚本从演示模式自动生成，不依赖远程服务。首页展示的不是概念图，而是当前仓库里可构建、可访问的 PC 与 H5 演示界面。
        </p>
      </div>
    </section>

    <section class="section showcase-section">
      <div class="container">
        <div class="showcase-grid">
          <article v-for="item in showcase" :key="item.title" class="showcase-card">
            <div class="phone-shot">
              <img :src="`${base}showcase/${item.image}`" :alt="`${item.title}截图`" />
            </div>
            <div class="showcase-copy">
              <h3>{{ item.title }}</h3>
              <p>{{ item.desc }}</p>
            </div>
          </article>
        </div>
      </div>
    </section>

    <section class="section capability-section">
      <div class="container">
        <div class="section-heading">
          <span class="section-kicker">核心能力</span>
          <h2>从交易闭环到运营增长，一套系统覆盖商城主流程</h2>
          <p>LYShop 以插件组织业务能力，以驱动层隔离第三方服务，适合开源学习、项目二开和中小型商城快速落地。</p>
        </div>
        <div class="capability-grid">
          <article v-for="item in capabilities" :key="item.title" class="capability-card">
            <h3>{{ item.title }}</h3>
            <p>{{ item.desc }}</p>
          </article>
        </div>
      </div>
    </section>

    <section class="section architecture-section">
      <div class="container architecture-layout">
        <div class="architecture-copy">
          <span class="section-kicker">技术架构</span>
          <h2>分层单体部署，插件化扩展</h2>
          <p>
            后端以 Go 单体服务承载核心接口，通过插件注册表启停业务模块；前端按 PC、移动端和后台独立构建，共用同一套 API 与 mock 演示数据。
          </p>
          <div class="stack-list">
            <span>Go</span>
            <span>Gin</span>
            <span>GORM</span>
            <span>Vue 3</span>
            <span>uni-app</span>
            <span>MySQL</span>
            <span>Redis</span>
            <span>Docker</span>
          </div>
        </div>
        <div class="architecture-map" aria-label="LYShop 架构示意">
          <div class="map-row">
            <span>PC Web</span>
            <span>H5 / 小程序 / App</span>
            <span>Admin</span>
          </div>
          <div class="map-core">
            <strong>Go Server</strong>
            <div>
              <span>product</span><span>order</span><span>marketing</span><span>decor</span><span>im</span><span>wms</span>
            </div>
          </div>
          <div class="map-row">
            <span>Payment</span>
            <span>Storage</span>
            <span>SMS</span>
            <span>AI</span>
          </div>
        </div>
      </div>
    </section>

    <section class="section channel-section">
      <div class="container">
        <div class="section-heading">
          <span class="section-kicker">交付入口</span>
          <h2>面向用户、商家和开发者的完整入口</h2>
        </div>
        <div class="channel-grid">
          <a v-for="item in channels" :key="item.name" class="channel-card" :href="`${base}${item.href}`">
            <strong>{{ item.name }}</strong>
            <span>{{ item.desc }}</span>
          </a>
        </div>
      </div>
    </section>

    <section class="final-cta">
      <div class="container final-cta-inner">
        <div>
          <span class="section-kicker">部署与二次开发</span>
          <h2>用开源商城底座启动你的业务系统</h2>
          <p>支持 Docker 部署、插件扩展和多端演示构建，适合直接试用、学习架构或作为商业项目二开基础。</p>
        </div>
        <div class="hero-actions">
          <a class="btn primary light" :href="`${base}deploy/`">部署文档</a>
          <a class="btn secondary light" href="https://github.com/ijry/lyshop" target="_blank" rel="noreferrer">GitHub 仓库</a>
        </div>
      </div>
    </section>
  </main>
</template>

<style scoped>
.lyshop-home {
  --ly-bg: #f7fafc;
  --ly-panel: #ffffff;
  --ly-ink: #0f172a;
  --ly-muted: #475569;
  --ly-line: #e2e8f0;
  --ly-brand: #b91c1c;
  --ly-blue: #0369a1;
  font-family: Inter, "PingFang SC", "Microsoft YaHei", system-ui, sans-serif;
  color: var(--ly-ink);
  background: var(--ly-bg);
}

.hero-shell {
  position: relative;
  overflow: hidden;
  padding: 92px 24px 72px;
  background:
    linear-gradient(135deg, rgba(15, 23, 42, 0.96), rgba(30, 41, 59, 0.96)),
    linear-gradient(90deg, #0f172a, #172554);
}

.hero-inner,
.container {
  max-width: 1180px;
  margin: 0 auto;
}

.hero-inner {
  display: grid;
  grid-template-columns: minmax(0, 0.9fr) minmax(420px, 1.1fr);
  gap: 48px;
  align-items: center;
}

.eyebrow,
.section-kicker {
  display: inline-flex;
  align-items: center;
  width: fit-content;
  min-height: 28px;
  padding: 5px 12px;
  border: 1px solid rgba(185, 28, 28, 0.18);
  border-radius: 999px;
  background: rgba(185, 28, 28, 0.08);
  color: #991b1b;
  font-size: 13px;
  font-weight: 700;
}

.hero-shell .eyebrow {
  border-color: rgba(248, 113, 113, 0.28);
  background: rgba(248, 113, 113, 0.12);
  color: #fecaca;
}

.hero-copy h1 {
  margin: 22px 0 18px;
  color: #fff;
  font-size: 58px;
  line-height: 1.08;
  font-weight: 850;
  letter-spacing: 0;
}

.hero-lead {
  margin: 0;
  max-width: 620px;
  color: rgba(255, 255, 255, 0.76);
  font-size: 18px;
  line-height: 1.8;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 32px;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 46px;
  padding: 0 20px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 750;
  text-decoration: none;
  transition: background-color 0.2s, border-color 0.2s, color 0.2s, box-shadow 0.2s;
}

.btn.primary {
  border: 1px solid #ef4444;
  background: #dc2626;
  color: #fff;
  box-shadow: 0 16px 34px rgba(220, 38, 38, 0.24);
}

.btn.primary:hover {
  background: #b91c1c;
  color: #fff;
}

.btn.secondary {
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.06);
  color: #fff;
}

.btn.secondary:hover {
  border-color: rgba(255, 255, 255, 0.4);
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.btn.light.secondary {
  border-color: rgba(255, 255, 255, 0.28);
}

.metric-row {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-top: 42px;
}

.metric-item {
  min-width: 0;
}

.metric-item strong {
  display: block;
  color: #fff;
  font-size: 30px;
  line-height: 1.1;
}

.metric-item span {
  display: block;
  margin-top: 6px;
  color: rgba(255, 255, 255, 0.56);
  font-size: 13px;
}

.hero-visual {
  position: relative;
  min-height: 520px;
}

.browser-frame,
.phone-frame,
.showcase-card,
.capability-card,
.channel-card,
.architecture-map {
  border: 1px solid var(--ly-line);
  background: var(--ly-panel);
  box-shadow: 0 20px 50px rgba(15, 23, 42, 0.08);
}

.browser-frame {
  overflow: hidden;
  border-radius: 14px;
}

.browser-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 38px;
  padding: 0 14px;
  border-bottom: 1px solid var(--ly-line);
  background: #f8fafc;
}

.browser-bar span {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #cbd5e1;
}

.browser-bar em {
  margin-left: 8px;
  color: #64748b;
  font-size: 12px;
  font-style: normal;
}

.browser-frame img {
  display: block;
  width: 100%;
  aspect-ratio: 16 / 10;
  object-fit: cover;
  object-position: top center;
}

.phone-frame {
  position: absolute;
  right: 22px;
  bottom: 0;
  width: 168px;
  overflow: hidden;
  border: 8px solid #0f172a;
  border-radius: 28px;
  background: #0f172a;
}

.phone-frame img,
.phone-shot img {
  display: block;
  width: 100%;
  aspect-ratio: 390 / 844;
  object-fit: cover;
  object-position: top center;
  background: #fff;
}

.section {
  padding: 76px 24px;
}

.split {
  display: grid;
  grid-template-columns: minmax(0, 0.9fr) minmax(0, 1.1fr);
  gap: 48px;
  align-items: end;
}

.split h2,
.section-heading h2,
.architecture-copy h2,
.final-cta h2 {
  margin: 16px 0 0;
  color: var(--ly-ink);
  font-size: 34px;
  line-height: 1.25;
  font-weight: 850;
  letter-spacing: 0;
}

.split p,
.section-heading p,
.architecture-copy p,
.final-cta p {
  margin: 0;
  color: var(--ly-muted);
  font-size: 16px;
  line-height: 1.8;
}

.intro-band {
  background: #fff;
}

.showcase-section {
  padding-top: 24px;
}

.showcase-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 22px;
}

.showcase-card {
  display: grid;
  grid-template-columns: 138px minmax(0, 1fr);
  gap: 18px;
  align-items: center;
  min-height: 260px;
  padding: 18px;
  border-radius: 12px;
}

.phone-shot {
  overflow: hidden;
  border: 7px solid #111827;
  border-radius: 24px;
  background: #111827;
}

.showcase-copy h3,
.capability-card h3 {
  margin: 0;
  color: var(--ly-ink);
  font-size: 18px;
  line-height: 1.35;
}

.showcase-copy p,
.capability-card p,
.channel-card span {
  margin: 10px 0 0;
  color: var(--ly-muted);
  font-size: 14px;
  line-height: 1.7;
}

.capability-section,
.channel-section {
  background: #fff;
}

.section-heading {
  max-width: 760px;
  margin-bottom: 34px;
}

.section-heading p {
  margin-top: 14px;
}

.capability-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.capability-card {
  min-height: 164px;
  padding: 24px;
  border-radius: 10px;
}

.architecture-section {
  background: #f1f5f9;
}

.architecture-layout {
  display: grid;
  grid-template-columns: minmax(0, 0.85fr) minmax(420px, 1.15fr);
  gap: 44px;
  align-items: center;
}

.stack-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 26px;
}

.stack-list span,
.map-row span,
.map-core span {
  display: inline-flex;
  align-items: center;
  min-height: 32px;
  padding: 0 12px;
  border-radius: 7px;
  background: #fff;
  color: #334155;
  font-size: 13px;
  font-weight: 700;
}

.architecture-map {
  padding: 24px;
  border-radius: 14px;
}

.map-row {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: center;
}

.map-core {
  margin: 20px 0;
  padding: 26px;
  border-radius: 12px;
  background: #0f172a;
  color: #fff;
  text-align: center;
}

.map-core strong {
  display: block;
  margin-bottom: 18px;
  font-size: 22px;
}

.map-core div {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: center;
}

.map-core span {
  background: rgba(255, 255, 255, 0.11);
  color: #fff;
}

.channel-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.channel-card {
  display: block;
  min-height: 128px;
  padding: 24px;
  border-radius: 10px;
  text-decoration: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.channel-card:hover {
  border-color: rgba(185, 28, 28, 0.32);
  box-shadow: 0 18px 42px rgba(15, 23, 42, 0.11);
}

.channel-card strong {
  display: block;
  color: var(--ly-ink);
  font-size: 18px;
}

.final-cta {
  padding: 76px 24px;
  background: #0f172a;
}

.final-cta-inner {
  display: flex;
  gap: 32px;
  align-items: center;
  justify-content: space-between;
}

.final-cta h2,
.final-cta p {
  color: #fff;
}

.final-cta p {
  max-width: 680px;
  margin-top: 14px;
  color: rgba(255, 255, 255, 0.72);
}

@media (max-width: 1040px) {
  .hero-inner,
  .architecture-layout {
    grid-template-columns: 1fr;
  }

  .hero-visual {
    min-height: 430px;
  }

  .showcase-grid,
  .capability-grid,
  .channel-grid {
    grid-template-columns: 1fr;
  }

  .showcase-card {
    grid-template-columns: 160px minmax(0, 1fr);
  }
}

@media (max-width: 720px) {
  .hero-shell {
    padding: 72px 18px 56px;
  }

  .hero-copy h1 {
    font-size: 42px;
  }

  .hero-lead {
    font-size: 16px;
  }

  .metric-row,
  .split {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .split {
    gap: 24px;
    align-items: start;
  }

  .section {
    padding: 56px 18px;
  }

  .split h2,
  .section-heading h2,
  .architecture-copy h2,
  .final-cta h2 {
    font-size: 28px;
  }

  .hero-visual {
    min-height: auto;
    padding-bottom: 0;
  }

  .phone-frame {
    display: none;
  }

  .showcase-card {
    grid-template-columns: 128px minmax(0, 1fr);
    min-height: 220px;
  }

  .final-cta-inner {
    display: block;
  }
}

@media (max-width: 460px) {
  .hero-actions,
  .metric-row,
  .split {
    grid-template-columns: 1fr;
  }

  .hero-actions {
    display: grid;
  }

  .btn {
    width: 100%;
  }

  .showcase-card {
    grid-template-columns: 1fr;
  }

  .phone-shot {
    max-width: 180px;
  }
}
</style>
