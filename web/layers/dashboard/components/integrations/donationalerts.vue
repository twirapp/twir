<script setup lang="ts">
import IconDonationAlerts from '@/assets/integrations/donationalerts.svg?use'
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue'
import OauthComponent from '@/components/integrations/variants/oauth.vue'
import { useIntegrationsPageData } from '@/api/integrations/integrations-page.ts'
import { useIntegrations } from '@/api/integrations/integrations.ts'

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
