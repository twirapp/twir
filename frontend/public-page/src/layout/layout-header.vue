<script setup lang="ts">
import LayoutHeaderMenu from './layout-header-menu.vue';
import LayoutHeaderStreamerProfile from './layout-streamer-profile.vue';

import LayoutHeaderUserProfile from '@/layout/layout-sidebar-user-profile.vue';
</script>

<template>
	<header class="sticky z-40 top-0 bg-background/80 backdrop-blur-lg border-b border-border">
		<div
			class="container flex flex-col sm:flex-col md:flex-row flex-wrap justify-between min-h-14 items-center gap-4 py-4 md:py-0"
		>
			<layout-header-streamer-profile />
			<layout-header-menu />
			<layout-header-user-profile />
		</div>
	</header>
</template>
