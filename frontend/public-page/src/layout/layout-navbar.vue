<script setup lang="ts">
import TwirLogo from '@twir/brand/src/logo.svg?url';
import {
	DialogRoot,
	DialogPortal,
	DialogOverlay,
	DialogContent,
	DialogTrigger,
} from 'radix-vue';
import { ref } from 'vue';

import LayoutSidebarMenu from './layout-sidebar-menu.vue';
import LayoutSidebarUserProfile from './layout-sidebar-user-profile.vue';

import BurgerToggle from '@/components/BurgerToggle.vue';
import { Separator } from '@/components/ui/separator';

const opened = ref(false);
</script>

<template>
	<header
		class="w-screen flex border-b border-border h-14 items-center justify-between px-6 z-20 absolute inset-0 pointer-events-auto"
	>
		<a class="flex h-16 w-16 gap-2 items-center" href="/">
			<img :src="TwirLogo" class="w-9 h-9" alt="twir-logo" />
			<span class="text-2xl font-semibold text-white" data-astro-cid-aygicng6="">Twir</span>
		</a>

		<DialogRoot v-model:open="opened" :modal="true">
			<DialogTrigger as-child>
				<BurgerToggle
					v-model="opened"
				/>
			</DialogTrigger>
			<DialogPortal>
				<DialogOverlay />
				<DialogContent
					class="fixed inset-y-0 left-0 z-10 h-[100dvh] w-screen gap-4 bg-zinc-900/90 px-6 pt-[3.5rem] backdrop-blur ease-out data-[state=closed]:duration-300 data-[state=open]:duration-300 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:slide-out-to-bottom-2 data-[state=open]:slide-in-from-bottom-2 md:px-10"
					style="outline: none"
					@interact-outside.prevent
					@focus-outside.prevent
					@pointer-down-outside.prevent
				>
					<layout-sidebar-menu class="mt-2" @navigate="opened = false" />
					<Separator class="mt-4 mb-4 px-[1rem]" />
					<layout-sidebar-user-profile class="my-4 px-[1rem]" />
				</DialogContent>
			</DialogPortal>
		</DialogRoot>
	</header>
</template>
