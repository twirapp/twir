<script setup lang="ts">
import { CheckIcon, SearchIcon, ListFilterIcon } from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';

import { useUsersTable } from '../composables/use-users-table';
import { FilterType, useUsersTableFilters } from '../composables/use-users-table-filters';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
	Command, CommandGroup, CommandItem, CommandList, CommandSeparator,
} from '@/components/ui/command';
import { Input } from '@/components/ui/input';
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover';
import { Separator } from '@/components/ui/separator';
import { cn } from '@/lib/utils';

const { t } = useI18n();
const usersTable = useUsersTable();
const usersTableFilters = useUsersTableFilters();

function applyFilter(filterKey: string, type: FilterType): void {
	usersTableFilters.setFilterValue(filterKey, type);
	usersTable.table.setPageIndex(0);
}

function isFilterApplied(filterKey: string, type: FilterType): boolean {
	if (type === 'status') {
		return filterKey in usersTableFilters.selectedStatuses;
	}

	if (type === 'badge') {
		return usersTableFilters.selectedBadges.includes(filterKey);
	}

	return false;
}
</script>

<template>
	<div class="flex gap-2 max-sm:flex-col">
		<div class="relative w-full items-center">
			<Input v-model="usersTableFilters.searchInput" type="text" :placeholder="t('sharedTexts.searchPlaceholder')" class="pl-10" />
			<span class="absolute start-2 inset-y-0 flex items-center justify-center px-2">
				<SearchIcon class="size-4 text-muted-foreground" />
			</span>
		</div>
		<Popover>
			<PopoverTrigger as-child>
				<Button variant="outline" size="sm" class="h-10">
					<ListFilterIcon class="mr-2 h-4 w-4" />
					{{ t('adminPanel.manageUsers.filters') }}

					<template v-if="usersTableFilters.selectedFiltersCount">
						<Separator orientation="vertical" class="mx-2 h-4" />
						<Badge
							variant="secondary"
							class="rounded-sm px-1 font-normal"
						>
							{{ t('adminPanel.manageUsers.countSelected', { count: usersTableFilters.selectedFiltersCount }) }}
						</Badge>
					</template>
				</Button>
			</PopoverTrigger>
			<PopoverContent class="w-[200px] p-0" align="end">
				<Command>
					<CommandList>
						<CommandGroup
							v-for="filters of usersTableFilters.filtersList"
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
									:class="[
										'mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary',
										isFilterApplied(filter.key, filters.type)
											? 'bg-primary text-primary-foreground'
											: 'opacity-50 [&_svg]:invisible'
									]"
								>
									<CheckIcon :class="cn('h-4 w-4')" />
								</div>
								<img v-if="filter.image" :src="filter.image" class="h-5 w-5 mr-2">
								<span>{{ filter.label }}</span>
							</CommandItem>
						</CommandGroup>

						<template v-if="usersTableFilters.selectedFiltersCount">
							<CommandSeparator />
							<CommandGroup>
								<CommandItem
									:value="{ label: 'Clear filters' }"
									class="justify-center text-center cursor-pointer"
									@select="usersTableFilters.clearFilters"
								>
									{{ t('adminPanel.manageUsers.clearFilters') }}
								</CommandItem>
							</CommandGroup>
						</template>
					</CommandList>
				</Command>
			</PopoverContent>
		</Popover>
	</div>
</template>
