<script setup lang="ts">
import ReinitEventsubChannels from '~~/layers/dashboard/features/admin-panel/actions/ui/reinit-eventsub-channels.vue'

import DropSessions from './ui/drop-sessions.vue'
import EventsubSubscribe from './ui/eventsub-subscribe.vue'
import RescheduleTimers from './ui/reschedule-timers.vue'
import SetupKickBot from './ui/setup-kick-bot.vue'

const { t } = useI18n()
</script>

<template>
	<div class="flex w-full flex-col gap-4">
		<EventsubSubscribe />
		<SetupKickBot />

		<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
			{{ t('adminPanel.adminActions.dangerZone.title') }}
		</h4>

		<DropSessions />
		<RescheduleTimers />
		<ReinitEventsubChannels />
	</div>
</template>
