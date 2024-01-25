<script lang="ts" setup>
import { TooltipProvider } from '@/components/ui/tooltip';
import LayoutSidebar from '@/layout/layout-sidebar.vue';
import LayoutStreamerProfile from '@/layout/layout-streamer-profile.vue';
</script>

<template>
	<div class="flex gap-4">
		<layout-sidebar />

		<div class="flex flex-nowrap container flex-col gap-4 py-5">
			<layout-streamer-profile />

			<div>
				<TooltipProvider :delayDuration="0">
					<router-view v-slot="{ Component }">
						<transition appear mode="out-in">
							<component :is="Component" />
						</transition>
					</router-view>
				</TooltipProvider>
			</div>
		</div>
	</div>
	<!--	<LayoutHeader />-->
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
