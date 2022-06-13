import path from 'path';
import { fileURLToPath, URL } from 'url';

import vueI18n from '@intlify/vite-plugin-vue-i18n';
import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import svgLoader from 'vite-svg-loader';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    svgLoader(),
    vueI18n({
      include: path.resolve(path.dirname(fileURLToPath(import.meta.url)), 'src/locales/**'),
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(path.dirname(fileURLToPath(import.meta.url)), 'src'),
      'vue-i18n': 'vue-i18n/dist/vue-i18n.runtime.esm-bundler.js',
    },
  },
  server: {
    port: 3005,
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL ?? 'http://localhost:3002',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
        ws: true,
      },
    },
  },
});
