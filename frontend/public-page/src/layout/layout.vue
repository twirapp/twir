<script lang="ts" setup>
import { useBreakpoints, breakpointsTailwind } from '@vueuse/core';
import { computed, watch } from 'vue';
import { useRouter } from 'vue-router';

import { TooltipProvider } from '@/components/ui/tooltip';
import { useStreamerProfile } from '@/composables/use-streamer-profile';
import LayoutNavbar from '@/layout/layout-navbar.vue';
import LayoutSidebar from '@/layout/layout-sidebar.vue';
import LayoutStreamerProfile from '@/layout/layout-streamer-profile.vue';
import { routeNames } from '@/router';

const router = useRouter();
const streamerName = computed(() => {
	const name = router.currentRoute.value.params.channelName;
	if (typeof name !== 'string') {
		return null;
	}

	return name;
});

const { isError } = useStreamerProfile(streamerName);

watch([streamerName, isError], ([name, error]) => {
	if (!name || error) {
		router.push({ name: routeNames.notFound });
	}
});

const breakpoints = useBreakpoints(breakpointsTailwind);
const mdAndSmaller = breakpoints.smallerOrEqual('md');
</script>

<template>
	<layout-navbar v-if="mdAndSmaller" />
	<div class="flex flex-row gap-4">
		<aside
			v-if="!mdAndSmaller"
			class="flex w-auto flex-col justify-between border-e h-screen sticky top-0"
		>
			<layout-sidebar />
		</aside>

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
