import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table';
import { defineStore } from 'pinia';
import { computed, h } from 'vue';
import { useI18n } from 'vue-i18n';

import { useNotificationsFilters } from './use-notifications-filters.js';
import UsersTableCellUser from '../../manage-users/components/users-table-cell-user.vue';
import CreatedAtTooltip from '../components/created-at-tooltip.vue';
import NotificationsTableActions from '../components/notifications-table-actions.vue';
import { useNotificationsForm } from '../composables/use-notifications-form.js';

import { useAdminNotifications } from '@/api/admin/notifications.js';
import { useLayout } from '@/composables/use-layout.js';
import { usePagination } from '@/composables/use-pagination.js';
import { NotificationType, type AdminNotification, type AdminNotificationsParams } from '@/gql/graphql.js';

export const useNotificationsTable = defineStore('admin-panel/notifications-table', () => {
	const layout = useLayout();
	const { t } = useI18n();

	const form = useNotificationsForm();

	const { pagination, setPagination } = usePagination();
	const filters = useNotificationsFilters();

	const reqParams = computed<AdminNotificationsParams>(() => ({
		perPage: pagination.value.pageSize,
		page: pagination.value.pageIndex,
		type: filters.filterInput,
		search: filters.debounceSearchInput,
	}));

	const notificationsApi = useAdminNotifications();
	const { data } = notificationsApi.useQueryNotifications(reqParams);
	const { executeMutation: deleteNotification } = notificationsApi.useMutationDeleteNotification();

	const notifications = computed(() => {
		if (!data.value) return [];
		return data.value.notificationsByAdmin.notifications as AdminNotification[];
	});

	const tableColumns = computed<ColumnDef<AdminNotification>[]>(() => {
		const columns: ColumnDef<AdminNotification>[] = [
			{
				accessorKey: 'message',
				size: 65,
				header: () => h('div', {}, t('adminPanel.notifications.messageLabel')),
				cell: ({ row }) => {
					return h('div', { class: 'break-words max-w-[450px]', innerHTML: row.original.text });
				},
			},
			{
				accessorKey: 'createdAt',
				size: 10,
				header: () => h('div', {}, t('adminPanel.notifications.createdAt')),
				cell: ({ row }) => {
					return h(CreatedAtTooltip, { time: new Date(row.original.createdAt) });
				},
			},
			{
				accessorKey: 'actions',
				size: 15,
				header: () => '',
				cell: ({ row }) => {
					return h(NotificationsTableActions, {
						onDelete: () => onDeleteNotification(row.original.id),
						onEdit: () => onEditNotification(row.original),
					});
				},
			},
		];

		if (filters.filterInput === NotificationType.User) {
			columns.unshift({
				accessorKey: 'id',
				size: 10,
				header: () => h('div', {}, t('adminPanel.notifications.userLabel')),
				cell: ({ row }) => {
					if (row.original.twitchProfile?.profileImageUrl && row.original.twitchProfile?.displayName) {
						return h('a',
							{
								class: 'flex flex-col',
								href: `https://twitch.tv/${row.original.twitchProfile.displayName.toLowerCase()}`,
								target: '_blank',
							},
							h(UsersTableCellUser, {
								userId: row.original.id,
								avatar: row.original.twitchProfile.profileImageUrl,
								name: row.original.twitchProfile.displayName,
							}),
						);
					}
				},
			});
		}

		return columns;
	});

	const totalNotifications = computed(() => data.value?.notificationsByAdmin.total ?? 0);

	const pageCount = computed(() => {
		return Math.ceil((data.value?.notificationsByAdmin.total ?? 0) / pagination.value.pageSize);
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

		await deleteNotification({ id: notificationId });
	}

	async function onEditNotification(notification: AdminNotification) {
		let isConfirmed = true;

		if (form.formValues.message || form.isEditableForm) {
			isConfirmed = confirm(t('adminPanel.notifications.confirmResetForm'));
		}

		if (isConfirmed) {
			form.editableMessageId = notification.id;
			form.userIdField.fieldModel = notification.userId ?? null;
			form.messageField.fieldModel = notification.text;
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
