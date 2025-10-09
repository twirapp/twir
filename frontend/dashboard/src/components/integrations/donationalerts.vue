<script setup lang="ts">
import IconDonationAlerts from '@/assets/integrations/donationalerts.svg?use'
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue'
import OauthComponent from '@/components/integrations/variants/oauth.vue'
import { useIntegrations } from '@/api/integrations/integrations.ts'

const integrationsManager = useIntegrations()
const { data } = integrationsManager.useQuery()

const logout = integrationsManager.donationAlertsLogout()
</script>

<template>
	<oauth-component
		title="DonationAlerts"
		:data="data?.donationAlerts"
		:logout="() => logout.executeMutation({})"
		:authLink="data?.donationAlertsAuthLink"
		:icon="IconDonationAlerts"
	>
		<template #description>
			<donate-description />
		</template>
	</oauth-component>
</template>
