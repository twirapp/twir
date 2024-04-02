import path from 'node:path';
import { fileURLToPath } from 'node:url';

import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite';
import { webUpdateNotice } from '@plugin-web-update-notification/vite';
import svgSprite from '@twirapp/vite-plugin-svg-spritemap';
import vue from '@vitejs/plugin-vue';
import autoprefixer from 'autoprefixer';
import tailwind from 'tailwindcss';
import { defineConfig, loadEnv } from 'vite';


// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, path.resolve(process.cwd(), '..', '..'), '');

	return {
		build: {
			minify: true,
		},
		css: {
			postcss: {
				plugins: [tailwind(), autoprefixer()],
			},
		},
		plugins: [
			vue({
				script: {
					defineModel: true,
				},
			}),
			svgSprite(['./src/assets/*/*.svg', './src/assets/*.svg']),
			webUpdateNotice({
				notificationProps: {
					title: 'New version',
					description: 'An update available, please refresh the page to get latest features and bug fixes!',
					buttonText: 'refresh',
					dismissButtonText: 'cancel',
				},
				checkInterval: 1 * 60 * 1000,
			}),
			VueI18nPlugin({
				include: [
					path.resolve(__dirname, './src/locales/**'),
				],
				strictMessage: false,
				escapeHtml: false,
				runtimeOnly: true,
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
			hmr: {
				protocol: env.USE_WSS === 'true' ? 'wss' : 'ws',
			},
		},
		clearScreen: false,
	};
});
