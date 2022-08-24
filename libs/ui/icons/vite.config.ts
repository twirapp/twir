import { resolve } from 'node:path';

import { defineConfig } from 'vite';
import svgPlugin from 'vite-svg-loader';

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
          preserveModules: false,
          esModule: true,
          exports: 'named',
          dir: 'dist',
          globals: {
            vue: 'Vue',
          },
          format: 'esm',
          entryFileNames: '[name].js',
        },
      ],
    },
  },
  resolve: {
    alias: [{ find: '@', replacement: resolve(__dirname, 'src') }],
  },
  plugins: [
    svgPlugin({
      svgo: false,
    }),
  ],
});
