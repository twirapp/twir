<script setup lang="ts">
import "vue-sonner/style.css";
import { computed, ref } from "vue";
import { type RouteLocationNormalized, useRoute } from "vue-router";

import { Toaster } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import Sidebar from "~~/layers/dashboard/layout/sidebar/sidebar.vue";
import Stats from "~~/layers/dashboard/layout/header/header.vue";

const isRouterReady = ref(false);
const route = useRoute();

if (import.meta.client) {
	isRouterReady.value = true;
}

const isFullScreen = computed(() => route.meta?.fullScreen === true);

interface HistoryState {
	noTransition?: boolean;
}

function getTransition(route: RouteLocationNormalized) {
	if (typeof window === 'undefined') return undefined;
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
			<slot />
			<Toaster />
		</template>
		<Sidebar v-else>
			<Stats />
			<div
				:style="{
					padding: route.meta?.noPadding ? undefined : '24px',
					height: '100%',
				}"
				class="bg-[#0b0b0c]"
			>
				<slot />
			</div>

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
