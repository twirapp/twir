import { resolve } from 'node:path';

import vuePlugin from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import declarationsPlugin from 'vite-plugin-dts';
import svgLoaderPlugin from 'vite-svg-loader';

export default defineConfig({
  build: {
    rollupOptions: {
      preserveEntrySignatures: 'exports-only',
      input: {
        components: resolve(__dirname, 'src/components/index.ts'),
        logos: resolve(__dirname, 'src/icons/logos/index.ts'),
        icons: resolve(__dirname, 'src/icons/index.ts'),
      },
      external: ['vue'],
      output: [
        {
          preserveModules: true,
          esModule: true,
          exports: 'named',
          dir: 'dist',
          globals: {
            vue: 'Vue',
          },
          format: 'esm',
          entryFileNames: (chunkInfo) => {
            if (chunkInfo.isEntry) {
              return 'index.js';
            }
            return '[name].js';
          },
        },
      ],
    },
  },
  resolve: {
    alias: [{ find: '@', replacement: resolve(__dirname, 'src') }],
  },
  plugins: [
    vuePlugin(),
    svgLoaderPlugin(),
    declarationsPlugin({
      beforeWriteFile: (filePath, content) => {
        return {
          filePath,
          content: content.replace(/(svg\?component)|(vue)/g, 'js'),
        };
      },
      exclude: './src/shims-vue.d.ts',
    }),
  ],
});
