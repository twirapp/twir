import tailwindcss from '@tailwindcss/vite'
import path from 'node:path'
import process from 'node:process'

import gqlcodegen from './modules/gql-codegen'

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0'

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
	],

	icon: {
		mode: 'svg',
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

	vite: {
		plugins: [tailwindcss()],
	},
	css: ['~/assets/css/tailwind.css', '~/assets/css/global.css'],

	nitro: {
		preset: 'bun',
	},

	app: {
		head: {
			script: [
				{
					src: 'https://rybbit.twir.app/api/script.js',
					async: true,
					defer: true,
					'data-site-id': '8eaa535a44ba',
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
})
