<script setup lang="ts">
import { InfoIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

const { t } = useI18n()
</script>

<template>
	<Alert>
		<InfoIcon class="size-5" />
		<AlertTitle>
			{{ t('greetings.info.title') }}
		</AlertTitle>
		<AlertDescription>
			{{ t('greetings.info.text') }}
		</AlertDescription>
	</Alert>
</template>
