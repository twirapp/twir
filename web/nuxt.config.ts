import path from 'node:path'
import process from 'node:process'

import { config, readEnv } from '@twir/config'
import { watch } from 'vite-plugin-watch'

readEnv(path.join(process.cwd(), '..', '.env'))

const https = config.TWITCH_CALLBACKURL.startsWith('https')

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0'
// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	compatibilityDate: '2024-04-03',

	devtools: {
		enabled: true,
		timeline: {
			enabled: true,
		},
	},

	devServer: {
		port: 3005,
		// https,
	},

	experimental: {
		inlineRouteRules: true,
	},

	srcDir: './app',

	modules: [
		'@nuxtjs/tailwindcss',
		'@nuxt/image',
		'@nuxt/fonts',
		'@bicou/nuxt-urql',
		'nuxt-svgo',
		'@vueuse/nuxt',
		'@nuxt-alt/proxy',
	],

	css: [
		'~/assets/styles/global.css',
	],

	vite: {
		plugins: [
			watch({
				onInit: true,
				pattern: '~/api/**/*.ts',
				command: 'graphql-codegen',
			}),
		],
		resolve: {
			alias: {
				'@': path.resolve(__dirname, 'app'),
			},
		},
		server: {
			hmr: {
				protocol: https ? 'wss' : 'ws',
			},
			proxy: {
				'/api-old': {
					target: 'http://127.0.0.1:3002',
					changeOrigin: true,
					rewrite: (path) => path.replace(/^\/api-old/, ''),
					ws: true,
				},
				'/api': {
					target: 'http://127.0.0.1:3009',
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

	proxy: {
		debug: true,
		experimental: {
			listener: true,
		},
		proxies: {
			'/api': {
				target: 'http://127.0.0.1:3009',
				changeHost: false,
				rewrite: (path) => path.replace(/^\/api/, ''),
				ws: true,
			},
		},
	},

	nitro: {
		devProxy: {
			'/api': {
				target: 'http://127.0.0.1:3009',
				changeOrigin: true,
				ws: true,
			},
		},
	},

	urql: {
		endpoint: process.env.NODE_ENV !== 'production'
			? `${https ? 'https' : 'http'}://${config.SITE_BASE_URL}/api/query`
			: 'https://twir.app/api/query',
		client: '~/configs/urql.ts',
	},
})
