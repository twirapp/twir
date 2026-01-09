<script setup lang="ts">
import { InfoIcon } from 'lucide-vue-next'




const { t } = useI18n()
</script>

<template>
	<UiAlert>
		<InfoIcon class="size-5" />
		<UiAlertDescription>
			{{ t('variables.info') }}
		</UiAlertDescription>
	</UiAlert>
</template>
