<script setup lang="ts">
import { toast } from 'vue-sonner'
import { useRoute } from 'vue-router'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert/index.js'
import { Button } from '@/components/ui/button/index.js'
import { useStreamElementsIntegration } from './composables/use-streamelements-integration.js'

const { t } = useI18n()
const route = useRoute()
const streamElements = useStreamElementsIntegration()
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
		const result = await streamElements.postCode.executeMutation({ input: { code, state } })
		if (result.error || !result.data?.streamelementsPostCode) {
			fail('imports.errors.callbackRequest')
			return
		}
		streamElements.broadcastRefresh()
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
			<AlertTitle>{{ t('imports.errors.callbackTitle', { provider: 'StreamElements' }) }}</AlertTitle>
			<AlertDescription class="flex flex-col gap-4">
				<p>{{ errorMessage }}</p>
				<Button variant="outline" @click="closeWindow">{{ t('imports.actions.close') }}</Button>
			</AlertDescription>
		</Alert>
		<div v-else class="flex items-center gap-3 text-muted-foreground">
			<Icon name="lucide:loader-circle" class="size-6 animate-spin" />
			<span>{{ t('imports.callback.connecting', { provider: 'StreamElements' }) }}</span>
		</div>
	</div>
</template>
