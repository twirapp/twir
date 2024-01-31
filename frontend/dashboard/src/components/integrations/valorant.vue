<script setup lang="ts">

import { useValorantIntegration } from '@/api/index.js';
import IconValorant from '@/assets/integrations/valorant.svg?use';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useValorantIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();
</script>

<template>
	<oauth-component
		title="Valorant"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconValorant"
	>
		<template #description>
			123
		</template>
	</oauth-component>
</template>
