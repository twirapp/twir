<script setup lang="ts">
import { VueMonacoEditor } from '@guolao/vue-monaco-editor'
import type { CodeEditorDialogProps } from '../composables/useCodeEditorDialog'
import { useCodeEditorDialog } from '../composables/useCodeEditorDialog'
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
import { Input } from '@/components/ui/input'

const props = withDefaults(defineProps<CodeEditorDialogProps>(), {
	html: '',
	css: '',
	js: '',
	refreshInterval: 5,
})
const { t } = useI18n()

const emit = defineEmits<{
	'update:open': [value: boolean]
	save: [data: { html: string; css: string; js: string; refreshInterval: number }]
}>()

const { localHtml, localCss, localJs, localRefreshInterval, showPreview, activeTab, isLoading, showVariablesPanel, variablesSearchQuery, htmlEditorRef, copiedVariableId, filteredVariables, parsedHtml, sanitizedHtml, previewContainer, previewContent, handleSave, handleCancel, copyVariable } = useCodeEditorDialog(props, emit)
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
					<Icon name="lucide:code-xml" class="h-5 w-5" />
					<span>{{ t('overlayBuilder.codeEditor.title') }}</span>
					<span v-if="layerName" class="text-muted-foreground font-normal">
						- {{ layerName }}
					</span>
				</DialogTitle>
				<DialogDescription>
					{{ t('overlayBuilder.codeEditor.description') }}
				</DialogDescription>
			</DialogHeader>

			<div class="flex-1 flex overflow-hidden">
				<!-- Code Editor Side -->
				<div class="flex flex-col border-r" style="flex: 1 1 0; min-width: 0;">
					<!-- Settings Bar -->
					<div class="flex items-center gap-4 px-4 py-3 border-b bg-muted/30">
						<div class="flex items-center gap-2">
							<Label for="refresh-interval" class="text-xs">{{ t('overlayBuilder.codeEditor.refreshInterval') }}</Label>
							<input
								id="refresh-interval"
								v-model.number="localRefreshInterval"
								type="number"
								min="1"
								max="60"
								class="w-16 px-2 py-1 text-xs border rounded bg-background"
								@keydown.stop
							/>
						</div>

						<div class="flex items-center gap-2 ml-auto">
							<Button
								variant="outline"
								size="sm"
								class="h-7 text-xs gap-1.5"
								@click="showVariablesPanel = !showVariablesPanel"
							>
								<Icon name="lucide:chevron-left" v-if="showVariablesPanel" class="h-3 w-3" />
								<Icon name="lucide:chevron-right" v-else class="h-3 w-3" />
								{{ t('overlayBuilder.codeEditor.variables') }}
							</Button>

							<Switch
								id="preview-toggle"
								:model-value="showPreview"
								@update:model-value="showPreview = $event"
							/>
							<Label for="preview-toggle" class="text-xs cursor-pointer flex items-center gap-1">
								<Icon name="lucide:eye" v-if="showPreview" class="h-3 w-3" />
								<Icon name="lucide:eye-off" v-else class="h-3 w-3" />
								{{ t('overlayBuilder.codeEditor.preview') }}
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
								@mount="(editor) => htmlEditorRef = editor"
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

				<!-- Variables Panel -->
				<div
					v-if="showVariablesPanel"
					class="flex flex-col border-r bg-background overflow-hidden"
					style="width: 400px;"
				>
					<div class="px-4 py-3 border-b shrink-0">
										<h3 class="text-sm font-semibold mb-2">{{ t('overlayBuilder.codeEditor.availableVariables') }}</h3>
						<div class="relative">
							<Icon name="lucide:search" class="absolute left-2 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
							<Input
								v-model="variablesSearchQuery"
														:placeholder="t('overlayBuilder.codeEditor.searchVariables')"
								class="pl-8 h-8 text-xs"
								@keydown.stop
							/>
						</div>
					</div>

					<div class="flex-1 overflow-y-auto p-2 space-y-1">
						<div
							v-for="variable in filteredVariables"
							:key="variable.name"
							class="group relative rounded-lg border p-3 hover:bg-accent/50 transition-colors"
						>
								<div class="flex items-start justify-between gap-2 mb-1">
									<div class="flex-1 min-w-0">
										<div class="flex items-center gap-2">
											<code class="text-xs font-mono font-semibold text-primary">
												{{ variable.name }}
											</code>
											<span
												v-if="'isBuiltIn' in variable && variable.isBuiltIn"
												class="text-[10px] px-1.5 py-0.5 rounded bg-blue-500/10 text-blue-500 font-medium"
											>
														{{ t('overlayBuilder.codeEditor.builtIn') }}
											</span>
											<span
												v-else
												class="text-[10px] px-1.5 py-0.5 rounded bg-purple-500/10 text-purple-500 font-medium"
											>
														{{ t('overlayBuilder.codeEditor.custom') }}
											</span>
										</div>
										<p v-if="variable.description" class="text-xs text-muted-foreground mt-1 line-clamp-2">
											{{ variable.description }}
										</p>
									</div>
								</div>

								<div class="flex items-center gap-1 mt-2">
									<code class="flex-1 text-[11px] px-2 py-1 rounded bg-muted/50 font-mono text-muted-foreground truncate">
										$({{ variable.example }})
									</code>
									<Button
										variant="ghost"
										size="icon"
										class="h-6 w-6 shrink-0"
										:title="copiedVariableId === variable.name ? t('overlayBuilder.codeEditor.copied') : t('overlayBuilder.codeEditor.copyToClipboard')"
										@click="copyVariable(variable)"
									>
										<Icon name="lucide:check" v-if="copiedVariableId === variable.name" class="h-3 w-3 text-green-500" />
										<Icon name="lucide:copy" v-else class="h-3 w-3" />
									</Button>
								</div>

								<div v-if="'links' in variable && variable.links && variable.links.length > 0" class="flex flex-wrap gap-1 mt-2">
									<a
										v-for="link in variable.links"
										:key="link.href"
										:href="link.href"
										target="_blank"
										rel="noopener noreferrer"
										class="inline-flex items-center gap-1 text-[10px] text-blue-500 hover:text-blue-600 hover:underline"
									>
										<Icon name="lucide:external-link" class="h-2.5 w-2.5" />
										{{ link.name }}
									</a>
								</div>
							</div>

						<div
							v-if="filteredVariables.length === 0"
							class="text-center py-8 text-sm text-muted-foreground"
						>
							<Icon name="lucide:search" class="h-8 w-8 mx-auto mb-2 opacity-50" />
							<p>{{ t('overlayBuilder.codeEditor.noVariables') }}</p>
						</div>
					</div>

					<div class="px-4 py-2 border-t text-xs text-muted-foreground shrink-0">
						{{ t('overlayBuilder.codeEditor.variableCount', filteredVariables.length) }}
					</div>
				</div>

				<!-- Preview Side -->
				<div v-if="showPreview" class="w-150 flex flex-col bg-slate-900">
					<div class="px-4 py-3 border-b bg-muted/30 flex items-center justify-between">
						<h3 class="text-sm font-medium">{{ t('overlayBuilder.codeEditor.preview') }}</h3>
						<div v-if="isLoading" class="flex items-center gap-2 text-xs text-muted-foreground">
							<div class="w-3 h-3 border-2 border-primary border-t-transparent rounded-full animate-spin" />
							<span>{{ t('overlayBuilder.codeEditor.parsing') }}</span>
						</div>
						<div v-else-if="parsedHtml" class="text-xs text-green-500">
							{{ t('overlayBuilder.codeEditor.live') }}
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
					{{ t('overlayBuilder.codeEditor.cancel') }}
				</Button>
				<Button @click="handleSave">
					{{ t('overlayBuilder.codeEditor.saveChanges') }}
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
