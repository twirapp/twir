<script setup lang="ts">

import { useNightbotIntegration } from '@/api/index.js';
import IconNightbot from '@/assets/integrations/nightbot.svg?use';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useNightbotIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();

const { data: authLink } = manager.useAuthLink();

</script>

<template>
	<oauth-component
		title="Nightbot"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconNightbot"
	>
		<template #description>
			<i18n-t keypath="integrations.nightbot.info" />
		</template>
	</oauth-component>
</template>
