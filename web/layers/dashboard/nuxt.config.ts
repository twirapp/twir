import path from 'node:path'
import url from 'node:url'

const currentDir = path.dirname(url.fileURLToPath(import.meta.url))

export default defineNuxtConfig({
	modules: ['@nuxtjs/tailwindcss', '@nuxtjs/color-mode', '@vueuse/nuxt', '@nuxt/icon'],

	tailwindcss: {
		exposeConfig: true,
		configPath: path.join(currentDir, 'tailwind.config.js'),
	},

	colorMode: {
		classSuffix: '',
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
})
