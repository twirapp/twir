<script setup lang="ts">
import { useSpotifyIntegration } from '~~/layers/dashboard/app/api/integrations/spotify.js'
import SongDescription from '~~/layers/dashboard/app/components/integrations/helpers/songDescription.vue'
import OauthComponent from '~~/layers/dashboard/app/components/integrations/variants/oauth.vue'
import { useIntegrationsPageData } from '~~/layers/dashboard/app/api/integrations/integrations-page.js'

const manager = useSpotifyIntegration()
const integrationsPage = useIntegrationsPageData()
</script>

<template>
	<OauthComponent
		title="Spotify"
		:data="integrationsPage.spotifyData.value"
		:logout="() => manager.logout.executeMutation({})"
		:authLink="integrationsPage.spotifyAuthLink.value"
		icon="twir-integrations:spotify"
	>
		<template #description>
			<SongDescription />
		</template>
	</OauthComponent>
</template>
