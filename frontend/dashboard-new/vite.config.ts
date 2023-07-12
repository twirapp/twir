import { fileURLToPath } from 'node:url';

import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
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
