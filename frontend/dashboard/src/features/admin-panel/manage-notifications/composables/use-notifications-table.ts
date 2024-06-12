import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import { useNotificationsFilters } from './use-notifications-filters.js'
import UsersTableCellUser from '../../manage-users/ui/users-table-cell-user.vue'
import { useNotificationsForm } from '../composables/use-notifications-form.js'
import CreatedAtTooltip from '../ui/created-at-tooltip.vue'
import NotificationsTableActions from '../ui/notifications-table-actions.vue'

import { useAdminNotifications } from '@/api/admin/notifications.js'
import BlocksRenderer from '@/components/ui/editorjs/blocks-render.vue'
import { useLayout } from '@/composables/use-layout.js'
import { usePagination } from '@/composables/use-pagination.js'
import {
	type AdminNotificationsParams,
	NotificationType,
	type NotificationsByAdminQuery,
} from '@/gql/graphql.js'
import { valueUpdater } from '@/helpers/value-updater.js'

type Notifications = NotificationsByAdminQuery['notificationsByAdmin']['notifications']

export const useNotificationsTable = createGlobalState(() => {
	const layout = useLayout()
	const { t } = useI18n()

	const form = useNotificationsForm()

	const { pagination } = usePagination()
	const { filterInput, debounceSearchInput } = useNotificationsFilters()

	const params = computed<AdminNotificationsParams>(() => ({
		perPage: pagination.value.pageSize,
		page: pagination.value.pageIndex,
		type: filterInput.value,
		search: debounceSearchInput.value,
	}))

	const notificationsApi = useAdminNotifications()
	const { data, fetching } = notificationsApi.useQueryNotifications(params)
	const { executeMutation: deleteNotification } = notificationsApi.useMutationDeleteNotification()

	const notifications = computed<Notifications>(() => {
		if (!data.value) return []
		return data.value.notificationsByAdmin.notifications
	})

	const tableColumns = computed(() => {
		const columns: ColumnDef<Notifications[0]>[] = [
			{
				accessorKey: 'message',
				size: 65,
				header: () => h('div', {}, t('adminPanel.notifications.messageLabel')),
				cell: ({ row }) => {
					if (row.original.text) {
						return h('div', { class: 'break-words max-w-[450px]', innerHTML: row.original.text })
					} else if (row.original.editorJsJson) {
						return h(BlocksRenderer, { data: row.original.editorJsJson })
					}
				},
			},
			{
				accessorKey: 'createdAt',
				size: 10,
				header: () => h('div', {}, t('adminPanel.notifications.createdAt')),
				cell: ({ row }) => {
					return h(CreatedAtTooltip, { time: new Date(row.original.createdAt) })
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
					})
				},
			},
		]

		if (filterInput.value === NotificationType.User) {
			columns.unshift({
				accessorKey: 'id',
				size: 10,
				header: () => h('div', {}, t('adminPanel.notifications.userLabel')),
				cell: ({ row }) => {
					if (row.original.twitchProfile?.profileImageUrl && row.original.twitchProfile?.displayName) {
						return h('a', {
							class: 'flex flex-col',
							href: `https://twitch.tv/${row.original.twitchProfile.displayName.toLowerCase()}`,
							target: '_blank',
						}, h(UsersTableCellUser, {
							avatar: row.original.twitchProfile.profileImageUrl,
							name: row.original.twitchProfile.login,
							displayName: row.original.twitchProfile.displayName,
						}))
					}
				},
			})
		}

		return columns
	})

	const totalNotifications = computed(() => data.value?.notificationsByAdmin.total ?? 0)

	const pageCount = computed(() => {
		return Math.ceil((data.value?.notificationsByAdmin.total ?? 0) / pagination.value.pageSize)
	})

	const table = useVueTable({
		get pageCount() {
			return pageCount.value
		},
		state: {
			pagination: pagination.value,
		},
		get data() {
			return notifications.value
		},
		get columns() {
			return tableColumns.value
		},
		manualPagination: true,
		getCoreRowModel: getCoreRowModel(),
		onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, pagination),
	})

	async function onDeleteNotification(notificationId: string) {
		if (form.editableMessageId.value === notificationId) {
			form.onReset()
		}

		await deleteNotification({ id: notificationId })
	}

	async function onEditNotification(notification: Notifications[0]) {
		let isConfirmed = true

		if (form.formValues.value.editorJsJson || form.isEditableForm) {
			// TODO: use confirm dialog from shadcn
			// eslint-disable-next-line no-alert
			isConfirmed = confirm(t('adminPanel.notifications.confirmResetForm'))
		}

		if (isConfirmed) {
			form.editableMessageId.value = notification.id
			form.userIdField.fieldModel.value = notification.userId ?? null
			form.editorJsJsonField.fieldModel.value = notification.editorJsJson
			layout.scrollToTop()
		}
	}

	return {
		isLoading: fetching,
		table,
		tableColumns,
		pagination,
		notifications,
		totalNotifications,
		onDeleteNotification,
		onEditNotification,
	}
})
