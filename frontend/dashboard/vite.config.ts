import path from 'node:path'
import process from 'node:process'
import { fileURLToPath } from 'node:url'
import { analyzer, unstableRolldownAdapter } from 'vite-bundle-analyzer'

import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import { webUpdateNotice } from '@plugin-web-update-notification/vite'
import svgSprite from '@twirapp/vite-plugin-svg-spritemap'
import vue from '@vitejs/plugin-vue'
import { type PluginOption, defineConfig, loadEnv } from 'vite'
import { watch } from 'vite-plugin-watch'
import tailwindcss from '@tailwindcss/vite'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
	const analyzeMode = process.env.ANALYZE
	const env = loadEnv(mode, path.resolve(process.cwd(), '..', '..'), '')

	const plugins: PluginOption[] = [
		vue(),
		svgSprite(['./src/assets/*/*.svg', './src/assets/*.svg']),
		webUpdateNotice({
			notificationProps: {
				title: 'New version',
				description:
					'An update available, please refresh the page to get latest features and bug fixes!',
				buttonText: 'refresh',
				dismissButtonText: 'cancel',
			},
			checkInterval: 1 * 60 * 1000,
		}),
		VueI18nPlugin({
			include: [path.resolve(__dirname, './src/locales/**')],
			strictMessage: false,
			escapeHtml: false,
			runtimeOnly: true,
		}),
		tailwindcss(),
		// https://github.com/nonzzz/vite-bundle-analyzer
		// ANALYZE=server bun run build
		analyzeMode && unstableRolldownAdapter(analyzer({
			analyzerMode: analyzeMode === 'json' ? 'json' : 'server',
		})),
	]

	if (mode === 'development') {
		plugins.push(
			watch({
				onInit: true,
				pattern: 'src/api/**/*.ts',
				command: 'graphql-codegen',
			})
		)
	}

	return {
		plugins,
		base: '/dashboard',
		resolve: {
			alias: {
				vue: 'vue/dist/vue.esm-bundler.js',
				'@': fileURLToPath(new URL('./src', import.meta.url)),
			},
		},
		server: {
			port: 3006,
			host: '0.0.0.0',
			hmr: {
				protocol: env.USE_WSS === 'true' ? 'wss' : 'ws',
			},
		},
		clearScreen: false,

		build: {
			sourcemap: true,
			rolldownOptions: {
				output: {
					chunkFileNames: 'assets/[name]-[hash].js',
					entryFileNames: 'assets/[name]-[hash].js',
					assetFileNames: 'assets/[name]-[hash].[ext]',

					manualChunks: (id) => {
						if (id.includes('node_modules')) {
							if (id.includes('node_modules/.bun/')) {
								const match = id.match(/\.bun\/([^@/]+)@/)
								if (match && match[1]) {
									return `vendor-${match[1]}`
								}

								const simpleMatch = id.match(/\.bun\/([^/]+)/)
								if (simpleMatch && simpleMatch[1] && simpleMatch[1] !== '.bun') {
									return `vendor-${simpleMatch[1]}`
								}

								return 'vendor-common'
							}

							const match = id.match(/node_modules\/(@[^/]+\/[^/]+|[^/@]+)/)
							if (match && match[1]) {
								const packageName = match[1].replace('@', '').replace('/', '-')
								return `vendor-${packageName}`
							}

							return 'vendor-common'
						}

						if (id.includes('src/gql/')) {
							return 'gql'
						}

						if (id.includes('src/plugins/')) {
							return 'plugins'
						}

						if (id.includes('src/api/')) {
							const match = id.match(/src\/api\/([^/]+)/)
							if (match && match[1]) {
								const fileName = match[1].replace(/\.ts$/, '')
								return `api-${fileName}`
							}
							return 'api-common'
						}

						if (id.includes('src/composables/')) {
							return 'composables'
						}

						if (id.includes('src/features/')) {
							const match = id.match(/src\/features\/([^/]+)/)
							if (match) {
								return `feature-${match[1]}`
							}
						}

						if (id.includes('src/components/')) {
							const match = id.match(/src\/components\/([^/]+(?:\/[^/]+)?)/)
							if (match && match[1]) {
								const path = match[1].replace(/\.(ts|vue)$/, '').replace(/\//g, '-')
								return `components-${path}`
							}
							return 'components-common'
						}

						if (id.includes('src/pages/')) {
							const match = id.match(/src\/pages\/([^/]+)/)
							if (match) {
								return `page-${match[1].replace('.vue', '').toLowerCase()}`
							}
						}
					},
				},
			},
		},
	}
})
