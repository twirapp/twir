<script setup lang="ts">

import { useDonationAlertsIntegration, useValorantIntegration } from '@/api/index.js';
import IconDonationAlerts from '@/assets/integrations/donationalerts.svg?use';
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue';
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
		:icon="IconDonationAlerts"
	>
		<template #description>
			123
		</template>
	</oauth-component>
</template>
