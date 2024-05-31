<script setup lang="ts">
import { useThemeVars } from 'naive-ui'
import { computed } from 'vue'

import ButtonDiscord from '../buttons/button-discord.vue'
import ButtonGithub from '../buttons/button-github.vue'
import ButtonPublicPage from '../buttons/button-public-page.vue'
import ButtonToggleTheme from '../buttons/button-toggle-theme.vue'
import DropdownLanguage from '../dropdowns/dropdown-language.vue'
import DropdownProfile from '../dropdowns/dropdown-profile'
import NotificationsButton from '../notifications/notifications-button.vue'
import StatsScrollable from '../stats/stats-scrollable.vue'

const themeVars = useThemeVars()
const blockColor = computed(() => themeVars.value.buttonColor2)
</script>

<template>
	<div style="flex-grow: 1; overflow-y: hidden; overflow-x: auto;">
		<StatsScrollable />
	</div>

	<div style="display: flex; gap: 12px;">
		<div class="block">
			<ButtonGithub />
			<ButtonDiscord />
		</div>

		<div class="block">
			<NotificationsButton />
			<DropdownLanguage />
			<ButtonToggleTheme />
			<ButtonPublicPage />
			<DropdownProfile />
		</div>
	</div>
</template>

<style scoped>
.block {
	background-color: v-bind(blockColor);
	display: flex;
	gap: 16px;
	padding: 16px;
	border-radius: 10px;
	align-items: center;
	height: 45px;
}
</style>
