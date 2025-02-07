<script setup lang="ts">
defineProps<{
	cooldown: number
	cooldownType: string
}>()
</script>

<template>
	<div class="flex items-center justify-center">
		<UiTooltip>
			<UiTooltipTrigger>
				<div class="flex gap-1 items-center">
					<span>{{ cooldown }}</span>
					<Icon :name="cooldownType === 'GLOBAL' ? 'lucide:globe' : 'lucide:user'" />
				</div>
			</UiTooltipTrigger>
			<UiTooltipContent>
				{{ cooldownType }}
			</UiTooltipContent>
		</UiTooltip>
	</div>
</template>

<style module>
.cooldownIcon {
	height: 18px;
	width: 18px;
}
</style>
