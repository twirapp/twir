<script setup lang="ts">
import { IconPencil, IconTrash } from '@tabler/icons-vue';
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
import { storeToRefs } from 'pinia';
import { computed, h, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api/index.js';
import { useVariablesApi } from '@/api/variables';
import { type CustomVariable, type EditableCustomVariable } from '@/api/variables';
import Modal from '@/components/variables/modal.vue';
import { VariableType } from '@/gql/graphql';
import { renderIcon } from '@/helpers/index.js';

const variablesApi = useVariablesApi();
const { customVariables, isLoading } = storeToRefs(variablesApi);
const deleter = variablesApi.useMutationRemoveVariable();

const showModal = ref(false);

const userCanManageVariables = useUserAccessFlagChecker('MANAGE_VARIABLES');

const { t } = useI18n();

const columns = computed<DataTableColumns<CustomVariable>>(() => [
	{
		title: t('sharedTexts.name'),
		key: 'name',
		render(row) {
			return row.name;
		},
	},
	{
		title: t('variables.type'),
		key: 'type',
		render(row) {
			return h(NTag, { type: 'info', bordered: true }, {
				default: () => row.type,
			});
		},
	},
	{
		title: t('sharedTexts.response'),
		key: 'response',
		render(row) {
			return h(NTag, { type: 'info', bordered: true }, {
				default: () => row.type === VariableType.Script ? 'Script' : row.response,
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
					onPositiveClick: () => deleter.executeMutation({ id: row.id! }),
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

			return h(NSpace, {}, { default: () => [editButton, deleteButton] });
		},
	},
]);

const editableVariable = ref<EditableCustomVariable | null>(null);

function openModal(t: EditableCustomVariable | null) {
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
		:isLoading="isLoading"
		:columns="columns"
		:data="customVariables"
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
