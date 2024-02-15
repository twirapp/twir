<script setup lang="ts">
import { type Command } from '@twir/api/messages/commands/commands';
import { NDataTable, NButton, NSpace, NModal, NInput } from 'naive-ui';
import { ref, toRaw, computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api/index.js';
import { createColumns } from '@/components/commands/list/createColumns';
import ManageGroups from '@/components/commands/manageGroups.vue';
import Modal from '@/components/commands/modal.vue';
import type { EditableCommand, ListRowData } from '@/components/commands/types.js';


const { t } = useI18n();

const props = withDefaults(defineProps<{
	commands: Command[]
	showHeader?: boolean
	showCreateButton?: boolean,
	groupByModule?: boolean,
}>(), {
	showHeader: false,
	showCreateButton: false,
	groupByModule: false,
});

const commandsWithGroups = computed<ListRowData[]>(() => {
	const commands = props.commands;
	let i = 0;
	const groups: Record<string, ListRowData> = {
		'no-group': {
			name: 'no-group',
			children: [],
			isGroup: true,
			groupColor: '',
			index: i,
		} as any as ListRowData,
	};

	if (props.groupByModule) {
		for (const command of commands) {
			i++;
			const group = command.module;
			if (!groups[group]) {
				groups[group] = {
					name: group,
					children: [],
					isGroup: true,
					groupColor: '',
					index: i,
				} as any as ListRowData;
			}

			groups[group]!.children!.push(command as ListRowData);
		}
	} else {
		for (const command of commands) {
			i++;
			const group = command.group?.name ?? 'no-group';
			if (!groups[group]) {
				groups[group] = {
					name: group,
					children: [],
					isGroup: true,
					groupColor: command.group!.color,
					index: i,
				} as any as ListRowData;
			}

			groups[group]!.children!.push(command as ListRowData);
		}
	}

	return [
		...groups['no-group']!.children!,
		...Object.entries(groups)
			.filter(([groupName]) => groupName !== 'no-group').map(([, group]) => group),
	];
});

const commandsFilter = ref('');
const filteredCommands = computed<ListRowData[]>(() => {
	return commandsWithGroups.value.filter(c => {
		return c.name.includes(commandsFilter.value) || c.aliases.some(a => a.includes(commandsFilter.value));
	});
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

const userCanManageCommands = useUserAccessFlagChecker('MANAGE_COMMANDS');

const columns = createColumns(editCommand);
</script>

<template>
	<div>
		<div v-if="showHeader" class="header">
			<div>
				<n-space align="center">
					<h2>{{ t('commands.name', 0) }}</h2>
					<n-input
						v-model:value="commandsFilter"
						:placeholder="t('commands.searchPlaceholder')"
					/>
				</n-space>
			</div>
			<div>
				<n-space>
					<n-button
						:disabled="!userCanManageCommands" secondary type="info"
						@click="showManageGroupsModal = true"
					>
						{{ t('commands.groups.manageButton') }}
					</n-button>

					<n-button
						v-if="showCreateButton"
						secondary
						type="success"
						:disabled="!userCanManageCommands"
						@click="() => {
							editableCommand = null;
							showCommandEditModal = true;
						}"
					>
						{{ t('sharedButtons.create') }}
					</n-button>
				</n-space>
			</div>
		</div>

		<n-data-table
			size="small"
			:columns="columns"
			:data="filteredCommands"
			:bordered="false"
			:row-key="r => r.index"
		/>

		<n-modal
			v-model:show="showCommandEditModal"
			:mask-closable="false"
			:segmented="true"
			preset="card"
			:title="editableCommand?.name ?? t('commands.newCommandTitle')"
			class="modal"
			:style="{
				width: '800px',
				top: '50px',
			}"
			:on-close="onModalClose"
		>
			<modal
				:command="editableCommand"
				@close="() => {
					showCommandEditModal = false;
					onModalClose()
				}"
			/>
		</n-modal>

		<n-modal
			v-model:show="showManageGroupsModal"
			:mask-closable="false"
			:segmented="true"
			preset="card"
			:title="t('commands.groups.manageButton')"
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
