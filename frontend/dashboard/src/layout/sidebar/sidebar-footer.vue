<script setup lang="ts">
import { ExternalLink } from "lucide-vue-next";
import { computed } from "vue";
import { useI18n } from "vue-i18n";

import { usePublicPageHref } from "../use-public-page-href";

import DiscordLogo from "@/assets/integrations/discord.svg?use";
import GithubLogo from "@/assets/integrations/github.svg?use";
import Badge from "@/components/ui/badge/Badge.vue";
import {
	SidebarFooter,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	useSidebar,
} from "@/components/ui/sidebar";
import { footerNavigationItems } from "@/config/navigation";
import { useNotifications } from "@/composables/use-notifications";

const { t } = useI18n();
const { setOpenMobile } = useSidebar();
const publicPageHref = usePublicPageHref();
const { notificationsCounter } = useNotifications();

// Filter and prepare footer items
const visibleFooterItems = computed(() => {
	return footerNavigationItems
		.filter((item) => {
			// Filter out public page dependent items if no public page
			if (item.isPublicPageDependent && !publicPageHref.value) {
				return false;
			}
			return true;
		})
		.map((item) => {
			// Compute dynamic hrefs
			let href = item.href;
			if (item.isPublicPageDependent && item.translationKey === "sidebar.publicPage") {
				href = publicPageHref.value || "";
			} else if (item.href.startsWith("/") && item.isExternal) {
				href = `${window.location.origin}${item.href}`;
			}

			return {
				...item,
				href,
			};
		});
});

// Separate Discord and GitHub for special layout
const socialItems = computed(() =>
	visibleFooterItems.value.filter((item) => item.icon === "discord" || item.icon === "github"),
);

const regularItems = computed(() =>
	visibleFooterItems.value.filter((item) => item.icon !== "discord" && item.icon !== "github"),
);
</script>

<template>
	<SidebarFooter>
		<SidebarMenu>
			<!-- Discord and GitHub in special layout -->
			<div class="flex gap-2 group-data-[collapsible=icon]:flex-col">
				<template v-for="item in socialItems" :key="item.name">
					<SidebarMenuButton
						class="flex justify-center"
						variant="active"
						as-child
						:tooltip="item.name"
					>
						<a :href="item.href" target="_blank">
							<DiscordLogo v-if="item.icon === 'discord'" />
							<GithubLogo v-else-if="item.icon === 'github'" />
						</a>
					</SidebarMenuButton>
				</template>
			</div>

			<!-- Regular footer items -->
			<SidebarMenuItem v-for="item in regularItems" :key="item.href">
				<SidebarMenuButton
					as-child
					:tooltip="item.translationKey ? t(item.translationKey) : item.name"
					@click="!item.isExternal && setOpenMobile(false)"
				>
					<component
						:is="item.isExternal ? 'a' : 'RouterLink'"
						:href="item.isExternal ? item.href : undefined"
						:to="!item.isExternal ? item.href : undefined"
						:target="item.isExternal ? '_blank' : undefined"
					>
						<component :is="item.icon" />
						<span>{{ item.translationKey ? t(item.translationKey) : item.name }}</span>
						<Badge
							v-if="item.showNotificationsBadge && notificationsCounter.counter > 0"
							variant="success"
							class="ml-auto"
						>
							{{ notificationsCounter.counter }}
						</Badge>
						<ExternalLink v-else-if="item.isExternal" class="ml-auto" />
					</component>
				</SidebarMenuButton>
			</SidebarMenuItem>
		</SidebarMenu>
	</SidebarFooter>
</template>
