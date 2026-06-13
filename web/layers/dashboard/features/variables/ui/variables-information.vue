<script setup lang="ts">
import { InfoIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { Alert, AlertDescription } from '@/components/ui/alert'

const { t } = useI18n()
</script>

<template>
	<Alert>
		<InfoIcon class="size-5" />
		<AlertDescription>
			{{ t('variables.info') }}
		</AlertDescription>
	</Alert>
</template>
