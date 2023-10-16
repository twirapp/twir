<script setup lang='ts'>
import { useI18n } from 'vue-i18n';

import { useSpotifyIntegration } from '@/api/index.js';
import IconSpotify from '@/assets/icons/integrations/spotify.svg?component';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useSpotifyIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();

const { t } = useI18n();
</script>

<template>
	<oauth-component
		title="Spotify"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconSpotify"
		:description="t('integrations.songServicesInfo')"
	/>
</template>
