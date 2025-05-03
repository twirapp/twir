<script setup lang="ts">
import { InfoIcon } from 'lucide-vue-next'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
</script>

<template>
	<Alert>
		<InfoIcon class="h-4 w-4" />
		<AlertTitle>{{ $t('events.information.title') }}</AlertTitle>
		<AlertDescription>
			{{ $t('events.information.description') }}
		</AlertDescription>
	</Alert>
</template>
