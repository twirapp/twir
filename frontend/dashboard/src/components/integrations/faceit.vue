<script setup lang='ts'>
import { useFaceitIntegration } from '@/api/index.js';
import IconFaceit from '@/assets/icons/integrations/faceit.svg?component';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useFaceitIntegration();
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
		:icon="IconFaceit"
	/>
</template>
