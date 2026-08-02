/// <reference types="@plugin-web-update-notification/nuxt" />

export function useAutoUpdate() {
	const where = useRoute()

	const handler = () => {
		if (!import.meta.client || !globalThis.window) return
		if (where.path.startsWith('/o/')) {
			window.location.reload()
			return
		}
	}

	onMounted(() => {
		if (!import.meta.client) return
		document.body.addEventListener('plugin_web_update_notice', handler)
	})

	onUnmounted(() => {
		if (!import.meta.client) return

		document.body.removeEventListener('plugin_web_update_notice', handler)
	})
}
