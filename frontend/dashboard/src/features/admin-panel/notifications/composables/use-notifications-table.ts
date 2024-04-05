import {
	type ColumnDef,
	getCoreRowModel,
	useVueTable,
	getPaginationRowModel,
} from '@tanstack/vue-table';
import type { Notification } from '@twir/api/messages/admin_notifications/admin_notifications';
import { defineStore } from 'pinia';
import { computed, h, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import CreatedAtTooltip from '../components/created-at-tooltip.vue';
import NotificationsTableActions from '../components/notifications-table-actions.vue';
import { useNotificationsForm } from '../composables/use-notifications-form.js';

import { useAdminNotifications } from '@/api/notifications.js';
import { useLayout } from '@/composables/use-layout.js';

export const useNotificationsTable = defineStore('admin-panel/notifications-table', () => {
	const layout = useLayout();
	const { t } = useI18n();

	const notificationsForm = useNotificationsForm();
	const notificationsCrud = useAdminNotifications();
	const notificationsData = notificationsCrud.getAll({});

	const tableColumns = computed<ColumnDef<Notification>[]>(() => [
		{
			accessorKey: 'message',
			size: 80,
			header: () => h('div', {}, t('adminPanel.notifications.message')),
			cell: ({ row }) => {
				return h('span', { innerHTML: row.original.message });
			},
		},
		{
			accessorKey: 'createdAt',
			size: 5,
			header: () => h('div', {}, t('adminPanel.notifications.createdAt')),
			cell: ({ row }) => {
				return h(CreatedAtTooltip, { time: new Date(row.original.createdAt) });
			},
		},
		{
			accessorKey: 'actions',
			size: 10,
			header: () => '',
			cell: ({ row }) => {
				return h(NotificationsTableActions, {
					onDelete: () => onDeleteNotification(row.original.id),
					onEdit: () => onEditNotification(row.original),
				});
			},
		},
	]);

	const notifications = computed<Notification[]>(() => {
		return notificationsData.data.value?.notifications ?? [];
	});

	const table = useVueTable({
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		get data() {
			return notifications.value;
		},
		get columns() {
			return tableColumns.value;
		},
	});

	function onDeleteNotification(notificationId: string) {
		notificationsCrud.deleteOne.mutate({ id: notificationId });
	}

	async function onEditNotification(notification: Notification) {
		notificationsForm.editableMessageId = notification.id;
		notificationsForm.form.setValues(notification);
		layout.scrollToTop();
	}

	// globals or users
	const notificationsFilter = ref('globals');

	return {
		table,
		tableColumns,
		notificationsFilter,
		onDeleteNotification,
		onEditNotification,
	};
});
