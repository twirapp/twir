import {
	type ColumnDef,
	getCoreRowModel,
	getPaginationRowModel,
	useVueTable,
} from '@tanstack/vue-table';
import type { Notification } from '@twir/api/messages/admin_notifications/admin_notifications';
import { NText } from 'naive-ui';
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
	const notificationsData = notificationsCrud.getAll({
		perPage: 50,
		page: 0,
	});

	const tableColumns = computed<ColumnDef<Notification>[]>(() => [
		{
			accessorKey: 'message',
			size: 80,
			header: () => h('div', {}, t('adminPanel.notifications.messageLabel')),
			cell: ({ row }) => {
				return h(NText, { class: 'break-words w-[450px]', innerHTML: row.original.message });
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

	async function onDeleteNotification(notificationId: string) {
		if (notificationsForm.editableMessageId === notificationId) {
			notificationsForm.onReset();
		}

		await notificationsCrud.deleteOne.mutateAsync({ id: notificationId });
	}

	async function onEditNotification(notification: Notification) {
		const confirmed = !notificationsForm.isFormDirty
			|| notificationsForm.editableMessageId
			&& confirm(t('adminPanel.notifications.confirmResetForm'));
		if (!confirmed) return;

		notificationsForm.editableMessageId = notification.id;
		notificationsForm.form.setValues(notification);
		layout.scrollToTop();
	}

	// globals or users
	const notificationsFilter = ref('globals');

	return {
		table,
		tableColumns,
		notifications,
		notificationsFilter,
		onDeleteNotification,
		onEditNotification,
	};
});
