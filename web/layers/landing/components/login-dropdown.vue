<script setup lang="ts">
interface Props {
	variant?: 'hero' | 'header'
}

const props = withDefaults(defineProps<Props>(), {
	variant: 'header',
})

const userStore = useAuth()

const platforms = [
	{
		label: 'Twitch',
		icon: 'simple-icons:twitch',
		iconClass: 'text-[#9146FF]',
		login: () => userStore.login(),
	},
	{
		label: 'Kick',
		icon: 'simple-icons:kick',
		iconClass: 'text-[#53FC18]',
		login: () => userStore.loginWithKick(),
	},
	{
		label: 'VK Live',
		icon: 'simple-icons:vk',
		iconClass: 'text-[#0077FF]',
		login: () => userStore.loginWithVk(),
	},
]

const triggerClass = computed(() => {
	return props.variant === 'hero'
		? 'xs:py-4 inline-flex items-center justify-center gap-2 rounded-lg bg-[#5D58F5] px-7 py-3 text-center text-base font-semibold whitespace-nowrap text-white transition-[background,box-shadow] hover:bg-[#6964FF] focus-visible:ring-4 focus-visible:ring-[#5D58F5]/50 focus-visible:outline-none sm:text-lg'
		: 'flex flex-row px-4 py-2 items-center gap-2 bg-[#5D58F5] text-white rounded-lg font-medium focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-[#5D58F5]/50 cursor-pointer hover:bg-[#6964FF] transition-shadow'
})
</script>

<template>
	<UiDropdownMenu>
		<UiDropdownMenuTrigger as="button" :class="triggerClass">
			{{ variant === 'hero' ? 'Start now' : 'Login' }}
			<Icon name="lucide:chevron-down" class="h-4 w-4" />
		</UiDropdownMenuTrigger>

		<UiDropdownMenuContent align="end" side="bottom" :side-offset="4" class="w-44">
			<UiDropdownMenuItem
				v-for="platform in platforms"
				:key="platform.label"
				as="button"
				class="flex w-full items-center cursor-pointer"
				@click="platform.login"
			>
				<Icon :name="platform.icon" :class="['mr-2 h-4 w-4', platform.iconClass]" />
				{{ platform.label }}
			</UiDropdownMenuItem>
		</UiDropdownMenuContent>
	</UiDropdownMenu>
</template>
