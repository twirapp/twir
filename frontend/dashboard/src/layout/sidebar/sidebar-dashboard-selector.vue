<script lang="ts" setup>
import { Check, ChevronsUpDown } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ListboxVirtualizer } from 'reka-ui'

import { resolveUserName } from '../../helpers'

import type { PopoverContentProps } from 'reka-ui'

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
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
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
	const dashboard = profile.value?.availableDashboards.find(
		(dashboard) => dashboard.id === profile.value?.selectedDashboardId
	)
	if (!dashboard) return null

	return dashboard
})

function selectDashboard(id: string) {
	setDashboard(id)
	open.value = false
}

const popoverProps = computed((): PopoverContentProps & { class?: string } => {
	if (sidebarOpen.value) return { class: 'w-(--radix-popper-anchor-width)' }
	return {
		class: 'w-[300px]',
		alignOffset: -4,
		align: 'start',
		sideOffset: 12,
		side: 'right',
	}
})

const options = computed(() => {
	return profile.value?.availableDashboards ?? []
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
							<span class="truncate font-semibold">{{
								currentDashboard.twitchProfile.displayName
							}}</span>
							<span class="truncate text-xs">{{ t(`dashboard.header.managingUser`) }}</span>
						</div>
						<ChevronsUpDown class="ml-auto" />
					</SidebarMenuButton>
				</PopoverTrigger>
				<PopoverContent class="p-0 min-h-20!" v-bind="popoverProps">
					<Command>
						<CommandInput class="h-9" placeholder="Search user" />
						<CommandEmpty> No user found </CommandEmpty>
						<CommandList class="max-h-full!">
							<ListboxVirtualizer
								v-slot="{ option }"
								:options="options"
								:text-content="(opt) => opt.twitchProfile.login"
								:estimateSize="32"
								class="min-h-28!"
							>
								<CommandGroup :heading="t(`dashboard.header.channelsAccess`)" class="w-full">
									<CommandItem
										style="height: 32px"
										:value="option.id"
										@select="
											(ev) => {
												if (typeof ev.detail.value === 'string') {
													selectDashboard(ev.detail.value)
												}
											}
										"
										class="cursor-pointer"
										:data-highligted="profile.selectedDashboardId === option.id"
									>
										<div class="flex items-center gap-2">
											<Avatar class="size-4">
												<AvatarImage
													:src="option.twitchProfile.profileImageUrl"
													:alt="option.twitchProfile.displayName"
												/>
												<AvatarFallback>
													{{ option.twitchProfile.displayName.slice(0, 2).toUpperCase() }}
												</AvatarFallback>
											</Avatar>
											<span>
												{{
													resolveUserName(
														option.twitchProfile.login,
														option.twitchProfile.displayName
													)
												}}
											</span>
										</div>
										<Check
											:class="
												cn(
													'ml-auto h-4 w-4',
													profile.selectedDashboardId === option.id ? 'opacity-100' : 'opacity-0'
												)
											"
										/>
									</CommandItem>
								</CommandGroup>
							</ListboxVirtualizer>
						</CommandList>
					</Command>
				</PopoverContent>
			</Popover>
		</SidebarMenuItem>
	</SidebarMenu>
</template>
