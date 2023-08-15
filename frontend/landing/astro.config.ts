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
		plugins: [svg({ defaultImport: 'url' }) as unknown as PluginOption],
		clearScreen: false,
		define: {
			'import.meta.env.HOST': JSON.stringify(config.HOSTNAME || 'localhost:3005'),
			'import.meta.env.DISCORD_FEEDBACK_URL': JSON.stringify(config.DISCORD_FEEDBACK_URL),
		},
	},
	server: {
		port: 3005,
		host: true,
	},
});

process.on('uncaughtException', console.error).on('unhandledRejection', console.error);
