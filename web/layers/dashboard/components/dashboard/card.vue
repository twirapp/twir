<script lang="ts" setup>
import { useAttrs } from 'vue'

import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'

import { type WidgetItem, useWidgets } from './widgets.js'

defineProps<{
	popup?: boolean
}>()

defineSlots<{
	default: any
	action?: any
	'header-extra'?: any
}>()

const widgets = useWidgets()

const attrs = useAttrs() as { item: WidgetItem; [x: string]: unknown } | undefined

function hideItem() {
	if (!attrs) return

	const item = widgets.value.find((item) => item.i === attrs.item.i)
	if (!item) return
	item.visible = false
}

const { t } = useI18n()
</script>

<template>
	<Card
		v-bind="$attrs"
		class="gap-0 p-0"
	>
		<CardContent class="mt-0 h-10 border-b px-2">
			<div class="flex h-10 flex-row items-center justify-between">
				<div class="widgets-draggable-handle flex items-center">
					<Icon
						name="lucide:grip-vertical"
						class="h-5 w-5"
					/>
					<template v-if="attrs?.item.displayName">
						{{ attrs.item.displayName }}
					</template>
					<template v-else>
						{{ t(`dashboard.widgets.${attrs?.item.i}.title`) }}
					</template>
				</div>

				<div
					v-if="!popup"
					class="flex gap-2"
				>
					<slot name="header-extra" />
					<Button
						size="sm"
						variant="ghost"
						@click="hideItem"
						class="cursor-pointer"
					>
						<Icon name="lucide:eye-off" />
					</Button>
				</div>
			</div>
		</CardContent>

		<slot />

		<template #action>
			<slot name="action" />
		</template>
	</Card>
</template>
