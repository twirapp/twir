<script setup lang="ts">
import { UserStoreKey } from '~/stores/user'
import KickIcon from '~~/layers/landing/components/kick-icon.vue'

const userStore = useAuth()

const isKickUser = computed(() => {
	return userStore.userWithoutDashboards?.currentPlatform === 'KICK'
})

await Promise.all([callOnce(UserStoreKey, () => userStore.getUserDataWithoutDashboards())])
</script>

<template>
	<div v-if="!userStore.userWithoutDashboards" class="flex flex-row items-center gap-2">
		<button
			class="flex flex-row px-4 py-2 items-center gap-2 bg-[#5D58F5] text-white rounded-lg font-medium focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-[#5D58F5]/50 cursor-pointer hover:bg-[#6964FF] transition-shadow"
			@click="() => userStore.login()"
		>
			Twitch
			<SvgoSocialTwitch :fontControlled="false" class="w-5 h-5 fill-white" />
		</button>
		<button
			class="flex flex-row px-4 py-2 items-center gap-2 bg-[#27272a] text-white rounded-lg font-medium focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-[#53FC18]/50 cursor-pointer hover:bg-[#27272a]/80 transition-shadow"
			@click="() => userStore.loginWithKick()"
		>
			Kick
			<KickIcon class="text-[#53FC18]" />
		</button>
	</div>

	<UiDropdownMenu v-else-if="userStore.userWithoutDashboards">
		<UiDropdownMenuTrigger
			class="inline-flex items-center gap-3 text-white/75 hover:text-white transition-colors"
			as="button"
		>
			<div class="flex items-center gap-3 min-w-0">
			<img
				:src="userStore.userWithoutDashboards.twitchProfile?.profileImageUrl ?? userStore.userWithoutDashboards.kickProfile?.profilePicture ?? ''"
				:alt="userStore.userWithoutDashboards.twitchProfile?.displayName ?? userStore.userWithoutDashboards.kickProfile?.displayName ?? ''"
				class="w-8 h-8 rounded-full shrink-0"
			/>
			<span class="max-[600px]:hidden truncate">
				{{ userStore.userWithoutDashboards?.twitchProfile?.login ?? userStore.userWithoutDashboards?.kickProfile?.displayName ?? '' }}
			</span>
				<KickIcon v-if="isKickUser" class="w-4 h-4 text-[#53FC18] shrink-0" />
				<Icon name="lucide:chevron-down" class="w-4 h-4 shrink-0" />
			</div>
		</UiDropdownMenuTrigger>

		<UiDropdownMenuContent align="end" side="bottom" :side-offset="4" class="w-50">
			<UiDropdownMenuItem as-child>
				<a href="/dashboard" class="flex w-full items-center">
					<Icon name="lucide:layout-dashboard" class="mr-2 h-4 w-4" />
					Dashboard
				</a>
			</UiDropdownMenuItem>

			<UiDropdownMenuSeparator />

			<UiDropdownMenuItem
				as="button"
				class="flex w-full items-center text-red-500"
				@click="userStore.logout"
			>
				<Icon name="lucide:log-out" class="mr-2 h-4 w-4" />
				Logout
			</UiDropdownMenuItem>
		</UiDropdownMenuContent>
	</UiDropdownMenu>
</template>
