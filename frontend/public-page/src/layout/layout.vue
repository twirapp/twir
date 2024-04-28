<script lang="ts" setup>
import { useBreakpoints, breakpointsTailwind } from '@vueuse/core';

import { useStreamerProfile } from '@/api/use-streamer-profile';
import { TooltipProvider } from '@/components/ui/tooltip';
import LayoutNavbar from '@/layout/layout-navbar.vue';
import LayoutSidebar from '@/layout/layout-sidebar.vue';
import LayoutStreamerProfile from '@/layout/layout-streamer-profile.vue';


useStreamerProfile();

const breakpoints = useBreakpoints(breakpointsTailwind);
const mdAndSmaller = breakpoints.smallerOrEqual('md');
</script>

<template>
	<layout-navbar v-if="mdAndSmaller" />
	<div class="flex flex-row gap-4">
		<aside
			v-if="!mdAndSmaller"
			class="flex min-w-[300px] w-[300px] flex-col justify-between border-e h-screen sticky top-0"
		>
			<layout-sidebar />
		</aside>

		<div
			class="flex flex-nowrap container flex-col gap-4 py-5"
			:class="{ 'mt-[3.5rem]': mdAndSmaller }"
		>
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
