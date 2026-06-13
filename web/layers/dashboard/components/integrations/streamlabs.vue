<script setup lang="ts">
import IconStreamLabs from '~~/layers/dashboard/assets/integrations/streamlabs.svg'
import DonateDescription from '~~/layers/dashboard/components/integrations/helpers/donateDescription.vue'
import OauthComponent from '~~/layers/dashboard/components/integrations/variants/oauth.vue'
import { useIntegrationsPageData } from '~~/layers/dashboard/api/integrations/integrations-page.js'
import { useIntegrations } from '~~/layers/dashboard/api/integrations/integrations.js'

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
		:icon="IconStreamLabs"
	>
		<template #description>
			<donate-description />
		</template>
	</oauth-component>
</template>
