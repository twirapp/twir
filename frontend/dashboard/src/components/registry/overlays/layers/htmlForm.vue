<!-- eslint-disable no-undef -->
<script lang="ts" setup>
import {
	NModal,
	NInputNumber,
	NFormItem,
	NAlert,
	NSelect,
	useMessage,
	NSwitch,
	NTabs,
	NTabPane,
} from 'naive-ui';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import BaseLayer, { type LayerProps } from './layer.vue';

import { useAllVariables, useKeywordsManager, useCommandsManager } from '@/api/index.js';
import { copyToClipBoard } from '@/helpers/index.js';

defineProps<LayerProps>();

defineEmits<{
	focus: [index: number]
	remove: [index: number]
}>();

const html = defineModel('html');
const css = defineModel('css');
const js = defineModel('js');
const pollInterval = defineModel('pollInterval', { default: 5 });
const periodicallyRefetchData = defineModel<boolean>('periodicallyRefetchData');

const showModal = ref(false);

const allVariables = useAllVariables();
const keywordsManager = useKeywordsManager();
const { data: keywords } = keywordsManager.getAll({});
const commandsManager = useCommandsManager();
const { data: commands } = commandsManager.getAll({});

const messages = useMessage();
const copyVariable = async (v: string) => {
	await copyToClipBoard(v);
	messages.success('Copied');
};

const variables = computed(() => {
	const vars = allVariables.data.value ?? [];
	const k = keywords.value?.keywords ?? [];
	const cmds = commands.value?.commands ?? [];

	return [
		...vars
			.filter(v => v.canBeUsedInRegistry)
			.map(v => {
				const name = `$(${v.isBuiltIn ? v.name : `customvar|${v.name}`})`;
				return {
					label: `${name} - ${v.description || 'Your custom variable'}`,
					value: name,
				};
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
	];
});

const { t } = useI18n();
</script>

<template>
	<base-layer
		:is-focused="isFocused"
		:layer-index="layerIndex"
		:type="type"
		@open-settings="showModal = true"
		@focus="$emit('focus', layerIndex)"
		@remove="$emit('remove', layerIndex)"
	/>

	<n-modal
		v-model:show="showModal"
		preset="card"
		title="Settings"
		style="width: 50vw"
	>
		<div class="flex flex-col gap-5">
			<n-form-item :label="t('overlaysRegistry.html.periodicallyRefetchData')">
				<n-switch v-model:value="periodicallyRefetchData" />
			</n-form-item>

			<n-form-item :label="t('overlaysRegistry.html.updateInterval')">
				<n-input-number v-model:value="pollInterval" :min="5" :max="300" />
			</n-form-item>

			<n-alert type="info" :title="t('overlaysRegistry.html.variablesAlert.title')">
				{{ t('overlaysRegistry.html.variablesAlert.selectToCopy') }}
				<n-select filterable :options="variables" @update:value="(v) => copyVariable(v)" />
			</n-alert>

			<n-tabs type="segment">
				<n-tab-pane name="HTML">
					<vue-monaco-editor
						v-model:value="html"
						theme="vs-dark"
						height="500px"
						language="html"
						:options="{
							readOnlyMessage: 'test'
						}"
					/>
				</n-tab-pane>

				<n-tab-pane name="CSS">
					<vue-monaco-editor
						v-model:value="css"
						theme="vs-dark"
						height="500px"
						language="css"
					/>
				</n-tab-pane>

				<n-tab-pane name="JS">
					<vue-monaco-editor
						v-model:value="js"
						theme="vs-dark"
						height="500px"
						language="javascript"
					/>
				</n-tab-pane>
			</n-tabs>
		</div>
	</n-modal>
</template>
