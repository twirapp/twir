<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';

const { t } = useI18n();

defineEmits<{
	cancel: [],
	confirm: [],
}>();

const open = defineModel<boolean>({ default: false });

defineProps<{
	confirmText?: string
}>()
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
