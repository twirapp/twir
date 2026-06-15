<script setup lang="ts">
import { Alert, AlertDescription } from '@/components/ui/alert'

const { t } = useI18n()
</script>

<template>
	<Alert>
		<Icon
			name="lucide:info"
			class="size-5"
		/>
		<AlertDescription>
			{{ t('variables.info') }}
		</AlertDescription>
	</Alert>
</template>
