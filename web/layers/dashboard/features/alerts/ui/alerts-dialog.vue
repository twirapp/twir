<script setup lang="ts">

import { ref } from 'vue'


import AlertsModalContent from './alerts-dialog-content.vue'

import type { Alert } from '#layers/dashboard/api/alerts'

import DialogOrSheet from '#layers/dashboard/components/dialog-or-sheet.vue'



defineProps<{ alert?: Alert | null }>()

const { t } = useI18n()
const open = ref(false)
</script>

<template>
	<UiDialog v-model:open="open">
		<UiDialogTrigger as-child>
			<slot name="dialog-trigger" />
		</UiDialogTrigger>

		<DialogOrSheet class="p-0">
			<UiDialogHeader class="p-6 border-b">
				<UiDialogTitle>
					{{ alert ? t('alerts.editAlert') : t('alerts.createAlert') }}
				</UiDialogTitle>
			</UiDialogHeader>

			<UiScrollArea class="max-h-[85vh]">
				<AlertsModalContent
					:alert="alert"
					@close="() => open = false"
				/>
			</UiScrollArea>
		</DialogOrSheet>
	</UiDialog>
</template>
