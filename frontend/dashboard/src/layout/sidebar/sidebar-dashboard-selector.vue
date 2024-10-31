<script lang="ts" setup>
import { Check, ChevronsUpDown, Globe } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { resolveUserName } from '../../helpers'

import { useDashboard, useProfile } from '@/api'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover'
import {
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	useSidebar,
} from '@/components/ui/sidebar'
import { usePublicPageHref } from '@/layout/use-public-page-href'
import { cn } from '@/lib/utils'

const { t } = useI18n()
const publicPageHref = usePublicPageHref()
const { data: profile } = useProfile()
const { setDashboard } = useDashboard()

const open = ref(false)

const currentDashboard = computed(() => {
	const dashboard = profile.value?.availableDashboards.find(dashboard => dashboard.id === profile.value?.selectedDashboardId)
	if (!dashboard) return null

	return dashboard
})

function selectDashboard(id: string) {
	setDashboard(id)
	open.value = false
}

function filterFunction(_items: any, searchTerm: string): string[] {
	if (!profile.value?.availableDashboards) return []

	return profile.value.availableDashboards
		.filter((item) => {
			return item.twitchProfile.login.toLowerCase().includes(searchTerm.toLowerCase())
		})
		.map(item => item.id)
}

const { open: sidebarOpen } = useSidebar()
</script>

<template>
	<SidebarMenu>
		<SidebarMenuItem v-if="profile">
			<Popover v-model:open="open">
				<PopoverTrigger as-child>
					<SidebarMenuButton
						v-if="currentDashboard"
						size="lg"
						class="flex items-center justify-center data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
					>
						<Avatar :class="{ 'size-6': !sidebarOpen, 'size-8': sidebarOpen }">
							<AvatarImage :src="currentDashboard.twitchProfile.profileImageUrl" :alt="currentDashboard.twitchProfile.displayName" />
							<AvatarFallback>
								{{ currentDashboard.twitchProfile.displayName.slice(0, 2).toUpperCase() }}
							</AvatarFallback>
						</Avatar>
						<div class="grid flex-1 text-left text-sm leading-tight">
							<span class="truncate font-semibold">{{ currentDashboard.twitchProfile.displayName }}</span>
							<span class="truncate text-xs">{{ t(`dashboard.header.managingUser`) }}</span>
						</div>
						<ChevronsUpDown class="ml-auto" />
					</SidebarMenuButton>
				</PopoverTrigger>
				<PopoverContent class="w-full p-0">
					<Command :filter-function="filterFunction">
						<CommandInput class="h-9" placeholder="Search user" />
						<CommandEmpty>No user found</CommandEmpty>
						<CommandList>
							<CommandGroup :heading="t(`dashboard.header.channelsAccess`)">
								<CommandItem
									v-for="dashboard in profile.availableDashboards"
									:key="dashboard.id"
									:value="dashboard.id"
									@select="(ev) => {
										if (typeof ev.detail.value === 'string') {
											selectDashboard(ev.detail.value)
										}

									}"
								>
									<div class="flex items-center gap-2">
										<Avatar class="size-4">
											<AvatarImage :src="dashboard.twitchProfile.profileImageUrl" :alt="dashboard.twitchProfile.displayName" />
											<AvatarFallback>
												{{ dashboard.twitchProfile.displayName.slice(0, 2).toUpperCase() }}
											</AvatarFallback>
										</Avatar>
										<span>
											{{ resolveUserName(dashboard.twitchProfile.login, dashboard.twitchProfile.displayName) }}
										</span>
									</div>
									<Check
										:class="cn(
											'ml-auto h-4 w-4',
											profile.selectedDashboardId === dashboard.id ? 'opacity-100' : 'opacity-0',
										)"
									/>
								</CommandItem>
							</CommandGroup>
						</CommandList>
					</Command>
				</PopoverContent>
			</Popover>
		</SidebarMenuItem>
		<SidebarMenuItem v-if="publicPageHref">
			<SidebarMenuButton as-child>
				<a :href="publicPageHref" target="_blank">
					<Globe class="mr-0.5" />
					{{ t('navbar.publicPage') }}
				</a>
			</SidebarMenuButton>
		</SidebarMenuItem>
	</SidebarMenu>
</template>
