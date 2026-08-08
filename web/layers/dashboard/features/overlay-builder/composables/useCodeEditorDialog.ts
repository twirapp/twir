import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useClipboard } from '@vueuse/core'
import { useMonaco } from '@guolao/vue-monaco-editor'

import { useChannelOverlayParseHtml } from '~~/layers/dashboard/api/overlays/custom'
import { useVariablesApi } from '~~/layers/dashboard/api/variables'

export interface CodeEditorDialogProps {
	open: boolean
	layerId?: string
	layerName?: string
	html?: string
	css?: string
	js?: string
	refreshInterval?: number
}

export interface CodeEditorSaveData {
	html: string
	css: string
	js: string
	refreshInterval: number
}

interface CodeEditorDialogEmit {
	(event: 'update:open', value: boolean): void
	(event: 'save', data: CodeEditorSaveData): void
}

export function useCodeEditorDialog(props: CodeEditorDialogProps, emit: CodeEditorDialogEmit) {
	const { t } = useI18n()
	const { monacoRef } = useMonaco()
	const parseHtmlMutation = useChannelOverlayParseHtml()
	const variablesApi = useVariablesApi()
	const { copy } = useClipboard()

	const localHtml = ref(props.html ?? '')
	const localCss = ref(props.css ?? '')
	const localJs = ref(props.js ?? '')
	const localRefreshInterval = ref(props.refreshInterval ?? 5)
	const showPreview = ref(true)
	const activeTab = ref('html')
	const parsedHtml = ref('')
	const pollInterval = ref<ReturnType<typeof setInterval>>()
	const isLoading = ref(false)
	const showVariablesPanel = ref(false)
	const variablesSearchQuery = ref('')
	const htmlEditorRef = ref<unknown>(null)
	const copiedVariableId = ref<string | null>(null)
	const previewContainer = ref<HTMLDivElement>()
	const previewContent = ref<HTMLDivElement>()
	const styleElement = ref<HTMLStyleElement>()

	watch(() => props.html, (value) => { localHtml.value = value ?? '' })
	watch(() => props.css, (value) => { localCss.value = value ?? '' })
	watch(() => props.js, (value) => { localJs.value = value ?? '' })
	watch(() => props.refreshInterval, (value) => { localRefreshInterval.value = value ?? 5 })

	const filteredVariables = computed(() => {
		const query = variablesSearchQuery.value.toLowerCase().trim()
		if (!query) return variablesApi.allVariables.value
		return variablesApi.allVariables.value.filter((variable) => (
			variable.name.toLowerCase().includes(query) ||
			variable.description?.toLowerCase().includes(query) ||
			variable.example?.toLowerCase().includes(query)
		))
	})

	function formatVariableForInsertion(variable: typeof variablesApi.allVariables.value[number]) {
		return `$(${variable.example})`
	}

	async function copyVariable(variable: typeof variablesApi.allVariables.value[number]) {
		await copy(formatVariableForInsertion(variable))
		copiedVariableId.value = variable.name
		setTimeout(() => { copiedVariableId.value = null }, 2000)
	}

	const sanitizedHtml = computed(() => {
		const html = parsedHtml.value || localHtml.value
		return html || `<div style="display: flex; align-items: center; justify-content: center; height: 100%; color: rgba(255,255,255,0.5); font-size: 14px;">${t('overlayBuilder.codeEditor.emptyHtml')}</div>`
	})

	async function parseHtml() {
		if (!localHtml.value) {
			parsedHtml.value = ''
			return
		}

		isLoading.value = true
		try {
			const result = await parseHtmlMutation.executeMutation({ html: localHtml.value })
			parsedHtml.value = result.data?.channelOverlayParseHtml ?? localHtml.value
			executeScript()
		} catch (error) {
			console.error('[Preview] Failed to parse HTML:', error)
			parsedHtml.value = localHtml.value
		} finally {
			isLoading.value = false
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
		if (localRefreshInterval.value > 0) pollInterval.value = setInterval(() => void parseHtml(), localRefreshInterval.value * 1000)
	}

	function updateStyles() {
		if (!previewContainer.value) return
		styleElement.value?.remove()
		styleElement.value = undefined
		if (localCss.value) {
			const style = document.createElement('style')
			style.textContent = localCss.value
			styleElement.value = style
			previewContainer.value.appendChild(style)
		}
	}

	function executeScript() {
		if (!localJs.value) return
		try {
			// eslint-disable-next-line no-new-func
			const scriptFunc = new Function('container', localJs.value)
			scriptFunc(previewContent.value)
		} catch (error) {
			console.error('Preview JS Error:', error)
		}
	}

	watch([localHtml, localCss, localJs], () => {
		setTimeout(() => {
			updateStyles()
			if (showPreview.value) void parseHtml()
		}, 50)
	})
	watch(localRefreshInterval, () => {
		if (showPreview.value && props.open) startPolling()
	})
	watch(() => props.open, (isOpen) => {
		if (isOpen && showPreview.value) setTimeout(() => { updateStyles(); startPolling() }, 100)
		else stopPolling()
	})
	watch(showPreview, (show) => {
		if (show && props.open) startPolling()
		else stopPolling()
	})

	function handleSave() {
		emit('save', { html: localHtml.value, css: localCss.value, js: localJs.value, refreshInterval: localRefreshInterval.value })
		emit('update:open', false)
	}

	function handleCancel() {
		localHtml.value = props.html ?? ''
		localCss.value = props.css ?? ''
		localJs.value = props.js ?? ''
		localRefreshInterval.value = props.refreshInterval ?? 5
		emit('update:open', false)
	}

	onMounted(() => {
		if (monacoRef.value) monacoRef.value.editor.defineTheme('twir-dark', { base: 'vs-dark', inherit: true, rules: [], colors: { 'editor.background': '#0b0b0c' } })
		setTimeout(() => {
			if (showPreview.value && props.open) { updateStyles(); startPolling() }
		}, 200)
	})
	onUnmounted(() => {
		stopPolling()
		styleElement.value?.remove()
	})

	return {
		localHtml,
		localCss,
		localJs,
		localRefreshInterval,
		showPreview,
		activeTab,
		isLoading,
		showVariablesPanel,
		variablesSearchQuery,
		htmlEditorRef,
		copiedVariableId,
		filteredVariables,
		parsedHtml,
		sanitizedHtml,
		previewContainer,
		previewContent,
		handleSave,
		handleCancel,
		copyVariable,
	}
}
