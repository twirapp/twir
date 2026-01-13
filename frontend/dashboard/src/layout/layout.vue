<script setup lang="ts">
import "vue-sonner/style.css";
import { computed, ref } from "vue";
import { type RouteLocationNormalized, RouterView, useRoute, useRouter } from "vue-router";

import { Toaster } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import Sidebar from "@/layout/sidebar/sidebar.vue";
import Stats from "@/layout/header/header.vue";

const isRouterReady = ref(false);
const router = useRouter();
const route = useRoute();

router.isReady().finally(() => (isRouterReady.value = true));

const isFullScreen = computed(() => route.meta?.fullScreen === true);

interface HistoryState {
	noTransition?: boolean;
}

function getTransition(route: RouteLocationNormalized) {
	const state = window.history.state as HistoryState;
	if (state.noTransition) {
		return undefined;
	}

	return route.meta.transition || "router";
}
</script>

<template>
	<TooltipProvider :delay-duration="100">
		<template v-if="isFullScreen">
			<RouterView v-slot="{ Component, route }">
				<transition :name="getTransition(route)" mode="out-in">
					<div :key="route.path" class="w-full h-full">
						<component :is="Component" />
					</div>
				</transition>
			</RouterView>
			<Toaster />
		</template>
		<Sidebar v-else>
			<Stats />
			<RouterView v-slot="{ Component, route }">
				<transition :name="getTransition(route)" mode="out-in">
					<div
						:key="route.path"
						:style="{
							padding: route.meta?.noPadding ? undefined : '24px',
							height: '100%',
						}"
						class="bg-[#0b0b0c]"
					>
						<component :is="Component" />
					</div>
				</transition>
			</RouterView>

			<Toaster />
		</Sidebar>
	</TooltipProvider>
</template>

<style>
.router-enter-active,
.router-leave-active {
	transition: all 0.2s cubic-bezier(0, 0, 0.2, 1);
}

.router-enter-from,
.router-leave-to {
	opacity: 0;
	transform: scale(0.98);
}

.app-loader {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100%;
}
</style>
