<script setup lang="ts">
import { NScrollbar } from 'naive-ui'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import AlertsModalContent from './alerts-dialog-content.vue'

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
	<Dialog v-model:open="open">
		<DialogTrigger as-child>
			<slot name="dialog-trigger" />
		</DialogTrigger>

		<DialogOrSheet class="p-0">
			<DialogHeader class="p-6 border-b-[1px]">
				<DialogTitle>
					{{ alert ? t('alerts.editAlert') : t('alerts.createAlert') }}
				</DialogTitle>
			</DialogHeader>

			<NScrollbar style="max-height: 85vh" trigger="none">
				<AlertsModalContent
					:alert="alert"
					@close="() => open = false"
				/>
			</NScrollbar>
		</DialogOrSheet>
	</Dialog>
</template>
