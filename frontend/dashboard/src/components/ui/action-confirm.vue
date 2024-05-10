<script lang="ts" setup>
import { useI18n } from 'vue-i18n'

import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
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

const open = defineModel<boolean>({ default: false })
</script>

<template>
	<AlertDialog v-model:open="open">
		<AlertDialogContent>
			<AlertDialogHeader>
				<AlertDialogTitle>{{ confirmText ?? t('deleteConfirmation.text') }}</AlertDialogTitle>
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
