<script lang="ts" setup>
import { useLocalStorage } from '@vueuse/core'
import { computed } from 'vue'
import { useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth'
import { baseNavigationItems } from '~~/layers/dashboard/config/navigation'

import Badge from '@/components/ui/badge/Badge.vue'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import {
	SidebarGroup,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarMenuSub,
	SidebarMenuSubItem,
	useSidebar,
} from '@/components/ui/sidebar'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

const { t } = useI18n()
const localePath = useLocalePath()
const currentRoute = useRoute()
const sidebar = useSidebar()

const canViewIntegrations = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewIntegrations)
const canViewEvents = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewEvents)
const canViewOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewOverlays)
const canViewSongRequests = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewSongRequests)
const canViewCommands = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewCommands)
const canViewTimers = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewTimers)
const canViewKeywords = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewKeywords)
const canViewVariables = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewVariables)
const canViewGreetings = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewGreetings)
const canViewAlerts = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewAlerts)
const canViewGames = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewGames)
const canViewModeration = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewModeration)
const canViewModules = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewModules)
const canViewGiveaways = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewGiveaways)

const twirSidebarOpenedStates = useLocalStorage<Record<string, boolean>>(
	'twir-sidebar-opened-states',
	{
		commands: false,
		community: false,
	}
)

// Map permissions to paths
const permissionMap: Record<string, boolean> = {
	'/dashboard/integrations': canViewIntegrations.value,
	'/dashboard/events/chat-alerts': canViewEvents.value,
	'/dashboard/events': canViewEvents.value,
	'/dashboard/overlays': canViewOverlays.value,
	'/dashboard/song-requests': canViewSongRequests.value,
	'/dashboard/commands': canViewCommands.value,
	'/dashboard/timers': canViewTimers.value,
	'/dashboard/keywords': canViewKeywords.value,
	'/dashboard/variables': canViewVariables.value,
	'/dashboard/greetings': canViewGreetings.value,
	'/dashboard/alerts': canViewAlerts.value,
	'/dashboard/games': canViewGames.value,
	'/dashboard/moderation': canViewModeration.value,
	'/dashboard/modules': canViewModules.value,
	'/dashboard/giveaways': canViewGiveaways.value,
}

const links = computed(() => {
	return baseNavigationItems.map((item) => {
		const hasPermission = item.path ? (permissionMap[item.path] ?? true) : true

		return {
			name: item.translationKey ? t(item.translationKey) : item.name || '',
			icon: item.icon,
			disabled: !hasPermission,
			path: item.path,
			isNew: item.isNew,
			openStateKey: item.openStateKey,
			child: item.child?.map((c) => ({
				name: c.translationKey ? t(c.translationKey) : c.name || '',
				icon: c.icon,
				path: c.path,
				isNew: c.isNew,
			})),
		}
	})
})

function goToRoute() {
	if (sidebar.isMobile.value) {
		sidebar.setOpenMobile(false)
	}
}
</script>

<template>
	<SidebarGroup>
		<SidebarMenu>
			<SidebarMenuItem
				v-for="item in links"
				:key="item.name"
			>
				<SidebarMenuButton
					v-if="!item.child"
					as-child
					:tooltip="item.name"
					:variant="currentRoute.path === localePath(item.path) ? 'active' : 'default'"
					@click="goToRoute"
				>
					<RouterLink :to="localePath(item.path!)">
						<Icon :name="item.icon" />
						<span>{{ item.name }}</span>
						<Badge
							v-if="item.isNew"
							class="rounded-md px-1 py-0.5 text-[10px] uppercase"
						>
							New
						</Badge>
					</RouterLink>
				</SidebarMenuButton>
				<Collapsible
					v-else-if="item.openStateKey"
					v-model:open="twirSidebarOpenedStates[item.openStateKey]"
					as-child
					class="group/collapsible"
				>
					<SidebarMenuItem>
						<CollapsibleTrigger as-child>
							<SidebarMenuButton
								:tooltip="item.name"
								:variant="
									item.path && currentRoute.path.startsWith(localePath(item.path)) ? 'active' : 'default'
								"
							>
								<Icon :name="item.icon" />
								<span>{{ item.name }}</span>
								<Badge
									v-if="item.isNew"
									class="rounded-md px-1 py-0.5 text-[10px] uppercase"
								>
									New
								</Badge>
								<Icon
									name="lucide:chevron-right"
									class="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"
								/>
							</SidebarMenuButton>
						</CollapsibleTrigger>
						<CollapsibleContent>
							<SidebarMenuSub>
								<SidebarMenuSubItem
									v-for="child in item.child"
									:key="child.name"
								>
									<SidebarMenuButton
										as-child
									:variant="
										currentRoute.path === localePath(child.path) || currentRoute.fullPath === localePath(child.path)
											? 'active'
											: 'default'
									"
										@click="goToRoute"
									>
										<RouterLink :to="localePath(child.path!)">
											<Icon :name="child.icon" />
											<span>{{ child.name }}</span>
											<Badge
												v-if="'isNew' in child && child.isNew"
												class="rounded-md px-1 py-0.5 text-[10px] uppercase"
											>
												New
											</Badge>
										</RouterLink>
									</SidebarMenuButton>
								</SidebarMenuSubItem>
							</SidebarMenuSub>
						</CollapsibleContent>
					</SidebarMenuItem>
				</Collapsible>
			</SidebarMenuItem>
		</SidebarMenu>
	</SidebarGroup>
</template>
