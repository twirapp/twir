<script lang="ts" setup>
import { useVirtualList } from '@vueuse/core'
import { ChevronsUpDown } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useFilter } from 'reka-ui'

import type { PopoverContentProps } from 'reka-ui'

import { useDashboard, useProfile } from '@/api/auth'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import {
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	useSidebar,
} from '@/components/ui/sidebar'
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

function getDashboardName(dashboard: NonNullable<typeof profile.value>['availableDashboards'][number]) {
	if (dashboard.platform === 'kick') {
		return dashboard.kickProfile?.displayName ?? dashboard.kickProfile?.slug ?? ''
	}
	return dashboard.twitchProfile?.displayName ?? ''
}

function getDashboardLogin(dashboard: NonNullable<typeof profile.value>['availableDashboards'][number]) {
	if (dashboard.platform === 'kick') {
		return dashboard.kickProfile?.slug ?? ''
	}
	return dashboard.twitchProfile?.login ?? ''
}

function getDashboardAvatar(dashboard: NonNullable<typeof profile.value>['availableDashboards'][number]) {
	if (dashboard.platform === 'kick') {
		return dashboard.kickProfile?.profilePicture ?? ''
	}
	return dashboard.twitchProfile?.profileImageUrl ?? ''
}

const options = computed(() => {
	return (
		profile.value?.availableDashboards.filter(
			(p) =>
				contains(getDashboardLogin(p), search.value) ||
				contains(getDashboardName(p), search.value)
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
						<img v-if="getDashboardAvatar(currentDashboard)" :src="getDashboardAvatar(currentDashboard)" class="rounded-full" />
					</div>
					<div class="grid flex-1 text-left text-sm leading-tight">
						<span class="truncate font-semibold">{{ getDashboardName(currentDashboard) }}</span>
						<span class="truncate text-xs flex items-center gap-1">
							{{ t(`dashboard.header.managingUser`) }}
							<Badge v-if="currentDashboard.platform === 'kick'" variant="outline" class="uppercase text-[10px] px-1 py-0 h-4">K</Badge>
							<Badge v-else variant="outline" class="uppercase text-[10px] px-1 py-0 h-4">T</Badge>
						</span>
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
										:src="getDashboardAvatar(option.data)"
										:alt="getDashboardName(option.data)"
									/>
									<AvatarFallback>
										{{ getDashboardName(option.data).slice(0, 2).toUpperCase() }}
									</AvatarFallback>
								</Avatar>
								<span class="truncate">{{ getDashboardLogin(option.data) }}</span>
								<Badge
									v-if="option.data.platform === 'kick'"
									variant="outline"
									class="uppercase text-[10px] px-1 py-0 h-4 ml-auto"
								>K</Badge>
								<Badge
									v-else
									variant="outline"
									class="uppercase text-[10px] px-1 py-0 h-4 ml-auto"
								>T</Badge>
							</Button>
						</div>
					</div>
				</PopoverContent>
			</Popover>
		</SidebarMenuItem>
	</SidebarMenu>
</template>
