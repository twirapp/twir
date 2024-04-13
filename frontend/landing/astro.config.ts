import node from '@astrojs/node';
import tailwind from '@astrojs/tailwind';
import vue from '@astrojs/vue';
import { config } from '@twir/config';
import svgSprite from '@twirapp/vite-plugin-svg-spritemap';
import { defineConfig } from 'astro/config';

// eslint-disable-next-line no-undef
process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';

// https://astro.build/config
export default defineConfig({
	compressHTML: true,
	output: 'server',
	adapter: node({ mode: 'middleware' }),
	integrations: [
		tailwind({
			applyBaseStyles: false,
		}),
		vue(),
	],
	vite: {
		// ssr: {
		// 	noExternal: true,
		// },
		plugins: [
			svgSprite('./src/assets/*/*.svg'),
		],
		clearScreen: false,
		define: {
			'import.meta.env.HOST': JSON.stringify(config.SITE_BASE_URL || 'localhost:3005'),
			'import.meta.env.DISCORD_FEEDBACK_URL': JSON.stringify(config.DISCORD_FEEDBACK_URL),
		},
		server: {
			proxy: {
				'/api-new': {
					target: 'http://127.0.0.1:3009',
					changeOrigin: true,
					rewrite: (path) => path.replace(/^\/api-new/, ''),
					ws: true,
				},
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
				},
				'/overlays': {
					target: 'http://127.0.0.1:3008',
					changeOrigin: true,
					ws: true,
				},
				'/dashboard': {
					target: 'http://127.0.0.1:3006/dashboard',
					changeOrigin: true,
					rewrite: (path) => path.replace(/^\/dashboard/, ''),
					ws: true,
				},
			},
		},
	},
	server: {
		port: 3005,
		host: true,
	},
});

process.on('uncaughtException', console.error).on('unhandledRejection', console.error);
