<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed } from 'vue'


import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'

import { useTimersEdit } from '~/features/timers/composables/use-timers-edit'
import { ChannelRolePermissionEnum } from '~/gql/graphql'

const { t } = useI18n()
const userCanManageTimers = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageTimers)
const { user: profile } = storeToRefs(useDashboardAuth())

const { timers } = useTimersEdit()
const timersLength = computed(() => timers.data?.value?.timers.length ?? 0)

const maxTimers = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan?.maxTimers ?? 10
})
</script>

<template>
	<div class="flex gap-2">
		<RouterLink v-slot="{ navigate, href }" custom to="/dashboard/timers/create">
			<UiButton
				as="a"
				:href="href"
				:disabled="!userCanManageTimers || timersLength >= maxTimers"
				@click="navigate"
			>
				<PlusIcon class="size-4 mr-2" />
				{{ timersLength >= maxTimers ? t('timers.limitExceeded') : t('sharedButtons.create') }} ({{
					timersLength }}/{{ maxTimers }})
			</UiButton>
		</RouterLink>
	</div>
</template>
