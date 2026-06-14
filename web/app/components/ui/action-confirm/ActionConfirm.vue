<script lang="ts" setup>
import { useI18n } from 'vue-i18n'

import {
	AlertDialog,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'

defineProps<{
	confirmText?: string
}>()

const emit = defineEmits<{
	cancel: []
	confirm: []
}>()

const { t } = useI18n()

const open = defineModel<boolean>('open', { default: false })

function handleCancel() {
	emit('cancel')
	open.value = false
}

function handleConfirm() {
	emit('confirm')
	open.value = false
}
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
				<Button type="button" variant="outline" class="mt-2 sm:mt-0" @click="handleCancel">
					{{ t('deleteConfirmation.cancel') }}
				</Button>
				<Button type="button" @click="handleConfirm">
					{{ t('deleteConfirmation.confirm') }}
				</Button>
			</AlertDialogFooter>
		</AlertDialogContent>
	</AlertDialog>
</template>
