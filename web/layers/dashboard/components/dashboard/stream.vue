<script lang="ts" setup>
import { computed } from 'vue'

import Card from './card.vue'

import { useProfile } from '~~/layers/dashboard/api/auth'

const { data: profile } = useProfile()
const requestUrl = useRequestURL()

const selectedDashboardTwitchUser = computed(() => {
	const dashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId,
	)
	if (dashboard?.profile?.platform.toLowerCase() !== 'twitch') return undefined

	return dashboard.profile
})

const streamUrl = computed(() => {
	if (!selectedDashboardTwitchUser.value) return

	const user = selectedDashboardTwitchUser.value
	const url = `https://player.twitch.tv/?channel=${user.platformLogin}&parent=${requestUrl.host}&autoplay=false`

	return url
})
</script>

<template>
	<Card>
		<iframe
			v-if="streamUrl"
			:src="streamUrl"
			width="100%"
			height="100%"
			frameborder="0"
			scrolling="no"
			allowfullscreen="true"
		>
		</iframe>
	</Card>
</template>
