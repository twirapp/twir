<script setup lang="ts">
import { computed } from 'vue'
import SearchBar from '~~/layers/dashboard/components/search-bar.vue'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
	Command,
	CommandGroup,
	CommandItem,
	CommandList,
	CommandSeparator,
} from '@/components/ui/command'
import {
	DropdownMenu,
	DropdownMenuCheckboxItem,
	DropdownMenuContent,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Separator } from '@/components/ui/separator'
import { Platform } from '~/gql/graphql.js'

import { useCommunityTableActions } from '../composables/use-community-table-actions.js'
import {
	TABLE_ACCESSOR_KEYS,
	useCommunityUsersTable,
} from '../composables/use-community-users-table.js'

const { t } = useI18n()
const communityTableActions = useCommunityTableActions()
const communityUsersTable = useCommunityUsersTable()

const platformOptions = [
	{ label: 'Twitch', value: Platform.Twitch },
	{ label: 'Kick', value: Platform.Kick },
]

function togglePlatform(platform: Platform) {
	if (communityTableActions.selectedPlatforms.value.includes(platform)) {
		communityTableActions.selectedPlatforms.value =
			communityTableActions.selectedPlatforms.value.filter((item) => item !== platform)
		communityUsersTable.table.setPageIndex(0)
		return
	}

	communityTableActions.selectedPlatforms.value = [
		...communityTableActions.selectedPlatforms.value,
		platform,
	]
	communityUsersTable.table.setPageIndex(0)
}

function clearFilters() {
	communityTableActions.selectedPlatforms.value = []
	communityUsersTable.table.setPageIndex(0)
}

// TODO: column labels
const columns = computed(() => {
	return communityUsersTable.table.getAllColumns().filter((column) => {
		if (column.id !== TABLE_ACCESSOR_KEYS.user) {
			return typeof column.accessorFn !== 'undefined' && column.getCanHide()
		}

		return false
	})
})
</script>

<template>
	<div class="flex gap-2">
		<SearchBar
			v-model="communityTableActions.searchInput.value"
			:placeholder="t('community.users.searchPlaceholder')"
		/>
		<Popover>
			<PopoverTrigger as-child>
				<Button
					variant="outline"
					size="sm"
					class="h-9"
				>
					<Icon
						name="lucide:list-filter"
						class="mr-2 h-4 w-4"
					/>
					Filters

					<template v-if="communityTableActions.selectedFiltersCount.value">
						<Separator
							orientation="vertical"
							class="mx-2 h-4"
						/>
						<Badge
							variant="secondary"
							class="rounded-sm px-1 font-normal"
						>
							{{ communityTableActions.selectedFiltersCount.value }}
						</Badge>
					</template>
				</Button>
			</PopoverTrigger>
			<PopoverContent
				class="w-[200px] p-0"
				align="end"
			>
				<Command>
					<CommandList>
						<CommandGroup heading="Platforms">
							<CommandItem
								v-for="option in platformOptions"
								:key="option.value"
								:value="option.value"
								@select="togglePlatform(option.value)"
							>
								<div
									class="border-primary mr-2 flex h-4 w-4 items-center justify-center rounded-sm border"
									:class="
										communityTableActions.selectedPlatforms.value.includes(option.value)
											? 'bg-primary text-primary-foreground'
											: 'opacity-50 [&_svg]:invisible'
									"
								>
									<Icon
										name="lucide:check"
										class="h-4 w-4"
									/>
								</div>
								<span>{{ option.label }}</span>
							</CommandItem>
						</CommandGroup>

						<template v-if="communityTableActions.selectedFiltersCount.value">
							<CommandSeparator />
							<CommandGroup>
								<CommandItem
									value="clear-filters"
									class="cursor-pointer justify-center text-center"
									@select="clearFilters"
								>
									Clear filters
								</CommandItem>
							</CommandGroup>
						</template>
					</CommandList>
				</Command>
			</PopoverContent>
		</Popover>
		<DropdownMenu>
			<DropdownMenuTrigger as-child>
				<Button
					variant="outline"
					size="sm"
					class="ml-auto flex h-9"
				>
					<Icon
						name="lucide:settings2"
						class="mr-2 h-4 w-4"
					/>
					{{ t('sharedTexts.view') }}
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent
				align="end"
				class="w-[200px]"
			>
				<DropdownMenuLabel>
					{{ t('sharedTexts.toggleColumns') }}
				</DropdownMenuLabel>
				<DropdownMenuSeparator />

				<DropdownMenuCheckboxItem
					v-for="column in columns"
					:key="column.id"
					class="capitalize"
					:checked="column.getIsVisible()"
					@update:checked="(value: boolean | 'indeterminate') => column.toggleVisibility(!!value)"
				>
					{{ column.id }}
				</DropdownMenuCheckboxItem>
			</DropdownMenuContent>
		</DropdownMenu>
	</div>
</template>
