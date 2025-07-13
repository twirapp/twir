<script lang="ts" setup>
import { IconDog, IconMoon, IconSun } from '@tabler/icons-vue'
import { NButton, NTooltip } from 'naive-ui'
import { computed, ref } from 'vue'

import Card from './card.vue'

import { useProfile, useTwitchGetUsers } from '@/api/index.js'
import { useTheme } from '@/composables/use-theme.js'

const { data: profile } = useProfile()
const { theme: chatTheme, toggleTheme } = useTheme()

const openFrankerFaceZ = ref(false)
const selectedDashboardTwitchUser = computed(() => {
	return profile.value?.availableDashboards.find((d) => d.id === profile.value?.selectedDashboardId)
		?.twitchProfile
})

const chatUrl = computed(() => {
	if (!selectedDashboardTwitchUser.value) return

	const user = selectedDashboardTwitchUser.value

	let url = `https://www.twitch.tv/embed/${user.login}/chat?parent=${window.location.host}`

	if (chatTheme.value === 'dark') {
		url += '&darkpopout'
	}

	if (openFrankerFaceZ.value) {
		url += '&ffz-settings'
	}

	return url
})
</script>

<template>
	<Card :content-style="{ 'margin-bottom': '10px', padding: '0px' }">
		<template #header-extra>
			<NTooltip trigger="hover" placement="bottom">
				<template #trigger>
					<NButton size="small" text @click="openFrankerFaceZ = !openFrankerFaceZ">
						<IconDog />
					</NButton>
				</template>

				FrankerFaceZ Control Center
			</NTooltip>

			<NButton size="small" text @click="toggleTheme">
				<IconSun v-if="chatTheme === 'dark'" color="orange" />
				<IconMoon v-else />
			</NButton>
		</template>

		<iframe v-if="chatUrl" :src="chatUrl" frameborder="0" class="w-full h-full"> </iframe>
	</Card>
</template>
