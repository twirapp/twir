<script lang="ts" setup>
import LayoutHeader from './layout-header.vue';

import { TooltipProvider } from '@/components/ui/tooltip';
</script>

<template>
	<LayoutHeader />

	<div class="container py-5">
		<TooltipProvider :delayDuration="0">
			<router-view v-slot="{ Component }">
				<transition appear mode="out-in">
					<component :is="Component" />
				</transition>
			</router-view>
		</TooltipProvider>
	</div>
</template>

<script setup lang="ts">
</script>
<style scoped>
.v-enter-active,
.v-leave-active {
	transition: opacity 0.2s ease;
}

.v-enter-from,
.v-leave-to {
	opacity: 0;
}
</style>
