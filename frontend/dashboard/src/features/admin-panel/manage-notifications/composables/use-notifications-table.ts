import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table';
import type { Notification } from '@twir/api/messages/admin_notifications/admin_notifications';
import { defineStore } from 'pinia';
import { computed, h } from 'vue';
import { useI18n } from 'vue-i18n';

import { useNotificationsFilters } from './use-notifications-filters.js';
import CreatedAtTooltip from '../components/created-at-tooltip.vue';
import NotificationsTableActions from '../components/notifications-table-actions.vue';
import { useNotificationsForm } from '../composables/use-notifications-form.js';

import { useAdminNotifications } from '@/api/admin/notifications.js';
import { useLayout } from '@/composables/use-layout.js';
import { usePagination } from '@/composables/use-pagination.js';

export const useNotificationsTable = defineStore('admin-panel/notifications-table', () => {
	const layout = useLayout();
	const { t } = useI18n();

	const form = useNotificationsForm();
	const crud = useAdminNotifications();

	const { pagination, setPagination } = usePagination();
	const filters = useNotificationsFilters();

	const reqParams = computed(() => ({
		perPage: pagination.value.pageSize,
		page: pagination.value.pageIndex,
		isUser: filters.isUsersFilter,
		search: filters.debounceSearchInput,
	}));

	const { data } = crud.getAll(reqParams);
	const notifications = computed<Notification[]>(() => {
		return data.value?.notifications ?? [];
	});

	const tableColumns = computed<ColumnDef<Notification>[]>(() => [
		{
			accessorKey: 'message',
			size: 80,
			header: () => h('div', {}, t('adminPanel.notifications.messageLabel')),
			cell: ({ row }) => {
				return h('div', { class: 'break-words max-w-[450px]', innerHTML: row.original.message });
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

	const totalNotifications = computed(() => data.value?.total ?? 0);

	const pageCount = computed(() => {
		return Math.ceil((data.value?.total ?? 0) / pagination.value.pageSize);
	});

	const table = useVueTable({
		get pageCount() {
			return pageCount.value;
		},
		state: {
			pagination: pagination.value,
		},
		get data() {
			return notifications.value;
		},
		get columns() {
			return tableColumns.value;
		},
		manualPagination: true,
		getCoreRowModel: getCoreRowModel(),
		onPaginationChange: (updater) => {
			if (typeof updater === 'function') {
				setPagination(
					updater({
						pageIndex: pagination.value.pageIndex,
						pageSize: pagination.value.pageSize,
					}),
				);
			} else {
				setPagination(updater);
			}
		},
	});

	async function onDeleteNotification(notificationId: string) {
		if (form.editableMessageId === notificationId) {
			form.onReset();
		}

		await crud.deleteOne.mutateAsync({ id: notificationId });
	}

	async function onEditNotification(notification: Notification) {
		let isConfirmed = true;

		if (form.formValues.message || form.isEditableForm) {
			isConfirmed = confirm(t('adminPanel.notifications.confirmResetForm'));
		}

		if (isConfirmed) {
			form.editableMessageId = notification.id;
			form.userIdField.fieldModel = notification.userId ?? null;
			form.messageField.fieldModel = notification.message;
			layout.scrollToTop();
		}
	}

	return {
		table,
		tableColumns,
		notifications,
		totalNotifications,
		onDeleteNotification,
		onEditNotification,
	};
});
