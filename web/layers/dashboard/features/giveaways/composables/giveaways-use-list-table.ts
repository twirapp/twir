import type { Giveaway } from '#layers/dashboard/api/giveaways.ts'
import type { ColumnDef } from '@tanstack/vue-table'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { BanIcon, EyeIcon, PlayIcon } from 'lucide-vue-next'
import { computed, h, ref } from 'vue'

import { useGiveaways } from '~/features/giveaways/composables/giveaways-use-giveaways.ts'
import GiveawaysCreateDialog from '~/features/giveaways/ui/giveaways-create-dialog.vue'
import { ChannelRolePermissionEnum } from '~/gql/graphql.ts'

export const useGiveawaysListTable = createGlobalState(() => {
	const { t } = useI18n()
	const { activeGiveaways, giveawaysListFetching, viewGiveaway, startGiveaway, stopGiveaway } =
		useGiveaways()

	const canManageGiveaways = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageGiveaways)

	const showCreateDialog = ref(false)

	const tableColumns = computed<ColumnDef<Giveaway>[]>(() => {
		return [
			{
				accessorKey: 'keyword',
				size: 20,
				header: () => h('div', {}, t('giveaways.keyword')),
				cell: ({ row }) =>
					h('div', { class: 'flex items-center gap-2' }, [
						h('span', {}, row.original.keyword),
						// row.original.startedAt && !row.original.stoppedAt && !row.original.endedAt
						// 	? h(Badge, { variant: 'success' }, () => 'Active')
						// 	: row.original.stoppedAt
						// 		? h(Badge, { variant: 'secondary' }, () => 'Stopped')
						// 		: h(Badge, { variant: 'outline' }, () => 'Created'),
					]),
			},
			{
				accessorKey: 'createdAt',
				size: 20,
				header: () => h('div', {}, t('giveaways.createdAt')),
				cell: ({ row }) => h('span', {}, new Date(row.original.createdAt).toLocaleString()),
			},
			{
				accessorKey: 'startedAt',
				size: 20,
				header: () => h('div', {}, t('giveaways.startedAt')),
				cell: ({ row }) =>
					h(
						'span',
						{},
						row.original.startedAt ? new Date(row.original.startedAt).toLocaleString() : '-'
					),
			},
			{
				accessorKey: 'actions',
				size: 40,
				header: () => h('div', { class: 'flex justify-end' }, [h(GiveawaysCreateDialog)]),
				cell: ({ row }) =>
					h('div', { class: 'flex gap-2 justify-end' }, [
						// View button
						h(
							Button,
							{
								size: 'sm',
								variant: 'outline',
								class: 'flex gap-2 items-center',
								onClick: () => viewGiveaway(row.original.id),
							},
							{
								default: () => [h(EyeIcon, { class: 'size-4' }), t('giveaways.view')],
							}
						),

						// Start button (if not started)
						!row.original.startedAt && !row.original.stoppedAt
							? h(
									Button,
									{
										size: 'sm',
										variant: 'default',
										class: 'flex gap-2 items-center',
										disabled: !canManageGiveaways.value,
										onClick: () => startGiveaway(row.original.id),
									},
									{
										default: () => [h(PlayIcon, { class: 'size-4' }), t('giveaways.start')],
									}
								)
							: null,

						// Stop button (if started and not stopped)
						row.original.startedAt && !row.original.stoppedAt
							? h(
									Button,
									{
										size: 'sm',
										variant: 'secondary',
										class: 'flex gap-2 items-center',
										disabled: !canManageGiveaways.value,
										onClick: () => stopGiveaway(row.original.id),
									},
									{
										default: () => [h(BanIcon, { class: 'size-4' }), t('giveaways.stop')],
									}
								)
							: null,
					]),
			},
		]
	})

	const table = useVueTable({
		get data() {
			return activeGiveaways.value
		},
		get columns() {
			return tableColumns.value
		},
		getCoreRowModel: getCoreRowModel(),
	})

	return {
		isLoading: giveawaysListFetching,
		table,
		showCreateDialog,
	}
})
