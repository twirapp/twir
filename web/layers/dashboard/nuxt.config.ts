export default defineNuxtConfig({
  ssr: false,
  routeRules: {
    '/dashboard/**': { ssr: false },
  },
  vite: {
    resolve: {
      alias: {
        vue: 'vue/dist/vue.esm-bundler.js',
      },
    },
  },
})
