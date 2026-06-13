<script setup lang="ts">
import { useSpotifyIntegration } from '@/api/integrations/spotify.js'
import IconSpotify from '@/assets/integrations/spotify.svg?use'
import SongDescription from '@/components/integrations/helpers/songDescription.vue'
import OauthComponent from '@/components/integrations/variants/oauth.vue'
import { useIntegrationsPageData } from '@/api/integrations/integrations-page.js'

const manager = useSpotifyIntegration()
const integrationsPage = useIntegrationsPageData()
</script>

<template>
	<OauthComponent
		title="Spotify"
		:data="integrationsPage.spotifyData.value"
		:logout="() => manager.logout.executeMutation({})"
		:authLink="integrationsPage.spotifyAuthLink.value"
		:icon="IconSpotify"
	>
		<template #description>
			<SongDescription />
		</template>
	</OauthComponent>
</template>
