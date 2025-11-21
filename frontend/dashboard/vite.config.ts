import path from 'node:path'
import process from 'node:process'
import { fileURLToPath } from 'node:url'

import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import { webUpdateNotice } from '@plugin-web-update-notification/vite'
import svgSprite from '@twirapp/vite-plugin-svg-spritemap'
import vue from '@vitejs/plugin-vue'
import autoprefixer from 'autoprefixer'
import tailwind from 'tailwindcss'
import { type PluginOption, defineConfig, loadEnv } from 'vite'
import { watch } from 'vite-plugin-watch'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, path.resolve(process.cwd(), '..', '..'), '')

	console.log(mode)

	const plugins: PluginOption[] = [
		// @ts-ignore
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
		css: {
			postcss: {
				plugins: [tailwind(), autoprefixer()],
			},
		},
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
		},
	}
})
