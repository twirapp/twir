<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile, useUserAccessFlagChecker } from '@/api'
import { Button } from '@/components/ui/button'
import { useTimersEdit } from '@/features/timers/composables/use-timers-edit'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const { t } = useI18n()
const userCanManageTimers = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageTimers)
const { data: profile } = useProfile()

const { timers } = useTimersEdit()
const timersLength = computed(() => timers.data?.value?.timers.length ?? 0)

const maxTimers = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan.maxTimers ?? 10
})
</script>

<template>
	<div class="flex gap-2">
		<RouterLink v-slot="{ navigate, href }" custom to="/dashboard/timers/create">
			<Button
				as="a"
				:href="href"
				:disabled="!userCanManageTimers || timersLength >= maxTimers"
				@click="navigate"
			>
				<PlusIcon class="size-4 mr-2" />
				{{ timersLength >= maxTimers ? t('timers.limitExceeded') : t('sharedButtons.create') }} ({{
					timersLength }}/{{ maxTimers }})
			</Button>
		</RouterLink>
	</div>
</template>
