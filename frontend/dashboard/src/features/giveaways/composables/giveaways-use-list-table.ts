import { getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { ArchiveIcon, BanIcon, EyeIcon, PlayIcon, PlusIcon } from 'lucide-vue-next'
import { computed, h, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import type { Giveaway } from '@/api/giveaways.ts'
import type { ColumnDef } from '@tanstack/vue-table'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Dialog, DialogTrigger } from '@/components/ui/dialog'
import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'
import GiveawaysCreateDialog from '@/features/giveaways/ui/giveaways-create-dialog.vue'

export const useGiveawaysListTable = createGlobalState(() => {
	const { t } = useI18n()
	const {
		activeGiveaways,
		giveawaysListFetching,
		viewGiveaway,
		startGiveaway,
		stopGiveaway,
		archiveGiveaway,
	} = useGiveaways()

	const showCreateDialog = ref(false)

	const tableColumns = computed<ColumnDef<Giveaway>[]>(() => {
		return [
			{
				accessorKey: 'keyword',
				size: 20,
				header: () => h('div', {}, t('giveaways.keyword')),
				cell: ({ row }) => h('div', { class: 'flex items-center gap-2' }, [
					h('span', {}, row.original.keyword),
					row.original.startedAt && !row.original.stoppedAt && !row.original.endedAt
						? h(Badge, { variant: 'success' }, () => 'Active')
						: row.original.stoppedAt
							? h(Badge, { variant: 'secondary' }, () => 'Stopped')
							: h(Badge, { variant: 'outline' }, () => 'Created'),
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
				cell: ({ row }) => h('span', {}, row.original.startedAt ? new Date(row.original.startedAt).toLocaleString() : '-'),
			},
			{
				accessorKey: 'actions',
				size: 40,
				header: () => h(Dialog, { 'onUpdate:open': (val) => showCreateDialog.value = val, 'open': showCreateDialog.value }, {
					default: () => [
						h(DialogTrigger, { asChild: true }, {
							default: () => h(Button, { size: 'sm', class: 'flex gap-2 items-center' }, {
								default: () => [
									h(PlusIcon, { class: 'size-4' }),
									t('giveaways.createNew'),
								],
							}),
						}),
						h(GiveawaysCreateDialog, {
							'open': showCreateDialog.value,
							'onUpdate:open': (val) => showCreateDialog.value = val,
						}),
					],
				}),
				cell: ({ row }) => h('div', { class: 'flex gap-2 justify-end' }, [
					// View button
					h(Button, {
						size: 'sm',
						variant: 'outline',
						class: 'flex gap-2 items-center',
						onClick: () => viewGiveaway(row.original.id),
					}, {
						default: () => [
							h(EyeIcon, { class: 'size-4' }),
							t('giveaways.view'),
						],
					}),

					// Start button (if not started)
					!row.original.startedAt && !row.original.stoppedAt
						? h(Button, {
							size: 'sm',
							variant: 'default',
							class: 'flex gap-2 items-center',
							onClick: () => startGiveaway(row.original.id),
						}, {
							default: () => [
								h(PlayIcon, { class: 'size-4' }),
								t('giveaways.start'),
							],
						})
						: null,

					// Stop button (if started and not stopped)
					row.original.startedAt && !row.original.stoppedAt
						? h(Button, {
							size: 'sm',
							variant: 'secondary',
							class: 'flex gap-2 items-center',
							onClick: () => stopGiveaway(row.original.id),
						}, {
							default: () => [
								h(BanIcon, { class: 'size-4' }),
								t('giveaways.stop'),
							],
						})
						: null,

					// Archive button
					h(Button, {
						size: 'sm',
						variant: 'destructive',
						class: 'flex gap-2 items-center',
						onClick: () => archiveGiveaway(row.original.id),
					}, {
						default: () => [
							h(ArchiveIcon, { class: 'size-4' }),
							t('giveaways.archive'),
						],
					}),
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
