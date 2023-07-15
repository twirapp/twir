<script setup lang='ts'>
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import { type Timer } from '@twir/grpc/generated/api/api/timers';
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

import { useTimersManager } from '@/api/index.js';
import Modal from '@/components/timers/modal.vue';
import { type EditableTimer } from '@/components/timers/types.js';
import { renderIcon } from '@/helpers/index.js';

const timersManager = useTimersManager();
const timers = timersManager.getAll({});

const columns: DataTableColumns<Timer> = [
	{
		title: 'Name',
		key: 'name',
		render(row) {
			return h(NTag, { type: 'info', bordered: false }, { default: () => row.name });
		},
	},
	{
		title: 'Responses',
		key: 'responses',
		render(row) {
			return h(NSpace, { vertical: true }, {
				default: () => row.responses.map((response) => {
					return h('span', null, response.text);
				}),
			});
		},
	},
	{
		title: 'Interval in minutes',
		key: 'timeInterval',
		render(row) {
			return h(NTag, { type: 'info' }, { default: () => `${row.timeInterval} m.` });
		},
	},
	{
		title: 'Interval in messages',
		key: 'messageInterval',
		render(row) {
			return h(NTag, { type: 'info' }, { default: () => `${row.messageInterval}` });
		},
	},
	{
		title: 'Status',
		key: 'enabled',
		render(row) {
			return h(NSwitch, {
				modelValue: row.enabled,
				onUpdateValue: (value) => {
					console.log(value);
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
					onPositiveClick: () => console.log('delete'),
				},
				{
					trigger: () => h(NButton, {
						type: 'error',
						size: 'small',
						quaternary: true,
					}, {
						default: renderIcon(IconTrash),
					}),
					default: () => 'Are you sure you want to delete this timer?',
				},
			);

			return h(NSpace, { }, { default: () => [editButton, deleteButton] });
		},
	},
];

const showModal = ref(false);

const editableTimer = ref<EditableTimer | null>(null);
function openModal(t: EditableTimer | null) {
	editableTimer.value = t;
	showModal.value = true;
}
</script>

<template>
  <n-space justify="space-between" align="center">
    <h2>Timers</h2>
    <n-button secondary type="success" @click="openModal(null)">
      Create
    </n-button>
  </n-space>
  <n-data-table
    :isLoading="timers.isLoading.value"
    :columns="columns"
    :data="[{
      name: 'test',
      responses: [{text: 'qwe'}, {text: 'asd'}],
      timeInterval: 5,
      messageInterval: 5,
    }]"
  />

  <n-modal
    v-model:show="showModal"
    :mask-closable="false"
    :segmented="true"
    preset="card"
    :title="editableTimer?.name ?? 'New timer'"
    class="modal"
    :style="{
      width: '600px',
      top: '50px',
    }"
    :on-close="() => showModal = false"
  >
    <modal :timer="editableTimer" />
  </n-modal>
</template>
