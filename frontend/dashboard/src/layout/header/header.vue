<script setup lang="ts">
import { Menu } from 'lucide-vue-next'

import { useSidebar } from '@/components/ui/sidebar'
import { useIsMobile } from '@/composables/use-is-mobile'
import DiscordButton from '@/layout/buttons/button-discord.vue'
import GithubButton from '@/layout/buttons/button-github.vue'
import ButtonToggleTheme from '@/layout/buttons/button-toggle-theme.vue'
import DropdownLanguage from '@/layout/dropdowns/dropdown-language.vue'
import NotificationsButton from '@/layout/notifications/notifications-button.vue'

const { isMobile } = useIsMobile()

const sidebar = useSidebar()
const openSidebar = () => sidebar.setOpenMobile(true)
</script>

<template>
	<div class="w-full bg-sidebar h-[45px] border-b border-b-border">
		<div class="flex h-full items-center px-3.5 py-2 justify-end gap-3">
			<div class="flex gap-3">
				<div class="flex gap-2 items-center">
					<DiscordButton />
					<GithubButton />
					<NotificationsButton />
					<DropdownLanguage />
					<ButtonToggleTheme />
				</div>
			</div>

			<template v-if="isMobile">
				<Button @click="openSidebar">
					<Menu />
				</Button>
			</template>
		</div>
	</div>
</template>

<style scoped>
.header {
	padding-top: 10px;
	padding-left: 10px;
	padding-right: 10px;
}

.content {
	display: flex;
	align-items: center;
	text-decoration: none;
	padding-right: 14px;
}

.content {
	display: flex;
	justify-content: space-between;
	align-items: center;
	height: 45px;
	column-gap: 12px;
}
</style>
