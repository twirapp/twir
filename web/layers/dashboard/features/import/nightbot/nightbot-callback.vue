<script setup lang="ts">
import { toast } from 'vue-sonner'
import { useRoute } from 'vue-router'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert/index.js'
import { Button } from '@/components/ui/button/index.js'
import { useNightbotIntegration } from './composables/use-nightbot-integration.js'

const { t } = useI18n()
const route = useRoute()
const nightbot = useNightbotIntegration()
const errorMessage = ref('')

function fail(messageKey: string) {
	errorMessage.value = t(messageKey)
	toast.error(errorMessage.value)
}

function closeWindow() {
	window.close()
}

onMounted(async () => {
	if (typeof route.query.error === 'string') {
		fail('imports.errors.callbackRejected')
		return
	}

	const { code, state } = route.query
	if (typeof code !== 'string' || typeof state !== 'string') {
		fail('imports.errors.callbackMissing')
		return
	}

	try {
		const result = await nightbot.postCode.executeMutation({ input: { code, state } })
		if (result.error || !result.data?.nightbotPostCode) {
			fail('imports.errors.callbackRequest')
			return
		}
		nightbot.broadcastRefresh()
		closeWindow()
	} catch {
		fail('imports.errors.callbackRequest')
	}
})
</script>

<template>
	<div class="flex h-full items-center justify-center p-6">
		<Alert v-if="errorMessage" variant="destructive" class="max-w-md">
			<Icon name="lucide:circle-alert" />
			<AlertTitle>{{ t('imports.errors.callbackTitle', { provider: 'Nightbot' }) }}</AlertTitle>
			<AlertDescription class="flex flex-col gap-4">
				<p>{{ errorMessage }}</p>
				<Button variant="outline" @click="closeWindow">{{ t('imports.actions.close') }}</Button>
			</AlertDescription>
		</Alert>
		<div v-else class="flex items-center gap-3 text-muted-foreground">
			<Icon name="lucide:loader-circle" class="size-6 animate-spin" />
			<span>{{ t('imports.callback.connecting', { provider: 'Nightbot' }) }}</span>
		</div>
	</div>
</template>
