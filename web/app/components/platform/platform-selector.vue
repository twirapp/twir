<script setup lang="ts">
import type { Platform } from '~/gql/graphql.js'

import { PLATFORM_OPTIONS } from '~/utils/platforms.js'

const props = withDefaults(
	defineProps<{
		exclude?: Platform[]
	}>(),
	{
		exclude: () => [],
	},
)

const platforms = defineModel<Platform[]>({ default: () => [] })

const visibleOptions = computed(() =>
	PLATFORM_OPTIONS.filter((opt) => !props.exclude.includes(opt.platform)),
)

function toggle(id: Platform) {
	const current = new Set(platforms.value)
	if (current.has(id)) {
		current.delete(id)
	} else {
		current.add(id)
	}
	platforms.value = Array.from(current)
}
</script>

<template>
	<div class="flex gap-2">
		<button
			v-for="opt in visibleOptions"
			:key="opt.platform"
			type="button"
			:data-active="platforms.includes(opt.platform)"
			class="flex items-center gap-2 px-4 py-2 rounded-lg border border-border bg-card text-muted-foreground transition-all hover:bg-accent"
			:class="opt.activeClass"
			@click="toggle(opt.platform)"
		>
			<Icon :name="opt.icon" class="size-4" />
			<span class="text-sm font-medium">{{ opt.label }}</span>
		</button>
	</div>
</template>
