<script setup lang='ts'>
import { IconTrash, IconPencil } from '@tabler/icons-vue';
import { type Command } from '@twir/grpc/generated/api/api/commands';
import { NDataTable, DataTableColumns, NText, NSwitch, NButton, NSpace, NBadge } from 'naive-ui';
import { h, type Ref } from 'vue';

import { renderIcon } from '@/helpers/index.js';

defineProps<{
	commands: Ref<Command[]>
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
			const editButton = h(NButton, { type: 'primary', size: 'small' }, {
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
</script>

<template>
  <n-data-table
    :columns="columns"
    :data="commands"
    :bordered="false"
  />
</template>

<style scoped lang='postcss'>

</style>
