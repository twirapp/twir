<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import DropSessions from './ui/drop-sessions.vue'
import EventsubSubscribe from './ui/eventsub-subscribe.vue'
import RescheduleTimers from './ui/reschedule-timers.vue'

import ReinitEventsubChannels from '@/features/admin-panel/actions/ui/reinit-eventsub-channels.vue'

const { t } = useI18n()
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<EventsubSubscribe />

		<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
			{{ t('adminPanel.adminActions.dangerZone.title') }}
		</h4>

		<DropSessions />
		<RescheduleTimers />
		<ReinitEventsubChannels />
	</div>
</template>
