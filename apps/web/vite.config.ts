import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

import i18n from '@intlify/vite-plugin-vue-i18n';
import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import ssr from 'vite-plugin-ssr/plugin';

const __dirname = dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  plugins: [
    vue(),
    ssr({
      prerender: true,
    }),
    i18n({
      include: resolve(__dirname, 'src/locales/**'),
      runtimeOnly: false,
    }),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      'vue-i18n': 'vue-i18n/dist/vue-i18n.cjs.prod.js',
    },
  },
  server: {
    port: 3000,
  },
});
