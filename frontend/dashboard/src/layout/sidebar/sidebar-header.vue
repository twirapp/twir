<script setup lang="ts">
import TwirLogo from '@/components/twir-logo.vue'
import { SidebarHeader, SidebarSeparator, SidebarTrigger, useSidebar } from '@/components/ui/sidebar'

const { isMobile } = useSidebar()

const isDev = import.meta.env.DEV
</script>

<template>
	<SidebarHeader>
		<div v-if="!isMobile" class="flex items-center justify-between group-data-[collapsible=icon]:justify-center">
			<a href="/" class="flex flex-row gap-2 items-center justify-center group-data-[collapsible=icon]:hidden ml-2">
				<TwirLogo class="size-8" />
				<h1 class="text-2xl font-semibold group-data-[collapsible=icon]:hidden text-accent-foreground">
					Twir {{ isDev ? 'dev' : '' }}
				</h1>
			</a>
			<SidebarTrigger />
		</div>
	</SidebarHeader>

	<SidebarSeparator v-if="!isMobile" />
</template>
