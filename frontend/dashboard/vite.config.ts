import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import tailwindcss from '@tailwindcss/vite'
import { webUpdateNotice } from '@plugin-web-update-notification/vite'
import svgSprite from '@twirapp/vite-plugin-svg-spritemap'
import vue from '@vitejs/plugin-vue'
import path from 'node:path'
import process from 'node:process'
import { fileURLToPath } from 'node:url'
import { type PluginOption, defineConfig, loadEnv } from 'vite'
import { analyzer, unstableRolldownAdapter } from 'vite-bundle-analyzer'
import { watch } from 'vite-plugin-watch'

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
		analyzeMode &&
			unstableRolldownAdapter(
				analyzer({
					analyzerMode: analyzeMode === 'json' ? 'json' : 'server',
				})
			),
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
					advancedChunks: {
						groups: [
							{
								test: /@vue\+/,
								name: 'vue',
							},
							{
								test: /@vuepic/,
								name: 'vuepic',
							},
							{
								test: /vue-draggable-plus/,
								name: 'vue-draggable-plus',
							},
							{
								test: /vexip/,
								name: 'vexip',
							},
							{
								test: /vee-validate/,
								name: 'vee-validate',
							},
							{
								test: /@floating-ui/,
								name: 'floating-ui',
							},
							{
								test: /@formkit\+drag-and-drop/,
								name: 'formkit-drag-and-drop',
							},
							{
								test: /@tanstack\+table/,
								name: 'tanstack-table',
							},
							{
								test: /@tanstack\+query/,
								name: 'tanstack-query',
							},
							{
								test: /tailwind/,
								name: 'tailwind',
							},
							{
								test: /zod/,
								name: 'zod',
							},
							{
								test: /date-fns/,
								name: 'date-fns',
							},
							{
								test: /grid-layout-plus/,
								name: 'grid-layout-plus',
							},
							{
								test: /lucide-vue-next/,
								name: 'lucide-vue-next',
							},
							{
								test: /@tabler\+/,
								name: 'tabler',
							},
							{
								test: /@discord-message/,
								name: 'discord-message',
							},
							{
								test: /editorjs|editor-js/,
								name: 'editorjs',
							},
							{
								test: /lightweight-charts/,
								name: 'lightweight-charts',
							},
							{
								test: /croact-moveable/,
								name: 'croact-moveable',
							},
							{
								test: /reka-ui/,
								name: 'reka-ui',
							},
							{
								test: /interactjs/,
								name: 'interactjs',
							},
							{
								test: /\/src\/gql/,
								name: 'gql',
							},
							{
								test: /\/src\/api/,
								name: 'api',
							},
							{
								test: /\/src\/components\/ui/,
								name: 'components-ui',
							},
							{
								test: /\/src\/components\/dashboard/,
								name: 'components-dashboard',
							},
							{
								test: /\/src\/components/,
								name: 'components',
							},
							{
								test: /\/src\/features\/overlays/,
								name: 'features-overlays',
							},
							{
								test: /\/src\/features\/overlay-builder/,
								name: 'features-overlay-builder',
							},
							{
								test: /\/src\/features\/admin-panel/,
								name: 'features-admin-panel',
							},
							{
								test: /\/src\/features\/events/,
								name: 'features-events',
							},
							{
								test: /\/src\/features\/giveaways/,
								name: 'features-giveaways',
							},
							{
								test: /\/src\/features\/bot-settings/,
								name: 'features-bot-settings',
							},
							{
								test: /\/src\/features\/commands/,
								name: 'features-commands',
							},
							{
								test: /\/src\/features\/integrations/,
								name: 'features-integrations',
							},
							{
								test: /\/src\/features\/moderation/,
								name: 'features-moderation',
							},
							{
								test: /\/src\/features\/community/,
								name: 'features-community',
							},
							{
								test: /\/src\/features\/dudes-settings/,
								name: 'features-dudes-settings',
							},
							{
								test: /\/src\/features/,
								name: 'features-misc',
							},
						],
					},
				},
			},
		},
	}
})
