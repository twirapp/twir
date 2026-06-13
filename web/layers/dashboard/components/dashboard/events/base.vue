<script lang="ts" setup>
import { UseTimeAgo } from '@vueuse/components'
import { computed } from 'vue'

import type { SVGProps } from '@tabler/icons-vue'
import type { FunctionalComponent } from 'vue'

import { useTheme } from '@/composables/use-theme.js'

const props = defineProps<{
	icon: FunctionalComponent<SVGProps, Record<string, any>, any>
	iconColor?: [light: string, dark: string]
	createdAt: string
}>()

defineSlots<{
	leftContent: any
	rightContent: any
}>()

const locale = navigator.language
const theme = useTheme()
const color = computed(() => {
	if (!props.iconColor) return
	return props.iconColor.at(Number(theme.theme.value === 'dark'))
})
const date = computed(() => new Date(Number(props.createdAt)))
</script>

<template>
	<div
		class="flex min-h-12.5 gap-2.5 px-2.5 select-text border-b-border border-b border-solid py-2"
	>
		<div class="flex justify-between items-center w-full">
			<div class="flex gap-2.5 items-center">
				<component :is="icon" class="flex items-center min-h-9 min-w-9" :style="{ color }" />
				<div class="flex flex-col">
					<slot name="leftContent" />
				</div>
			</div>

			<div class="flex items-end text-xs h-full py-2 shrink-0">
				<UseTimeAgo v-slot="{ timeAgo }" :time="date" :update-interval="1000" show-second>
					<span class="tooltip" :data-utc="date.toLocaleString(locale)">
						{{ timeAgo }}
					</span>
				</UseTimeAgo>
			</div>
		</div>
	</div>
</template>

<style scoped>
.tooltip {
	position: relative;
}

.tooltip:after {
	content: attr(data-utc);
	position: absolute;
	top: 0px;
	right: 0;
	left: -100px;
	display: none;
	text-align: center;
	background-color: #000;
	border-radius: 4px;
	padding: 2px;
	white-space: nowrap;
}

.tooltip:hover:after {
	display: block;
}
</style>
