<script setup lang="ts">
import CodeTool from '@editorjs/code'
import Delimiter from '@editorjs/delimiter'
import EditorJS from '@editorjs/editorjs'
import Header from '@editorjs/header'
import InlineCode from '@editorjs/inline-code'
import List from '@editorjs/list'
import Paragraph from '@editorjs/paragraph'
import Quote from '@editorjs/quote'
import SimpleImage from '@editorjs/simple-image'
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import './notifications-form.css'

import { useNotificationsForm } from '../composables/use-notifications-form.js'
import { textareaButtons, useTextarea } from '../composables/use-textarea.js'

import TwitchUsersSelect from '@/components/twitchUsers/twitch-users-select.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

const { t } = useI18n()

const notificationsForm = useNotificationsForm()
const { textareaRef, applyModifier } = useTextarea()

const editor = ref<EditorJS | null>(null)
const editorHtmlElement = ref<HTMLElement | null>(null)

const d = ref()

onMounted(() => {
	if (!editorHtmlElement.value) return

	editor.value = new EditorJS({
		holder: editorHtmlElement.value,
		onChange(api, event) {
			editor.value?.save().then((outputData) => {
				d.value = outputData
			})
		},
		tools: {
			header: Header,
			image: SimpleImage,
			list: {
				class: List,
				inlineToolbar: true,
				config: {
					defaultStyle: 'unordered',
				},
			},
			inlineCode: {
				class: InlineCode,
				shortcut: 'CMD+SHIFT+M',
			},
			delimiter: Delimiter,
			code: CodeTool,
			paragraph: {
				class: Paragraph,
				inlineToolbar: true,
			},
			quote: Quote,
		},
	})
})
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
							v-model="notificationsForm.messageField.fieldModel.value"
							rows="8"
						/>
					</div>

					{{ d }}
					<div ref="editorHtmlElement" class="border border-border rounded-md"></div>
				</div>

				<template v-if="notificationsForm.formValues.value.message">
					<Label>{{ t('adminPanel.notifications.messagePreview') }}</Label>
					<div class="border rounded-md p-2" v-html="notificationsForm.formValues.value.message"></div>
				</template>
			</CardContent>

			<CardFooter class="flex justify-end gap-4">
				<Button
					:disabled="!notificationsForm.formValues.value.message && !notificationsForm.editableMessageId"
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
