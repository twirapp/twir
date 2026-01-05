<script setup lang="ts">
import { UseTimeAgo } from '@vueuse/components'

import {
	Tooltip,
	TooltipContent,
	TooltipTrigger
} from '@/components/ui/tooltip'

defineProps<{ time: Date }>()
</script>

<template>
	<Tooltip>
		<TooltipTrigger as-child>
			<UseTimeAgo v-slot="{ timeAgo }" :time="time">
				<span class="cursor-pointer text-nowrap">{{ timeAgo }}</span>
			</UseTimeAgo>
		</TooltipTrigger>
		<TooltipContent side="top">
			{{ time.toLocaleString() }}
		</TooltipContent>
	</Tooltip>
</template>
