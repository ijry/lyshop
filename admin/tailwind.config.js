/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,ts,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: { DEFAULT: '#1e40af', light: '#3b82f6' }
      }
    }
  },
  plugins: []
}
