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

const merchantAppShots = [
  { title: '商家工作台', image: 'eapp-dashboard.png', desc: '今日经营数据、待办事项与核心入口集中呈现' },
  { title: '订单管理', image: 'eapp-orders.png', desc: '订单状态筛选、快速进入发货与售后处理流程' },
  { title: '商品管理', image: 'eapp-products.png', desc: '商品搜索、上下架切换与高频编辑能力' },
]

const merchantFunctions = [
  { title: '一屏经营概览', desc: '在移动端快速查看订单、销售额、待办和库存预警。' },
  { title: '履约与售后联动', desc: '订单发货、物流同步、售后审核和退款动作打通。' },
  { title: '商品运营闭环', desc: '商品上下架、编辑与营销活动挂接统一在商家端处理。' },
  { title: '统一后台接口', desc: '直接复用 /admin/api/*，降低二次开发与维护成本。' },
]

const marketingAppShots = [
  { title: '营销首页', image: 'eapp-marketing-index.png', desc: '优惠券、秒杀、拼团、砍价功能入口统一收拢' },
  { title: '秒杀活动管理', image: 'eapp-marketing-seckill.png', desc: '活动切换、活动商品维护与状态管理移动化操作' },
]

const marketingFunctions = [
  { title: '活动创建与维护', desc: '可在商家 App 快速创建活动并调整活动时间和状态。' },
  { title: '活动商品管理', desc: '按活动绑定商品 SKU，并维护活动价、库存与限购。' },
  { title: '多营销类型统一', desc: '秒杀、拼团、砍价复用统一的移动端管理交互。' },
  { title: '运营即时响应', desc: '不依赖 PC，在手机端即可完成营销日常调整与巡检。' },
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
  { name: '商家 eapp', desc: '商家移动端运营工作台演示', href: 'eapp-demo/index.html#/pages/dashboard/index' },
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
          <span class="section-kicker">产品界面预览</span>
          <h2>首页直接展示可运行的商城能力</h2>
        </div>
        <p>
          PC 商城与移动端商城覆盖商品浏览、营销入口、订单管理和客服触达等核心流程，让访问者可以直观看到系统落地后的业务界面。
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

    <section class="section merchant-app-section">
      <div class="container">
        <div class="section-heading">
          <span class="section-kicker">商家 App 截图与功能</span>
          <h2>商家移动端 eapp：掌上运营工作台</h2>
          <p>围绕单店铺日常运营场景，把工作台、订单、商品等高频管理能力落到手机端。</p>
        </div>
        <div class="merchant-layout">
          <div class="merchant-shot-grid">
            <article v-for="item in merchantAppShots" :key="item.title" class="merchant-shot-card">
              <div class="phone-shot">
                <img :src="`${base}showcase/${item.image}`" :alt="`${item.title}截图`" />
              </div>
              <div class="showcase-copy">
                <h3>{{ item.title }}</h3>
                <p>{{ item.desc }}</p>
              </div>
            </article>
          </div>
          <aside class="merchant-feature-panel">
            <h3>商家端核心能力</h3>
            <ul class="feature-list">
              <li v-for="item in merchantFunctions" :key="item.title">
                <strong>{{ item.title }}</strong>
                <p>{{ item.desc }}</p>
              </li>
            </ul>
            <a class="btn primary" :href="`${base}guide/eapp-merchant`">查看商家 eapp 文档</a>
          </aside>
        </div>
      </div>
    </section>

    <section class="section marketing-app-section">
      <div class="container">
        <div class="section-heading">
          <span class="section-kicker">营销功能 App 端截图</span>
          <h2>营销能力在商家 App 的移动化呈现</h2>
          <p>从营销入口到活动商品配置，关键营销动作可在商家移动端直接完成。</p>
        </div>
        <div class="marketing-shot-grid">
          <article v-for="item in marketingAppShots" :key="item.title" class="marketing-shot-card">
            <div class="phone-shot">
              <img :src="`${base}showcase/${item.image}`" :alt="`${item.title}截图`" />
            </div>
            <div class="showcase-copy">
              <h3>{{ item.title }}</h3>
              <p>{{ item.desc }}</p>
            </div>
          </article>
        </div>
        <div class="marketing-feature-grid">
          <article v-for="item in marketingFunctions" :key="item.title" class="capability-card">
            <h3>{{ item.title }}</h3>
            <p>{{ item.desc }}</p>
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

.merchant-app-section {
  background: #fff;
}

.marketing-app-section {
  background: #f8fafc;
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

.merchant-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(280px, 0.8fr);
  gap: 22px;
  align-items: start;
}

.merchant-shot-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.merchant-shot-card,
.marketing-shot-card {
  border: 1px solid var(--ly-line);
  border-radius: 12px;
  background: #fff;
  padding: 16px;
  box-shadow: 0 12px 28px rgba(15, 23, 42, 0.06);
}

.merchant-feature-panel {
  border: 1px solid var(--ly-line);
  border-radius: 12px;
  background: #f8fafc;
  padding: 22px;
}

.merchant-feature-panel h3 {
  margin: 0;
  font-size: 20px;
}

.feature-list {
  margin: 16px 0 20px;
  padding: 0;
  list-style: none;
  display: grid;
  gap: 12px;
}

.feature-list li {
  padding: 12px;
  border-radius: 10px;
  border: 1px solid #e5e7eb;
  background: #fff;
}

.feature-list li strong {
  display: block;
  font-size: 15px;
  color: var(--ly-ink);
}

.feature-list li p {
  margin: 8px 0 0;
  color: var(--ly-muted);
  font-size: 13px;
  line-height: 1.6;
}

.marketing-shot-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.marketing-feature-grid {
  margin-top: 16px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
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
  .merchant-shot-grid,
  .marketing-shot-grid,
  .capability-grid,
  .channel-grid,
  .marketing-feature-grid {
    grid-template-columns: 1fr;
  }

  .merchant-layout {
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

  .merchant-shot-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
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

  .merchant-shot-grid {
    grid-template-columns: 1fr;
  }

  .phone-shot {
    max-width: 180px;
  }
}
</style>
