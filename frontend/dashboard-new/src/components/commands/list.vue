<script setup lang='ts'>
import { type Command } from '@twir/grpc/generated/api/api/commands';
import { NDataTable, NButton, NSpace, NModal } from 'naive-ui';
import { ref, toRaw, computed } from 'vue';

import { useCommandsManager } from '@/api/index.js';
import { createListColumns } from '@/components/commands/createListColumns.js';
import ManageGroups from '@/components/commands/manageGroups.vue';
import Modal from '@/components/commands/modal.vue';
import type { EditableCommand, ListRowData } from '@/components/commands/types.js';

const props = withDefaults(defineProps<{
	commands: Command[]
	showHeader: boolean
	showCreateButton: boolean,
}>(), {
	showHeader: false,
	showCreateButton: false,
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

const showCommandEditModal = ref(false);
const showManageGroupsModal = ref(false);

const editableCommand = ref<EditableCommand | null>(null);
function editCommand(command: EditableCommand) {
	editableCommand.value = structuredClone(toRaw(command));
	showCommandEditModal.value = true;
}

function onModalClose() {
	editableCommand.value = null;
}
</script>

<template>
  <div>
    <div class="header">
      <div v-if="showHeader">
        <h2>Commands</h2>
      </div>
      <div>
        <n-space>
          <n-button secondary type="info" @click="showManageGroupsModal = true">
            Manage groups
          </n-button>

          <n-button v-if="showCreateButton" secondary type="success">
            Create command
          </n-button>
        </n-space>
      </div>
    </div>

    <n-data-table
      :columns="createListColumns(editCommand, commandsDeleter)"
      :data="commandsWithGroups"
      :bordered="false"
    />

    <n-modal
      v-model:show="showCommandEditModal"
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

    <n-modal
      v-model:show="showManageGroupsModal"
      :mask-closable="false"
      :segmented="true"
      preset="card"
      title="Commands groups"
      class="modal"
      :style="{
        width: '600px',
      }"
      :on-close="onModalClose"
    >
      <manage-groups />
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
