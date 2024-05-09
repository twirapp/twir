<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import type { Alert } from '@/api/alerts.js'

import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import {
	Dialog,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'

defineProps<{ alert?: Alert | null }>()

const { t } = useI18n()
const open = ref(false)
</script>

<template>
	<Dialog
		v-model:open="open"
		@update:open="(state) => {
			if (!state && !alert) {
				// reset form
			}
		}"
	>
		<DialogTrigger as-child>
			<slot name="dialog-trigger" />
		</DialogTrigger>
		<DialogOrSheet class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>
					{{ alert ? t('alerts.editAlert') : t('alerts.createAlert') }}
				</DialogTitle>
			</DialogHeader>

			TODO: form
		</DialogOrSheet>
	</Dialog>
</template>
