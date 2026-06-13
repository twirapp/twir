<script setup lang="ts">
import { ref } from 'vue'

import { useMutationEventSubInitChannels } from '@/api/admin/actions'
import ActionConfirm from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

const mutation = useMutationEventSubInitChannels()

async function onResubscribe() {
	await mutation.executeMutation({})
}

const confirmOpened = ref(false)
</script>

<template>
	<Card class="rounded-lg border bg-card text-card-foreground shadow-xs border-red-500 p-4 flex flex-col gap-4">
		<div class="flex items-center">
			<div class="flex-auto">
				<small class="text-sm font-medium leading-none">
					Reinit eventsub subscriptions
				</small>
				<p class="text-sm text-muted-foreground">
					Will recreate eventsub subscriptions
				</p>
			</div>
			<Button variant="destructive" @click="confirmOpened = true">
				Resubscribe
			</Button>
		</div>
	</Card>

	<!-- TODO: reusable action confirm -->
	<ActionConfirm
		v-model:open="confirmOpened"
		confirm-text="Are you sure you want to resubscribe all channels?"
		@confirm="onResubscribe"
	/>
</template>
