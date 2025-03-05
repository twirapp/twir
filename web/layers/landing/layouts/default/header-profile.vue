<script setup lang="ts">
import type { DropdownMenuContentProps } from 'radix-vue'

import { useAuthLink, useLogout, useProfile } from '~~/layers/landing/api/user'

const pageUrl = useRequestURL()

const redirectUrl = computed(() => {
	return `${pageUrl.origin}/dashboard`
})

const [
	{ data: profile },
	{ data: authLinkData },
] = await Promise.all([
	useProfile(),
	useAuthLink(redirectUrl),
])

const dropdownProps = computed((): DropdownMenuContentProps & { class?: string } => {
	return {
		class: 'w-[200px]',
		side: 'bottom',
		align: 'end',
		sideOffset: 4,
	}
})

const logout = useLogout()
</script>

<template>
	<a
		v-if="!profile"
		class="flex flex-row px-4 py-2 items-center gap-2 bg-[#5D58F5] text-white rounded-lg font-medium focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-[#5D58F5]/50 cursor-pointer hover:bg-[#6964FF] transition-shadow"
		:href="authLinkData?.authLink"
	>
		Login
		<SvgoSocialTwitch :fontControlled="false" class="w-5 h-5 fill-white" />
	</a>

	<UiDropdownMenu v-else>
		<UiDropdownMenuTrigger class="inline-flex items-center gap-3 text-white/75 hover:text-white transition-colors" as="button">
			<div class="flex items-center gap-3 min-w-0">
				<img
					:src="profile.authenticatedUser.twitchProfile.profileImageUrl"
					:alt="profile.authenticatedUser.twitchProfile.displayName"
					class="w-8 h-8 rounded-full flex-shrink-0"
				/>
				<span class="max-[600px]:hidden truncate">
					{{ profile.authenticatedUser.twitchProfile.login }}
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

			<UiDropdownMenuItem as="button" class="flex w-full items-center text-destructive" @click="logout">
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
