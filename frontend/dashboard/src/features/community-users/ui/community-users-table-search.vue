<script setup lang="ts">
import { Settings2Icon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useCommunityTableActions } from '../composables/use-community-table-actions.js'
import {
	TABLE_ACCESSOR_KEYS,
	useCommunityUsersTable,
} from '../composables/use-community-users-table.js'

import SearchBar from '@/components/search-bar.vue'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuCheckboxItem,
	DropdownMenuContent,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

const { t } = useI18n()
const communityTableActions = useCommunityTableActions()
const communityUsersTable = useCommunityUsersTable()

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
		<DropdownMenu>
			<DropdownMenuTrigger as-child>
				<Button variant="outline" size="sm" class="flex ml-auto h-9">
					<Settings2Icon class="mr-2 h-4 w-4" />
					{{ t('sharedTexts.view') }}
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="end" class="w-[200px]">
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
