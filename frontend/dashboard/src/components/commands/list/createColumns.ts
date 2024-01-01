import { DataTableColumns } from 'naive-ui';
import { computed, h } from 'vue';
import { useI18n } from 'vue-i18n';

import ColumnActions from './column-actions.vue';
import ColumnName from './column-name.vue';
import ColumnResponses from './column-responses.vue';
import ColumnStatus from './column-status.vue';

import type { EditableCommand, ListRowData } from '@/components/commands/types';

export const createColumns = (editCommand: (command: EditableCommand) => void) => {
	const { t } = useI18n();

	return computed<DataTableColumns<ListRowData>>(() => [
		{
			title: t('sharedTexts.name'),
			key: 'name',
			width: 250,
			render(row) {
				return h(ColumnName, { row });
			},
		},
		{
			title: t('sharedTexts.responses'),
			key: 'responses',
			render(row) {
				return h(ColumnResponses, { row });
			},
		},
		{
			title: t('sharedTexts.status'),
			key: 'enabled',
			width: 100,
			render(row) {
				return h(ColumnStatus, { row });
			},
		},
		{
			title: t('sharedTexts.actions'),
			key: 'actions',
			width: 150,
			render(row) {
				return h(ColumnActions, { row, onEdit: editCommand });
			},
		},
	]);
};
