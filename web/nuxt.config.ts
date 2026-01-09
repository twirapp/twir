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

	extends: ['./layers/dashboard'],

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
		// '@nuxtjs/i18n',
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
					src: 'https://rybbit.a.twir.app/api/script.js',
					async: true,
					defer: true,
					'data-site-id': '1',
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
		disableI18nIntegration: true,
	},

	telemetry: {
		enabled: true,
		consent: 1,
	},

	i18n: {
		locales: [
			{ code: 'en', name: 'English', file: 'en.json' },
			{ code: 'ru', name: 'Russian', file: 'ru.json' },
			{ code: 'uk', name: 'Українська', file: 'uk.json' },
			{ code: 'de', name: 'Deutsch', file: 'de.json' },
			{ code: 'ja', name: '日本語', file: 'ja.json' },
			{ code: 'sk', name: 'Slovenčina', file: 'sk.json' },
			{ code: 'es', name: 'Español', file: 'es.json' },
			{ code: 'pt', name: 'Português', file: 'pt.json' },
		],
		defaultLocale: 'en',
		strategy: 'no_prefix',
		lazy: true,
		langDir: 'locales',
		detectBrowserLanguage: {
			useCookie: true,
			cookieKey: 'twir_locale',
			redirectOn: 'root',
		},
		compilation: {
			strictMessage: false,
			escapeHtml: false,
		},
	},
})
