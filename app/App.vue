<script>
import { useTheme } from '@/composables/useTheme'

export default {
  onLaunch() {
    useTheme().initTheme()
    uni.hideTabBar({ animation: false })
  },
  onShow() {
    if (typeof uni !== 'undefined' && typeof uni.$emit === 'function') {
      uni.$emit('up:tabbar-sync')
    }
  },
  onHide() {}
}
</script>
<style lang="scss">
@import 'uview-plus/index.scss';

/* App-level theme tokens (light defaults) */
:root,
page {
  --app-page-bg: #f5f5f5;
  --app-card-bg: #ffffff;
  --app-text-primary: #111111;
  --app-text-secondary: #333333;
  --app-text-tertiary: #666666;
  --app-text-placeholder: #999999;
  --app-text-inverse: #ffffff;
  --app-border-color: #f0f0f0;
  --app-divider-color: #f5f5f5;
  --app-accent: #dc2626;
  --app-accent-light: #fef2f2;
  --app-price-color: #dc2626;
  --app-shadow: 0 2px 12px rgba(0,0,0,0.06);
  --app-shadow-sm: 0 1px 6px rgba(0,0,0,0.04);
}

/* Auto dark mode (system preference) */
@media (prefers-color-scheme: dark) {
  :root,
  page {
    --app-page-bg: #1f1f1f;
    --app-card-bg: #1c1c1e;
    --app-text-primary: #f5f5f5;
    --app-text-secondary: #d1d5db;
    --app-text-tertiary: #9ca3af;
    --app-text-placeholder: #6b7280;
    --app-text-inverse: #111111;
    --app-border-color: #3a3a3c;
    --app-divider-color: #2c2c2e;
    --app-accent: #ef4444;
    --app-accent-light: #3a2222;
    --app-price-color: #f87171;
    --app-shadow: 0 2px 12px rgba(0,0,0,0.3);
    --app-shadow-sm: 0 1px 6px rgba(0,0,0,0.2);
  }
}

/* Explicit light */
[data-up-theme='light'],
[data-up-theme='light'] page {
  --app-page-bg: #f5f5f5;
  --app-card-bg: #ffffff;
  --app-text-primary: #111111;
  --app-text-secondary: #333333;
  --app-text-tertiary: #666666;
  --app-text-placeholder: #999999;
  --app-text-inverse: #ffffff;
  --app-border-color: #f0f0f0;
  --app-divider-color: #f5f5f5;
  --app-accent: #dc2626;
  --app-accent-light: #fef2f2;
  --app-price-color: #dc2626;
  --app-shadow: 0 2px 12px rgba(0,0,0,0.06);
  --app-shadow-sm: 0 1px 6px rgba(0,0,0,0.04);
}

/* Global UnoCSS class dark overrides */
@media (prefers-color-scheme: dark) {
  .bg-white { background-color: #1c1c1e !important; }
  .bg-gray-50 { background-color: #1f1f1f !important; }
  .bg-gray-100 { background-color: #2c2c2e !important; }
  .text-gray-800 { color: #f5f5f5 !important; }
  .text-gray-700 { color: #d1d5db !important; }
  .text-gray-600 { color: #9ca3af !important; }
  .text-gray-500 { color: #9ca3af !important; }
  .text-gray-400 { color: #6b7280 !important; }
  .border-gray-100 { border-color: #3a3a3c !important; }
  .border-gray-200 { border-color: #3a3a3c !important; }
  .border-slate-100 { border-color: #3a3a3c !important; }
  .shadow-sm { box-shadow: 0 1px 6px rgba(0,0,0,0.2) !important; }
}

[data-up-theme='dark'] .bg-white { background-color: #1c1c1e !important; }
[data-up-theme='dark'] .bg-gray-50 { background-color: #1f1f1f !important; }
[data-up-theme='dark'] .bg-gray-100 { background-color: #2c2c2e !important; }
[data-up-theme='dark'] .text-gray-800 { color: #f5f5f5 !important; }
[data-up-theme='dark'] .text-gray-700 { color: #d1d5db !important; }
[data-up-theme='dark'] .text-gray-600 { color: #9ca3af !important; }
[data-up-theme='dark'] .text-gray-500 { color: #9ca3af !important; }
[data-up-theme='dark'] .text-gray-400 { color: #6b7280 !important; }
[data-up-theme='dark'] .border-gray-100 { border-color: #3a3a3c !important; }
[data-up-theme='dark'] .border-gray-200 { border-color: #3a3a3c !important; }
[data-up-theme='dark'] .border-slate-100 { border-color: #3a3a3c !important; }
[data-up-theme='dark'] .shadow-sm { box-shadow: 0 1px 6px rgba(0,0,0,0.2) !important; }

/* Explicit dark */
[data-up-theme='dark'],
[data-up-theme='dark'] page {
  --app-page-bg: #1f1f1f;
  --app-card-bg: #1c1c1e;
  --app-text-primary: #f5f5f5;
  --app-text-secondary: #d1d5db;
  --app-text-tertiary: #9ca3af;
  --app-text-placeholder: #6b7280;
  --app-text-inverse: #111111;
  --app-border-color: #3a3a3c;
  --app-divider-color: #2c2c2e;
  --app-accent: #ef4444;
  --app-accent-light: #3a2222;
  --app-price-color: #f87171;
  --app-shadow: 0 2px 12px rgba(0,0,0,0.3);
  --app-shadow-sm: 0 1px 6px rgba(0,0,0,0.2);
}
</style>
