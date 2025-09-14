import path from 'node:path'
import process from 'node:process'

import gqlcodegen from './modules/gql-codegen'

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	compatibilityDate: '2024-04-03',
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
		'@nuxtjs/tailwindcss',
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
		'nuxt-shiki',
	],

	icon: {
		localApiEndpoint: '/_nuxt_icon',
		clientBundle: {
			includeCustomCollections: true,
		},
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

	css: ['~/assets/css/global.css'],

	nitro: {
		preset: 'bun',
	},

	app: {
		head: {
			script: [
				{
					src: 'https://rybbit.a.twir.app/api/script.js',
					async: true,
					defer: true,
					'data-site-id': '1',
					'data-session-replay': 'true',
					'data-mask-patterns': '["/overlays/**"]',
				},
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

	tailwindcss: {
		config: {
			content: {
				files: [
					path.join(
						path.dirname(require.resolve('@twir/frontend-valorant-stats')),
						'**/*.{js,vue,ts}'
					),
				],
			},
		},
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
	},

	shiki: {
		bundledThemes: ['dark-plus'],
		defaultTheme: 'dark-plus',
	},
	telemetry: {
		enabled: true,
		consent: 1,
	},
})
