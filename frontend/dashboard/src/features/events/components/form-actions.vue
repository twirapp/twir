<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'

const props = defineProps<{
	isNewEvent: boolean
	isSubmitting: boolean
	onCancel: () => void
}>()

const { t } = useI18n()
</script>

<template>
	<div class="flex justify-end gap-4">
		<Button variant="outline" type="button" @click="onCancel">
			{{ t('sharedTexts.cancel') }}
		</Button>

		<Button type="submit" :disabled="isSubmitting">
			{{ isSubmitting ? t('sharedTexts.saving') : (isNewEvent ? t('sharedTexts.create') : t('sharedTexts.save')) }}
		</Button>
	</div>
</template>
