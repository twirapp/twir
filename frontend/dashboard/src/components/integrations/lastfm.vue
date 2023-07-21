<script setup lang='ts'>
import { useLastfmIntegration } from '@/api/index.js';
import IconLastFM from '@/assets/icons/integrations/lastfm.svg?component';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useLastfmIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();
</script>

<template>
	<oauth-component
		name="Faceit"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconLastFM"
	/>
</template>
