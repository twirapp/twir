<script setup lang="ts">
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';

const props = defineProps<{
	name: string
	aliases: string[]
}>();

const mappedAliases = props.aliases.map((alias) => alias).join(', ');
</script>

<template>
	<Tooltip v-if="aliases.length">
		<TooltipTrigger>
			<span>!{{ name }}</span>
		</TooltipTrigger>
		<TooltipContent>
			{{ mappedAliases }}
		</TooltipContent>
	</Tooltip>

	<span v-else>!{{ name }}</span>
</template>
