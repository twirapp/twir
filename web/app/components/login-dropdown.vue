<script setup lang="ts">
import { Platform } from '~/gql/graphql.js'
import { PLATFORM_OPTIONS } from '~/utils/platforms.js'

interface Props {
	variant?: 'hero' | 'header' | 'sidebar'
}

const props = withDefaults(defineProps<Props>(), {
	variant: 'header',
})

const userStore = useAuth()

const loginHandlers: Record<Platform, () => void> = {
	[Platform.Twitch]: () => userStore.login(),
	[Platform.Kick]: () => userStore.loginWithKick(),
	[Platform.VkVideoLive]: () => userStore.loginWithVk(),
	[Platform.Youtube]: () => userStore.loginWithYoutube(),
}

const platforms = PLATFORM_OPTIONS.map((meta) => ({
	...meta,
	label: meta.platform === Platform.VkVideoLive ? 'VK Live' : meta.label,
	login: loginHandlers[meta.platform],
}))

const triggerClass = computed(() => {
	return props.variant === 'hero'
		? 'xs:py-4 inline-flex items-center justify-center gap-2 rounded-lg bg-[#5D58F5] px-7 py-3 text-center text-base font-semibold whitespace-nowrap text-white transition-[background,box-shadow] hover:bg-[#6964FF] focus-visible:ring-4 focus-visible:ring-[#5D58F5]/50 focus-visible:outline-none sm:text-lg'
		: 'flex flex-row px-4 py-2 items-center gap-2 bg-[#5D58F5] text-white rounded-lg font-medium focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-[#5D58F5]/50 cursor-pointer hover:bg-[#6964FF] transition-shadow'
})

const contentClass = computed(() => {
	return props.variant === 'sidebar' ? 'min-w-56 rounded-lg w-(--reka-dropdown-menu-trigger-width)' : 'w-44'
})
</script>

<template>
	<UiDropdownMenu>
		<UiDropdownMenuTrigger v-if="variant === 'sidebar'" as-child>
			<UiSidebarMenuButton
				size="lg"
				class="items-center justify-center bg-[#5D58F5] text-white hover:bg-[#6964FF] hover:text-white"
			>
				Login
				<Icon name="lucide:chevrons-up-down" class="ml-auto h-4 w-4" />
			</UiSidebarMenuButton>
		</UiDropdownMenuTrigger>

		<UiDropdownMenuTrigger v-else as="button" :class="triggerClass">
			{{ variant === 'hero' ? 'Start now' : 'Login' }}
			<Icon name="lucide:chevron-down" class="h-4 w-4" />
		</UiDropdownMenuTrigger>

		<UiDropdownMenuContent align="end" side="bottom" :side-offset="4" :class="contentClass">
			<UiDropdownMenuItem
				v-for="platform in platforms"
				:key="platform.platform"
				as="button"
				class="flex w-full items-center cursor-pointer"
				@click="platform.login"
			>
				<Icon :name="platform.icon" :class="['mr-2 h-4 w-4', platform.colorClass]" />
				{{ platform.label }}
			</UiDropdownMenuItem>
		</UiDropdownMenuContent>
	</UiDropdownMenu>
</template>
