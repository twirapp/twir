<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'


import AlertsDialog from './alerts-dialog.vue'

import { type Alert,useAlertsApi } from '#layers/dashboard/api/alerts'



const props = defineProps<{ alert: Alert, withSelect: boolean }>()
const emits = defineEmits<{ 'update:select-alert': [alert: Alert] }>()

const { t } = useI18n()
const showDelete = ref(false)

const manager = useAlertsApi()
const alertsDeleteMutation = manager.useAlertsDeleteMutation()

function deleteAlert() {
	alertsDeleteMutation.executeMutation({ id: props.alert.id })
}
</script>

<template>
	<template v-if="withSelect">
		<UiButton @click="emits('update:select-alert', alert)">
			{{ t('sharedButtons.select') }}
		</UiButton>
	</template>
	<template v-else>
		<div class="flex justify-end items-center gap-2">
			<AlertsDialog :alert="alert">
				<template #dialog-trigger>
					<UiButton variant="secondary" size="icon">
						<PencilIcon class="h-4 w-4" />
					</UiButton>
				</template>
			</AlertsDialog>

			<UiButton variant="destructive" size="icon" @click="showDelete = true">
				<TrashIcon class="h-4 w-4" />
			</UiButton>
		</div>

		<UiActionConfirm v-model:open="showDelete" @confirm="deleteAlert" />
	</template>
</template>
