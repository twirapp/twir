<script setup lang="ts">
import { useRouter } from 'vue-router';

import { routesNames } from '@/routes.js';

const router = useRouter();
</script>

<template>
	<div>
		<div class="flex flex-wrap items-end h-full gap-3">
			<button
				v-for="routeName of Object.values(routesNames)"
				:key="routeName"
				class="p-2 rounded-md text-slate-200 shadow-lg"
				:class="{
					'bg-neutral-600': router.currentRoute.value.name === routeName,
					'bg-neutral-700': router.currentRoute.value.name !== routeName
				}"
				@click="router.push({ name: routeName })"
			>
				{{ routeName }}
			</button>
		</div>
	</div>
</template>
