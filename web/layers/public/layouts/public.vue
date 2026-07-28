<script lang="ts" setup>
import PublicNavigation from './public/public-navigation.vue'
import PublicUserProfile from './public/public-user-profile.vue'

import { useStreamerProfile } from '~~/layers/public/api/use-streamer-profile'

const streamerProfile = useStreamerProfile()
await useAsyncData('streamerProfile', () => streamerProfile.fetchProfile().then(() => true))

const platformUrls: Record<string, { url: (login: string) => string; label: string }> = {
	twitch: { url: (login) => `https://twitch.tv/${login}`, label: 'Twitch' },
	kick: { url: (login) => `https://kick.com/${login}`, label: 'Kick' },
	vk_video_live: { url: (login) => `https://live.vkvideo.ru/${login}`, label: 'VK Video Live' },
}

const profile = computed(() => {
	const account = streamerProfile.profile?.channelBySlug?.profile
	if (!account) return null

	const platform = platformUrls[account.platform] ?? {
		url: () => '',
		label: account.platform,
	}

	return {
		avatar: account.platformAvatar ?? undefined,
		displayName: account.platformDisplayName,
		login: account.platformLogin,
		url: platform.url(account.platformLogin),
		platform: platform.label,
	}
})
</script>

<template>
	<UiSidebarProvider>
		<UiSidebar collapsible="icon" variant="inset">
			<UiSidebarHeader>
				<div class="flex items-center justify-between group-data-[collapsible=icon]:justify-center">
					<NuxtLink
						to="/"
						class="flex flex-row gap-2 items-center justify-center group-data-[collapsible=icon]:hidden ml-2"
					>
						<TwirLogo class="w-8 h-8" />
						<h1
							class="text-2xl font-semibold group-data-[collapsible=icon]:hidden text-accent-foreground"
						>
							Twir
						</h1>
					</NuxtLink>
				</div>
			</UiSidebarHeader>

			<UiSidebarContent>
				<PublicNavigation />
			</UiSidebarContent>

			<UiSidebarFooter>
				<PublicUserProfile />
			</UiSidebarFooter>
		</UiSidebar>

		<UiSidebarInset class="p-4">
			<UiCard>
				<UiCardContent class="p-6">
					<div class="flex flex-row flex-wrap justify-between w-full gap-4">
						<div class="flex gap-4 flex-row flex-1">
								<img
									:src="profile?.avatar ?? undefined"
									class="w-16 h-16 rounded-full"
									:alt="`${profile?.login ?? ''}-avatar`"
								/>
								<div class="flex flex-col gap-2">
									<span class="text-4xl">{{ profile?.displayName }}</span>
								<span class="text-sm text-muted-foreground break-all">
									{{ streamerProfile.publicProfile?.userPublicSettings.description }}
								</span>
								</div>
							</div>
							<div class="flex flex-col gap-2 flex-none">
								<a
									v-if="profile"
									class="underline"
									:href="profile.url"
								>
									{{ profile.platform }}
								</a>
							<a
								v-for="(link, idx) of streamerProfile.publicProfile?.userPublicSettings.socialLinks"
								:key="idx"
								class="underline"
								:href="link.href"
							>
								{{ link.title }}
							</a>
						</div>
					</div>
				</UiCardContent>
			</UiCard>

			<div class="mt-4">
				<slot />
			</div>
		</UiSidebarInset>
	</UiSidebarProvider>
</template>
