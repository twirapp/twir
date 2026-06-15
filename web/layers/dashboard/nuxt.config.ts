export default defineNuxtConfig({
	routeRules: {
		'/dashboard': { ssr: false },
		'/dashboard/**': { ssr: false },
		'/**/dashboard': { ssr: false },
		'/**/dashboard/**': { ssr: false },
	},
})
