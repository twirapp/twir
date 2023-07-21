<script setup lang='ts'>
import { useStreamlabsIntegration } from '@/api/index.js';
import IconStreamLabs from '@/assets/icons/integrations/streamlabs.svg?component';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useStreamlabsIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();
</script>

<template>
	<oauth-component
		name="StreamLabs"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconStreamLabs"
	/>
</template>
