<script setup lang="ts">

import { useStreamlabsIntegration } from '@/api/integrations/oauth';
import IconStreamLabs from '@/assets/integrations/streamlabs.svg?use';
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useStreamlabsIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();
</script>

<template>
	<oauth-component
		title="StreamLabs"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconStreamLabs"
	>
		<template #description>
			<donate-description />
		</template>
	</oauth-component>
</template>
