<script lang="ts" setup>
import { useVirtualList } from '@vueuse/core'
import { Check, ChevronsUpDown } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useFilter } from 'reka-ui'

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
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

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
	search.value = ''
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

const search = ref('')

const { contains } = useFilter({ sensitivity: 'base' })
const options = computed(() => {
	return (
		profile.value?.availableDashboards.filter((p) =>
			contains(p.twitchProfile.login, search.value)
		) ?? []
	)
})

const {
	list: virtualizedList,
	containerProps,
	wrapperProps,
} = useVirtualList(options, {
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
							<span class="truncate font-semibold">{{
								currentDashboard.twitchProfile.displayName
							}}</span>
							<span class="truncate text-xs">{{ t(`dashboard.header.managingUser`) }}</span>
						</div>
						<ChevronsUpDown class="ml-auto" />
					</SidebarMenuButton>
				</PopoverTrigger>
				<PopoverContent class="p-0 min-h-20!" v-bind="popoverProps">
					<div class="p-2">
						<Input v-model="search" placeholder="Search user..." />
					</div>

					<div
						v-bind="containerProps"
						class="max-h-72 w-full px-2 overflow-y-auto [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-zinc-600 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-800"
					>
						<div v-bind="wrapperProps" class="w-full">
							<Button
								v-for="option in virtualizedList"
								:key="option.data.id"
								style="height: 32px"
								class="flex justify-start w-full"
								variant="ghost"
								@click="selectDashboard(option.data.id)"
							>
								<Avatar class="size-4">
									<AvatarImage
										:src="option.data.twitchProfile.profileImageUrl"
										:alt="option.data.twitchProfile.displayName"
									/>
									<AvatarFallback>
										{{ option.data.twitchProfile.displayName.slice(0, 2).toUpperCase() }}
									</AvatarFallback>
								</Avatar>
								{{ option.data.twitchProfile.login }}
							</Button>
						</div>
					</div>
				</PopoverContent>
			</Popover>
		</SidebarMenuItem>
	</SidebarMenu>
</template>
