<script setup lang="ts">
import { ArrowDown, ArrowDownUp, ArrowUp, EyeOff } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import type { CommunityUsersResetType } from '@/gql/graphql.js'
import type { Column } from '@tanstack/vue-table'

import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { cn } from '@/lib/utils.js'

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
	<div v-if="column.getCanSort()" :class="cn('flex items-center space-x-2', $attrs.class ?? '')">
		<DropdownMenu>
			<DropdownMenuTrigger as-child>
				<Button
					variant="ghost"
					size="sm"
					class="-ml-3 h-8 data-[state=open]:bg-accent"
				>
					<span>{{ title }}</span>
					<ArrowDown v-if="column.getIsSorted() === 'desc'" class="ml-2 h-4 w-4" />
					<ArrowUp v-else-if=" column.getIsSorted() === 'asc'" class="ml-2 h-4 w-4" />
					<ArrowDownUp v-else class="ml-2 h-4 w-4" />
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="start">
				<DropdownMenuItem @click="column.toggleSorting(false)">
					<ArrowUp class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
					{{ t('sharedTexts.asc') }}
				</DropdownMenuItem>
				<DropdownMenuItem @click="column.toggleSorting(true)">
					<ArrowDown class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
					{{ t('sharedTexts.desc') }}
				</DropdownMenuItem>
				<DropdownMenuSeparator />
				<DropdownMenuItem @click="column.toggleVisibility(false)">
					<EyeOff class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
					{{ t('sharedTexts.hide') }}
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	</div>

	<div v-else :class="$attrs.class">
		{{ title }}
	</div>
</template>
