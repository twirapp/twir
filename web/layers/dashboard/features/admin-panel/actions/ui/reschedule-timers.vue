<script setup lang="ts">
import { ref } from 'vue'

import { useMutationRescheduleTimers } from '#layers/dashboard/api/admin/actions'




const mutation = useMutationRescheduleTimers()

async function onReschedule() {
	await mutation.executeMutation({})
}

const confirmOpened = ref(false)
</script>

<template>
	<UiCard class="rounded-lg border bg-card text-card-foreground shadow-xs border-red-500 p-4 flex flex-col gap-4">
		<div class="flex items-center">
			<div class="flex-auto">
				<small class="text-sm font-medium leading-none">
					Reschedule timers
				</small>
				<p class="text-sm text-muted-foreground">
					Will drop all timers from queue and create them again
				</p>
			</div>
			<UiButton variant="destructive" @click="confirmOpened = true">
				Reschedule
			</UiButton>
		</div>
	</UiCard>

	<!-- TODO: reusable action confirm -->
	<UiActionConfirm
		v-model:open="confirmOpened"
		confirm-text="Are you sure you want to reschedule all timers?"
		@confirm="onReschedule"
	/>
</template>
