// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	extends: ['../urql'],

	modules: [
		'@nuxtjs/tailwindcss',
		'@nuxt/image',
		'@nuxt/fonts',
		'nuxt-svgo',
		'@vueuse/nuxt',
		'@nuxt-alt/proxy',
		'nuxt-typed-router',
	],
})
