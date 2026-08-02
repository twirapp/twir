<script setup lang="ts">
import type { ImportReportData } from '../components/types.js'

import { toast } from 'vue-sonner'

import { useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth.js'
import { useIntegrationsPageData } from '~~/layers/dashboard/api/integrations/integrations-page.js'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'
import ImportProviderCard from '../components/import-provider-card.vue'
import ImportSettings from '../components/import-settings.vue'
import { useStreamElementsIntegration } from './composables/use-streamelements-integration.js'

const { t } = useI18n()
const integrationsPage = useIntegrationsPageData()
const streamElements = useStreamElementsIntegration()

const commandsReport = ref<ImportReportData | null>(null)
const timersReport = ref<ImportReportData | null>(null)
const commandsImporting = ref(false)
const timersImporting = ref(false)
const authLoading = ref(false)

const connected = computed(() => Boolean(integrationsPage.streamelementsData.value?.userName))
const canManageIntegrations = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageIntegrations)
const canManageCommands = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageCommands)
const canManageTimers = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageTimers)

async function importCommands() {
	commandsImporting.value = true
	try {
		const result = await streamElements.importCommands.executeMutation({})
		if (result.error || !result.data?.streamelementsImportCommands) {
			toast.error(t('imports.errors.import'))
			return
		}
		commandsReport.value = result.data.streamelementsImportCommands
	} catch {
		toast.error(t('imports.errors.import'))
	} finally {
		commandsImporting.value = false
	}
}

async function importTimers() {
	timersImporting.value = true
	try {
		const result = await streamElements.importTimers.executeMutation({})
		if (result.error || !result.data?.streamelementsImportTimers) {
			toast.error(t('imports.errors.import'))
			return
		}
		timersReport.value = result.data.streamelementsImportTimers
	} catch {
		toast.error(t('imports.errors.import'))
	} finally {
		timersImporting.value = false
	}
}

async function logout() {
	authLoading.value = true
	try {
		const result = await streamElements.logout.executeMutation({})
		if (result.error || !result.data?.streamelementsLogout) {
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
		title="StreamElements"
		icon="twir-integrations:streamelements"
		:connected="connected"
		:description="t('imports.providers.streamelements.description')"
		:donation-description="connected ? t('imports.providers.streamelements.donationsEnabled') : ''"
		:account="integrationsPage.streamelementsData.value"
		:auth-link="integrationsPage.streamelementsAuthLink.value"
		:is-loading="integrationsPage.fetching.value || authLoading"
		:can-manage-integration="canManageIntegrations"
		@logout="logout"
	>
		<template #settings>
			<ImportSettings
				:connected="connected"
				:can-manage-commands="canManageCommands"
				:can-manage-timers="canManageTimers"
				:commands-importing="commandsImporting"
				:timers-importing="timersImporting"
				:commands-report="commandsReport"
				:timers-report="timersReport"
				@import-commands="importCommands"
				@import-timers="importTimers"
			/>
		</template>
	</ImportProviderCard>
</template>
