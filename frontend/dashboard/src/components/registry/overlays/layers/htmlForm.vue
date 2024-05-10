<script lang="ts" setup>
import {
	NAlert,
	NFormItem,
	NInputNumber,
	NModal,
	NSelect,
	NSwitch,
	NTabPane,
	NTabs,
	useMessage,
} from 'naive-ui'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import BaseLayer, { type LayerProps } from './layer.vue'

import { useCommandsApi } from '@/api/commands/commands'
import { useKeywordsApi } from '@/api/keywords'
import { useVariablesApi } from '@/api/variables'
import { copyToClipBoard } from '@/helpers/index.js'

defineProps<LayerProps>()

defineEmits<{
	focus: [index: number]
	remove: [index: number]
}>()

const html = defineModel('html')
const css = defineModel('css')
const js = defineModel('js')
const pollInterval = defineModel('pollInterval', { default: 5 })
const periodicallyRefetchData = defineModel<boolean>('periodicallyRefetchData')

const showModal = ref(false)

const { allVariables } = useVariablesApi()
const keywordsManager = useKeywordsApi()
const { data: keywords } = keywordsManager.useQueryKeywords()

const commandsManager = useCommandsApi()
const { data: commands } = commandsManager.useQueryCommands()

const messages = useMessage()
async function copyVariable(v: string) {
	await copyToClipBoard(v)
	messages.success('Copied')
}

const variables = computed(() => {
	const k = keywords.value?.keywords ?? []
	const cmds = commands.value?.commands ?? []

	return [
		...allVariables.value
			.filter(v => v.canBeUsedInRegistry)
			.map(v => {
				const name = `$(${v.isBuiltIn ? v.name : `customvar|${v.name}`})`
				return {
					label: `${name} - ${v.description || 'Your custom variable'}`,
					value: name,
				}
			}),
		...k
			.map(k => ({
				label: `$(keywords.counter|${k.id}) - How many times "${k.text}" was used`,
				value: `$(keywords.counter|${k.id})`,
			})),
		...cmds
			.map(c => ({
				label: `$(command.counter.fromother|${c.name}) - How many times "${c.name}" was used`,
				value: `$(command.counter.fromother|${c.name})`,
			})),
	]
})

const { t } = useI18n()
</script>

<template>
	<BaseLayer
		:is-focused="isFocused"
		:layer-index="layerIndex"
		:type="type"
		@open-settings="showModal = true"
		@focus="$emit('focus', layerIndex)"
		@remove="$emit('remove', layerIndex)"
	/>

	<NModal
		v-model:show="showModal"
		preset="card"
		title="Settings"
		style="width: 50vw"
	>
		<div class="flex flex-col gap-5">
			<NFormItem :label="t('overlaysRegistry.html.periodicallyRefetchData')">
				<NSwitch v-model:value="periodicallyRefetchData" />
			</NFormItem>

			<NFormItem :label="t('overlaysRegistry.html.updateInterval')">
				<NInputNumber v-model:value="pollInterval" :min="5" :max="300" />
			</NFormItem>

			<NAlert type="info" :title="t('overlaysRegistry.html.variablesAlert.title')">
				{{ t('overlaysRegistry.html.variablesAlert.selectToCopy') }}
				<NSelect filterable :options="variables" @update:value="(v) => copyVariable(v)" />
			</NAlert>

			<NTabs type="segment">
				<NTabPane name="HTML">
					<vue-monaco-editor
						v-model:value="html"
						theme="vs-dark"
						height="500px"
						language="html"
						:options="{
							readOnlyMessage: 'test',
						}"
					/>
				</NTabPane>

				<NTabPane name="CSS">
					<vue-monaco-editor
						v-model:value="css"
						theme="vs-dark"
						height="500px"
						language="css"
					/>
				</NTabPane>

				<NTabPane name="JS">
					<vue-monaco-editor
						v-model:value="js"
						theme="vs-dark"
						height="500px"
						language="javascript"
					/>
				</NTabPane>
			</NTabs>
		</div>
	</NModal>
</template>
