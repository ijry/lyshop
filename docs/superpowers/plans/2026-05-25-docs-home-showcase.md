# Docs Home Showcase Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build an automated docs-site showcase screenshot workflow and use the generated images in a more commercial LYShop homepage.

**Architecture:** Keep all changes inside `docs-site` plus documentation under `docs/superpowers`. A Node script serves built demo assets locally and uses Playwright to capture stable PNG files into VitePress public assets. The VitePress homepage reads those static assets and renders a polished product marketing page.

**Tech Stack:** VitePress, Vue 3, Node.js, Playwright, existing LYShop mock demos.

---

### Task 1: Screenshot Automation

**Files:**
- Create: `docs-site/scripts/capture-showcase.js`
- Modify: `docs-site/package.json`

- [ ] Add a Node script that starts a local static server for `docs-site/docs/public`.
- [ ] Use Playwright Chromium to capture:
  - `/demo/index.html#/pages/index/index` as `h5-home.png`
  - `/demo/index.html#/pages/product/list` as `h5-products.png`
  - `/demo/index.html#/pages/order/list` as `h5-orders.png`
  - `/web-demo/index.html#/` as `web-home.png`
- [ ] Add `showcase:capture` and `showcase:build` npm scripts.

### Task 2: Homepage Redesign

**Files:**
- Modify: `docs-site/docs/.vitepress/theme/HomePage.vue`
- Modify: `docs-site/docs/.vitepress/theme/style.css`

- [ ] Replace emoji-heavy cards with SVG/icon-like text treatments and screenshot-led product sections.
- [ ] Add a commercial hero with value proposition, CTAs, proof metrics, and showcase image composition.
- [ ] Add sections for commerce workflow, plugin architecture, multi-terminal delivery, and deployment CTA.
- [ ] Keep responsive CSS for 375px, 768px, 1024px, and 1440px widths.

### Task 3: Documentation

**Files:**
- Modify: `docs-site/docs/README.md`
- Modify: `docs-site/docs/guide/features.md`

- [ ] Document screenshot generation commands.
- [ ] Document that homepage images live in `docs-site/docs/public/showcase/`.
- [ ] Note deployment impact: static assets only, no runtime screenshot service.

### Task 4: Verification

**Files:**
- Generated: `docs-site/docs/public/showcase/*.png`

- [ ] Install missing docs-site dependencies if needed.
- [ ] Run `npm run showcase:build`.
- [ ] Run `npm run docs:build`.
- [ ] Inspect generated screenshots and homepage layout locally.
