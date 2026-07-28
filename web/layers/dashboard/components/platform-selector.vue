<script setup lang="ts">
import { Platform } from '~/gql/graphql.js'

const props = withDefaults(
	defineProps<{
		exclude?: string[]
	}>(),
	{
		exclude: () => [],
	},
)

const platforms = defineModel<string[]>({ default: () => [] })

const options = [
	{
		id: 'twitch',
		label: 'Twitch',
		icon: 'simple-icons:twitch',
		colorClass:
			'data-[active=true]:border-[#9146FF] data-[active=true]:bg-[#9146FF]/10 data-[active=true]:text-[#9146FF]',
	},
	{
		id: 'kick',
		label: 'Kick',
		icon: 'simple-icons:kick',
		colorClass:
			'data-[active=true]:border-[#53FC18] data-[active=true]:bg-[#53FC18]/10 data-[active=true]:text-[#53FC18]',
	},
	{
		id: Platform.VkVideoLive.toLowerCase(),
		label: 'VK Video Live',
		icon: 'simple-icons:vk',
		colorClass:
			'data-[active=true]:border-[#0077FF] data-[active=true]:bg-[#0077FF]/10 data-[active=true]:text-[#0077FF]',
	},
] as const

const visibleOptions = computed(() => options.filter((opt) => !props.exclude.includes(opt.id)))

function toggle(id: string) {
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
			:key="opt.id"
			type="button"
			:data-active="platforms.includes(opt.id)"
			@click="toggle(opt.id)"
			class="flex items-center gap-2 px-4 py-2 rounded-lg border border-border bg-card text-muted-foreground transition-all hover:bg-accent"
			:class="opt.colorClass"
		>
			<Icon :name="opt.icon" class="size-4" />
			<span class="text-sm font-medium">{{ opt.label }}</span>
		</button>
	</div>
</template>
