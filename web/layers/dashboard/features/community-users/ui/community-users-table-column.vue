<script setup lang="ts">
import type { Column } from '@tanstack/vue-table'

import { cn } from '~~/layers/dashboard/lib/utils.js'

import type { CommunityUsersResetType } from '~/gql/graphql.js'

import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

defineOptions({
	inheritAttrs: false,
})

defineProps<{
	columnType?: CommunityUsersResetType
	column: Column<any, any>
	title: string
}>()

const { t } = useI18n()
</script>

<template>
	<div
		v-if="column.getCanSort()"
		:class="cn('flex items-center space-x-2', $attrs.class ?? '')"
	>
		<DropdownMenu>
			<DropdownMenuTrigger as-child>
				<Button
					variant="ghost"
					size="sm"
					class="data-[state=open]:bg-accent -ml-3 h-8"
				>
					<span>{{ title }}</span>
					<Icon
						name="lucide:arrow-down"
						v-if="column.getIsSorted() === 'desc'"
						class="ml-2 h-4 w-4"
					/>
					<Icon
						name="lucide:arrow-up"
						v-else-if="column.getIsSorted() === 'asc'"
						class="ml-2 h-4 w-4"
					/>
					<Icon
						name="lucide:arrow-down-up"
						v-else
						class="ml-2 h-4 w-4"
					/>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="start">
				<DropdownMenuItem @click="column.toggleSorting(false)">
					<Icon
						name="lucide:arrow-up"
						class="text-muted-foreground/70 mr-2 h-3.5 w-3.5"
					/>
					{{ t('sharedTexts.asc') }}
				</DropdownMenuItem>
				<DropdownMenuItem @click="column.toggleSorting(true)">
					<Icon
						name="lucide:arrow-down"
						class="text-muted-foreground/70 mr-2 h-3.5 w-3.5"
					/>
					{{ t('sharedTexts.desc') }}
				</DropdownMenuItem>
				<DropdownMenuSeparator />
				<DropdownMenuItem @click="column.toggleVisibility(false)">
					<EyeOff class="text-muted-foreground/70 mr-2 h-3.5 w-3.5" />
					{{ t('sharedTexts.hide') }}
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	</div>

	<div
		v-else
		:class="$attrs.class"
	>
		{{ title }}
	</div>
</template>
