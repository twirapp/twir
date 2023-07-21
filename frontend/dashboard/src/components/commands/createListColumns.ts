import { IconTrash, IconPencil } from '@tabler/icons-vue';
import { DataTableColumns, NButton, NPopconfirm, NSpace, NSwitch, NTag, NText } from 'naive-ui';
import { computed, h } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCommandsManager, useUserAccessFlagChecker } from '@/api/index.js';
import type { ListRowData, EditableCommand } from '@/components/commands/types.js';
import { renderIcon } from '@/helpers/index.js';

type Deleter = ReturnType<typeof useCommandsManager>['deleteOne']

const rgbaPattern = /rgba?\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*(?:,\s*([\d.]+)\s*)?\)/;
const computeGroupTextColor = (color?: string) => {
	const result = rgbaPattern.exec(color ?? '');
	if (!result) return '#c2b7b7';
	const [r, g, b] = result.slice(1).map(i => parseInt(i, 10));

	const bright = (
		(((r * 299) + (g * 587) + (b * 114)) / 1000) - 128
	) * -1000;

	return `rgba(${bright},${bright},${bright})`;
};

export const createListColumns = (
	editCommand: (command: EditableCommand) => void,
	deleter: Deleter,
	patcher: (id: string, value: boolean) => any,
) => {
	const userCanManageCommands = useUserAccessFlagChecker('MANAGE_COMMANDS');
	const { t } = useI18n();

	return computed<DataTableColumns<ListRowData>>(() => [
		{
			title: t('tables.columns.name'),
			key: 'name',
			width: 250,
			render(row) {
				return h(
					NTag,
					{
						bordered: false,
						color: { color: row.isGroup ? row.groupColor : 'rgba(112, 192, 232, 0.16)' },
					},
					{ default: () => h(
						'p',
						{
							style: `color: ${computeGroupTextColor(row.groupColor)}`,
						},
						row.name),
					},
				);
			},
		},
		{
			title: t('commands.table.responses.name'),
			key: 'responses',
			render(row) {
				if (row.isGroup) return;
				return h(
					NText,
					{},
					{
						default: () => {
							if (row.module !== 'CUSTOM') return row.description ?? 'No description';
							return row.responses.length
								? h(NSpace, { vertical: true }, {
									default: () => row.responses?.map(r => h('span', null, `${r.text}`)),
								})
								: t('commands.table.responses.empty');
						},
					},
				);
			},
		},
		{
			title: t('tables.columns.status'),
			key: 'enabled',
			width: 100,
			render(row) {
				if (row.isGroup) return;

				return h(
					NSwitch,
					{
						value: row.enabled,
						disabled: !userCanManageCommands.value,
						onUpdateValue: (value: boolean) => {
							row.enabled = value;
							patcher(row.id, value);
						},
					},
					{ default: () => row.enabled },
				);
			},
		},
		{
			title: t('tables.columns.actions'),
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
						quaternary: true,
						disabled: !userCanManageCommands.value,
					}, {
						icon: renderIcon(IconPencil),
					});

				const deleteButton = h(
					NPopconfirm,
					{
						onPositiveClick: () => deleter.mutate({ commandId: row.id }),
						positiveText: t('deleteConfirmation.confirm'),
						negativeText: t('deleteConfirmation.cancel'),
					},
					{
						trigger: () => h(NButton, {
							type: 'error',
							size: 'small',
							quaternary: true,
							disabled: row.default || !userCanManageCommands.value,
						}, {
							default: renderIcon(IconTrash),
						}),
						default: () => t('deleteConfirmation.text'),
					},
				);

				return h(NSpace, { }, { default: () => [editButton, deleteButton] });
			},
		},
	]);
};
