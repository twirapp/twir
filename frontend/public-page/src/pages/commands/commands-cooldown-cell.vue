<script setup lang="ts">
import { IconUser, IconWorld } from '@tabler/icons-vue';

import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';

defineProps<{
	cooldown: bigint,
	cooldownType: string,
}>();
</script>

<template>
	<div class="flex items-center justify-center">
		<Tooltip>
			<TooltipTrigger>
				<div class="flex gap-0.5 items-center">
					<span>{{ cooldown }}</span>
					<IconWorld v-if="cooldownType === 'GLOBAL'" :class="$style.cooldownIcon" />
					<IconUser v-else :class="$style.cooldownIcon" />
				</div>
			</TooltipTrigger>
			<TooltipContent>
				{{ cooldownType }}
			</TooltipContent>
		</Tooltip>
	</div>
</template>

<style module>
.cooldownIcon {
	height: 18px;
	width: 18px;
}
</style>
