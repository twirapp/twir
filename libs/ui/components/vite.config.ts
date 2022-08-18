import { resolve } from 'node:path';

import vuePlugin from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import declarationsPlugin from 'vite-plugin-dts';

export default defineConfig({
  build: {
    cssCodeSplit: false,
    lib: {
      entry: resolve(__dirname, 'src/index.ts'),
      formats: ['es'],
    },
    rollupOptions: {
      external: ['vue', '@tsuwari/ui-icons/icons', 'vee-validate'],
      output: {
        esModule: true,
        exports: 'named',
        dir: 'dist',
        globals: {
          vue: 'Vue',
        },
        format: 'esm',
        entryFileNames: 'index.js',
      },
    },
  },
  resolve: {
    alias: [{ find: '@', replacement: resolve(__dirname, 'src') }],
  },
  plugins: [
    vuePlugin(),
    declarationsPlugin({
      cleanVueFileName: true,
      exclude: './src/shims-vue.d.ts',
    }),
  ],
});
