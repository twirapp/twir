<script setup lang="ts">
import { ArrowDown, ArrowDownUp, ArrowUp } from 'lucide-vue-next'


import type { Column } from '@tanstack/vue-table'



import { cn } from '~/lib/utils.js'

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
		<UiDropdownMenu>
			<UiDropdownMenuTrigger as-child>
				<UiButton
					variant="ghost"
					size="sm"
					class="-ml-3 h-8 data-[state=open]:bg-accent"
				>
					<span>{{ title }}</span>
					<ArrowDown v-if="column.getIsSorted() === 'desc'" class="ml-2 h-4 w-4" />
					<ArrowUp v-else-if=" column.getIsSorted() === 'asc'" class="ml-2 h-4 w-4" />
					<ArrowDownUp v-else class="ml-2 h-4 w-4" />
				</UiButton>
			</UiDropdownMenuTrigger>
			<UiDropdownMenuContent align="start">
				<UiDropdownMenuItem @click="column.toggleSorting(false)">
					<ArrowUp class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
					{{ t('sharedTexts.asc') }}
				</UiDropdownMenuItem>
				<UiDropdownMenuItem @click="column.toggleSorting(true)">
					<ArrowDown class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
					{{ t('sharedTexts.desc') }}
				</UiDropdownMenuItem>
			</UiDropdownMenuContent>
		</UiDropdownMenu>
	</div>
</template>
