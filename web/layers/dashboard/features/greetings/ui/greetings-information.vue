<script setup lang="ts">
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

const { t } = useI18n()
</script>

<template>
	<Alert>
		<Icon
			name="lucide:info"
			class="size-5"
		/>
		<AlertTitle>
			{{ t('greetings.info.title') }}
		</AlertTitle>
		<AlertDescription>
			{{ t('greetings.info.text') }}
		</AlertDescription>
	</Alert>
</template>
