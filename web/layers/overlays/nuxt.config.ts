export default defineNuxtConfig({
  routeRules: {
    '/o/:apiKey/**': { ssr: false },
    '/overlays/:apiKey/**': {
      redirect: '/o/:apiKey/**',
    },
  },
})
