<script setup lang="ts">
import { CheckIcon, ListFilterIcon } from 'lucide-vue-next'


import { useUsersTableFilters } from '../composables/use-users-table-filters.js'
import { useUsersTable } from '../composables/use-users-table.js'

import type { FilterType } from '../composables/use-users-table-filters.js'

import SearchBar from '#layers/dashboard/components/search-bar.vue'






const { t } = useI18n()
const usersTable = useUsersTable()
const usersTableFilters = useUsersTableFilters()

function applyFilter(filterKey: string, type: FilterType): void {
	usersTableFilters.setFilterValue(filterKey, type)
	usersTable.table.setPageIndex(0)
}
</script>

<template>
	<div class="flex gap-2">
		<SearchBar v-model="usersTableFilters.searchInput.value" />
		<UiPopover>
			<UiPopoverTrigger as-child>
				<UiButton variant="outline" size="sm" class="h-9">
					<ListFilterIcon class="mr-2 h-4 w-4" />
					{{ t('adminPanel.manageUsers.filters') }}

					<template v-if="usersTableFilters.selectedFiltersCount">
						<UiSeparator orientation="vertical" class="mx-2 h-4" />
						<UiBadge
							variant="secondary"
							class="rounded-sm px-1 font-normal"
						>
							{{ t('adminPanel.manageUsers.countSelected', { count: usersTableFilters.selectedFiltersCount.value }) }}
						</UiBadge>
					</template>
				</UiButton>
			</UiPopoverTrigger>
			<UiPopoverContent class="w-[200px] p-0" align="end">
				<UiCommand>
					<UiCommandList>
						<UiCommandGroup
							v-for="filters of usersTableFilters.filtersList.value"
							:key="filters.group"
							:heading="filters.group"
						>
							<UiCommandItem
								v-for="filter of filters.list"
								:key="filter.key"
								:value="filter.key"
								@select="applyFilter(filter.key, filters.type)"
							>
								<div
									class="mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary"
									:class="[usersTableFilters.isFilterApplied(filter.key, filters.type)
										? 'bg-primary text-primary-foreground'
										: 'opacity-50 [&_svg]:invisible',
									]"
								>
									<CheckIcon class="h-4 w-4" />
								</div>
								<img v-if="filter.image" :src="filter.image" class="h-5 w-5 mr-2">
								<span>{{ filter.label }}</span>
							</UiCommandItem>
						</UiCommandGroup>

						<template v-if="usersTableFilters.selectedFiltersCount">
							<UiCommandSeparator />
							<UiCommandGroup>
								<UiCommandItem
									:value="{ label: 'Clear filters' }"
									class="justify-center text-center cursor-pointer"
									@select="usersTableFilters.clearFilters"
								>
									{{ t('adminPanel.manageUsers.clearFilters') }}
								</UiCommandItem>
							</UiCommandGroup>
						</template>
					</UiCommandList>
				</UiCommand>
			</UiPopoverContent>
		</UiPopover>
	</div>
</template>
