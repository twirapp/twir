<script setup lang="ts">
import { IconMenu } from '@tabler/icons-vue';
import TwirLogo from '@twir/brand/src/logo.svg?url';
import { ref } from 'vue';

import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import {
	Sheet,
	SheetContent,
} from '@/components/ui/sheet';
import LayoutSidebarMenu from '@/layout/layout-sidebar-menu.vue';
import LayoutSidebarUserProfile from '@/layout/layout-sidebar-user-profile.vue';

const sheetOpened = ref(false);
</script>

<template>
	<header class="w-full flex border-b border-border h-14 items-center px-2 justify-between">
		<a class="flex h-16 w-16 gap-2 items-center" href="/">
			<img :src="TwirLogo" class="w-10 h-10" alt="twir-logo" />
			<span class="text-2xl font-semibold text-white" data-astro-cid-aygicng6="">Twir</span>
		</a>

		<Button size="icon" variant="ghost" @click="sheetOpened = true">
			<IconMenu />
		</Button>
	</header>

	<Sheet v-model:open="sheetOpened">
		<SheetContent class="sm:w-[100%] w-[100%]">
			<layout-sidebar-menu class="mt-2" />
			<Separator class="my-4" />
			<layout-sidebar-user-profile />
		</SheetContent>
	</Sheet>
</template>

<style scoped>

</style>
