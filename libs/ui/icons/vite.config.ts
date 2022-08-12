import { resolve } from 'node:path';

import { defineConfig } from 'vite';
import declarationsPlugin from 'vite-plugin-dts';
import svgLoaderPlugin from 'vite-svg-loader';

export default defineConfig({
  build: {
    rollupOptions: {
      preserveEntrySignatures: 'exports-only',
      input: {
        logos: resolve(__dirname, 'src/logos/index.ts'),
        icons: resolve(__dirname, 'src/index.ts'),
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
          assetFileNames: (chunkInfo) => {
            if (chunkInfo.name?.includes('.css')) {
              return '[name].css';
            }
            return 'assets/[name].[ext]';
          },
        },
      ],
    },
  },
  resolve: {
    alias: [{ find: '@', replacement: resolve(__dirname, 'src') }],
  },
  plugins: [
    svgLoaderPlugin(),
    declarationsPlugin({
      beforeWriteFile: (filePath, content) => ({
        filePath,
        content: content.replace(/(\.svg\?component)|(\.vue)/g, '.js'),
      }),
    }),
  ],
});
