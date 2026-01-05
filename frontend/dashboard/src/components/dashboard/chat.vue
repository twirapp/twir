<script lang="ts" setup>
import { IconDog, IconMoon, IconSun } from '@tabler/icons-vue'
import { computed, ref } from 'vue'

import Card from './card.vue'

import { useProfile } from '@/api/index.js'
import { useTheme } from '@/composables/use-theme.js'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

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
	<Card :content-style="{ 'margin-bottom': '10px', 'padding': '0px' }">
		<template #header-extra>
			<TooltipProvider>
				<Tooltip>
					<TooltipTrigger as-child>
						<Button size="sm" variant="ghost" @click="openFrankerFaceZ = !openFrankerFaceZ">
							<IconDog />
						</Button>
					</TooltipTrigger>
					<TooltipContent>
						FrankerFaceZ Control Center
					</TooltipContent>
				</Tooltip>
			</TooltipProvider>

			<Button size="sm" variant="ghost" @click="toggleTheme">
				<IconSun v-if="chatTheme === 'dark'" color="orange" />
				<IconMoon v-else />
			</Button>
		</template>

		<iframe v-if="chatUrl" :src="chatUrl" frameborder="0" class="w-full h-full"> </iframe>
	</Card>
</template>
