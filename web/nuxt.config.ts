import fs from 'node:fs'
import path from 'node:path'
import process from 'node:process'

import tailwindcss from '@tailwindcss/vite'
import { createResolver } from 'nuxt/kit'

import gqlcodegen from './modules/gql-codegen'

const { resolve } = createResolver(import.meta.url)

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0'

function buildDiagnosticsPlugin(): any {
	if (process.env.TWIR_BUILD_DIAGNOSTICS !== '1') return null
	const logPath = process.env.TWIR_BUILD_DIAG_LOG || '/tmp/nuxt-build-diag.log'
	fs.writeFileSync(logPath, '')
	let count = 0
	return {
		name: 'twir-build-diagnostics',
		enforce: 'pre' as const,
		transform(_code: string, id: string) {
			count++
			const env = (this as any).environment?.name ?? 'unknown'
			if (env === 'ssr') {
				const mem = process.memoryUsage()
				fs.appendFileSync(
					logPath,
					`${new Date().toISOString()} [${env}] #${count} START ${id} rss=${(mem.rss / 1024 / 1024).toFixed(1)} heap=${(mem.heapTotal / 1024 / 1024).toFixed(1)}\n`
				)
			}
			return null
		},
	}
}

const diagnosticsPlugin = buildDiagnosticsPlugin()

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	compatibilityDate: '2025-12-14',
	future: {
		compatibilityVersion: 4,
	},
	devtools: {
		enabled: true,
		timeline: {
			enabled: true,
		},
	},

	site: { indexable: true },

	modules: [
		'@pinia/nuxt',
		'@bicou/nuxt-urql',
		'reka-ui/nuxt',
		'@nuxtjs/color-mode',
		'shadcn-nuxt',
		'@nuxt/image',
		'@nuxt/icon',
		'@nuxt/fonts',
		'nuxt-svgo',
		'@vueuse/nuxt',
		'@nuxtjs/seo',
		gqlcodegen,
		'@nuxtjs/fontaine',
		'@nuxtjs/i18n',
		'@xiaobailong/web-update-notice-plugin/nuxt',
	],

	i18n: {
		locales: [
			{ code: 'en', file: 'en.json' },
			{ code: 'ru', file: 'ru.json' },
			{ code: 'de', file: 'de.json' },
			{ code: 'es', file: 'es.json' },
			{ code: 'ja', file: 'ja.json' },
			{ code: 'pt', file: 'pt.json' },
			{ code: 'sk', file: 'sk.json' },
			{ code: 'uk', file: 'uk.json' },
		],
		defaultLocale: 'en',
		langDir: 'locales',
		compilation: {
			strictMessage: false,
			escapeHtml: false,
		},
	},

	icon: {
		mode: 'svg',
		localApiEndpoint: '/_nuxt_icon',
		clientBundle: {
			includeCustomCollections: true,
		},
		customCollections: [
			{
				prefix: 'twir-overlays',
				dir: resolve('./layers/dashboard/assets/overlays'),
			},
			{
				prefix: 'twir-integrations',
				dir: resolve('./layers/dashboard/assets/integrations'),
			},
		],
	},

	devServer: {
		host: '0.0.0.0',
		port: 3010,
	},

	experimental: {
		inlineRouteRules: true,
		typedPages: true,
		renderJsonPayloads: true,
		// asyncContext: true,
	},

	vite: {
		plugins: [diagnosticsPlugin, tailwindcss()],
		optimizeDeps: {
			include: [
				'@urql/vue',
				'graphql-ws',
				'vue-sonner',
				'clsx',
				'tailwind-merge',
				'class-variance-authority',
				'zod',
				'vee-validate',
				'vue3-moveable',
				'grid-layout-plus',
				'@tanstack/vue-table',
				'@zero-dependency/utils',
				'@unovis/vue',
				'vue-draggable-plus',
			],
		},
		server: {
			allowedHosts: ['dev.twir.app', 'localhost'],
			fs: {
				allow: [
					'/home/satont/Documents/Projects/twir',
					'/home/satont/Documents/Projects/twir/node_modules/.bun',
				],
			},
		},
	},
	css: ['~/assets/css/tailwind.css', '~/assets/css/global.css'],

	nitro: {
		preset: 'bun',
		devProxy: {
			'/api': {
				target: 'http://127.0.0.1:3009',
				changeOrigin: true,
			},
		},
	},

	app: {
		head: {
			script: [
				// {
				// 	src: 'https://rybbit.twir.app/api/script.js',
				// 	async: true,
				// 	defer: true,
				// 	'data-site-id': '8eaa535a44ba',
				// 	'data-mask-patterns': '["/overlays/**"]',
				// },
			],
		},
	},

	shadcn: {
		/**
		 * Prefix for all the imported component
		 */
		prefix: 'Ui',
		/**
		 * Directory that the component lives in.
		 * @default "./components/ui"
		 */
		componentDir: './app/components/ui',
	},

	imports: {
		imports: [
			{
				from: 'tailwind-variants',
				name: 'tv',
			},
			{
				from: 'tailwind-variants',
				name: 'VariantProps',
				type: true,
			},
		],
	},

	urql: {
		endpoint: `/api/query`,
		client: path.join(process.cwd(), 'urql.ts'),
		ssr: {
			endpoint:
				process.env.NODE_ENV !== 'production'
					? // ? `${https ? 'https' : 'http'}://${config.SITE_BASE_URL}/api/query`
						'http://localhost:3009/query'
					: 'http://api-gql:3009/query',
		},
	},

	robots: {
		blockAiBots: true,
		disallow: [
			'/s',
			'/dashboard',
			'/dashboard/**',
			'/s/**',
			'/h',
			'/h/**',
			'/overlays/**',
			'/overlays',
		],
	},

	telemetry: {
		enabled: true,
		consent: 1,
	},

	webUpdateNotice: {
		checkInterval: 5 * 60 * 1000,
		base: '/',
		autoRefresh: false,
		text: {
			title: '✨ Update available',
			desc: 'New update available, you MUST update to keep site work correctly.',
			cancel: 'Cancel',
			confirm: 'Confirm',
		},
	},
})
