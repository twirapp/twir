import path from 'node:path';
import { fileURLToPath } from 'node:url';

import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite';
import { webUpdateNotice } from '@plugin-web-update-notification/vite';
import vue from '@vitejs/plugin-vue';
import { defineConfig, loadEnv } from 'vite';
// import { VitePWA } from 'vite-plugin-pwa';
import svg from 'vite-svg-loader';

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '');

	return {
		plugins: [
			vue({
				script: {
					defineModel: true,
				},
			}),
			svg({ svgo: false }),
			// VitePWA({
			// 	registerType: 'autoUpdate',
			// 	devOptions: {
			// 		enabled: true,
			// 	},
			// 	includeAssets: ['favicon.ico', 'apple-touch-icon.png', 'mask-icon.svg'],
			// 	manifest: {
			// 		name: 'Twir - The perfect bot for your stream',
			// 		short_name: 'TwirApp',
			// 		description: 'The perfect bot for your stream. Created by streamers. Made for streamers. Used by streamers. For streamers with love.',
			// 		theme_color: '#101014',
			// 		icons: [
			// 			{
			// 				src: 'pwa-192x192.png',
			// 				sizes: '192x192',
			// 				type: 'image/png',
			// 			},
			// 			{
			// 				src: 'pwa-512x512.png',
			// 				sizes: '512x512',
			// 				type: 'image/png',
			// 			},
			// 			{
			// 				src: 'pwa-512x512.png',
			// 				sizes: '512x512',
			// 				type: 'image/png',
			// 				purpose: 'any',
			// 			},
			// 			{
			// 				src: 'pwa-512x512.png',
			// 				sizes: '512x512',
			// 				type: 'image/png',
			// 				purpose: 'maskable',
			// 			},
			// 		],
			// 	},
			// }),
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
				'/p': {
					target: 'http://127.0.0.1:3007',
					changeOrigin: true,
					ws: true,
					// rewrite: (path) => path.replace(/^\/p/, ''),
				},
				'/overlays': {
					target: 'http://127.0.0.1:3008',
					changeOrigin: true,
					ws: true,
				},
			},
		},
		clearScreen: false,
	};
});
