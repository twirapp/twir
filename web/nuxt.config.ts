import fs from 'node:fs'
import path from 'node:path'
import process from 'node:process'

import tailwindcss from '@tailwindcss/vite'
import { createResolver } from 'nuxt/kit'

const { resolve } = createResolver(import.meta.url)

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0'

const localeCodes = fs
	.readdirSync(resolve('./locales'))
	.filter((f) => f.endsWith('.json'))
	.map((f) => f.replace('.json', ''))

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

const siteUrl = process.env.NUXT_PUBLIC_SITE_URL || 'https://twir.app'

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

	site: {
		indexable: true,
		url: siteUrl,
		name: 'Twir',
		description:
			'Powerful and useful Twitch bot that helps manage chat on big channels. Developed from streamers for streamers with love.',
		defaultLocale: 'en',
	},

	routeRules: {
		'/dashboard': { ssr: false },
		'/dashboard/**': { ssr: false },
		...Object.fromEntries(
			localeCodes.flatMap((l) => [
				[`/${l}/dashboard`, { ssr: false }],
				[`/${l}/dashboard/**`, { ssr: false }],
			])
		),
	},

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
		'@nuxtjs/i18n',
		'@nuxtjs/seo',
		'@nuxtjs/fontaine',
	],

	i18n: {
		baseUrl: siteUrl,
		locales: localeCodes.map((code) => ({
			code,
			file: `${code}.json`,
			language: code,
		})) as any, // TODO: remove any, no ai written xd
		defaultLocale: 'en',
		langDir: 'locales',
		compilation: {
			strictMessage: false,
			escapeHtml: false,
		},
		detectBrowserLanguage: {
			useCookie: true,
			cookieKey: 'i18n_redirected',
			redirectOn: 'root',
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
			{
				prefix: 'twir-compare',
				dir: resolve('./layers/landing/assets/compare'),
			},
		],
	},

	devServer: {
		host: '0.0.0.0',
		port: 3010,
	},

	build: {
		// frontend-chat ships Vue SFC source, including scoped message styles.
		transpile: ['@twir/frontend-chat'],
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
			exclude: ['@twir/frontend-chat'],
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
				'@vuepic/vue-datepicker',
				'@guolao/vue-monaco-editor',
				'@unhead/schema-org/vue',
				'@tanstack/vue-virtual',
				'tinycolor2',
			],
		},
		server: {
			allowedHosts: ['dev.twir.app', 'localhost'],
			// fs: {
			// 	allow: [
			// 		'/home/satont/Documents/Projects/twir',
			// 		'/home/satont/Documents/Projects/twir/node_modules/.bun',
			// 	],
			// },
		},
		css: {
			lightningcss: {
				errorRecovery: true,
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
			'/s/',
			'/s/**',
			'/dashboard',
			'/dashboard/**',
			'/h/',
			'/h/**',
			'/overlays/',
			'/overlays/**',
		],
	},

	ogImage: {
		fontSubsets: ['latin', 'latin-ext', 'cyrillic', 'cyrillic-ext'],
	},

	sitemap: {
		exclude: [
			'/dashboard',
			'/dashboard/**',
			'/s',
			'/s/**',
			'/h/**',
			'/o',
			'/o/**',
			'/overlays',
			'/overlays/**',
			'/login',
			'/login/**',
			'/import',
			'/import/**',
			'/settings',
			'/settings/**',
			'/en/dashboard',
			'/en/dashboard/**',
			'/**/dashboard',
			'/**/dashboard/**',
			'/**/s/**',
			'/**/h/**',
			'/**/o',
			'/**/o/**',
			'/**/overlays',
			'/**/overlays/**',
			'/**/login',
			'/**/login/**',
			'/**/url-shortener/profile',
			'/**/import',
			'/**/import/**',
			'/**/settings',
			'/**/settings/**',
		],
	},

	telemetry: {
		enabled: true,
		consent: 1,
	},
})
