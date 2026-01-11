<script setup lang="ts">
import IconStreamLabs from '@/assets/integrations/streamlabs.svg?use'
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue'
import OauthComponent from '@/components/integrations/variants/oauth.vue'
import { useIntegrationsPageData } from '@/api/integrations/integrations-page.ts'
import { useIntegrations } from '@/api/integrations/integrations.ts'

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
