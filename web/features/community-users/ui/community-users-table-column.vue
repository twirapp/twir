<script setup lang="ts">
import type { CommunityUsersResetType } from '@/gql/graphql.js'
import type { Column } from '@tanstack/vue-table'

import {
	DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu'
import { cn } from '@/lib/utils.js'

defineOptions({
	inheritAttrs: false,
})

defineProps<{
	columnType?: CommunityUsersResetType
	column: Column<any, any>
	title: string
	hideReset?: boolean
}>()

// const communityUsersApi = useCommunityUsersApi()
// const communityResetMutation = communityUsersApi.useMutationCommunityReset()

const showConfirm = ref(false)
// async function resetColumn() {
// 	if (!props.columnType) return

// 	await communityResetMutation.executeMutation({
// 		type: props.columnType,
// 	})
// }
</script>

<template>
	<div v-if="column.getCanSort()" :class="cn('flex items-center space-x-2', $attrs.class ?? '')">
		<UiDropdownMenu>
			<UiDropdownMenuTrigger as-child>
				<UiButton
					variant="ghost"
					size="sm"
					class="-ml-3 h-8 data-[state=open]:bg-accent"
				>
					<span>{{ title }}</span>
					<Icon v-if="column.getIsSorted() === 'desc'" name="lucide:arrow-down" class="ml-2 h-4 w-4" />
					<Icon v-else-if=" column.getIsSorted() === 'asc'" name="lucide:arrow-up" class="ml-2 h-4 w-4" />
					<Icon v-else name="lucide:arrow-down-up" class="ml-2 h-4 w-4" />
				</UiButton>
			</UiDropdownMenuTrigger>
			<UiDropdownMenuContent align="start">
				<UiDropdownMenuItem @click="column.toggleSorting(false)">
					<Icon name="lucide:arrow-up" class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
					TASC
				</UiDropdownMenuItem>
				<UiDropdownMenuItem @click="column.toggleSorting(true)">
					<Icon name="lucide:arrow-up-down" class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
					TDESC
				</UiDropdownMenuItem>
				<DropdownMenuSeparator />
				<UiDropdownMenuItem @click="column.toggleVisibility(false)">
					<Icon name="lucide:eye-off" class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
					THIDE
				</UiDropdownMenuItem>
				<UiDropdownMenuItem v-if="columnType && !hideReset" @click="showConfirm = true">
					<Icon name="lucide:trash" class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
					Reset
				</UiDropdownMenuItem>
			</UiDropdownMenuContent>
		</UiDropdownMenu>
	</div>

	<div v-else :class="$attrs.class">
		{{ title }}
	</div>

	<!-- <ActionConfirm
		v-model:open="showConfirm"
		:confirm-text="t('community.users.reset.resetQuestion', { title })"
		@confirm="resetColumn"
	/> -->
</template>
