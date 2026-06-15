<script setup lang="ts">


import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

const { t } = useI18n()
</script>

<template>
	<Alert>
		<Icon name="lucide:info" class="size-5" />
		<AlertTitle>
			{{ t('keywords.title') }}
		</AlertTitle>
		<AlertDescription>
			{{ t('keywords.info.text') }}
		</AlertDescription>
	</Alert>
</template>
