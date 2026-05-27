import { defineConfig, presetWind, presetIcons, transformerDirectives } from 'unocss'

export default defineConfig({
  presets: [
    presetWind({ preflight: false }),
    presetIcons({
      scale: 1.2,
      cdn: 'https://esm.sh/',
    }),
  ],
  transformers: [transformerDirectives()],
  rules: [
    [/^border-(l|r|t|b)-3$/, ([, direction]) => ({ [`border-${direction}-width`]: '3px' })],
    [/^(min-w|w|h)-4\.5$/, ([, property]) => ({ [property]: '1.125rem' })],
  ],
  shortcuts: {
    'flex-center': 'flex justify-center items-center',
    'flex-between': 'flex justify-between items-center',
  },
  theme: {
    colors: {
      primary: { DEFAULT: '#1e40af', light: '#3b82f6' },
    },
  },
})
