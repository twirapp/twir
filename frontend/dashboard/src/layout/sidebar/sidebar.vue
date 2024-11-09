<script setup lang="ts">
import SidebarNavigation from '@/layout/sidebar/sidebar-navigation.vue'

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
</script>

<template>
	<SidebarProvider>
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
