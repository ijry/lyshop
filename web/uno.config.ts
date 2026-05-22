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
    'btn-primary': 'px-6 py-2.5 bg-blue-700 text-white rounded-lg text-sm font-medium hover:bg-blue-800 transition-colors cursor-pointer disabled:opacity-50',
    'btn-outline': 'px-6 py-2.5 border border-gray-300 text-gray-700 rounded-lg text-sm font-medium hover:border-blue-700 hover:text-blue-700 transition-colors cursor-pointer',
    'card': 'bg-white rounded-xl shadow-sm border border-gray-100',
    'input-base': 'w-full px-4 py-2.5 border border-gray-200 rounded-lg text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/20 transition-all',
  },
  theme: {
    colors: {
      primary: { DEFAULT: '#1e40af', light: '#3b82f6', dark: '#1e3a8a' },
    },
  },
})
