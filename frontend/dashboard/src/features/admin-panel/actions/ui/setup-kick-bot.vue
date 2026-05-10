<script setup lang="ts">
import { ref } from 'vue'

import { useMutationKickBotSetupLink } from '@/api/admin/actions'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

const mutation = useMutationKickBotSetupLink()
const loading = ref(false)

async function onSetup() {
	loading.value = true
	const result = await mutation.executeMutation({})
	loading.value = false

	if (result.error) {
		toast.error(result.error.message)
		return
	}

	const url = result.data?.kickBotSetupLink
	if (url) {
		window.open(url, '_blank')
	}
}
</script>

<template>
	<Card class="rounded-lg border bg-card text-card-foreground shadow-xs p-4 flex flex-col gap-4">
		<div class="flex items-center">
			<div class="flex-auto">
				<small class="text-sm font-medium leading-none">
					Setup Kick Bot
				</small>
				<p class="text-sm text-muted-foreground">
					Authorize the default Kick bot account that will be used to send messages in chat.
				</p>
			</div>
			<Button :disabled="loading" @click="onSetup">
				Authorize
			</Button>
		</div>
	</Card>
</template>
