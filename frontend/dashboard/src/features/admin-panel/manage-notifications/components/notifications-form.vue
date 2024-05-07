<script setup lang="ts">
import { NCard } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import { useNotificationsForm } from '../composables/use-notifications-form.js'
import { textareaButtons, useTextarea } from '../composables/use-textarea.js'

import TwitchUserSingle from '@/components/twitchUsers/single.vue'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

const { t } = useI18n()

const notificationsForm = useNotificationsForm()
const { textareaRef, applyModifier } = useTextarea()
</script>

<template>
	<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
		{{ t('adminPanel.notifications.formTitle') }}
	</h4>
	<NCard size="small" bordered>
		<form class="flex flex-col gap-4" @submit="notificationsForm.onSubmit">
			<div class="space-y-2">
				<Label for="userId">
					{{ t('adminPanel.notifications.userLabel') }}
				</Label>
				<TwitchUserSingle v-model="notificationsForm.userIdField.fieldModel" twir-only />
			</div>

			<div class="space-y-2">
				<Label for="message">
					{{ t('adminPanel.notifications.messageLabel') }}
				</Label>

				<div class="flex flex-col gap-2">
					<div class="flex gap-2 flex-wrap">
						<TooltipProvider>
							<Tooltip v-for="button in textareaButtons" :key="button.name">
								<TooltipTrigger as-child>
									<Button
										type="button"
										variant="secondary"
										size="icon"
										@click="applyModifier(button.name)"
									>
										<component :is="button.icon" class="h-4 w-4" />
									</Button>
								</TooltipTrigger>
								<TooltipContent>
									<p>{{ button.title }}</p>
								</TooltipContent>
							</Tooltip>
						</TooltipProvider>
					</div>

					<Textarea
						ref="textareaRef"
						v-model="notificationsForm.messageField.fieldModel"
						rows="8"
					/>
				</div>
			</div>

			<template v-if="notificationsForm.formValues.message">
				<Label>{{ t('adminPanel.notifications.messagePreview') }}</Label>
				<div class="border rounded-md p-2" v-html="notificationsForm.formValues.message"></div>
			</template>

			<div class="flex justify-end gap-4">
				<Button
					:disabled="!notificationsForm.formValues.message && !notificationsForm.editableMessageId"
					type="button"
					variant="secondary"
					@click="notificationsForm.onReset"
				>
					<template v-if="notificationsForm.editableMessageId">
						{{ t('sharedButtons.cancel') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.reset') }}
					</template>
				</Button>
				<Button type="submit">
					<template v-if="notificationsForm.editableMessageId">
						{{ t('sharedButtons.edit') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.send') }}
					</template>
				</Button>
			</div>
		</form>
	</NCard>
</template>
