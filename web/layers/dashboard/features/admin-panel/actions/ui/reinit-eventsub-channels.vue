<script setup lang="ts">
import { ref } from 'vue'

import { useMutationEventSubInitChannels } from '#layers/dashboard/api/admin/actions'




const mutation = useMutationEventSubInitChannels()

async function onResubscribe() {
	await mutation.executeMutation({})
}

const confirmOpened = ref(false)
</script>

<template>
	<UiCard class="rounded-lg border bg-card text-card-foreground shadow-xs border-red-500 p-4 flex flex-col gap-4">
		<div class="flex items-center">
			<div class="flex-auto">
				<small class="text-sm font-medium leading-none">
					Reinit eventsub subscriptions
				</small>
				<p class="text-sm text-muted-foreground">
					Will recreate eventsub subscriptions
				</p>
			</div>
			<UiButton variant="destructive" @click="confirmOpened = true">
				Resubscribe
			</UiButton>
		</div>
	</UiCard>

	<!-- TODO: reusable action confirm -->
	<UiActionConfirm
		v-model:open="confirmOpened"
		confirm-text="Are you sure you want to resubscribe all channels?"
		@confirm="onResubscribe"
	/>
</template>
