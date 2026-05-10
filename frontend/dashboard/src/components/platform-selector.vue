<script setup lang="ts">
const platforms = defineModel<string[]>({ default: () => [] })

const options = [
	{
		id: 'twitch',
		label: 'Twitch',
		colorClass:
			'data-[active=true]:border-[#9146FF] data-[active=true]:bg-[#9146FF]/10 data-[active=true]:text-[#9146FF]',
	},
	{
		id: 'kick',
		label: 'Kick',
		colorClass:
			'data-[active=true]:border-[#53FC18] data-[active=true]:bg-[#53FC18]/10 data-[active=true]:text-[#53FC18]',
	},
] as const

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
			v-for="opt in options"
			:key="opt.id"
			type="button"
			:data-active="platforms.includes(opt.id)"
			@click="toggle(opt.id)"
			class="flex items-center gap-2 px-4 py-2 rounded-lg border border-border bg-card text-muted-foreground transition-all hover:bg-accent"
			:class="opt.colorClass"
		>
			<svg
				v-if="opt.id === 'twitch'"
				xmlns="http://www.w3.org/2000/svg"
				fill="currentColor"
				viewBox="0 0 24 24"
				class="size-4"
			>
				<path
					d="M1.3 4.6 2.8.8h19.9v14.5L17 21h-4.6l-3 2.9H6V21H1.3V4.6Zm15.8 12.6 3.7-3.8V2.7H5v14.5h3.7v3l2.8-3h5.6Z"
				/>
				<path d="M17.1 7h-1.8v5.5H17V7Zm-4.6 0h-1.9v5.5h1.9V7Z" />
			</svg>
			<svg
				v-else
				xmlns="http://www.w3.org/2000/svg"
				fill="currentColor"
				viewBox="0 0 24 24"
				class="size-4"
			>
				<path d="M3 5h3.5l5 6.5-5 6.5H3V5z" />
				<path d="M15 5h3v13h-3V5z" />
			</svg>
			<span class="text-sm font-medium">{{ opt.label }}</span>
		</button>
	</div>
</template>
