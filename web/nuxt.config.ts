import path from 'node:path'

import { watch } from 'vite-plugin-watch'

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
		endpoint:
		// eslint-disable-next-line node/prefer-global/process
												process.env.NODE_ENV !== 'production'
													? 'https://dev.twir.app/api/query'
													: 'https://twir.app/api/query',
		client: '~/configs/urql.ts',
	},
})
