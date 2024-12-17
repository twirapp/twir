<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useMutationRescheduleTimers } from '@/api/admin/actions'
import ActionConfirm from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

const { t } = useI18n()
const mutation = useMutationRescheduleTimers()

async function onReschedule() {
	await mutation.executeMutation({})
}

const confirmOpened = ref(false)
</script>

<template>
	<Card class="rounded-lg border bg-card text-card-foreground shadow-sm border-red-500 p-4 flex flex-col gap-4">
		<div class="flex items-center">
			<div class="flex-auto">
				<small class="text-sm font-medium leading-none">
					Reschedule timers
				</small>
				<p class="text-sm text-muted-foreground">
					Will drop all timers from queue and create them again
				</p>
			</div>
			<Button variant="destructive" @click="confirmOpened = true">
				Reschedule
			</Button>
		</div>
	</Card>

	<!-- TODO: reusable action confirm -->
	<ActionConfirm
		v-model:open="confirmOpened"
		:confirm-text="t('adminPanel.adminActions.dangerZone.revokeAllSessionsConfirm')"
		@confirm="onReschedule"
	/>
</template>
