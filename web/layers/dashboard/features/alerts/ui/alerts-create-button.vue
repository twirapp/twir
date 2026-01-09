<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed } from 'vue'


import AlertsDialog from './alerts-dialog.vue'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useAlertsApi } from '#layers/dashboard/api/alerts'

import { ChannelRolePermissionEnum } from '~/gql/graphql'

const { t } = useI18n()
const { user: profile } = storeToRefs(useDashboardAuth())
const alertsApi = useAlertsApi()
const userCanManageAlerts = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageAlerts)

const { data: alertsData } = alertsApi.useAlertsQuery()

const maxAlerts = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan?.maxAlerts ?? 50
})

const alertsLength = computed(() => alertsData.value?.channelAlerts.length ?? 0)

const isCreateDisabled = computed(() => {
	return alertsLength.value >= maxAlerts.value || !userCanManageAlerts.value
})
</script>

<template>
	<div class="flex gap-2">
		<AlertsDialog>
			<template #dialog-trigger>
				<UiButton :disabled="isCreateDisabled">
					<PlusIcon class="size-4 mr-2" />
					{{ alertsLength >= maxAlerts ? t('alerts.limitExceeded') : t('alerts.createAlert') }} ({{
						alertsLength }}/{{ maxAlerts }})
				</UiButton>
			</template>
		</AlertsDialog>
	</div>
</template>
