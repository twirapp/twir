<script setup lang='ts'>
import { IconTrash, IconPencil } from '@tabler/icons-vue';
import { type Command } from '@twir/grpc/generated/api/api/commands';
import { NDataTable, DataTableColumns, NText, NSwitch, NButton, NSpace, NModal, NTag, NPopconfirm } from 'naive-ui';
import { h, ref, toRaw, VNode, computed } from 'vue';

import { useCommandsManager } from '@/api/index.js';
import Modal from '@/components/commands/modal.vue';
import { renderIcon } from '@/helpers/index.js';

const props = withDefaults(defineProps<{
	commands: Command[]
	showHeader: boolean
}>(), {
	showHeader: false,
});

const commandsManager = useCommandsManager();
const commandsDeleter = commandsManager.deleteOne;

type RowData = Command & { isGroup?: boolean, groupColor?: string, children?: RowData[] };

const commandsWithGroups = computed<RowData[]>(() => {
	const commands = props.commands;
	const groups: Record<string, RowData> = {
		'no-group': {
			name: 'no-group',
			children: [],
			isGroup: true,
			groupColor: '',
		} as any as RowData,
	};

	for (const command of commands) {
		const group = command.group?.name ?? 'no-group';
		if (!groups[group]) {
			groups[group] = {
				name: group,
				children: [],
				isGroup: true,
				groupColor: command.group!.color,
			} as any as RowData;
		}

		groups[group]!.children!.push(command as RowData);
	}

	return [
		...groups['no-group']!.children!,
		...Object.entries(groups)
			.filter(([groupName]) => groupName !== 'no-group').map(([, group]) => group),
	];
});

const columns: DataTableColumns<RowData> = [
	{
		title: 'Name',
		key: 'name',
		width: 150,
		render(row) {
			return h(
				NTag,
				{
					bordered: false,
					color: { color: row.isGroup ? row.groupColor : 'rgba(112, 192, 232, 0.16)' },
				},
				{ default: () => row.name },
			);
		},
	},
	{
		title: 'Responses',
		key: 'responses',
		render(row) {
			if (row.isGroup) return;
			return h(
				NText,
				{},
				{ default: () => row.responses?.map(r => `${r.text}`).join('\n') },
			);
		},
	},
	{
		title: 'Status',
		key: 'enabled',
		width: 100,
		render(row) {
			if (row.isGroup) return;

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
			if (row.isGroup) return;
			const editButton = h(
				NButton,
				{
					type: 'primary',
					size: 'small',
					onClick: () => editCommand(row),
				}, {
				icon: renderIcon(IconPencil),
			});

			const deleteButton = h(
				NPopconfirm,
				{
					onPositiveClick: () => commandsDeleter.mutate({ commandId: row.id }),
				},
				{
					trigger: h(NButton, {
						type: 'error',
						size: 'small',
					}, {
						default: renderIcon(IconTrash),
					}),
					default: () => 'Are you sure you want to delete this command?',
				},
			);

			const buttons: VNode[] = [editButton];

			if (!row.default) {
				buttons.push(deleteButton);
			}

			return h(NSpace, {  }, { default: () => buttons });
		},
	},
];

const showModal = ref(false);

const editableCommand = ref<Command | null>(null);
function editCommand(command: Command) {
	editableCommand.value = structuredClone(toRaw(command));
	showModal.value = true;
}

function onModalClose() {
	editableCommand.value = null;
}
</script>

<template>
  <div>
    <div v-if="showHeader" class="header">
      <div>
        <h2>Commands</h2>
      </div>
      <div>
        <n-button type="primary">
          Create
        </n-button>
      </div>
    </div>

    <n-data-table
      :columns="columns"
      :data="commandsWithGroups"
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
        top: '50px',
      }"
      :on-close="onModalClose"
    >
      <modal :command="editableCommand" />
    </n-modal>
  </div>
</template>

<style scoped>
.header {
	display: flex;
	justify-content: space-between;
	align-items: center;
}
</style>
