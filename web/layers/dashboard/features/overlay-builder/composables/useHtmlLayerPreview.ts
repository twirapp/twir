import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

import { useChannelOverlayParseHtml } from '~~/layers/dashboard/api/overlays/custom'

export interface HtmlLayerPreviewProps {
	html?: string
	css?: string
	js?: string
	width: number
	height: number
	refreshInterval?: number
}

export function useHtmlLayerPreview(props: HtmlLayerPreviewProps) {
	const { t } = useI18n()
	const containerRef = ref<HTMLDivElement>()
	const shadowRoot = ref<ShadowRoot>()
	const renderKey = ref(0)
	const parsedHtml = ref('')
	const pollInterval = ref<ReturnType<typeof setInterval>>()
	const parseHtmlMutation = useChannelOverlayParseHtml()

	const sanitizedHtml = computed(() => {
		const html = parsedHtml.value || props.html
		return html || `<div style="display: flex; align-items: center; justify-content: center; height: 100%; color: rgba(255,255,255,0.5); font-size: 14px;">${t('overlayBuilder.codeEditor.emptyHtml')}</div>`
	})

	async function parseHtml() {
		if (!props.html) {
			parsedHtml.value = ''
			return
		}

		try {
			const result = await parseHtmlMutation.executeMutation({ html: props.html })
			parsedHtml.value = result.data?.channelOverlayParseHtml ?? props.html
			executeScript()
		} catch (error) {
			console.error('Canvas Layer: Failed to parse HTML:', error)
			parsedHtml.value = props.html
		}
	}

	function stopPolling() {
		if (pollInterval.value) {
			clearInterval(pollInterval.value)
			pollInterval.value = undefined
		}
	}

	function startPolling() {
		stopPolling()
		void parseHtml()
		if (props.refreshInterval && props.refreshInterval > 0) {
			pollInterval.value = setInterval(() => void parseHtml(), props.refreshInterval * 1000)
		}
	}

	function executeScript() {
		if (!props.js || !shadowRoot.value) return
		try {
			const contentElement = shadowRoot.value.querySelector('.html-content')
			if (!contentElement) return
			// eslint-disable-next-line no-new-func
			const scriptFunc = new Function('container', props.js)
			scriptFunc(contentElement)
		} catch (error) {
			console.error('Canvas Layer JS Error:', error)
		}
	}

	function renderInShadowDOM() {
		if (!shadowRoot.value) return
		const baseStyles = `
			* { box-sizing: border-box; }
			:host { display: block; width: 100%; height: 100%; overflow: hidden; background: transparent; color: #fff; font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; }
			.html-content { width: 100%; height: 100%; pointer-events: none; }
		`
		shadowRoot.value.innerHTML = `<style>${baseStyles + (props.css || '')}</style><div class="html-content">${sanitizedHtml.value}</div>`
		setTimeout(executeScript, 0)
	}

	function forceUpdate() {
		renderKey.value++
		setTimeout(renderInShadowDOM, 0)
	}

	function initializeShadowDOM() {
		if (!containerRef.value || shadowRoot.value) return
		try {
			shadowRoot.value = containerRef.value.attachShadow({ mode: 'open' })
		} catch (error) {
			console.error('Failed to create shadow DOM:', error)
			return
		}
		renderInShadowDOM()
	}

	watch(() => props.html, () => {
		forceUpdate()
		startPolling()
	})
	watch(() => props.css, renderInShadowDOM)
	watch(() => props.js, renderInShadowDOM)
	watch(() => props.refreshInterval, startPolling)
	watch(sanitizedHtml, renderInShadowDOM)

	onMounted(() => {
		initializeShadowDOM()
		startPolling()
	})
	onUnmounted(() => {
		stopPolling()
		if (shadowRoot.value) shadowRoot.value.innerHTML = ''
	})

	return { containerRef, renderKey }
}
