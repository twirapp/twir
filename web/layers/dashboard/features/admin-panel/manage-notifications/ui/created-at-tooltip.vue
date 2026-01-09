<script setup lang="ts">
import { UseTimeAgo } from '@vueuse/components'



defineProps<{ time: Date }>()
</script>

<template>
	<UiTooltip>
		<UiTooltipTrigger as-child>
			<UseTimeAgo v-slot="{ timeAgo }" :time="time">
				<span class="cursor-pointer text-nowrap">{{ timeAgo }}</span>
			</UseTimeAgo>
		</UiTooltipTrigger>
		<UiTooltipContent side="top">
			{{ time.toLocaleString() }}
		</UiTooltipContent>
	</UiTooltip>
</template>
