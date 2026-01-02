<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { VueMonacoEditor, useMonaco } from '@guolao/vue-monaco-editor'
import { Code2, Eye, EyeOff } from 'lucide-vue-next'

import { useChannelOverlayParseHtml } from '@/api/overlays/custom'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'

interface Props {
	open: boolean
	layerId?: string
	layerName?: string
	html?: string
	css?: string
	js?: string
	refreshInterval?: number
}

const props = withDefaults(defineProps<Props>(), {
	html: '',
	css: '',
	js: '',
	refreshInterval: 5,
})

const emit = defineEmits<{
	'update:open': [value: boolean]
	save: [data: { html: string; css: string; js: string; refreshInterval: number }]
}>()

const { monacoRef } = useMonaco()
const parseHtmlMutation = useChannelOverlayParseHtml()

// Local state
const localHtml = ref(props.html)
const localCss = ref(props.css)
const localJs = ref(props.js)
const localRefreshInterval = ref(props.refreshInterval)
const showPreview = ref(true)
const activeTab = ref('html')
const parsedHtml = ref('')
const pollInterval = ref<ReturnType<typeof setInterval>>()
const isLoading = ref(false)

// Watch props changes
watch(() => props.html, (newVal) => { localHtml.value = newVal })
watch(() => props.css, (newVal) => { localCss.value = newVal })
watch(() => props.js, (newVal) => { localJs.value = newVal })
watch(() => props.refreshInterval, (newVal) => { localRefreshInterval.value = newVal })

// Preview refs
const previewContainer = ref<HTMLDivElement>()
const previewContent = ref<HTMLDivElement>()
const styleElement = ref<HTMLStyleElement>()

const sanitizedHtml = computed(() => {
	// Use parsed HTML if available, otherwise use raw HTML
	const html = parsedHtml.value || localHtml.value
	return html || '<div style="display: flex; align-items: center; justify-content: center; height: 100%; color: rgba(255,255,255,0.5); font-size: 14px;">Enter HTML to preview</div>'
})

// Parse HTML with variables
async function parseHtml() {
	console.log('[Preview] parseHtml called')

	if (!localHtml.value) {
		console.log('[Preview] No HTML to parse')
		parsedHtml.value = ''
		return
	}

	console.log('[Preview] Parsing HTML:', localHtml.value.substring(0, 100))

	isLoading.value = true
	try {
		const result = await parseHtmlMutation.executeMutation({ html: localHtml.value })
		console.log('[Preview] Parse result:', result.data?.channelOverlayParseHtml?.substring(0, 100))

		parsedHtml.value = result.data?.channelOverlayParseHtml ?? localHtml.value

		console.log('[Preview] Calling executeScript after parse')
		// Call onDataUpdate after parsing
		executeScript()
	} catch (e) {
		console.error('[Preview] Failed to parse HTML:', e)
		parsedHtml.value = localHtml.value
	} finally {
		isLoading.value = false
	}
}

// Start periodic polling
function startPolling() {
	console.log('[Preview] startPolling called, interval:', localRefreshInterval.value)
	stopPolling()

	// Initial parse
	parseHtml()

	// Set up interval
	if (localRefreshInterval.value > 0) {
		console.log('[Preview] Setting up interval:', localRefreshInterval.value * 1000, 'ms')
		pollInterval.value = setInterval(() => {
			console.log('[Preview] Interval tick - parsing HTML')
			parseHtml()
		}, localRefreshInterval.value * 1000)
	}
}

// Stop polling
function stopPolling() {
	if (pollInterval.value) {
		console.log('[Preview] Stopping polling')
		clearInterval(pollInterval.value)
		pollInterval.value = undefined
	}
}

// Apply CSS by injecting a style element
function updateStyles() {
	if (!previewContainer.value) return

	// Remove old style element if exists
	if (styleElement.value) {
		styleElement.value.remove()
		styleElement.value = undefined
	}

	// Create new style element with scoped styles
	if (localCss.value) {
		const style = document.createElement('style')
		style.textContent = localCss.value
		styleElement.value = style
		previewContainer.value.appendChild(style)
	}
}

// Execute JavaScript
function executeScript() {
	if (!localJs.value) return

	try {
		// eslint-disable-next-line no-new-func
		const scriptFunc = new Function('container', localJs.value)
		scriptFunc(previewContent.value)
	} catch (e) {
		console.error('Preview JS Error:', e)
	}
}

// Watch for code changes
watch([localHtml, localCss, localJs], () => {
	setTimeout(() => {
		updateStyles()
		// Re-parse HTML when it changes
		if (showPreview.value) {
			parseHtml()
		}
	}, 50)
})

// Watch refresh interval changes
watch(localRefreshInterval, () => {
	if (showPreview.value && props.open) {
		startPolling()
	}
})

// Watch when dialog opens to initialize preview
watch(() => props.open, (isOpen) => {
	if (isOpen && showPreview.value) {
		setTimeout(() => {
			updateStyles()
			startPolling()
		}, 100)
	} else {
		stopPolling()
	}
})

// Watch preview toggle
watch(showPreview, (show) => {
	if (show && props.open) {
		startPolling()
	} else {
		stopPolling()
	}
})

function handleSave() {
	emit('save', {
		html: localHtml.value,
		css: localCss.value,
		js: localJs.value,
		refreshInterval: localRefreshInterval.value,
	})
	emit('update:open', false)
}

function handleCancel() {
	// Reset to props values
	localHtml.value = props.html
	localCss.value = props.css
	localJs.value = props.js
	localRefreshInterval.value = props.refreshInterval
	emit('update:open', false)
}

onMounted(() => {
	// Configure Monaco themes if needed
	if (monacoRef.value) {
		monacoRef.value.editor.defineTheme('twir-dark', {
			base: 'vs-dark',
			inherit: true,
			rules: [],
			colors: {
				'editor.background': '#0b0b0c',
			},
		})
	}

	// Initialize preview
	setTimeout(() => {
		if (showPreview.value && props.open) {
			updateStyles()
			startPolling()
		}
	}, 200)
})

onUnmounted(() => {
	stopPolling()
	if (styleElement.value) {
		styleElement.value.remove()
	}
})
</script>

<template>
	<Dialog :open="open" @update:open="emit('update:open', $event)">
		<DialogContent
			class="h-[90vh] flex flex-col p-0"
			:style="{ maxWidth: '95vw', width: '95vw' }"
			@keydown.stop
			@keyup.stop
			@keypress.stop
		>
			<DialogHeader class="px-6 pt-6 pb-4 border-b">
				<DialogTitle class="flex items-center gap-2">
					<Code2 class="h-5 w-5" />
					<span>Edit HTML Layer</span>
					<span v-if="layerName" class="text-muted-foreground font-normal">
						- {{ layerName }}
					</span>
				</DialogTitle>
				<DialogDescription>
					Edit HTML, CSS, and JavaScript for this layer. Changes are previewed in real-time.
				</DialogDescription>
			</DialogHeader>

			<div class="flex-1 flex overflow-hidden">
				<!-- Code Editor Side -->
				<div class="flex-1 flex flex-col border-r">
					<!-- Settings Bar -->
					<div class="flex items-center gap-4 px-4 py-3 border-b bg-muted/30">
						<div class="flex items-center gap-2">
							<Label for="refresh-interval" class="text-xs">Refresh Interval (seconds):</Label>
							<input
								id="refresh-interval"
								v-model.number="localRefreshInterval"
								type="number"
								min="1"
								max="60"
								class="w-16 px-2 py-1 text-xs border rounded bg-background"
							/>
						</div>

						<div class="flex items-center gap-2 ml-auto">
							<Switch
								id="preview-toggle"
								:checked="showPreview"
								@update:checked="showPreview = $event"
							/>
							<Label for="preview-toggle" class="text-xs cursor-pointer flex items-center gap-1">
								<Eye v-if="showPreview" class="h-3 w-3" />
								<EyeOff v-else class="h-3 w-3" />
								Preview
							</Label>
						</div>
					</div>

					<!-- Tabs -->
					<Tabs v-model="activeTab" class="flex-1 flex flex-col" @keydown.stop @keyup.stop>
						<TabsList class="w-full justify-start rounded-none border-b bg-muted/30 px-4">
							<TabsTrigger value="html">HTML</TabsTrigger>
							<TabsTrigger value="css">CSS</TabsTrigger>
							<TabsTrigger value="js">JavaScript</TabsTrigger>
						</TabsList>

						<TabsContent value="html" class="flex-1 mt-0 p-0" @keydown.stop @keyup.stop>
							<VueMonacoEditor
								v-model:value="localHtml"
								language="html"
								theme="vs-dark"
								:options="{
									automaticLayout: true,
									minimap: { enabled: false },
									fontSize: 14,
									lineNumbers: 'on',
									scrollBeyondLastLine: false,
									wordWrap: 'on',
									tabSize: 2,
									contextmenu: true,
									selectOnLineNumbers: true,
									quickSuggestions: true,
									suggest: {
										snippetsPreventQuickSuggestions: false
									},
									readOnly: false,
									domReadOnly: false
								}"
								class="h-full"
							/>
						</TabsContent>

						<TabsContent value="css" class="flex-1 mt-0 p-0" @keydown.stop @keyup.stop>
							<VueMonacoEditor
								v-model:value="localCss"
								language="css"
								theme="vs-dark"
								:options="{
									automaticLayout: true,
									minimap: { enabled: false },
									fontSize: 14,
									lineNumbers: 'on',
									scrollBeyondLastLine: false,
									wordWrap: 'on',
									tabSize: 2,
									contextmenu: true,
									selectOnLineNumbers: true,
									quickSuggestions: true,
									suggest: {
										snippetsPreventQuickSuggestions: false
									},
									readOnly: false,
									domReadOnly: false
								}"
								class="h-full"
							/>
						</TabsContent>

						<TabsContent value="js" class="flex-1 mt-0 p-0" @keydown.stop @keyup.stop>
							<VueMonacoEditor
								v-model:value="localJs"
								language="javascript"
								theme="vs-dark"
								:options="{
									automaticLayout: true,
									minimap: { enabled: false },
									fontSize: 14,
									lineNumbers: 'on',
									scrollBeyondLastLine: false,
									wordWrap: 'on',
									tabSize: 2,
									contextmenu: true,
									selectOnLineNumbers: true,
									quickSuggestions: true,
									suggest: {
										snippetsPreventQuickSuggestions: false
									},
									readOnly: false,
									domReadOnly: false
								}"
								class="h-full"
							/>
						</TabsContent>
					</Tabs>
				</div>

				<!-- Preview Side -->
				<div v-if="showPreview" class="w-150 flex flex-col bg-slate-900">
					<div class="px-4 py-3 border-b bg-muted/30 flex items-center justify-between">
						<h3 class="text-sm font-medium">Preview</h3>
						<div v-if="isLoading" class="flex items-center gap-2 text-xs text-muted-foreground">
							<div class="w-3 h-3 border-2 border-primary border-t-transparent rounded-full animate-spin" />
							<span>Parsing...</span>
						</div>
						<div v-else-if="parsedHtml" class="text-xs text-green-500">
							✓ Live
						</div>
					</div>
					<div class="flex-1 p-4 overflow-auto">
						<div
							ref="previewContainer"
							class="w-full h-full bg-[#1a1a1a] rounded border border-slate-700 p-4 overflow-auto"
						>
							<div
								ref="previewContent"
								class="preview-content w-full h-full"
								v-html="sanitizedHtml"
							/>
						</div>
					</div>
				</div>
			</div>

			<DialogFooter class="px-6 py-4 border-t">
				<Button variant="outline" @click="handleCancel">
					Cancel
				</Button>
				<Button @click="handleSave">
					Save Changes
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>

<style scoped>
:deep(.monaco-editor) {
	height: 100%;
}

:deep(.tabs-content) {
	height: 100%;
}

.preview-content {
	color: #fff;
	font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.preview-content :deep(*) {
	box-sizing: border-box;
}
</style>
