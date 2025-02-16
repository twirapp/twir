function getRequestLanguage() {
	const nuxtApp = useNuxtApp()

	if (import.meta.server) {
		const reqLocale = nuxtApp.ssrContext?.event.node.req.headers['accept-language']?.split(',')[0]
		return reqLocale || 'en-US'
	}

	return navigator.language
}

export function formatNumber(num: number): string {
	const reqLocale = getRequestLanguage()
	const intl = new Intl.NumberFormat(reqLocale)

	return intl.format(num)
}
