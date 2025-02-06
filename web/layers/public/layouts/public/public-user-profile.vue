<script setup lang="ts">
const userStore = useUser()
await useAsyncData('authedUser', async () => userStore.refetch().then(() => true))
</script>

<template>
	<UiSidebarMenuButton
		v-if="userStore.user"
		size="lg"
		class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
	>
		<img :src="userStore.user.twitchProfile.profileImageUrl" class="size-8 rounded-full" />
		<div class="grid flex-1 text-left text-sm leading-tight">
			<span class="truncate font-semibold">{{ userStore.user.twitchProfile.displayName }}</span>
			<span class="truncate text-xs">Logged as</span>
		</div>
		<Icon name="lucide:chevrons-up-down" class="ml-auto size-4" />
	</UiSidebarMenuButton>
	<UiSidebarMenuButton v-else size="lg" class="items-center justify-center">
		Login with Twitch
	</UiSidebarMenuButton>
</template>
