<script setup lang="ts">
import { CheckIcon, ListFilterIcon } from 'lucide-vue-next'
import { computed } from 'vue'


import { useAuditFilters } from '../composables/use-audit-filters'
import { useAuditTable } from '../composables/use-audit-table'

import type { AuditFilterType } from '../composables/use-audit-filters'

import TwitchUserSearch from '#layers/dashboard/components/twitchUsers/twitch-user-select.vue'






const { t } = useI18n()
const auditTable = useAuditTable()
const {
	searchType,
	searchUserId,
	searchOptions,
	filtersList,
	selectedFiltersCount,
	clearFilters,
	setFilterValue,
	isFilterApplied,
} = useAuditFilters()

function applyFilter(filterKey: string, type: AuditFilterType): void {
	setFilterValue(type, filterKey)
	auditTable.table.setPageIndex(0)
}

const usersSearchPlaceholder = computed(() => {
	return searchType.value === 'channel'
		? t('dashboard.widgets.audit-logs.search.channel')
		: t('dashboard.widgets.audit-logs.search.actor')
})
</script>

<template>
	<div class="flex gap-2">
		<TwitchUserSearch
			v-model="searchUserId"
			:placeholder="usersSearchPlaceholder"
		/>

		<UiPopover>
			<UiPopoverTrigger as-child>
				<UiButton variant="outline" size="sm" class="h-10">
					<ListFilterIcon class="mr-2 h-4 w-4" />
					{{ t('adminPanel.manageUsers.filters') }}

					<UiSeparator orientation="vertical" class="mx-2 h-4" />
					<UiBadge
						variant="secondary"
						class="rounded-sm px-1 font-normal"
					>
						{{ t('adminPanel.manageUsers.countSelected', {
							count: selectedFiltersCount,
						}) }}
					</UiBadge>
				</UiButton>
			</UiPopoverTrigger>
			<UiPopoverContent class="w-[250px] p-0" align="end">
				<UiCommand>
					<UiCommandList class="mb-10">
						<UiCommandGroup :heading="t('dashboard.widgets.audit-logs.search-label')">
							<UiCommandItem
								v-for="option of searchOptions"
								:key="option.value"
								:value="option.value"
								@select="searchType = option.value"
							>
								<div
									class="mr-2 flex h-4 w-4 items-center justify-center rounded-full border border-primary"
									:class="[searchType === option.value
										? 'bg-primary text-primary-foreground'
										: 'opacity-50 [&_svg]:invisible',
									]"
								>
									<CheckIcon class="h-3 w-3" />
								</div>
								<span>{{ option.label }}</span>
							</UiCommandItem>
						</UiCommandGroup>

						<UiCommandGroup
							v-for="filters of filtersList"
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
									:class="[isFilterApplied(filters.type, filter.key)
										? 'bg-primary text-primary-foreground'
										: 'opacity-50 [&_svg]:invisible',
									]"
								>
									<CheckIcon class="h-4 w-4" />
								</div>
								<span>{{ filter.label }}</span>
							</UiCommandItem>
						</UiCommandGroup>

						<div class="absolute bottom-0 w-full bg-background">
							<UiCommandGroup class="border-t border-border">
								<UiCommandItem
									:value="{ label: 'Clear filters' }"
									class="justify-center text-center cursor-pointer"
									@select="clearFilters"
								>
									{{ t('adminPanel.manageUsers.clearFilters') }}
								</UiCommandItem>
							</UiCommandGroup>
						</div>
					</UiCommandList>
				</UiCommand>
			</UiPopoverContent>
		</UiPopover>
	</div>
</template>
