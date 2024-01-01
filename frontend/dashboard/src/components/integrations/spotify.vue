<script setup lang="ts">
import { useSpotifyIntegration } from '@/api/index.js';
import IconSpotify from '@/assets/icons/integrations/spotify.svg?component';
import SongDescription from '@/components/integrations/helpers/songDescription.vue';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useSpotifyIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();
</script>

<template>
	<oauth-component
		title="Spotify"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconSpotify"
	>
		<template #description>
			<song-description />
		</template>
	</oauth-component>
</template>
