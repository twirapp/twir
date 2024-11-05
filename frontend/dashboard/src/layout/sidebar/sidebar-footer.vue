<script setup lang="ts">
import { DISCORD_INVITE_URL, GITHUB_REPOSITORY_URL } from '@twir/brand'
import { BellIcon, ExternalLink, Globe } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import SidebarProfile from './sidebar-profile.vue'
import { usePublicPageHref } from '../use-public-page-href'

import DiscordLogo from '@/assets/integrations/discord.svg?use'
import GithubLogo from '@/assets/integrations/github.svg?use'
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

			<div class="flex flex-row group-data-[collapsible=icon]:flex-col w-full">
				<SidebarMenuButton as-child tooltip="Discord">
					<a :href="DISCORD_INVITE_URL" target="_blank">
						<DiscordLogo />
						<span>Discord</span>
						<ExternalLink class="ml-auto" />
					</a>
				</SidebarMenuButton>

				<SidebarMenuButton as-child tooltip="GitHub">
					<a :href="GITHUB_REPOSITORY_URL" target="_blank">
						<GithubLogo />
						<span>GitHub</span>
						<ExternalLink class="ml-auto" />
					</a>
				</SidebarMenuButton>
			</div>

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
