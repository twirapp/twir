<script setup lang='ts'>
import { type Command } from '@twir/grpc/generated/api/api/commands';
import { NDataTable, DataTableColumns, NText } from 'naive-ui';
import { h } from 'vue';

defineProps<{
	commands: Command[]
}>();

const columns: DataTableColumns<Command> = [
	{
		title: 'Name',
		key: 'name',
	},
	{
		title: 'Responses',
		key: 'responses',
		render(row) {
			return h(
				NText,
				{

				},
				{ default: () => row.responses.map(r => `${r.text}`).join('\n') },
			);
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
