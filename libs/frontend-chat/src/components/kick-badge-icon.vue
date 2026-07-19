<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
	type: string
	text?: string
}>()

type BadgeIcon = 'crown' | 'shield' | 'star' | 'gem' | 'verified' | 'award' | 'gift' | 'wrench'

interface BadgeMeta {
	color: string
	icon: BadgeIcon
}

const badgeMap: Record<string, BadgeMeta> = {
	broadcaster: { color: '#53FC18', icon: 'crown' },
	moderator: { color: '#00E701', icon: 'shield' },
	subscriber: { color: '#FFD700', icon: 'star' },
	founder: { color: '#FFD700', icon: 'star' },
	vip: { color: '#FF4F9A', icon: 'gem' },
	verified: { color: '#1D9BF0', icon: 'verified' },
	og: { color: '#B14AED', icon: 'award' },
	sub_gifter: { color: '#FF9F1C', icon: 'gift' },
	staff: { color: '#FF4444', icon: 'wrench' },
}

const meta = computed<BadgeMeta | undefined>(() => badgeMap[props.type])
</script>

<template>
	<svg
		v-if="meta"
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 24 24"
		width="1em"
		height="1em"
		fill="none"
		:stroke="meta.color"
		stroke-width="2"
		stroke-linecap="round"
		stroke-linejoin="round"
		class="kick-badge"
	>
		<g v-if="meta.icon === 'crown'">
			<path
				d="M11.562 3.266a.5.5 0 0 1 .876 0L15.39 8.87a1 1 0 0 0 1.516.294L21.183 5.5a.5.5 0 0 1 .798.519l-2.834 10.246a1 1 0 0 1-.956.735H5.81a1 1 0 0 1-.957-.735L2.02 6.02a.5.5 0 0 1 .798-.519l4.276 3.664a1 1 0 0 0 1.516-.294z"
			/>
			<path d="M5 21h14" />
		</g>
		<g v-else-if="meta.icon === 'shield'">
			<path
				d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z"
			/>
		</g>
		<g v-else-if="meta.icon === 'star'">
			<path
				d="M11.525 2.295a.53.53 0 0 1 .95 0l2.31 4.679a2.123 2.123 0 0 0 1.595 1.16l5.166.756a.53.53 0 0 1 .294.904l-3.736 3.638a2.123 2.123 0 0 0-.611 1.878l.882 5.14a.53.53 0 0 1-.771.56l-4.618-2.428a2.122 2.122 0 0 0-1.973 0L6.396 21.01a.53.53 0 0 1-.77-.56l.881-5.139a2.122 2.122 0 0 0-.611-1.879L2.16 9.795a.53.53 0 0 1 .294-.906l5.165-.755a2.122 2.122 0 0 0 1.597-1.16z"
			/>
		</g>
		<g v-else-if="meta.icon === 'gem'">
			<path d="M6 3h12l4 6-10 13L2 9Z" />
			<path d="M11 3 8 9l4 13 4-13-3-6" />
			<path d="M2 9h20" />
		</g>
		<g v-else-if="meta.icon === 'verified'">
			<path
				d="M3.85 8.62a4 4 0 0 1 4.78-4.77 4 4 0 0 1 6.74 0 4 4 0 0 1 4.78 4.78 4 4 0 0 1 0 6.74 4 4 0 0 1-4.77 4.78 4 4 0 0 1-6.75 0 4 4 0 0 1-4.78-4.77 4 4 0 0 1 0-6.76Z"
			/>
			<path d="m9 12 2 2 4-4" />
		</g>
		<g v-else-if="meta.icon === 'award'">
			<circle cx="12" cy="8" r="6" />
			<path d="M15.477 12.89 17 22l-5-3-5 3 1.523-9.11" />
		</g>
		<g v-else-if="meta.icon === 'gift'">
			<rect x="3" y="8" width="18" height="4" rx="1" />
			<path d="M12 8v13" />
			<path d="M19 12v7a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2v-7" />
			<path d="M7.5 8a2.5 2.5 0 0 1 0-5A4.8 8 0 0 1 12 8a4.8 8 0 0 1 4.5-5 2.5 2.5 0 0 1 0 5" />
		</g>
		<g v-else-if="meta.icon === 'wrench'">
			<path
				d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"
			/>
		</g>
	</svg>
	<span v-else class="text-badge">{{ text || type }}</span>
</template>

<style scoped>
.kick-badge {
	height: 1em;
	width: 1em;
	display: inline-block;
	vertical-align: middle;
}

.text-badge {
	padding: 0 4px;
	font-size: 0.6em;
	background-color: #6d6767;
	border-radius: 4px;
	text-transform: uppercase;
}
</style>
