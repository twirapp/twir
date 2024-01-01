<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import { useDonationAlertsIntegration } from '@/api/index.js';
import IconDonationAlerts from '@/assets/icons/integrations/donationalerts.svg?component';
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useDonationAlertsIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();

const { t } = useI18n();
</script>

<template>
	<oauth-component
		title="DonationAlerts"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconDonationAlerts"
	>
		<template #description>
			<donate-description />
		</template>
	</oauth-component>
</template>
