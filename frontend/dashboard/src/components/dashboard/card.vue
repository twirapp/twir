<script lang="ts" setup>
import { IconEyeOff, IconGripVertical } from '@tabler/icons-vue'
import { NButton, NCard } from 'naive-ui'
import { useAttrs } from 'vue'
import { useI18n } from 'vue-i18n'

import { type WidgetItem, useWidgets } from './widgets.js'

import type { CSSProperties } from 'vue'

withDefaults(defineProps<{
	contentStyle?: CSSProperties
	popup?: boolean
}>(), {
	contentStyle: () => ({ padding: '0px' }),
})

defineSlots<{
	'default': any
	'action'?: any
	'header-extra'?: any
}>()

const widgets = useWidgets()

const attrs = useAttrs() as { item: WidgetItem, [x: string]: unknown } | undefined

function hideItem() {
	if (!attrs) return

	const item = widgets.value.find(item => item.i === attrs.item.i)
	if (!item) return
	item.visible = false
}

const { t } = useI18n()
</script>

<template>
	<NCard
		:segmented="{
			content: true,
			footer: 'soft',
		}"
		header-style="padding: 5px;"
		:content-style="contentStyle"
		style="width: 100%; height: 100%"
		v-bind="$attrs"
	>
		<template v-if="!popup" #header>
			<div class="widgets-draggable-handle flex items-center">
				<IconGripVertical class="w-5 h-5" />
				{{ t(`dashboard.widgets.${attrs?.item.i}.title`) }}
			</div>
		</template>

		<template v-if="!popup" #header-extra>
			<div class="flex gap-1">
				<slot name="header-extra" />
				<NButton text size="small" @click="hideItem">
					<IconEyeOff />
				</NButton>
			</div>
		</template>

		<slot />

		<template #action>
			<slot name="action" />
		</template>
	</NCard>
</template>
