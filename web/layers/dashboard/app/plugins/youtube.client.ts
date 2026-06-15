export default defineNuxtPlugin(() => {
	if (typeof window !== 'undefined' && !(window as any).YT) {
		const tag = document.createElement('script')
		tag.src = 'https://www.youtube.com/iframe_api'
		const firstScriptTag = document.getElementsByTagName('script')[0]
		firstScriptTag?.parentNode?.insertBefore(tag, firstScriptTag)
	}
})
