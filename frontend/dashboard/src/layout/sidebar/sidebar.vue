<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'

import SidebarDashboardSelector from './sidebar-dashboard-selector.vue'
import SidebarFooter from './sidebar-footer.vue'
import SidebarHeader from './sidebar-header.vue'

import {
	Sidebar,
	SidebarContent,
	SidebarInset,
	SidebarProvider,
	SidebarSeparator,
} from '@/components/ui/sidebar'
import SidebarNavigation from '@/layout/sidebar/sidebar-navigation.vue'

const open = useLocalStorage('twirSidebarIsOpen', true)
</script>

<template>
	<SidebarProvider v-model:open="open">
		<Sidebar collapsible="icon">
			<SidebarHeader />

			<SidebarDashboardSelector />

			<SidebarSeparator />

			<SidebarContent>
				<SidebarNavigation />
			</SidebarContent>

			<SidebarSeparator />

			<SidebarFooter />
		</Sidebar>
		<SidebarInset>
			<slot />
		</SidebarInset>
	</SidebarProvider>
</template>
