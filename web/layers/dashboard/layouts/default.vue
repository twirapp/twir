<script setup lang="ts">
import 'vue-sonner/style.css'

const route = useRoute()
const isFullScreen = computed(() => route.meta?.fullScreen === true)
</script>

<template>
	<UiTooltipProvider :delay-duration="100">
		<template v-if="isFullScreen">
			<div class="w-full h-full">
				<slot />
			</div>
			<UiToaster />
		</template>
		<template v-else>
			<DashboardSidebar>
				<DashboardHeader />
				<div
					:style="{
						padding: route.meta?.noPadding ? undefined : '24px',
						height: '100%',
					}"
					class="bg-[#0b0b0c]"
				>
					<slot />
				</div>
				<UiToaster />
			</DashboardSidebar>
		</template>
	</UiTooltipProvider>
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
</style>
