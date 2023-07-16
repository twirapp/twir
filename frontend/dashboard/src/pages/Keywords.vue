<script setup lang='ts'>
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import type { Keyword } from '@twir/grpc/generated/api/api/keywords';
import {
	type DataTableColumns,
	NDataTable,
	NSpace,
	NTag,
	NSwitch,
	NButton,
	NPopconfirm,
	NModal,
} from 'naive-ui';
import { h, ref } from 'vue';

import { useKeywordsManager } from '@/api/index.js';
import Modal from '@/components/keywords/modal.vue';
import { type EditableKeyword } from '@/components/keywords/types.js';
import { renderIcon } from '@/helpers/index.js';

const keywordsManager = useKeywordsManager();
const keywords = keywordsManager.getAll({});
const keywordsDeleter = keywordsManager.deleteOne;
const keywordsPatcher = keywordsManager.patch;

const showModal = ref(false);

const columns: DataTableColumns<Keyword> = [
	{
		title: 'Trigger',
		key: 'text',
		render(row) {
			return h(NTag, { type: 'info', bordered: false }, { default: () => row.text });
		},
	},
	{
		title: 'Response',
		key: 'response',
		render(row) {
			return h(NTag, { type: 'info', bordered: true }, { default: () => row.response || 'No response' });
		},
	},
	{
		title: 'Usages',
		key: 'usages',
		render(row) {
			return h(NTag, { type: 'info', bordered: true }, { default: () => row.usages });
		},
	},
	{
		title: 'Status',
		key: 'enabled',
		render(row) {
			return h(NSwitch, {
				value: row.enabled,
				onUpdateValue: (value) => {
					keywordsPatcher!.mutateAsync({
						id: row.id,
						enabled: value,
					}).then(() => row.enabled = value);
				},
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
					onPositiveClick: () => keywordsDeleter.mutate({ id: row.id }),
				},
				{
					trigger: () => h(NButton, {
						type: 'error',
						size: 'small',
						quaternary: true,
					}, {
						default: renderIcon(IconTrash),
					}),
					default: () => 'Are you sure you want to delete this keyword?',
				},
			);

			return h(NSpace, { }, { default: () => [editButton, deleteButton] });
		},
	},
];

const editableKeyword = ref<EditableKeyword | null>(null);
function openModal(t: EditableKeyword | null) {
	editableKeyword.value = t;
	showModal.value = true;
}
function closeModal() {
	showModal.value = false;
}
</script>

<template>
  <n-space justify="space-between" align="center">
    <h2>Keywords</h2>
    <n-button secondary type="success" @click="openModal(null)">
      Create
    </n-button>
  </n-space>
  <n-data-table
    :isLoading="keywords.isLoading.value"
    :columns="columns"
    :data="keywords.data.value?.keywords ?? []"
  />

  <n-modal
    v-model:show="showModal"
    :mask-closable="false"
    :segmented="true"
    preset="card"
    :title="editableKeyword?.id ? 'Edit keyword' : 'New keyword'"
    class="modal"
    :style="{
      width: '600px',
      top: '50px',
    }"
    :on-close="closeModal"
  >
    <modal :keyword="editableKeyword" @close="closeModal" />
  </n-modal>
</template>
