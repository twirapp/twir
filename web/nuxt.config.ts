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

	modules: [
		'@nuxtjs/tailwindcss',
		'@nuxt/image',
		'@nuxt/fonts',
		'nuxt-svgo',
		'@vueuse/nuxt',
		'@nuxt-alt/proxy',
	],

	devServer: {
		port: 3005,
	},

	experimental: {
		inlineRouteRules: true,
		typedPages: true,
	},

	vite: {
		plugins: [
			watch({
				onInit: true,
				pattern: '~/layers/**/*.ts',
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
})
