<script setup lang="ts">
import { useLastfmIntegration } from '@/api/index.js';
import IconLastFM from '@/assets/integrations/lastfm.svg?use';
import SongDescription from '@/components/integrations/helpers/songDescription.vue';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useLastfmIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();
</script>

<template>
	<oauth-component
		title="Last.fm"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconLastFM"
	>
		<template #description>
			<song-description />
		</template>
	</oauth-component>
</template>
