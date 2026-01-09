<script setup lang="ts">


import { useNotificationsForm } from '../composables/use-notifications-form.js'

import TwitchUserSelect from '#layers/dashboard/components/twitchUsers/twitch-user-select.vue'






const { t } = useI18n()

const notificationsForm = useNotificationsForm()
</script>

<template>
	<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
		{{ t('adminPanel.notifications.formTitle') }}
	</h4>

	<UiCard>
		<form class="flex flex-col gap-4" @submit="notificationsForm.onSubmit">
			<UiCardContent class="flex flex-col gap-4 p-4">
				<div class="space-y-2">
					<UiLabel for="userId">
						{{ t('adminPanel.notifications.userLabel') }}
					</UiLabel>
					<TwitchUserSelect
						v-model="notificationsForm.userIdField.fieldModel.value as string"
						twir-only
					/>
				</div>

				<div class="space-y-2">
					<UiLabel for="message">
						{{ t('adminPanel.notifications.messageLabel') }}
					</UiLabel>

					<UiEditorJS
						v-model:model-value="notificationsForm.editorJsJsonField.fieldModel.value"
						@update:model-value="(v: string) => notificationsForm.editorJsJsonField.fieldModel.value = v"
					/>
				</div>

				<template v-if="notificationsForm.editorJsJsonField.fieldModel.value">
					<UiLabel>{{ t('adminPanel.notifications.messagePreview') }}</UiLabel>
					<UiBlocksRender :data="notificationsForm.editorJsJsonField.fieldModel.value" />
				</template>
			</UiCardContent>

			<UiCardFooter class="flex justify-end gap-4">
				<UiButton
					:disabled="!notificationsForm.editorJsJsonField.fieldModel.value && !notificationsForm.editableMessageId"
					type="button"
					variant="secondary"
					@click="notificationsForm.onReset"
				>
					<template v-if="notificationsForm.editableMessageId.value">
						{{ t('sharedButtons.cancel') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.reset') }}
					</template>
				</UiButton>
				<UiButton type="submit">
					<template v-if="notificationsForm.editableMessageId.value">
						{{ t('sharedButtons.edit') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.send') }}
					</template>
				</UiButton>
			</UiCardFooter>
		</form>
	</UiCard>
</template>
