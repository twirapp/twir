<script lang="ts" setup>
import PublicNavigation from './public/public-navigation.vue'

import { useStreamerProfile, useStreamerPublicSettings } from '~/layers/public/api/use-streamer-profile'

const { data } = await useStreamerProfile()
const profileId = computed(() => data.value?.twitchGetUserByName?.id ?? '')

const { data: publicSettings } = await useStreamerPublicSettings(profileId)
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
		</UiSidebar>

		<UiSidebarInset class="p-4 container">
			<UiCard style="background-color: rgb(24, 24, 28)">
				<UiCardContent class="p-6">
					<div class="flex flex-row flex-wrap justify-between w-full">
						<div class="flex gap-4 flex-row">
							<img :src="data?.twitchGetUserByName?.profileImageUrl" class="size-16 rounded-full" />
							<div class="flex flex-col gap-2">
								<span class="text-4xl">{{ data?.twitchGetUserByName?.displayName }}</span>
								<span class="text-sm text-muted-foreground"> {{ publicSettings?.userPublicSettings.description || data?.twitchGetUserByName?.description }}</span>
							</div>
						</div>
						<div class="flex flex-col gap-2">
							<a
								class="underline"
								:href="`https://twitch.tv/${data?.twitchGetUserByName?.login}`"
							>
								Twitch
							</a>
							<a
								v-for="link of publicSettings?.userPublicSettings.socialLinks"
								:key="link"
								class="underline"
								:href="link"
							>
								{{ link.title }}
							</a>
						</div>
					</div>
				</UiCardContent>
			</UiCard>

			<div class="mt-2">
				<slot />
			</div>
		</UiSidebarInset>
	</UiSidebarProvider>
</template>
