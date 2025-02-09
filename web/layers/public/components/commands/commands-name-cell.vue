<script setup lang="ts">
const props = defineProps<{
	name: string
	aliases: string[]
}>()

const mappedAliases = props.aliases.join(', ')
</script>

<template>
	<UiTooltip v-if="aliases.length">
		<UiTooltipTrigger>
			<span>{{ name }}</span>
		</UITooltipTrigger>
		<UiTooltipContent>
			{{ mappedAliases }}
		</UiTooltipContent>
	</UiTooltip>

	<span v-else>{{ name }}</span>
</template>
