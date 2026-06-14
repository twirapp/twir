export default defineNuxtPlugin(() => {
	if (typeof window === 'undefined') return

	const origCreate = document.createElement.bind(document)
	document.createElement = function (name: string, options?: ElementCreationOptions) {
		if (!name || typeof name !== 'string' || name.includes(':') || name.includes(' ') || /^[0-9]/.test(name) || name.length > 50) {
			console.error('[DEBUG] Invalid createElement name:', JSON.stringify(name?.slice?.(0, 100)), new Error().stack)
		}
		return origCreate(name, options)
	}
})
