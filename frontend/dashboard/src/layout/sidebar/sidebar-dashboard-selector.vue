<script lang="ts" setup>
import { useVirtualList } from '@vueuse/core'
import { Check, ChevronsUpDown } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { resolveUserName } from '../../helpers'

import type { PopoverContentProps } from 'radix-vue'

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
import { cn } from '@/lib/utils'

const { t } = useI18n()
const { open: sidebarOpen } = useSidebar()
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

const search = ref('')

function filterFunction(_items: any, searchTerm: string): string[] {
	if (!profile.value?.availableDashboards) return []

	search.value = searchTerm

	return profile.value.availableDashboards
		.filter((item) => {
			return item.twitchProfile.login.toLowerCase().includes(searchTerm.toLowerCase())
		})
		.map(item => item.id)
}

const popoverProps = computed((): PopoverContentProps & { class?: string } => {
	if (sidebarOpen.value) return { class: 'w-[--radix-popper-anchor-width]' }
	return {
		class: 'w-[300px]',
		alignOffset: -4,
		align: 'start',
		sideOffset: 12,
		side: 'right',
	}
})

const options = computed(() => {
	return profile.value?.availableDashboards.filter((item) => {
		return item.twitchProfile.login.toLowerCase().includes(search.value.toLowerCase())
	}) ?? []
})

const { list: virtualizedList, containerProps, wrapperProps } = useVirtualList(options, {
	itemHeight: 32,
})
</script>

<template>
	<SidebarMenu class="p-2">
		<SidebarMenuItem v-if="profile">
			<Popover v-model:open="open">
				<PopoverTrigger as-child>
					<SidebarMenuButton
						v-if="currentDashboard"
						size="lg"
						class="flex justify-start items-center data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
					>
						<div class="flex aspect-square size-8 items-center justify-center">
							<img :src="currentDashboard.twitchProfile.profileImageUrl" class="rounded-full" />
						</div>
						<div class="grid flex-1 text-left text-sm leading-tight">
							<span class="truncate font-semibold">{{ currentDashboard.twitchProfile.displayName }}</span>
							<span class="truncate text-xs">{{ t(`dashboard.header.managingUser`) }}</span>
						</div>
						<ChevronsUpDown class="ml-auto" />
					</SidebarMenuButton>
				</PopoverTrigger>
				<PopoverContent class="p-0" v-bind="popoverProps">
					<Command :filter-function="filterFunction">
						<CommandInput class="h-9" placeholder="Search user" />
						<CommandEmpty>
							No user found
						</CommandEmpty>
						<CommandList class="!max-h-full">
							<CommandGroup :heading="t(`dashboard.header.channelsAccess`)">
								<div v-bind="containerProps" class="max-h-72">
									<div v-bind="wrapperProps">
										<CommandItem
											v-for="option in virtualizedList"
											:key="option.index"
											style="height: 32px"
											:value="option.data.id"
											@select="(ev) => {
												if (typeof ev.detail.value === 'string') {
													selectDashboard(ev.detail.value)
												}
											}"
										>
											<div class="flex items-center gap-2">
												<Avatar class="size-4">
													<AvatarImage :src="option.data.twitchProfile.profileImageUrl" :alt="option.data.twitchProfile.displayName" />
													<AvatarFallback>
														{{ option.data.twitchProfile.displayName.slice(0, 2).toUpperCase() }}
													</AvatarFallback>
												</Avatar>
												<span>
													{{ resolveUserName(option.data.twitchProfile.login, option.data.twitchProfile.displayName) }}
												</span>
											</div>
											<Check
												:class="cn(
													'ml-auto h-4 w-4',
													profile.selectedDashboardId === option.data.id ? 'opacity-100' : 'opacity-0',
												)"
											/>
										</CommandItem>
									</div>
								</div>
							</CommandGroup>
						</CommandList>
					</Command>
				</PopoverContent>
			</Popover>
		</SidebarMenuItem>
	</SidebarMenu>
</template>
