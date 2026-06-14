<script setup lang="ts">
import { computed, onMounted } from 'vue'

import SongDescription from '~~/layers/dashboard/components/integrations/helpers/songDescription.vue'
import OauthComponent from '~~/layers/dashboard/components/integrations/variants/oauth.vue'
import {
	lastfmBroadcaster,
	useLastfmIntegration,
} from '~~/layers/dashboard/features/integrations/composables/lastfm/use-lastfm-integration.js'

const { userName, avatar, logout, authLink, isDataFetching, refetchData } = useLastfmIntegration()

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
	lastfmBroadcaster.onmessage = async (event) => {
		if (event.data !== 'refresh') return

		await refetchData()
	}
})
</script>

<template>
	<OauthComponent
		title="Last.fm"
		:data="oauthData"
		:logout="logout"
		:authLink="authLink"
		icon="twir-integrations:lastfm"
		:is-loading="isDataFetching"
	>
		<template #description>
			<SongDescription />
		</template>
	</OauthComponent>
</template>
