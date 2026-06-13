<script setup lang="ts">
import IconDonationAlerts from '~~/layers/dashboard/assets/integrations/donationalerts.svg?use'
import DonateDescription from '~~/layers/dashboard/components/integrations/helpers/donateDescription.vue'
import OauthComponent from '~~/layers/dashboard/components/integrations/variants/oauth.vue'
import { useIntegrationsPageData } from '~~/layers/dashboard/api/integrations/integrations-page.js'
import { useIntegrations } from '~~/layers/dashboard/api/integrations/integrations.js'

const integrationsPage = useIntegrationsPageData()
const integrationsManager = useIntegrations()
const logout = integrationsManager.donationAlertsLogout()
</script>

<template>
	<oauth-component
		title="DonationAlerts"
		:data="integrationsPage.donationAlertsData.value"
		:logout="() => logout.executeMutation({})"
		:authLink="integrationsPage.donationAlertsAuthLink.value"
		:icon="IconDonationAlerts"
	>
		<template #description>
			<donate-description />
		</template>
	</oauth-component>
</template>
