<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'

defineProps<{
	isNewEvent: boolean
	onCancel: () => void
}>()

const { t } = useI18n()
</script>

<template>
	<div class="flex justify-end gap-4">
		<Button variant="outline" type="button" @click="onCancel">
			{{ t('sharedTexts.cancel') }}
		</Button>

		<Button type="submit">
			{{ isNewEvent ? t('sharedTexts.create') : t('sharedTexts.save') }}
		</Button>
	</div>
</template>
