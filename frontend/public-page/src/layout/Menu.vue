<script setup lang="ts">
import { useRouter } from 'vue-router';

const router = useRouter();

const routes = [
	{
		name: 'Commands',
	},
	{
		name: 'Song requests',
	},
];

</script>

<template>
	<div class="flex items-center h-full">
		<button
			v-for="r of routes"
			:key="r.name"
			class="p-2 rounded text-slate-200 m-1 shadow-lg"
			:class="{
				'bg-neutral-600': router.currentRoute.value.name === r.name,
				'bg-neutral-700': router.currentRoute.value.name !== r.name
			}"
			@click="() => {
				console.log(r.name)
				router.push({ name: r.name })
			}"
		>
			{{ r.name }}
		</button>
	</div>
</template>
