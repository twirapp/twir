<script setup lang="ts">
import { useSpotifyIntegration } from '@/api/integrations/spotify.ts'
import IconSpotify from '@/assets/integrations/spotify.svg?use'
import SongDescription from '@/components/integrations/helpers/songDescription.vue'
import OauthComponent from '@/components/integrations/variants/oauth.vue'
import { useIntegrations } from '@/api/integrations/integrations.ts'

const manager = useSpotifyIntegration()

const integrationsManager = useIntegrations()
const { data } = integrationsManager.useQuery()
</script>

<template>
	<OauthComponent
		title="Spotify"
		:data="data?.spotifyData"
		:logout="() => manager.logout.executeMutation({})"
		:authLink="data?.spotifyAuthLink"
		:icon="IconSpotify"
	>
		<template #description>
			<SongDescription />
		</template>
	</OauthComponent>
</template>
