import node from '@astrojs/node';
import tailwind from '@astrojs/tailwind';
import vue from '@astrojs/vue';
import { config } from '@twir/config';
import { defineConfig } from 'astro/config';
import { PluginOption } from 'vite';
import svg from 'vite-svg-loader';

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
		plugins: [svg({ defaultImport: 'url', svgo: false }) as unknown as PluginOption],
		clearScreen: false,
		define: {
			'import.meta.env.HOST': JSON.stringify(config.HOSTNAME || 'localhost:3005'),
			'import.meta.env.DISCORD_FEEDBACK_URL': JSON.stringify(config.DISCORD_FEEDBACK_URL),
		},
		server: {
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
				'/dashboard': {
					target: 'http://127.0.0.1:3006/dashboard',
					changeOrigin: true,
					rewrite: (path) => path.replace(/^\/dashboard/, ''),
					ws: true,
				},
				'/cdn': {
					target: 'http://127.0.0.1:8000',
					changeOrigin: true,
					rewrite: (path) => path.replace(/^\/cdn/, ''),
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
