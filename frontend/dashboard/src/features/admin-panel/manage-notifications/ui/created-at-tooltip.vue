<script setup lang="ts">
import { NTime } from 'naive-ui'

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
			<NTime class="cursor-pointer text-nowrap" :time="time" type="relative" />
		</TooltipTrigger>
		<TooltipContent side="top">
			{{ time.toLocaleString() }}
		</TooltipContent>
	</Tooltip>
</template>
