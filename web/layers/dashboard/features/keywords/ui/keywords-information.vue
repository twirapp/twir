<script setup lang="ts">
import { InfoIcon } from 'lucide-vue-next'




const { t } = useI18n()
</script>

<template>
	<UiAlert>
		<InfoIcon class="size-5" />
		<UiAlertTitle>
			{{ t('keywords.title') }}
		</UiAlertTitle>
		<UiAlertDescription>
			{{ t('keywords.info.text') }}
		</UiAlertDescription>
	</UiAlert>
</template>
