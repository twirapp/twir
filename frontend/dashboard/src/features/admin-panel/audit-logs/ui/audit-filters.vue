<script setup lang="ts">
import { CheckIcon, ListFilterIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useAuditFilters } from '../composables/use-audit-filters'
import { useAuditTable } from '../composables/use-audit-table'

import type { AuditFilterType } from '../composables/use-audit-filters'

import TwitchUserSearch from '@/components/twitchUsers/twitch-user-select.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
	Command,
	CommandGroup,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover'
import { Separator } from '@/components/ui/separator'

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

		<Popover>
			<PopoverTrigger as-child>
				<Button variant="outline" size="sm" class="h-10">
					<ListFilterIcon class="mr-2 h-4 w-4" />
					{{ t('adminPanel.manageUsers.filters') }}

					<Separator orientation="vertical" class="mx-2 h-4" />
					<Badge
						variant="secondary"
						class="rounded-sm px-1 font-normal"
					>
						{{ t('adminPanel.manageUsers.countSelected', {
							count: selectedFiltersCount,
						}) }}
					</Badge>
				</Button>
			</PopoverTrigger>
			<PopoverContent class="w-[250px] p-0" align="end">
				<Command>
					<CommandList class="mb-10">
						<CommandGroup :heading="t('dashboard.widgets.audit-logs.search-label')">
							<CommandItem
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
							</CommandItem>
						</CommandGroup>

						<CommandGroup
							v-for="filters of filtersList"
							:key="filters.group"
							:heading="filters.group"
						>
							<CommandItem
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
							</CommandItem>
						</CommandGroup>

						<div class="absolute bottom-0 w-full bg-background">
							<CommandGroup class="border-t-[1px] border-border">
								<CommandItem
									:value="{ label: 'Clear filters' }"
									class="justify-center text-center cursor-pointer"
									@select="clearFilters"
								>
									{{ t('adminPanel.manageUsers.clearFilters') }}
								</CommandItem>
							</CommandGroup>
						</div>
					</CommandList>
				</Command>
			</PopoverContent>
		</Popover>
	</div>
</template>
