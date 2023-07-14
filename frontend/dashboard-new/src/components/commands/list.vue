<script setup lang='ts'>
import { IconTrash, IconPencil } from '@tabler/icons-vue';
import { type Command } from '@twir/grpc/generated/api/api/commands';
import { NDataTable, DataTableColumns, NText, NSwitch, NButton, NSpace, NModal, NTag, NPopconfirm } from 'naive-ui';
import { h, ref, toRaw, VNode, computed } from 'vue';

import { useCommandsManager } from '@/api/index.js';
import { createListColumns } from '@/components/commands/createListColumns.js';
import Modal from '@/components/commands/modal.vue';
import type { EditableCommand, ListRowData } from '@/components/commands/types.js';
import { renderIcon } from '@/helpers/index.js';

const props = withDefaults(defineProps<{
	commands: Command[]
	showHeader: boolean
}>(), {
	showHeader: false,
});

const commandsManager = useCommandsManager();
const commandsDeleter = commandsManager.deleteOne;

const commandsWithGroups = computed<ListRowData[]>(() => {
	const commands = props.commands;
	const groups: Record<string, ListRowData> = {
		'no-group': {
			name: 'no-group',
			children: [],
			isGroup: true,
			groupColor: '',
		} as any as ListRowData,
	};

	for (const command of commands) {
		const group = command.group?.name ?? 'no-group';
		if (!groups[group]) {
			groups[group] = {
				name: group,
				children: [],
				isGroup: true,
				groupColor: command.group!.color,
			} as any as ListRowData;
		}

		groups[group]!.children!.push(command as ListRowData);
	}

	return [
		...groups['no-group']!.children!,
		...Object.entries(groups)
			.filter(([groupName]) => groupName !== 'no-group').map(([, group]) => group),
	];
});

const showModal = ref(false);

const editableCommand = ref<EditableCommand | null>(null);
function editCommand(command: EditableCommand) {
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
      :columns="createListColumns(editCommand, commandsDeleter)"
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
