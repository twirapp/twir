<script setup lang="ts">
import type { DropdownMenuContentProps } from 'radix-vue'

import { UserStoreKey } from '~/stores/user'

const userStore = useAuth()

await Promise.all([
	callOnce(UserStoreKey, () => userStore.getUserData()),
])

const dropdownProps = computed((): DropdownMenuContentProps & { class?: string } => {
	return {
		class: 'w-[200px]',
		side: 'bottom',
		align: 'end',
		sideOffset: 4,
	}
})
</script>

<template>
	<button
		v-if="!userStore.user"
		class="flex flex-row px-4 py-2 items-center gap-2 bg-[#5D58F5] text-white rounded-lg font-medium focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-[#5D58F5]/50 cursor-pointer hover:bg-[#6964FF] transition-shadow"
		@click="() => userStore.login()"
	>
		Login
		<SvgoSocialTwitch :fontControlled="false" class="w-5 h-5 fill-white" />
	</button>

	<UiDropdownMenu v-else-if="userStore.user">
		<UiDropdownMenuTrigger class="inline-flex items-center gap-3 text-white/75 hover:text-white transition-colors" as="button">
			<div class="flex items-center gap-3 min-w-0">
				<img
					:src="userStore.user.twitchProfile.profileImageUrl"
					:alt="userStore.user.twitchProfile.displayName"
					class="w-8 h-8 rounded-full flex-shrink-0"
				/>
				<span class="max-[600px]:hidden truncate">
					{{ userStore.user.twitchProfile.login }}
				</span>
				<Icon name="lucide:chevron-down" class="w-4 h-4 flex-shrink-0" />
			</div>
		</UiDropdownMenuTrigger>

		<UiDropdownMenuContent v-bind="dropdownProps">
			<UiDropdownMenuItem as-child>
				<a href="/dashboard" class="flex w-full items-center">
					<Icon name="lucide:layout-dashboard" class="mr-2 h-4 w-4" />
					Dashboard
				</a>
			</UiDropdownMenuItem>

			<UiDropdownMenuSeparator />

			<UiDropdownMenuItem as="button" class="flex w-full items-center text-destructive" @click="userStore.logout">
				<Icon name="lucide:log-out" class="mr-2 h-4 w-4" />
				Logout
			</UiDropdownMenuItem>
		</UiDropdownMenuContent>
	</UiDropdownMenu>
</template>

<style scoped>
.text-destructive {
	color: rgb(239 68 68);
}
</style>
