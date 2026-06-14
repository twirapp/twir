<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useProfile, useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth'
import { useEventsApi } from '~~/layers/dashboard/api/events.js'
import { Button } from '@/components/ui/button'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

const { t } = useI18n()
const router = useRouter()
const { data: profile } = useProfile()
const eventApi = useEventsApi()
const userCanManageEvents = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageEvents)

const { data: events } = eventApi.useQueryEvents()

const maxEvents = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan?.maxEvents ?? 50
})

const eventsLength = computed(() => events.value?.events?.length ?? 0)

const isCreateDisabled = computed(() => {
	return eventsLength.value >= maxEvents.value || !userCanManageEvents.value
})

function createEvent() {
	router.push('/dashboard/events/new')
}
</script>

<template>
	<div class="flex gap-2">
		<Button type="button" :disabled="isCreateDisabled" @click="createEvent">
			<Icon name="lucide:plus" class="size-4 mr-2" />
			{{ eventsLength >= maxEvents ? t('events.limitExceeded') : t('sharedTexts.create') }} ({{
				eventsLength }}/{{ maxEvents }})
		</Button>
	</div>
</template>
