<script lang="ts" setup>
import { useI18n } from 'vue-i18n'

import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle
} from '@/components/ui/alert-dialog'

defineProps<{
	confirmText?: string
}>()

defineEmits<{
	cancel: []
	confirm: []
}>()

const { t } = useI18n()

const open = defineModel<boolean>('open', { default: false })
</script>

<template>
	<AlertDialog v-model:open="open">
		<AlertDialogContent>
			<AlertDialogHeader>
				<AlertDialogTitle>{{ confirmText ?? t('deleteConfirmation.text') }}</AlertDialogTitle>
				<AlertDialogDescription class="sr-only">
					{{ confirmText ?? t('deleteConfirmation.text') }}
				</AlertDialogDescription>
			</AlertDialogHeader>
			<AlertDialogFooter>
				<AlertDialogCancel @click="$emit('cancel')">
					{{ t('deleteConfirmation.cancel') }}
				</AlertDialogCancel>
				<AlertDialogAction @click="$emit('confirm')">
					{{ t('deleteConfirmation.confirm') }}
				</AlertDialogAction>
			</AlertDialogFooter>
		</AlertDialogContent>
	</AlertDialog>
</template>
