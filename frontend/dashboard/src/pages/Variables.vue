<script setup lang='ts'>
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import { type Variable, VariableType } from '@twir/grpc/generated/api/api/variables';
import {
	type DataTableColumns,
  NDataTable,
  NSpace,
  NTag,
  NAlert,
  NButton,
  NPopconfirm,
  NModal,
} from 'naive-ui';
import { computed, h, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker, useVariablesManager } from '@/api/index.js';
import Modal from '@/components/variables/modal.vue';
import { type EditableVariable } from '@/components/variables/types.js';
import { renderIcon } from '@/helpers/index.js';

const variablesManager = useVariablesManager();
const variables = variablesManager.getAll({});
const variablesDeleter = variablesManager.deleteOne;

const showModal = ref(false);

const userCanManageVariables = useUserAccessFlagChecker('MANAGE_VARIABLES');

const { t } = useI18n();

const columns = computed<DataTableColumns<Variable>>(() => [
	{
		title: t('sharedTexts.name'),
		key: 'name',
		render(row) {
			return h(NTag, { type: 'info', bordered: false }, { default: () => row.name });
		},
	},
	{
		title: t('variables.type'),
		key: 'type',
		render(row) {
			return h(NTag, { type: 'info', bordered: true }, {
				default: () => {
					switch(row.type) {
						case VariableType.SCRIPT:
							return 'Script';
						case VariableType.TEXT:
							return 'Text';
						case VariableType.NUMBER:
							return 'Number';
						default:
							return 'Unknown';
					}
				},
			});
		},
	},
	{
		title: t('sharedTexts.response'),
		key: 'response',
		render(row) {
			return h(NTag, { type: 'info', bordered: true }, {
				default: () => row.type === VariableType.SCRIPT ? 'Script' : row.response,
			});
		},
	},
	{
		title: t('sharedTexts.actions'),
		key: 'actions',
		width: 150,
		render(row) {
			const editButton = h(
				NButton,
				{
					type: 'primary',
					size: 'small',
					onClick: () => openModal(row),
					quaternary: true,
					disabled: !userCanManageVariables.value,
				}, {
					icon: renderIcon(IconPencil),
				});

			const deleteButton = h(
				NPopconfirm,
				{
					onPositiveClick: () => variablesDeleter.mutate({ id: row.id! }),
					positiveText: t('deleteConfirmation.confirm'),
					negativeText: t('deleteConfirmation.cancel'),
				},
				{
					trigger: () => h(NButton, {
						type: 'error',
						size: 'small',
						quaternary: true,
						disabled: !userCanManageVariables.value,
					}, {
						default: renderIcon(IconTrash),
					}),
					default: () => t('deleteConfirmation.text'),
				},
			);

			return h(NSpace, { }, { default: () => [editButton, deleteButton] });
		},
	},
]);

const editableVariable = ref<EditableVariable | null>(null);
function openModal(t: EditableVariable | null) {
	editableVariable.value = t;
	showModal.value = true;
}
function closeModal() {
	showModal.value = false;
}
</script>

<template>
	<n-space justify="space-between" align="center">
		<h2>{{ t('variables.title') }}</h2>
		<n-button :disabled="!userCanManageVariables" secondary type="success" @click="openModal(null)">
			{{ t('sharedButtons.create') }}
		</n-button>
	</n-space>
	<n-alert type="info">
		{{ t('variables.info') }}
	</n-alert>
	<n-data-table
		:isLoading="variables.isLoading.value"
		:columns="columns"
		:data="variables.data.value?.variables ?? []"
		style="margin-top: 20px;"
	/>

	<n-modal
		v-model:show="showModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="editableVariable?.name ?? 'New variable'"
		class="modal"
		:style="{
			width: 'auto',
			top: '50px',
		}"
		:on-close="closeModal"
	>
		<modal :variable="editableVariable" @close="closeModal" />
	</n-modal>
</template>
