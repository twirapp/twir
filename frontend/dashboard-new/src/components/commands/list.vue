<script setup lang='ts'>
import { IconTrash, IconPencil } from '@tabler/icons-vue';
import { type Command } from '@twir/grpc/generated/api/api/commands';
import { NDataTable, DataTableColumns, NText, NSwitch, NButton, NSpace, NBadge, NModal } from 'naive-ui';
import { h, ref } from 'vue';

import Modal from '@/components/commands/modal.vue';
import { renderIcon } from '@/helpers/index.js';

defineProps<{
	commands: Command[]
}>();

const columns: DataTableColumns<Command> = [
	{
		title: 'Name',
		key: 'name',
		width: 150,
		render(row) {
			return h(
				NBadge,
				{ type: 'info', value: row.name },
				{},
			);
		},
	},
	{
		title: 'Responses',
		key: 'responses',
		render(row) {
			return h(
				NText,
				{},
				{ default: () => row.responses.map(r => `${r.text}`).join('\n') },
			);
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
						row.enabled = value;
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
					onClick: () => editCommand(row),
				}, {
				icon: renderIcon(IconPencil),
			});
			const deleteButton = h(NButton, { type: 'error', size: 'small' }, {
				icon: renderIcon(IconTrash),
			});
			return h(NSpace, {  }, {
				default: () => [
					editButton,
					deleteButton,
				],
			});
		},
	},
];

const showModal = ref(false);

const editableCommand = ref<Command | null>(null);
function editCommand(command: Command) {
	editableCommand.value = command;
	showModal.value = true;
}

function onModalClose() {
	editableCommand.value = null;
}
</script>

<template>
  <div>
    <n-data-table
      :columns="columns"
      :data="commands"
      :bordered="false"
    />

    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :segmented="true"
      preset="card"
      :title="editableCommand?.name ?? 'New command'"
      class="modal"
      :style="{
        width: '800px',
        position: 'fixed',
        left: `calc(50% - 800px/2)`,
        top: '50px',
      }"
      :on-close="onModalClose"
    >
      <modal />
    </n-modal>
  </div>
</template>

<style scoped lang='postcss'>

</style>
