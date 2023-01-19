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
    }),
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
        target: 'http://127.0.0.1:3002',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
        ws: true,
      },
      '/dashboard': {
        target: 'http://127.0.0.1:3006/dashboard',
        changeOrigin: false,
        ws: true,
        rewrite: (path) => path.replace(/^\/dashboard/, ''),
      },
      '/socket.io': {
        target: 'http://127.0.0.1:3004',
        changeOrigin: true,
        ws: true,
      },
      '/p': {
        target: 'http://127.0.0.1:3007/p',
        changeOrigin: true,
        ws: true,
        rewrite: (path) => path.replace(/^\/p/, ''),
      },
    },
  },
});
