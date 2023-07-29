import node from '@astrojs/node';
import tailwind from '@astrojs/tailwind';
import * as config from '@twir/config';
import { defineConfig } from 'astro/config';

console.log(config);
// process.env.HOST = config.HOSTNAME ?? 'localhost:3005';

import.meta.env.HOST = config.HOSTNAME ?? 'localhost:3005';

// https://astro.build/config
export default defineConfig({
	compressHTML: true,
	output: 'server',
	adapter: node({ mode: 'middleware' }),
	integrations: [
		tailwind({
			applyBaseStyles: false,
		}),
	],
	vite: {
		ssr: {
			noExternal: true,
		},
	},
	server: {
		port: 3005,
		host: true,
		proxy: {
			'/api': {
				target: 'http://127.0.0.1:3002',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, ''),
				ws: true,
			},
			'/dashboard': {
				target: 'http://127.0.0.1:3006',
				changeOrigin: true,
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
});
