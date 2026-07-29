<script setup lang="ts">
import { UserStoreKey } from '~/stores/user'

const userStore = useAuth()

const currentPlatform = computed(() => {
	return userStore.userWithoutDashboards?.currentPlatform.toLowerCase() ?? ''
})

await Promise.all([callOnce(UserStoreKey, () => userStore.getUserDataWithoutDashboards())])
</script>

<template>
	<div v-if="!userStore.userWithoutDashboards" class="flex flex-row items-center gap-2">
		<LoginDropdown />
	</div>

	<UiDropdownMenu v-else-if="userStore.userWithoutDashboards">
		<UiDropdownMenuTrigger
			class="inline-flex items-center gap-3 text-white/75 hover:text-white transition-colors"
			as="button"
		>
			<div class="flex items-center gap-3 min-w-0">
			<img
				:src="userStore.currentAccount?.platformAvatar ?? ''"
				:alt="userStore.currentAccount?.platformDisplayName ?? ''"
				class="w-8 h-8 rounded-full shrink-0"
			/>
			<span class="max-[600px]:hidden truncate">
				{{ userStore.currentAccount?.platformDisplayName ?? userStore.currentAccount?.platformLogin ?? '' }}
			</span>
				<Icon v-if="currentPlatform === 'kick'" name="simple-icons:kick" class="w-4 h-4 text-[#53FC18] shrink-0" />
				<Icon v-else-if="currentPlatform === 'vk_video_live'" name="simple-icons:vk" class="w-4 h-4 text-[#0077FF] shrink-0" />
			<Icon v-else-if="currentPlatform === 'youtube'" name="simple-icons:youtube" class="w-4 h-4 text-[#FF0000] shrink-0" />
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
