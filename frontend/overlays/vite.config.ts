import path from 'node:path'

import { webUpdateNotice } from '@plugin-web-update-notification/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import { watch } from 'vite-plugin-watch'

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [
		watch({
			onInit: true,
			pattern: 'src/**/*.ts',
			command: 'graphql-codegen',
		}),
		vue(),
		webUpdateNotice({
			hiddenDefaultNotification: true,
			checkInterval: 1 * 60 * 1000,
		}),
	],
	resolve: {
		alias: {
			'@': path.resolve(__dirname, './src'),
		},
	},
	base: '/overlays',
	server: {
		host: true,
		port: 3008,
	},
	clearScreen: false,
})
