<script setup lang="ts">
import { useSpotifyIntegration } from '@/api/integrations/spotify.ts'
import IconSpotify from '@/assets/integrations/spotify.svg?use'
import SongDescription from '@/components/integrations/helpers/songDescription.vue'
import OauthComponent from '@/components/integrations/variants/oauth.vue'

const manager = useSpotifyIntegration()
const { data } = manager.spotifyData
</script>

<template>
	<OauthComponent
		title="Spotify"
		:data="data?.profile"
		:logout="() => manager.logout.executeMutation({})"
		:authLink="data?.spotifyAuthLink"
		:icon="IconSpotify"
	>
		<template #description>
			<SongDescription />
		</template>
	</OauthComponent>
</template>
