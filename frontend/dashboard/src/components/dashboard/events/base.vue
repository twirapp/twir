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

const theme = useTheme()
const color = computed(() => {
	if (!props.iconColor) return
	return props.iconColor.at(Number(theme.theme.value === 'dark'))
})
</script>

<template>
	<div class="flex min-h-[50px] gap-2.5 px-2.5 select-text border-b-[color:var(--n-border-color)] border-b border-solid">
		<div class="flex justify-between items-center w-full">
			<div class="flex gap-2.5 items-center">
				<component
					:is="icon"
					class="flex items-center min-h-9 min-w-9"
					:style="{ color }"
				/>
				<div class="flex flex-col">
					<slot name="leftContent" />
				</div>
			</div>

			<div class="flex items-end text-xs h-full py-2 flex-shrink-0">
				<UseTimeAgo
					v-slot="{ timeAgo }"
					:time="new Date(Number(createdAt))"
					:update-interval="1000"
					show-second
				>
					{{ timeAgo }}
				</UseTimeAgo>
			</div>
		</div>
	</div>
</template>
