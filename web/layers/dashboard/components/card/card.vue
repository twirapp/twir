<script setup lang="ts">
import type { FunctionalComponent } from 'vue'

import { Card, CardContent, CardDescription, CardFooter, CardHeader } from '@/components/ui/card'

withDefaults(
	defineProps<{
		title: string
		description?: string
		icon?: FunctionalComponent
		iconStroke?: number
		withStroke?: boolean
		iconFill?: string
		iconWidth?: string
		iconHeight?: string
		isLoading?: boolean
	}>(),
	{
		withStroke: true,
		iconWidth: '48px',
		iconHeight: '48px',
	}
)

defineEmits<{
	openSettings: []
}>()

defineSlots<{
	content?: FunctionalComponent
	footer?: FunctionalComponent
	headerExtra?: FunctionalComponent
}>()


</script>

<template>
	<Card class="flex flex-col h-full">
		<CardHeader class="space-y-4">
			<div class="flex gap-2 items-center">
				<component
					:is="icon"
					v-if="icon"
					:style="{
						color: iconFill,
						fill: iconFill ? 'currentColor' : null,
						stroke: withStroke ? '#61e8bb' : null,
						strokeWidth: iconStroke,
						width: iconWidth,
						height: iconHeight,
					}"
				/>
				<h2 class="text-xl font-semibold text-foreground">
					{{ title }}
				</h2>
				<slot name="headerExtra" />
			</div>

			<CardDescription v-if="description">
				{{ description }}
			</CardDescription>
		</CardHeader>

		<CardContent class="text-muted-foreground">
			<slot name="content" />
		</CardContent>

		<CardFooter class="mt-auto">
			<div class="flex gap-2 flex-wrap w-full">
				<slot name="footer" />
			</div>
		</CardFooter>
	</Card>
</template>

<style scoped>
@reference '@/assets/index.css';

:deep(button span) {
	@apply text-sm;
}

@media (max-width: 568px) {
	:deep(button) {
		@apply w-full;
	}
}
</style>
