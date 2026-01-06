<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import AlertsDialog from './alerts-dialog.vue'

import { useProfile, useUserAccessFlagChecker } from '@/api/auth'
import { useAlertsApi } from '@/api/alerts'
import { Button } from '@/components/ui/button'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const { t } = useI18n()
const { data: profile } = useProfile()
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
				<Button :disabled="isCreateDisabled">
					<PlusIcon class="size-4 mr-2" />
					{{ alertsLength >= maxAlerts ? t('alerts.limitExceeded') : t('alerts.createAlert') }} ({{
						alertsLength }}/{{ maxAlerts }})
				</Button>
			</template>
		</AlertsDialog>
	</div>
</template>
