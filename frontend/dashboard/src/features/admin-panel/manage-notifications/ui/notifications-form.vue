<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { useNotificationsForm } from '../composables/use-notifications-form.js'

import TwitchUsersSelect from '@/components/twitchUsers/twitch-users-select.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import BlocksRender from '@/components/ui/editorjs/blocks-render.vue'
import EditorJS from '@/components/ui/editorjs/editorjs.vue'
import { Label } from '@/components/ui/label'

const { t } = useI18n()

const notificationsForm = useNotificationsForm()
</script>

<template>
	<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
		{{ t('adminPanel.notifications.formTitle') }}
	</h4>

	<Card>
		<form class="flex flex-col gap-4" @submit="notificationsForm.onSubmit">
			<CardContent class="flex flex-col gap-4 p-4">
				<div class="space-y-2">
					<Label for="userId">
						{{ t('adminPanel.notifications.userLabel') }}
					</Label>
					<TwitchUsersSelect
						v-model="notificationsForm.userIdField.fieldModel.value"
						twir-only
					/>
				</div>

				<div class="space-y-2">
					<Label for="message">
						{{ t('adminPanel.notifications.messageLabel') }}
					</Label>

					<EditorJS v-model:model-value="notificationsForm.editorJsJsonField.fieldModel.value" />
				</div>

				<template v-if="notificationsForm.editorJsJsonField.fieldModel.value">
					<Label>{{ t('adminPanel.notifications.messagePreview') }}</Label>
					<BlocksRender :data="notificationsForm.editorJsJsonField.fieldModel.value" />
				</template>
			</CardContent>

			<CardFooter class="flex justify-end gap-4">
				<Button
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
				</Button>
				<Button type="submit">
					<template v-if="notificationsForm.editableMessageId.value">
						{{ t('sharedButtons.edit') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.send') }}
					</template>
				</Button>
			</CardFooter>
		</form>
	</Card>
</template>
