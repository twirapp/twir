<script setup lang="ts">
import { DISCORD_INVITE_URL, GITHUB_REPOSITORY_URL } from '@twir/brand'
import { BellIcon, ClipboardPenLine, ExternalLink, Globe, LinkIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { usePublicPageHref } from '../use-public-page-href'

import DiscordLogo from '@/assets/integrations/discord.svg?use'
import GithubLogo from '@/assets/integrations/github.svg?use'
import Badge from '@/components/ui/badge/Badge.vue'
import {
	SidebarFooter,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	useSidebar,
} from '@/components/ui/sidebar'
import { useNotifications } from '@/composables/use-notifications'

const { t } = useI18n()
const { setOpenMobile } = useSidebar()
const publicPageHref = usePublicPageHref()
const { notificationsCounter } = useNotifications()

const hastebinLink = computed(() => {
	return `${window.location.origin}/h`
})

const urlShortenerLink = computed(() => {
	return `${window.location.origin}/url-shortener`
})
</script>

<template>
	<SidebarFooter>
		<SidebarMenu>
			<div class="flex gap-2 group-data-[collapsible=icon]:flex-col">
				<SidebarMenuButton class="flex justify-center" variant="active" as-child tooltip="Discord">
					<a :href="DISCORD_INVITE_URL" target="_blank">
						<DiscordLogo />
					</a>
				</SidebarMenuButton>

				<SidebarMenuButton class="flex justify-center" variant="active" as-child tooltip="GitHub">
					<a :href="GITHUB_REPOSITORY_URL" target="_blank">
						<GithubLogo />
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
						<Badge v-if="notificationsCounter.counter > 0" variant="success" class="ml-auto">
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
				<SidebarMenuButton as-child tooltip="Hastebin">
					<a :href="urlShortenerLink" target="_blank">
						<LinkIcon />
						<span>URL Shortener</span>
						<ExternalLink class="ml-auto" />
					</a>
				</SidebarMenuButton>
				<SidebarMenuButton as-child tooltip="Hastebin">
					<a :href="hastebinLink" target="_blank">
						<ClipboardPenLine />
						<span>Hastebin</span>
						<ExternalLink class="ml-auto" />
					</a>
				</SidebarMenuButton>
			</SidebarMenuItem>
		</SidebarMenu>
	</SidebarFooter>
</template>
