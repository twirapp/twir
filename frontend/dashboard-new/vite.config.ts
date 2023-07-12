import { fileURLToPath } from 'node:url';

import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import svg from 'vite-svg-loader';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), svg({ svgo: false })],
	base: '/dashboard',
	resolve: {
		alias: {
			vue: 'vue/dist/vue.esm-bundler.js',
			'@': fileURLToPath(new URL('./src', import.meta.url)),
		},
	},
	server: {
		port: 3006,
		host: true,
		proxy: {
			'/api': {
				target: 'http://127.0.0.1:3002',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, ''),
				ws: true,
			},
		},
	},
	clearScreen: false,
});
