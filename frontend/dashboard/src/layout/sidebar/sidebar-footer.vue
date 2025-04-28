<script setup lang="ts">
import { DISCORD_INVITE_URL, GITHUB_REPOSITORY_URL } from '@twir/brand'
import { BellIcon, ExternalLink, Globe } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import SidebarProfile from './sidebar-profile.vue'
import { usePublicPageHref } from '../use-public-page-href'

// import GithubLogo from '@/assets/integrations/github.svg?use'
import Badge from '@/components/ui/badge/Badge.vue'
import { SidebarFooter, SidebarMenu, SidebarMenuButton, SidebarMenuItem, useSidebar } from '@/components/ui/sidebar'
import { useNotifications } from '@/composables/use-notifications'

const { t } = useI18n()
const { setOpenMobile } = useSidebar()
const publicPageHref = usePublicPageHref()
const { notificationsCounter } = useNotifications()
</script>

<template>
	<SidebarFooter>
		<SidebarMenu>
			<div class="flex gap-2 group-data-[collapsible=icon]:flex-col">
				<SidebarMenuButton class="flex justify-center" variant="active" as-child tooltip="Discord">
					<a :href="DISCORD_INVITE_URL" target="_blank">
						<!--						<DiscordLogo /> -->
					</a>
				</SidebarMenuButton>

				<SidebarMenuButton class="flex justify-center" variant="active" as-child tooltip="GitHub">
					<a :href="GITHUB_REPOSITORY_URL" target="_blank">
						<!--						<GithubLogo /> -->
					</a>
				</SidebarMenuButton>
			</div>

			<SidebarMenuItem>
				<SidebarMenuButton
					as-child
					:tooltip="t('sidebar.notifications')"
					@click="setOpenMobile(false)"
				>
					<RouterLink to="/dashboard/notifications">
						<BellIcon />
						<span>{{ t('sidebar.notifications') }}</span>
						<Badge
							v-if="notificationsCounter.counter > 0"
							variant="success"
							class="ml-auto"
						>
							{{ notificationsCounter.counter }}
						</Badge>
					</RouterLink>
				</SidebarMenuButton>
			</SidebarMenuItem>

			<SidebarMenuItem v-if="publicPageHref">
				<SidebarMenuButton as-child :tooltip="t('sidebar.publicPage')">
					<a :href="publicPageHref" target="_blank">
						<Globe />
						<span>{{ t('sidebar.publicPage') }}</span>
						<ExternalLink class="ml-auto" />
					</a>
				</SidebarMenuButton>
			</SidebarMenuItem>

			<SidebarMenuItem class="mt-1">
				<SidebarProfile />
			</SidebarMenuItem>
		</SidebarMenu>
	</SidebarFooter>
</template>
