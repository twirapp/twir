<script setup lang="ts">
import { Import } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import TwirLogo from '@/components/twir-logo.vue'
import {
	Sidebar,
	SidebarContent,
	SidebarFooter,
	SidebarHeader,
	SidebarInset,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarProvider,
	SidebarRail,
	SidebarSeparator,
	SidebarTrigger,
} from '@/components/ui/sidebar'
import SidebarDashboardSelector from '@/layout/sidebar/sidebar-dashboard-selector.vue'
import SidebarNavigation from '@/layout/sidebar/sidebar-navigation.vue'
import SidebarProfile from '@/layout/sidebar/sidebar-profile.vue'

const { t } = useI18n()
</script>

<template>
	<SidebarProvider>
		<Sidebar collapsible="icon">
			<SidebarHeader>
				<div class="flex justify-between group-data-[collapsible=icon]:justify-center">
					<a href="/" class="flex flex-row gap-2 items-center justify-center group-data-[collapsible=icon]:hidden ml-2">
						<TwirLogo class="size-6" />
						<h1 class="text-xl group-data-[collapsible=icon]:hidden">
							Twir
						</h1>
					</a>
					<SidebarTrigger class="size-[24px]" />
				</div>

				<SidebarSeparator />

				<SidebarDashboardSelector />
			</SidebarHeader>
			<SidebarSeparator />
			<SidebarContent>
				<SidebarNavigation />
			</SidebarContent>
			<SidebarSeparator />
			<SidebarFooter>
				<SidebarMenu>
					<SidebarMenuItem>
						<SidebarMenuButton as-child variant="outline" :tooltip="t('sidebar.import')">
							<RouterLink to="/dashboard/import">
								<Import />
								<span>{{ t('sidebar.import') }}</span>
							</RouterLink>
						</SidebarMenuButton>
					</SidebarMenuItem>
					<SidebarMenuItem class="mt-1">
						<SidebarProfile />
					</SidebarMenuItem>
				</SidebarMenu>
			</SidebarFooter>
			<SidebarRail />
		</Sidebar>
		<SidebarInset>
			<slot />
		</SidebarInset>
	</SidebarProvider>
</template>
