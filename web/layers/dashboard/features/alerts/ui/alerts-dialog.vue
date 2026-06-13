<script setup lang="ts">

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
import { ScrollArea } from '@/components/ui/scroll-area'

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
			<DialogHeader class="p-6 border-b">
				<DialogTitle>
					{{ alert ? t('alerts.editAlert') : t('alerts.createAlert') }}
				</DialogTitle>
			</DialogHeader>

			<ScrollArea class="max-h-[85vh]">
				<AlertsModalContent
					:alert="alert"
					@close="() => open = false"
				/>
			</ScrollArea>
		</DialogOrSheet>
	</Dialog>
</template>
