<script lang="ts" setup>
import { EyeOff, GripVertical } from 'lucide-vue-next'
import { useAttrs } from 'vue'
import { useI18n } from 'vue-i18n'

import { type WidgetItem, useWidgets } from './widgets.js'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

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
	<Card v-bind="$attrs" class="p-0">
		<CardContent class="h-10 border-b">
			<div class="flex flex-row justify-between items-center h-10">
				<div class="widgets-draggable-handle flex items-center">
					<GripVertical class="w-5 h-5" />
					{{ t(`dashboard.widgets.${attrs?.item.i}.title`) }}
				</div>

				<div v-if="!popup" class="flex gap-2">
					<slot name="header-extra" />
					<Button size="sm" variant="outline" @click="hideItem">
						<EyeOff />
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
