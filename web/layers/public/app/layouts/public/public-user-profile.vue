<script setup lang="ts">
import { computed } from 'vue'

import type { DropdownMenuContentProps } from 'radix-vue'

import { UserStoreKey } from '~/stores/user'

const userStore = useAuth()

await Promise.all([callOnce(UserStoreKey, () => userStore.getUserDataWithoutDashboards())])

const dropdownProps = computed((): DropdownMenuContentProps & { class?: string } => {
	return {
		class: 'w-(--reka-dropdown-menu-trigger-width)',
		side: 'bottom',
		align: 'end',
		sideOffset: 4,
	}
})
</script>

<template>
	<div v-if="!userStore.userWithoutDashboards" class="flex flex-col gap-2">
		<UiSidebarMenuButton
			size="lg"
			class="items-center justify-center bg-[#5D58F5] text-white hover:bg-[#6964FF] hover:text-white"
			@click="userStore.login"
		>
			Login with Twitch
		</UiSidebarMenuButton>
		<UiSidebarMenuButton
			size="lg"
			class="items-center justify-center bg-[#53FC18] text-black hover:bg-[#53FC18]/80 hover:text-black"
			@click="userStore.loginWithKick"
		>
			Login with Kick
		</UiSidebarMenuButton>
	</div>

	<UiDropdownMenu v-else>
		<UiDropdownMenuTrigger as-child>
			<UiSidebarMenuButton
				v-if="userStore.userWithoutDashboards"
				size="lg"
				class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
			>
				<img
					:src="userStore.userWithoutDashboards.twitchProfile.profileImageUrl"
					class="w-8 h-8 rounded-full"
				/>
				<div class="grid flex-1 text-left text-sm leading-tight">
					<span class="truncate font-semibold">{{
						userStore.userWithoutDashboards.twitchProfile.displayName
					}}</span>
					<span class="truncate text-xs">Logged as</span>
				</div>
				<Icon name="lucide:chevrons-up-down" class="ml-auto w-4 h-4" />
			</UiSidebarMenuButton>
		</UiDropdownMenuTrigger>

		<UiDropdownMenuContent class="min-w-56 rounded-lg" v-bind="dropdownProps">
			<!--			<UiDropdownMenuGroup> -->
			<!--				<div class="p-2"> -->
			<!--					There is nothing, yet... -->
			<!--				</div> -->
			<!--			</UiDropdownMenuGroup> -->

			<!--			<UiDropdownMenuSeparator /> -->

			<UiDropdownMenuItem @click="() => userStore.logout()">
				<Icon name="lucide:log-out" class="mr-2 w-4 h-4" />
				Logout
			</UiDropdownMenuItem>
		</UiDropdownMenuContent>
	</UiDropdownMenu>
</template>
