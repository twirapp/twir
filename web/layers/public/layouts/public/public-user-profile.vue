<script setup lang="ts">
import { computed } from 'vue'

import type { DropdownMenuContentProps } from 'radix-vue'

const authStore = useAuth()
await useAsyncData('user', async () => authStore.refetch().then(() => true))
await useAsyncData('authLink', async () => authStore.fetchAuthLink().then(() => true))

const dropdownProps = computed((): DropdownMenuContentProps & { class?: string } => {
	return {
		class: 'w-[--radix-dropdown-menu-trigger-width]',
		side: 'bottom',
		align: 'end',
		sideOffset: 4,
	}

	// return {
	// 	class: 'w-[300px]',
	// 	alignOffset: -4,
	// 	align: 'start',
	// 	sideOffset: 12,
	// 	side: 'right',
	// }
})
</script>

<template>
	<UiSidebarMenuButton v-if="!authStore.user" size="lg" class="items-center justify-center" :href="authStore.authLink" as="a">
		Login with Twitch
	</UiSidebarMenuButton>

	<UiDropdownMenu v-else>
		<UiDropdownMenuTrigger as-child>
			<UiSidebarMenuButton
				v-if="authStore.user"
				size="lg"
				class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
			>
				<img :src="authStore.user.twitchProfile.profileImageUrl" class="size-8 rounded-full" />
				<div class="grid flex-1 text-left text-sm leading-tight">
					<span class="truncate font-semibold">{{ authStore.user.twitchProfile.displayName }}</span>
					<span class="truncate text-xs">Logged as</span>
				</div>
				<Icon name="lucide:chevrons-up-down" class="ml-auto size-4" />
			</UiSidebarMenuButton>
		</UiDropdownMenuTrigger>

		<UiDropdownMenuContent
			class="min-w-56 rounded-lg"
			v-bind="dropdownProps"
		>
			<!--			<UiDropdownMenuGroup> -->
			<!--				<div class="p-2"> -->
			<!--					There is nothing, yet... -->
			<!--				</div> -->
			<!--			</UiDropdownMenuGroup> -->

			<!--			<UiDropdownMenuSeparator /> -->

			<UiDropdownMenuItem @click="() => authStore.logout()">
				<Icon name="lucide:log-out" class="mr-2 size-4" />
				Logout
			</UiDropdownMenuItem>
		</UiDropdownMenuContent>
	</UiDropdownMenu>
</template>
