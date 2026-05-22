import { defineConfig, presetWind, presetIcons, transformerDirectives } from 'unocss'

export default defineConfig({
  presets: [
    presetWind(),
    presetIcons({
      scale: 1.2,
      cdn: 'https://esm.sh/',
    }),
  ],
  transformers: [transformerDirectives()],
  shortcuts: {
    'flex-center': 'flex justify-center items-center',
    'flex-between': 'flex justify-between items-center',
    'btn-primary': 'px-6 py-2.5 bg-red-600 text-white rounded-lg text-sm font-medium hover:bg-red-700 transition-colors cursor-pointer disabled:opacity-50',
    'btn-outline': 'px-6 py-2.5 border border-gray-300 text-gray-700 rounded-lg text-sm font-medium hover:border-red-600 hover:text-red-600 transition-colors cursor-pointer',
    'card': 'bg-white rounded-xl shadow-sm border border-gray-100',
    'input-base': 'w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:border-red-500 focus:ring-1 focus:ring-red-500/20 transition-all',
  },
  theme: {
    colors: {
      primary: { DEFAULT: '#dc2626', light: '#ef4444', dark: '#b91c1c' },
    },
  },
})
