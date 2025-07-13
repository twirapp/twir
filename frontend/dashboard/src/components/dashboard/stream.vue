<script lang="ts" setup>
import { computed } from 'vue'

import Card from './card.vue'

import { useProfile } from '@/api/index.js'

const { data: profile } = useProfile()

const selectedDashboardTwitchUser = computed(() => {
	return profile.value?.availableDashboards.find((d) => d.id === profile.value?.selectedDashboardId)
		?.twitchProfile
})

const streamUrl = computed(() => {
	if (!selectedDashboardTwitchUser.value) return

	const user = selectedDashboardTwitchUser.value
	const url = `https://player.twitch.tv/?channel=${user.login}&parent=${window.location.host}&autoplay=false`

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
