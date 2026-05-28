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

:root,
page {
  --eapp-primary: #2563eb;
  --eapp-success: #16a34a;
  --eapp-warning: #f59e0b;
  --eapp-danger: #dc2626;
  --eapp-accent: #f97316;
  --eapp-bg: #f8fafc;
  --eapp-card: #ffffff;
  --eapp-text: #1e293b;
  --eapp-text-muted: #64748b;
  --eapp-border: #e2e8f0;
  --eapp-primary-soft: #eff6ff;
  --eapp-success-soft: #dcfce7;
  --eapp-warning-soft: #fef3c7;
  --eapp-danger-soft: #fee2e2;
  --eapp-accent-soft: #ffedd5;
  --eapp-text-faint: #94a3b8;
  --eapp-border-strong: #cbd5e1;
}

page {
  background: var(--eapp-bg);
  color: var(--eapp-text);
}
</style>
