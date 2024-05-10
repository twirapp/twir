<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import AlertsDialog from './alerts-dialog.vue'

import { type Alert, useAlertsDeleteMutation } from '@/api/alerts.js'
import ActionConfirm from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'

const props = defineProps<{ alert: Alert, withSelect: boolean }>()
const emits = defineEmits<{ 'update:select-alert': [alert: Alert] }>()

const { t } = useI18n()
const showDelete = ref(false)

const alertsDeleteMutation = useAlertsDeleteMutation()

function deleteAlert() {
	alertsDeleteMutation.executeMutation({ id: props.alert.id })
}
</script>

<template>
	<template v-if="withSelect">
		<Button @click="emits('update:select-alert', alert)">
			{{ t('sharedButtons.select') }}
		</Button>
	</template>
	<template v-else>
		<div class="flex justify-end items-center gap-2">
			<AlertsDialog :alert="alert">
				<template #dialog-trigger>
					<Button variant="secondary" size="icon">
						<PencilIcon class="h-4 w-4" />
					</Button>
				</template>
			</AlertsDialog>

			<Button variant="destructive" size="icon" @click="showDelete = true">
				<TrashIcon class="h-4 w-4" />
			</Button>
		</div>

		<ActionConfirm v-model:open="showDelete" @confirm="deleteAlert" />
	</template>
</template>
