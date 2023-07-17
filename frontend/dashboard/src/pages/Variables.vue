<script setup lang='ts'>
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import { type Variable, VariableType } from '@twir/grpc/generated/api/api/variables';
import {
	type DataTableColumns, NTag, NButton, NPopconfirm, NSpace,
} from 'naive-ui';
import { h, ref } from 'vue';

import { useVariablesManager } from '@/api/index.js';
import Modal from '@/components/variables/modal.vue';
import { type EditableVariable } from '@/components/variables/types.js';
import { renderIcon } from '@/helpers/index.js';

const variablesManager = useVariablesManager();
const variables = variablesManager.getAll({});
const variablesDeleter = variablesManager.deleteOne;

const showModal = ref(false);

const columns: DataTableColumns<Variable> = [
	{
		title: 'Name',
		key: 'name',
		render(row) {
			return h(NTag, { type: 'info', bordered: false }, { default: () => row.name });
		},
	},
	{
		title: 'Type',
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
		title: 'Response',
		key: 'response',
		render(row) {
			return h(NTag, { type: 'info', bordered: true }, {
				default: () => row.type === VariableType.SCRIPT ? 'Script' : row.response,
			});
		},
	},
	{
		title: 'Actions',
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
				}, {
					icon: renderIcon(IconPencil),
				});

			const deleteButton = h(
				NPopconfirm,
				{
					onPositiveClick: () => variablesDeleter.mutate({ id: row.id! }),
				},
				{
					trigger: () => h(NButton, {
						type: 'error',
						size: 'small',
						quaternary: true,
					}, {
						default: renderIcon(IconTrash),
					}),
					default: () => 'Are you sure you want to delete this variable?',
				},
			);

			return h(NSpace, { }, { default: () => [editButton, deleteButton] });
		},
	},
];

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
    <h2>Variables</h2>
    <n-button secondary type="success" @click="openModal(null)">
      Create
    </n-button>
  </n-space>
  <n-alert>
    When you create variable, then you can use them in bot responses.
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
