<script setup lang="ts">
import DonateDescription from '~~/layers/dashboard/app/components/integrations/helpers/donateDescription.vue'
import OauthComponent from '~~/layers/dashboard/app/components/integrations/variants/oauth.vue'
import { useIntegrationsPageData } from '~~/layers/dashboard/app/api/integrations/integrations-page.js'
import { useIntegrations } from '~~/layers/dashboard/app/api/integrations/integrations.js'

const integrationsPage = useIntegrationsPageData()
const integrationsManager = useIntegrations()
const logout = integrationsManager.streamlabsLogout()
</script>

<template>
	<oauth-component
		title="StreamLabs"
		:data="integrationsPage.streamlabsData.value"
		:logout="() => logout.executeMutation({})"
		:authLink="integrationsPage.streamlabsAuthLink.value"
		icon="twir-integrations:streamlabs"
	>
		<template #description>
			<donate-description />
		</template>
	</oauth-component>
</template>
