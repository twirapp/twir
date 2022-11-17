import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import ssr from 'vite-plugin-ssr/plugin';
import svg from 'vite-svg-loader';

const __dirname = dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  define: {
    __VUE_I18N_FULL_INSTALL__: false,
    __VUE_I18N_LEGACY_API__: false,
    __INTLIFY_PROD_DEVTOOLS__: false,
  },
  plugins: [
    vue(),
    svg({
      svgo: false,
      defaultImport: 'url',
    }),
    ssr({
      prerender: true,
    })
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  // base: 'dashboard',
  server: {
    host: true,
    port: Number(process.env.VITE_PORT ?? 3005),
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL ?? 'http://localhost:3002',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
        ws: true,
      },
      '/dashboard': {
        target: process.env.VITE_API_URL ?? 'http://localhost:3006/dashboard',
        changeOrigin: false,
        rewrite: (path) => path.replace(/^\/dashboard/, ''),
      },
    },
  },
});
