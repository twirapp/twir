<script setup lang="ts">
import { toast } from 'vue-sonner'

import { useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth.js'
import { useIntegrations } from '~~/layers/dashboard/api/integrations/integrations.js'
import { useIntegrationsPageData } from '~~/layers/dashboard/api/integrations/integrations-page.js'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'
import ImportProviderCard from '../components/import-provider-card.vue'

const { t } = useI18n()
const integrationsPage = useIntegrationsPageData()
const streamlabsLogout = useIntegrations().streamlabsLogout()
const canManageIntegrations = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageIntegrations)
const authLoading = ref(false)

async function logout() {
	authLoading.value = true
	try {
		const result = await streamlabsLogout.executeMutation({})
		if (result.error || !result.data?.streamlabsLogout) {
			toast.error(t('imports.errors.logout'))
			return
		}
		await integrationsPage.refetch()
	} catch {
		toast.error(t('imports.errors.logout'))
	} finally {
		authLoading.value = false
	}
}
</script>

<template>
	<ImportProviderCard
		title="Streamlabs"
		icon="twir-integrations:streamlabs"
		:description="t('imports.providers.streamlabs.description')"
		:donation-description="integrationsPage.streamlabsData.value?.enabled ? t('imports.providers.streamlabs.donationsEnabled') : ''"
		:unavailable-description="t('imports.providers.streamlabs.importUnavailable')"
		:account="integrationsPage.streamlabsData.value"
		:auth-link="integrationsPage.streamlabsAuthLink.value"
		:is-loading="integrationsPage.fetching.value || authLoading"
		:can-manage-integration="canManageIntegrations"
		:import-available="false"
		@logout="logout"
	/>
</template>
