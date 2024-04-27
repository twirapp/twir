import { fileURLToPath } from 'node:url';

import svgSprite from '@twirapp/vite-plugin-svg-spritemap';
import vue from '@vitejs/plugin-vue';
import autoprefixer from 'autoprefixer';
import tailwind from 'tailwindcss';
import { defineConfig } from 'vite';
import { watch } from 'vite-plugin-watch';

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [
		vue(),
		svgSprite(['./src/assets/icons/*/*.svg']),
		watch({
			onInit: true,
			pattern: 'src/api/**/*.ts',
			command: 'graphql-codegen',
		}),
	],
	clearScreen: false,
	base: '/p',
	css: {
		postcss: {
			plugins: [tailwind(), autoprefixer()],
		},
	},
	resolve: {
		alias: {
			vue: 'vue/dist/vue.esm-bundler.js',
			'@': fileURLToPath(new URL('./src', import.meta.url)),
		},
	},
	server: {
		port: 3007,
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
});
