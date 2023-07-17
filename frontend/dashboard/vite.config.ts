import { fileURLToPath } from 'node:url';

import vue from '@vitejs/plugin-vue';
import AutoImport from 'unplugin-auto-import/vite';
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers';
import Components from 'unplugin-vue-components/vite';
import { defineConfig } from 'vite';
import svg from 'vite-svg-loader';


// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
		vue({
			script: {
				defineModel: true,
			},
		}),
		svg({ svgo: false }),
		AutoImport({
			imports: [
				'vue',
				{
					'naive-ui': [
						'useDialog',
						'useMessage',
						'useNotification',
						'useLoadingBar',
					],
				},
			],
		}),
		Components({
			resolvers: [NaiveUiResolver()],
		}),
	],
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
			'/socket': {
				target: 'http://127.0.0.1:3004',
				changeOrigin: true,
				ws: true,
				rewrite: (path) => path.replace(/^\/socket/, ''),
			},
		},
	},
	clearScreen: false,
});
