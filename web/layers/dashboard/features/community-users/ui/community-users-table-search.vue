<script setup lang="ts">
import { Settings2Icon } from 'lucide-vue-next'
import { computed } from 'vue'


import { useCommunityTableActions } from '../composables/use-community-table-actions.js'
import {
	TABLE_ACCESSOR_KEYS,
	useCommunityUsersTable,
} from '../composables/use-community-users-table.js'

import SearchBar from '#layers/dashboard/components/search-bar.vue'



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
		<UiDropdownMenu>
			<UiDropdownMenuTrigger as-child>
				<UiButton variant="outline" size="sm" class="flex ml-auto h-9">
					<Settings2Icon class="mr-2 h-4 w-4" />
					{{ t('sharedTexts.view') }}
				</UiButton>
			</UiDropdownMenuTrigger>
			<UiDropdownMenuContent align="end" class="w-[200px]">
				<UiDropdownMenuLabel>
					{{ t('sharedTexts.toggleColumns') }}
				</UiDropdownMenuLabel>
				<UiDropdownMenuSeparator />

				<UiDropdownMenuCheckboxItem
					v-for="column in columns"
					:key="column.id"
					class="capitalize"
					:checked="column.getIsVisible()"
					@update:checked="(value: boolean | 'indeterminate') => column.toggleVisibility(!!value)"
				>
					{{ column.id }}
				</UiDropdownMenuCheckboxItem>
			</UiDropdownMenuContent>
		</UiDropdownMenu>
	</div>
</template>
