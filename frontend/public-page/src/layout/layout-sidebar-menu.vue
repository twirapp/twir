<script setup lang="ts">
import { useRouter } from 'vue-router';

import { channelRoutes } from '@/router';

const router = useRouter();
const emits = defineEmits<{
	navigate: []
}>();

function goToRoute(name: string) {
	router.push({ name });
	emits('navigate');
}
</script>

<template>
	<ul class="flex flex-col flex-nowrap gap-2 justify-center items-center [list-style:none] px-4 py-[0]">
		<li
			v-for="route of Object.values(channelRoutes)"
			:key="route.name"
			:class="[
				'flex items-center gap-2 w-full text-left px-4 py-2 rounded-[var(--radius)] text-[0.9rem] cursor-pointer',
				{ 'item-active': router.currentRoute.value.name === route.name }
			]"
			@click="goToRoute(route.name as string)"
		>
			<component
				:is="route.meta?.icon"
				v-if="route.meta?.icon"
				class="w-[1.3rem] h-[1.3rem] stroke-[2] item-icon"
			/>
			<span>
				{{ route.name }}
			</span>
		</li>
	</ul>
</template>

<style scoped>
.item:hover {
	background-color: hsl(var(--accent));
}

.item-active {
	background-color: hsl(var(--accent) / 50);
}

.item-icon {
	stroke: hsl(var(--accent-foreground));
	color: hsl(var(--accent-foreground));
}
</style>
