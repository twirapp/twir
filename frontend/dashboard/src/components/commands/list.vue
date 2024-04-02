<script setup lang="ts">
import { IconChevronRight, IconChevronDown, IconSearch } from '@tabler/icons-vue';
import type { ColumnDef } from '@tanstack/vue-table';
import {
    FlexRender,
    getCoreRowModel,
    useVueTable,
		getExpandedRowModel,
} from '@tanstack/vue-table';
import { type Command } from '@twir/api/messages/commands/commands';
import { rgbToHex, hexToRgb, colorBrightness, type Rgb } from '@zero-dependency/utils';
import { NButton, NSpace, NModal, NInput, useThemeVars, NIcon } from 'naive-ui';
import { ref, h, computed } from 'vue';
import { useI18n } from 'vue-i18n';

import ColumnActions from './list/column-actions.vue';
import { type Group, isCommand, createGroups } from './list/create-groups';

import { useUserAccessFlagChecker } from '@/api/index.js';
import ManageGroups from '@/components/commands/manageGroups.vue';
import Modal from '@/components/commands/modal.vue';
import type { EditableCommand } from '@/components/commands/types.js';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';

const { t } = useI18n();
const themeVars = useThemeVars();

const props = withDefaults(defineProps<{
	commands: Command[]
	showHeader?: boolean
	showCreateButton?: boolean,
	enableGroups?: boolean,
	showBackground?: boolean
}>(), {
	showHeader: false,
	showCreateButton: false,
	enableGroups: false,
});

const commandsFilter = ref('');

const showCommandEditModal = ref(false);
const showManageGroupsModal = ref(false);

const editableCommand = ref<EditableCommand | null>(null);

function onModalClose() {
	editableCommand.value = null;
}

const userCanManageCommands = useUserAccessFlagChecker('MANAGE_COMMANDS');

const columns: ColumnDef<Command | Group>[] = [
  {
    accessorKey: 'name',
		size: 10,
    header: () => h('div', {}, 'Name'),
    cell: ({ row }) => {
			const chevron = row.getCanExpand() ? h(row.getIsExpanded() ? IconChevronDown : IconChevronRight) : null;

			if (isCommand(row.original)) {
				return h(
					'div',
					{ class: 'flex gap-2 items-center select-none' },
					[chevron, '!' + row.getValue('name') as string],
				);
			}

			let rgbColor: Rgb | null = null;
			if (row.original.color) {
				rgbColor = hexToRgb(rgbToHex(row.original.color));
			}

			const color = rgbColor
				? (colorBrightness(rgbColor) >= 128 ? '#000' : '#fff')
				: 'var(--n-text-color)';

      return h(
				'div',
				{ class: 'flex gap-2 items-center select-none' },
				[
					chevron,
					h(
						'span',
						{
							class: 'p-1 rounded',
							style: `background-color: ${row.original.color}; color: ${color}`,
						},
						row.original.name.charAt(0).toLocaleUpperCase() + row.original.name.slice(1),
					),
				],
			);
    },
  },
	{
    accessorKey: 'responses',
    header: () => h('div', {  }, 'Responses'),
		size: 85,
    cell: ({ row }) => {
			if (!isCommand(row.original)) {
				return;
			}

			const responses: Command['responses'] = row.getValue('responses');
			if (!responses.length) {
				return row.original.description;
			}

			const mappedResponses = responses.map((r) => h('span', {}, r.text));
      return h('div', { class: 'flex flex-col' }, mappedResponses);
    },
  },
	{
		id: 'actions',
		size: 5,
    cell: ({ row }) => {
			if (!isCommand(row.original)) {
				return;
			}

			return h(
				ColumnActions,
				{
					row: row.original,
					onEdit: () => {
						// typescript hi-hi + ha-ha
						editableCommand.value = row.original as Command;
						showCommandEditModal.value = true;
					},
				},
			);
    },
  },
];

const filteredCommands = computed(() => {
	return props.commands
		.filter(c => c.name.includes(commandsFilter.value) || c.aliases.some(a => a.includes(commandsFilter.value)));
});

const tableValue = computed(() => props.enableGroups ? createGroups(filteredCommands.value) : filteredCommands.value);

const table = useVueTable({
	get data() {
		return tableValue.value;
	},
	get columns() {
		return columns;
	},
	getCoreRowModel: getCoreRowModel(),
	getExpandedRowModel: getExpandedRowModel(),
	getSubRows: (original) => {
		if ('commands' in original) {
			return original.commands;
		}
	},
});
</script>

<template>
	<div>
		<div v-if="showHeader" class="flex justify-between items-center flex-wrap gap-2">
			<div>
				<n-input
					v-model:value="commandsFilter"
					:placeholder="t('commands.searchPlaceholder')"
				>
					<template #prefix>
						<n-icon :component="IconSearch"></n-icon>
					</template>
				</n-input>
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

		<n-modal
			v-model:show="showCommandEditModal"
			:mask-closable="false"
			:segmented="true"
			preset="card"
			:title="editableCommand?.name ?? t('commands.newCommandTitle')"
			class="modal"
			:style="{
				width: '800px',
				height: '90dvh',
			}"
			:on-close="onModalClose"
			content-style="padding: 5px;"
		>
			<modal
				:command="editableCommand"
				@close="() => {
					showCommandEditModal = false;
					onModalClose()
				}"
			/>
		</n-modal>

		<div
			class="border rounded-md"
			:class="{ 'mt-5': showHeader }"
			:style="{
				backgroundColor: props.showBackground ? themeVars.cardColor : 'inherit',
				color: themeVars.textColor2
			}"
		>
			<Table>
				<TableHeader>
					<TableRow
						v-for="headerGroup in table.getHeaderGroups()"
						:key="headerGroup.id"
						class="border-b"
					>
						<TableHead
							v-for="header in headerGroup.headers"
							:key="header.id"
							:style="{ width: `${header.getSize()}%` }"
						>
							<FlexRender
								v-if="!header.isPlaceholder"
								:render="header.column.columnDef.header"
								:props="header.getContext()"
							/>
						</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					<template v-if="table.getRowModel().rows?.length">
						<TableRow
							v-for="row in table.getRowModel().rows" :key="row.id"
							:data-state="row.getIsSelected() ? 'selected' : undefined"
							class="border-b"
							:class="{ 'cursor-pointer': !isCommand(row.original) }"
						>
							<TableCell
								v-for="cell in row.getVisibleCells()"
								:key="cell.id"
								@click="() => {
									if (row.getCanExpand()) {
										row.getToggleExpandedHandler()()
									}
								}"
							>
								<FlexRender
									:render="cell.column.columnDef.cell"
									:props="cell.getContext()"
								/>
							</TableCell>
						</TableRow>
					</template>
					<template v-else>
						<TableRow>
							<TableCell :colSpan="columns.length" class="h-24 text-center">
								No commands
							</TableCell>
						</TableRow>
					</template>
				</TableBody>
			</Table>
		</div>

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
