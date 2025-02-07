<script lang="ts" setup>
import PublicNavigation from './public/public-navigation.vue'
import PublicUserProfile from './public/public-user-profile.vue'

import { useStreamerProfile } from '~/layers/public/api/use-streamer-profile'

const streamerProfile = useStreamerProfile()
await useAsyncData('streamerProfile', () => streamerProfile.fetchProfile().then(() => true))
</script>

<template>
	<UiSidebarProvider>
		<UiSidebar collapsible="icon">
			<UiSidebarHeader>
				<div class="flex items-center justify-between group-data-[collapsible=icon]:justify-center">
					<a href="/" class="flex flex-row gap-2 items-center justify-center group-data-[collapsible=icon]:hidden ml-2">
						<TwirLogo class="size-8" />
						<h1 class="text-2xl font-semibold group-data-[collapsible=icon]:hidden text-accent-foreground">
							Twir
						</h1>
					</a>
				</div>
			</UiSidebarHeader>

			<UiSidebarSeparator />

			<UiSidebarContent>
				<PublicNavigation />
			</UiSidebarContent>

			<UiSidebarSeparator />

			<UiSidebarFooter>
				<div class="min-h-12">
					<ClientOnly>
						<template #fallback>
							<div class="h-full w-full flex items-center justify-center">
								Loading...
							</div>
						</template>
						<PublicUserProfile />
					</ClientOnly>
				</div>
			</UiSidebarFooter>
		</UiSidebar>

		<UiSidebarInset class="p-4 container">
			<UiCard style="background-color: rgb(24, 24, 28)">
				<UiCardContent class="p-6">
					<div class="flex flex-row flex-wrap justify-between w-full gap-4">
						<div class="flex gap-4 flex-row flex-1">
							<img :src="streamerProfile.profile?.twitchGetUserByName?.profileImageUrl" class="size-16 rounded-full" />
							<div class="flex flex-col gap-2">
								<span class="text-4xl">{{ streamerProfile.profile?.twitchGetUserByName?.displayName }}</span>
								<span class="text-sm text-muted-foreground break-all">
									{{ streamerProfile.publicProfile?.userPublicSettings.description || streamerProfile.profile?.twitchGetUserByName?.description }}
								</span>
							</div>
						</div>
						<div class="flex flex-col gap-2 flex-none">
							<a
								class="underline"
								:href="`https://twitch.tv/${streamerProfile.profile?.twitchGetUserByName?.login}`"
							>
								Twitch
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

<style>
html,
body {
	@apply bg-background-main text-foreground;
}
</style>
