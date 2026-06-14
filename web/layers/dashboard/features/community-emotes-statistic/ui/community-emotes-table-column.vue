<script setup lang="ts">

import type { Column } from '@tanstack/vue-table'

import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import { cn } from '~~/layers/dashboard/lib/utils.js'

defineOptions({
	inheritAttrs: false
})

defineProps<{
	title: string
	column: Column<any, any>
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
					<Icon name="lucide:arrow-down" v-if="column.getIsSorted() === 'desc'" class="ml-2 h-4 w-4" />
					<Icon name="lucide:arrow-up" v-else-if=" column.getIsSorted() === 'asc'" class="ml-2 h-4 w-4" />
					<Icon name="lucide:arrow-down-up" v-else class="ml-2 h-4 w-4" />
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="start">
				<DropdownMenuItem @click="column.toggleSorting(false)">
					<Icon name="lucide:arrow-up" class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
					{{ t('sharedTexts.asc') }}
				</DropdownMenuItem>
				<DropdownMenuItem @click="column.toggleSorting(true)">
					<Icon name="lucide:arrow-down" class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
					{{ t('sharedTexts.desc') }}
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	</div>
</template>
