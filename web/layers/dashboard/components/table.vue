<script setup lang="ts" generic="T extends RowData">
import { FlexRender, type RowData, type Table } from '@tanstack/vue-table'
import { ListX } from 'lucide-vue-next'

import { Card } from '@/components/ui/card'
import { useIsMobile } from '#layers/dashboard/composables/use-is-mobile'
import ShadcnLayout from '#layers/dashboard/shadcn-layout.vue'

defineProps<{
	table: Table<T>
	isLoading: boolean
	hideHeader?: boolean
}>()

const { isDesktop } = useIsMobile()
const { t } = useI18n()
</script>

<template>
	<ShadcnLayout v-if="isDesktop">
		<UiTable>
			<UiTableHeader v-if="!hideHeader">
				<UiTableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id" class="border-b">
					<UiTableHead v-for="header in headerGroup.headers" :key="header.id" :style="{ width: `${header.getSize()}%` }">
						<FlexRender
							v-if="!header.isPlaceholder"
							:render="header.column.columnDef.header"
							:props="header.getContext()"
						/>
					</UiTableHead>
				</UiTableRow>
			</UiTableHeader>
			<UiTableBody :class="[isLoading ? 'animate-pulse' : '']">
				<template v-if="table.getRowModel().rows?.length">
					<UiTableRow
						v-for="row in table.getRowModel().rows" :key="row.id"
						:data-state="row.getIsSelected() ? 'selected' : undefined" class="border-b"
					>
						<UiTableCell
							v-for="cell in row.getVisibleCells()"
							:key="cell.id"
							class="md:break-all"
							:class="{
								'cursor-pointer': row.getCanExpand(),
							}"
							@click="() => {
								if (row.getCanExpand()) {
									row.getToggleExpandedHandler()()
								}
							}"
						>
							<FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
						</UiTableCell>
					</UiTableRow>
				</template>
				<template v-else>
					<UiTableRow>
						<UiTableCell :colSpan="table.getAllColumns().length" class="h-24 text-center">
							<slot name="empty-message">
								<div class="flex items-center flex-col justify-center">
									<ListX class="size-12" />
									<span class="font-medium text-2xl">{{ t('sharedTexts.noData') }}</span>
								</div>
							</slot>
						</UiTableCell>
					</UiTableRow>
				</template>
			</UiTableBody>
		</UiTable>
	</ShadcnLayout>

	<div v-else class="grid grid-cols-1 gap-4">
		<Card v-for="row in table.getRowModel().rows" :key="row.id">
			<div
				v-for="cell in row.getVisibleCells()"
				:key="cell.id"
			>
				<div v-if="row.getCanExpand()" class="px-2 my-2 cursor-pointer">
					<FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" @click="() => row.getToggleExpandedHandler()()" />
				</div>

				<div v-else-if="cell.column.id !== 'actions'" class="px-4 py-2 border-b-2">
					<FlexRender :render="cell.column.columnDef.header" class="text-sm text-zinc-400/80" />
					<FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
				</div>

				<div v-else class="flex h-auto py-2 px-2 justify-end">
					<FlexRender
						:render="cell.column.columnDef.cell"
						:props="cell.getContext()"
					/>
				</div>
			</div>
		</Card>
	</div>
</template>
