<script setup lang="ts">
import { computed, onMounted } from 'vue'

import IconVk from '@/assets/integrations/vk.svg?use'
import SongDescription from '@/components/integrations/helpers/songDescription.vue'
import OauthComponent from '@/components/integrations/variants/oauth.vue'
import {
	useVKIntegration,
	vkBroadcaster,
} from '@/features/integrations/composables/vk/use-vk-integration.ts'

const { userName, avatar, logout, authLink, isDataFetching, refetchData } = useVKIntegration()

const oauthData = computed(() => {
	if (!userName.value && !avatar.value) {
		return null
	}

	return {
		userName: userName.value,
		avatar: avatar.value,
	}
})

onMounted(async () => {
	vkBroadcaster.onmessage = async (event) => {
		if (event.data !== 'refresh') return

		await refetchData()
	}
})
</script>

<template>
	<OauthComponent
		title="VK"
		:data="oauthData"
		:logout="logout"
		:authLink="authLink"
		:icon="IconVk"
		:is-loading="isDataFetching"
	>
		<template #description>
			<SongDescription />
		</template>
	</OauthComponent>
</template>
