<script setup lang='ts'>
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import type { Greeting } from '@twir/grpc/generated/api/api/greetings';
import { type Variable, VariableType } from '@twir/grpc/generated/api/api/variables';
import {
	type DataTableColumns,
	NDataTable,
	NSpace,
	NTag,
	NAlert,
	NButton,
	NPopconfirm,
	NModal, NSwitch,
} from 'naive-ui';
import { h, ref } from 'vue';

import { useGreetingsManager } from '@/api/index.js';
import Modal from '@/components/greetings/modal.vue';
import { EditableGreeting } from '@/components/greetings/types.js';
import { renderIcon } from '@/helpers/index.js';

const greetingsManager = useGreetingsManager();
const greetings = greetingsManager.getAll({});
const greetingsDeleter = greetingsManager.deleteOne;
const greetingsPatcher = greetingsManager.patch!;

const showModal = ref(false);

const columns: DataTableColumns<Greeting> = [
	{
		title: 'User name',
		key: 'userName',
		render(row) {
			return h(NTag, { type: 'info', bordered: false }, { default: () => row.userId });
		},
	},
	{
		title: 'Text',
		key: 'text',
		render(row) {
			return h(NTag, { type: 'info', bordered: true }, { default: () => row.text });
		},
	},
	{
		title: 'Status',
		key: 'enabled',
		width: 100,
		render(row) {
			return h(
				NSwitch,
				{
					value: row.enabled,
					onUpdateValue: (value: boolean) => {
						greetingsPatcher.mutateAsync({ id: row.id, enabled: value }).then(() => {
							row.enabled = value;
						});
					},
				},
				{ default: () => row.enabled },
			);
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
					onPositiveClick: () => greetingsDeleter.mutate({ id: row.id! }),
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

const editableGreeting = ref<EditableGreeting | null>(null);
function openModal(t: EditableGreeting | null) {
	editableGreeting.value = t;
	showModal.value = true;
}
function closeModal() {
	showModal.value = false;
}
</script>

<template>
  <n-space justify="space-between" align="center">
    <h2>Greetings</h2>
    <n-button secondary type="success" @click="openModal(null)">
      Create
    </n-button>
  </n-space>
  <n-alert>
    <p>Greeting system used for welcoming new users typed their first message on stream.</p>
    <p>
      If you wanna greet every user in chat, not only specified - then you can use events system.
    </p>
  </n-alert>
  <n-data-table
    :isLoading="greetings.isLoading.value"
    :columns="columns"
    :data="greetings.data.value?.variables ?? []"
    style="margin-top: 20px;"
  />

  <n-modal
    v-model:show="showModal"
    :mask-closable="false"
    :segmented="true"
    preset="card"
    :title="editableGreeting?.id ?? 'New greeting'"
    class="modal"
    :style="{
      width: 'auto',
      top: '50px',
    }"
    :on-close="closeModal"
  >
    <modal :greeting="editableVariable" @close="closeModal" />
  </n-modal>
</template>
