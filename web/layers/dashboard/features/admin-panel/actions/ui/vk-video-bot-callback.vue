<script setup lang="ts">
import * as z from 'zod'
import { useRoute } from 'vue-router'
import { toast } from 'vue-sonner'

import {
	useMutationVKVideoBotSetupComplete,
	vkVideoBotSetupBroadcastChannelName,
} from '~~/layers/dashboard/api/admin/actions'
import Alert from '@/components/ui/alert/Alert.vue'
import AlertDescription from '@/components/ui/alert/AlertDescription.vue'
import Card from '@/components/ui/card/Card.vue'
import CardContent from '@/components/ui/card/CardContent.vue'
import CardDescription from '@/components/ui/card/CardDescription.vue'
import CardHeader from '@/components/ui/card/CardHeader.vue'
import CardTitle from '@/components/ui/card/CardTitle.vue'

const callbackQuerySchema = z.object({
	code: z.string().min(1).optional(),
	error: z.string().min(1).optional(),
	state: z.string().min(1).optional(),
})

const route = useRoute()
const completeSetupMutation = useMutationVKVideoBotSetupComplete()
const failureMessage = ref<string | null>(null)

function showFailure(message: string) {
	failureMessage.value = message
	toast.error(message)
}

onMounted(async () => {
	const callback = callbackQuerySchema.safeParse(route.query)
	if (!callback.success || callback.data.error || !callback.data.code || !callback.data.state) {
		showFailure('VK Video Live authorization was not completed')
		return
	}

	const result = await completeSetupMutation.executeMutation({
		code: callback.data.code,
		state: callback.data.state,
	})
	if (result.error || !result.data?.vkVideoBotSetupComplete) {
		showFailure(result.error?.message ?? 'Unable to complete VK Video Live bot setup')
		return
	}

	new BroadcastChannel(vkVideoBotSetupBroadcastChannelName).postMessage('refresh')
	toast.success('VK Video Live bot connected')
	window.close()
})
</script>

<template>
	<main class="flex min-h-[100dvh] items-center justify-center p-4">
		<Card class="w-full max-w-md">
			<CardHeader>
				<CardTitle>VK Video Live bot</CardTitle>
				<CardDescription>Completing authorization for the global bot account.</CardDescription>
			</CardHeader>
			<CardContent>
				<Alert v-if="failureMessage" variant="destructive">
					<AlertDescription>{{ failureMessage }}</AlertDescription>
				</Alert>
				<div v-else class="flex items-center gap-2 text-sm text-muted-foreground">
					<Icon name="lucide:loader-circle" class="size-4 animate-spin" />
					Completing authorization...
				</div>
			</CardContent>
		</Card>
	</main>
</template>
