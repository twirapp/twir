<script setup lang="ts">
import TwirLogo from '@twir/brand/src/logo.svg?url';

import { Separator } from '@/components/ui/separator';
import LayoutSidebarMenu from '@/layout/layout-sidebar-menu.vue';
import LayoutSidebarUserProfile from '@/layout/layout-sidebar-user-profile.vue';
</script>

<template>
	<div class="flex flex-col gap-x-2">
		<a class="flex h-16 w-16 gap-2 items-center justify-center w-full" href="/">
			<img :src="TwirLogo" class="w-10 h-10" alt="twir-logo" />
			<span class="text-2xl font-semibold text-white" data-astro-cid-aygicng6="">Twir</span>
		</a>

		<Separator />

		<layout-sidebar-menu class="mt-2" />
	</div>

	<div class="sticky inset-x-0 bottom-0">
		<Separator />
		<div class="p-2">
			<layout-sidebar-user-profile />
		</div>
	</div>
</template>

<style scoped>

</style>
