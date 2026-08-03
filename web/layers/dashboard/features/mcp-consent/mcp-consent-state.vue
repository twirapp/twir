<script setup lang="ts">
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

defineProps<{
	state: 'loading' | 'expired' | 'permission' | 'network'
}>()

const emit = defineEmits<{
	(e: 'retry'): void
}>()

const { t } = useI18n()
</script>

<template>
	<Card v-if="state === 'loading'">
		<CardHeader>
			<Skeleton class="h-6 w-48" />
			<Skeleton class="h-4 w-full" />
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			<Skeleton class="h-20 w-full" />
			<Skeleton class="h-10 w-full" />
		</CardContent>
	</Card>

	<Alert v-else-if="state === 'expired'" variant="destructive">
		<Icon name="lucide:clock-alert" />
		<AlertTitle>{{ t('mcpConsent.errors.expired.title') }}</AlertTitle>
		<AlertDescription>{{ t('mcpConsent.errors.expired.description') }}</AlertDescription>
	</Alert>

	<Alert v-else-if="state === 'permission'" variant="destructive">
		<Icon name="lucide:shield-alert" />
		<AlertTitle>{{ t('mcpConsent.errors.permission.title') }}</AlertTitle>
		<AlertDescription>{{ t('mcpConsent.errors.permission.description') }}</AlertDescription>
	</Alert>

	<Alert v-else-if="state === 'network'" variant="destructive">
		<Icon name="lucide:triangle-alert" />
		<AlertTitle>{{ t('mcpConsent.errors.network.title') }}</AlertTitle>
		<AlertDescription class="flex flex-col gap-3">
			<span>{{ t('mcpConsent.errors.network.description') }}</span>
			<Button class="w-fit" type="button" variant="outline" @click="emit('retry')">
				{{ t('mcpConsent.retry') }}
			</Button>
		</AlertDescription>
	</Alert>
</template>
