<script setup lang="ts">
import { useRouter } from 'vue-router';

import { channelRoutes } from '@/router';

const router = useRouter();

function goToRoute(name: string) {
	router.push({ name });
}
</script>

<template>
	<ul :class="$style.menu">
		<li
			v-for="route of Object.values(channelRoutes)"
			:key="route.name"
			:class="[
				$style.item,
				{
					[$style.itemActive]: router.currentRoute.value.name === route.name
				}
			]"
			@click="goToRoute(route.name as string)"
		>
			<component
				:is="route.meta?.icon"
				v-if="route.meta?.icon"
				:class="$style.itemIcon"
			/>
			<span>
				{{ route.name }}
			</span>
		</li>
	</ul>
</template>

<style module>
.menu {
	display: flex;
	flex-direction: column;
	flex-wrap: nowrap;
	gap: 0.5rem;
	justify-content: center;
	align-items: center;
	list-style: none;
	padding: 0 1rem;
}

.item {
	display: flex;
	align-items: center;
	gap: 0.5rem;
	width: 100%;
	text-align: left;
	padding: 0.5rem 1rem;
	border-radius: var(--radius);
	font-size: 0.9rem;
	cursor: pointer;
}

.item:hover {
	background-color: hsl(var(--accent));
}

.itemActive {
	background-color: hsl(var(--accent)/50);
}

.itemIcon {
	width: 1.3rem;
	height: 1.3rem;
	stroke: hsl(var(--accent-foreground));
	color: hsl(var(--accent-foreground));
	stroke-width: 2;
}

</style>
