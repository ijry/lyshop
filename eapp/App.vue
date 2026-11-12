<script>
export default {
  onLaunch() {},
  onShow() {
    const pages = getCurrentPages()
    const route = pages[pages.length - 1]?.route || ''
    const isMock = import.meta.env.VITE_MOCK === 'true'
    const token = String(uni.getStorageSync('eapp_admin_token') || '')
    const isLogin = route === 'pages/login/index'

    if (isMock) {
      if (!token) {
        uni.setStorageSync('eapp_admin_token', 'demo_admin_token')
        uni.setStorageSync('eapp_admin_username', 'admin')
      }
      if (isLogin) {
        uni.switchTab({ url: '/pages/dashboard/index' })
      }
      return
    }

    if (!token && !isLogin) {
      uni.reLaunch({ url: '/pages/login/index' })
      return
    }
    if (token && isLogin) {
      uni.switchTab({ url: '/pages/dashboard/index' })
    }
  },
  onHide() {},
}
</script>

<style lang="scss">
@import 'uview-plus/index.scss';

/* ==========================================================================
   LYShop 商家移动端 · 设计令牌
   单一事实来源：页面与组件只引用这里的变量，不硬编码颜色/圆角/阴影。
   uview-plus 通过 --up-* 变量桥接，自动继承同一套主题。
   ========================================================================== */
:root,
page {
  /* 品牌 */
  --eapp-primary: #1d4ed8;
  --eapp-primary-deep: #173b9e;
  --eapp-primary-soft: #eef3fd;
  --eapp-primary-border: #d3def7;
  --eapp-primary-on: #ffffff;

  /* 语义色：低饱和，适合长时间办公阅读 */
  --eapp-success: #157f4b;
  --eapp-success-soft: #eaf6ef;
  --eapp-warning: #a75a08;
  --eapp-warning-soft: #fdf3e6;
  --eapp-danger: #b42318;
  --eapp-danger-soft: #fbeeec;
  --eapp-accent: #0e7490;
  --eapp-accent-soft: #e9f4f7;
  --eapp-info: #1d4ed8;
  --eapp-info-soft: #eef3fd;

  /* 中性色阶 */
  --eapp-bg: #f5f7fa;
  --eapp-card: #ffffff;
  --eapp-surface-subtle: #fafbfd;
  --eapp-text: #101828;
  --eapp-text-secondary: #475467;
  --eapp-text-muted: #667085;
  --eapp-text-faint: #98a2b3;
  --eapp-border: #e7eaf0;
  --eapp-border-strong: #d0d5dd;
  --eapp-divider: #eef0f4;

  /* 圆角：克制、统一 */
  --eapp-radius-sm: 8rpx;
  --eapp-radius: 12rpx;
  --eapp-radius-lg: 16rpx;
  --eapp-radius-xl: 20rpx;
  --eapp-radius-pill: 999rpx;

  /* 阴影：极淡，靠边框而非投影划分层次 */
  --eapp-shadow: 0 1rpx 2rpx rgba(16, 24, 40, 0.04), 0 4rpx 16rpx rgba(16, 24, 40, 0.04);
  --eapp-shadow-md: 0 8rpx 28rpx rgba(16, 24, 40, 0.1);

  /* 字阶 */
  --eapp-fs-xs: 20rpx;
  --eapp-fs-sm: 24rpx;
  --eapp-fs-base: 26rpx;
  --eapp-fs-md: 28rpx;
  --eapp-fs-lg: 32rpx;
  --eapp-fs-xl: 36rpx;
  --eapp-fs-2xl: 44rpx;

  /* 间距 */
  --eapp-space-1: 8rpx;
  --eapp-space-2: 12rpx;
  --eapp-space-3: 16rpx;
  --eapp-space-4: 20rpx;
  --eapp-space-5: 24rpx;
  --eapp-space-6: 32rpx;

  /* uview-plus 主题桥接 */
  --up-primary: #1d4ed8;
  --up-primary-dark: #173b9e;
  --up-primary-disabled: #aabff0;
  --up-primary-light: #eef3fd;
  --up-success: #157f4b;
  --up-success-dark: #116a3f;
  --up-success-disabled: #9ccdb2;
  --up-success-light: #eaf6ef;
  --up-warning: #a75a08;
  --up-warning-dark: #8b4a06;
  --up-warning-disabled: #e2bd92;
  --up-warning-light: #fdf3e6;
  --up-error: #b42318;
  --up-error-dark: #911b12;
  --up-error-disabled: #e4aaa4;
  --up-error-light: #fbeeec;
  --up-info: #667085;
  --up-info-dark: #475467;
  --up-info-disabled: #cdd3dd;
  --up-info-light: #f2f4f7;
  --up-main-color: #101828;
  --up-content-color: #475467;
  --up-tips-color: #98a2b3;
  --up-light-color: #d0d5dd;
  --up-border-color: #e7eaf0;
  --up-bg-color: #f5f7fa;

  /* legacy 变量兜底 */
  --u-primary: #1d4ed8;
  --u-success: #157f4b;
  --u-warning: #a75a08;
  --u-error: #b42318;
  --u-main-color: #101828;
  --u-content-color: #475467;
  --u-tips-color: #98a2b3;
  --u-border-color: #e7eaf0;
  --u-bg-color: #f5f7fa;
}

page {
  background: var(--eapp-bg);
  color: var(--eapp-text);
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Helvetica Neue',
    'Microsoft YaHei', 'Source Han Sans SC', sans-serif;
  font-size: var(--eapp-fs-base);
  line-height: 1.5;
  -webkit-font-smoothing: antialiased;
}

button::after {
  border: none;
}

/* 数字等宽对齐：金额、指标、编号 */
.eapp-num {
  font-variant-numeric: tabular-nums;
  font-feature-settings: 'tnum';
}

/* 通用卡片表面 */
.eapp-surface {
  background: var(--eapp-card);
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius-lg);
  box-shadow: var(--eapp-shadow);
  box-sizing: border-box;
}

/* UnoCSS eapp-card 快捷类补上边框与阴影（不与其原子类属性冲突） */
.eapp-card {
  border: 1rpx solid var(--eapp-border);
  box-shadow: var(--eapp-shadow);
}

/* 区块结构 */
.eapp-section {
  margin-bottom: var(--eapp-space-5);
}
.eapp-section__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--eapp-space-3);
}
.eapp-section__title {
  font-size: var(--eapp-fs-md);
  font-weight: 600;
  color: var(--eapp-text);
  letter-spacing: 0.5rpx;
}
.eapp-section__extra {
  font-size: var(--eapp-fs-sm);
  color: var(--eapp-text-muted);
}

.eapp-divider {
  height: 1rpx;
  background: var(--eapp-divider);
}

.eapp-muted {
  color: var(--eapp-text-muted);
}
.eapp-faint {
  color: var(--eapp-text-faint);
}

/* 交互反馈：小程序用 hover-class，H5 用 :active */
.eapp-hover {
  opacity: 0.68;
}
.eapp-press {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.eapp-press:active {
  opacity: 0.68;
}

/* 页头图标按钮：统一尺寸、边框与点击热区 */
.eapp-icon-btn {
  position: relative;
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1rpx solid var(--eapp-border);
  border-radius: var(--eapp-radius);
  background: var(--eapp-card);
  box-sizing: border-box;
  flex-shrink: 0;
}

.eapp-icon-btn__dot {
  position: absolute;
  top: 10rpx;
  right: 10rpx;
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
  background: var(--eapp-danger);
  border: 2rpx solid var(--eapp-card);
  box-sizing: border-box;
}

.u-button--mini {
  width: auto !important;
  flex-shrink: 0;
}
</style>
