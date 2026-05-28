import presetWeapp from 'unocss-preset-weapp'
import { extractorAttributify, transformerClass } from 'unocss-preset-weapp/transformer'
import { defineConfig } from 'unocss'
import transformerDirectives from '@unocss/transformer-directives'

const { presetWeappAttributify, transformerAttributify } = extractorAttributify()

export default defineConfig({
  presets: [
    // @ts-expect-error type mismatch
    presetWeapp({ dark: 'media' }),
    // @ts-expect-error type mismatch
    presetWeappAttributify(),
  ],
  shortcuts: {
    'flex-center': 'flex justify-center items-center',
    'flex-between': 'flex justify-between items-center',
    'eapp-page': 'min-h-screen px-24rpx pt-24rpx box-border',
    'eapp-card': 'rounded-24rpx bg-white p-24rpx shadow-sm',
  },
  transformers: [
    transformerDirectives({ enforce: 'pre' }),
    // @ts-expect-error type mismatch
    transformerAttributify(),
    // @ts-expect-error type mismatch
    transformerClass(),
  ],
})
