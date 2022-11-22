import path from 'path';
import { fileURLToPath } from 'url';

import vueI18n from '@intlify/vite-plugin-vue-i18n';
import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import svgLoader from 'vite-svg-loader';

// https://vitejs.dev/config/
export default defineConfig({
  clearScreen: false,
  plugins: [
    vue(),
    svgLoader(),
    vueI18n({
      include: path.resolve(path.dirname(fileURLToPath(import.meta.url)), 'src/locales/**'),
    }),
  ],
  resolve: {
    alias: {
      'xmlhttprequest-ssl': './node_modules/engine.io-client/lib/xmlhttprequest.js',
      '@': path.resolve(path.dirname(fileURLToPath(import.meta.url)), 'src'),
      '@elements': path.resolve(
        path.dirname(fileURLToPath(import.meta.url)),
        'src',
        'components',
        'elements',
      ),
      'vue-i18n': 'vue-i18n/dist/vue-i18n.runtime.esm-bundler.js',
    },
  },
  optimizeDeps: {
    exclude: [
      'amqplib',
      'mqtt',
      'kafkajs',
      '@grpc/grpc-js',
      'amqp-connection-manager',
      '@grpc/proto-loader',
      '@nestjs/websockets/socket-module',
    ],
  },
  server: {
    host: '0.0.0.0',
    port: Number(process.env.VITE_PORT ?? 3005),
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:3002',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
        ws: true,
      },
      '/socket.io': {
        target: 'http://127.0.0.1:3004',
        changeOrigin: true,
        ws: true,
      },
      // '/socket.io/:path*/': {
      //   target: process.env.VITE_API_URL ?? 'http://localhost:3003',
      //   changeOrigin: true,
      //   ws: true,
      // },
    },
  },
});
