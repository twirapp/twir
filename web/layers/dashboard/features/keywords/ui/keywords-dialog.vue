<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { InfoIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { ref, toRaw, watch } from 'vue'

import * as z from 'zod'

import type { KeywordResponse } from '#layers/dashboard/api/keywords'

import { useKeywordsApi } from '#layers/dashboard/api/keywords'
import DialogOrSheet from '#layers/dashboard/components/dialog-or-sheet.vue'








import { toast } from 'vue-sonner'
import VariableInput from '#layers/dashboard/components/variable-input.vue'
import FormRolesSelector from '~/features/commands/ui/form-roles-selector.vue'

const props = defineProps<{
	keyword?: Omit<KeywordResponse, 'id'> & { id?: string }
}>()

const emits = defineEmits<{
	close: []
}>()

const { t } = useI18n()

const open = ref(false)

const keywordsForm = useForm({
	validationSchema: toTypedSchema(
		z.object({
			id: z.string().optional(),
			text: z.string().min(1),
			response: z.string().optional().nullable(),
			isRegularExpression: z.boolean(),
			isReply: z.boolean(),
			cooldown: z.number().min(0).optional(),
			usageCount: z.number().min(0).optional(),
			rolesIds: z.array(z.string()).optional(),
			enabled: z.boolean().optional().default(true),
		})
	),
	initialValues: {
		text: '',
		usageCount: 0,
		cooldown: 0,
		enabled: true,
		isRegularExpression: false,
		isReply: true,
		response: null,
		rolesIds: [],
	},
	keepValuesOnUnmount: true,
})

function resetFormValue() {
	keywordsForm.resetForm({
		values: {
			text: '',
			usageCount: 0,
			cooldown: 0,
			enabled: true,
			isRegularExpression: false,
			isReply: true,
			response: null,
			rolesIds: [],
		},
	})
}

watch(
	() => props.keyword,
	(k) => {
		console.log(k)
		if (!k) return

		keywordsForm.setValues(structuredClone(toRaw(k)))
	},
	{ immediate: true }
)

const keywordsApi = useKeywordsApi()
const updateMutation = keywordsApi.useMutationUpdateKeyword()
const createMutation = keywordsApi.useMutationCreateKeyword()

const save = keywordsForm.handleSubmit(async (values) => {
	try {
		if (props.keyword?.id) {
			delete values.id

			await updateMutation.executeMutation({
				id: props.keyword.id,
				opts: values,
			})
		} else {
			await createMutation.executeMutation({ opts: values })
		}
		emits('close')
		open.value = false
		resetFormValue()
	} catch (e) {
		toast.error('Error occured while saving keyword')
		console.error(e)
	}
})
</script>

<template>
	<UiDialog
		v-model:open="open"
		@update:open="
			(state) => {
				if (!state && !keyword) resetFormValue()
			}
		"
	>
		<UiDialogTrigger as-child>
			<slot name="dialog-trigger" />
		</UiDialogTrigger>
		<DialogOrSheet class="sm:max-w-[424px]">
			<UiDialogHeader>
				<UiDialogTitle>
					{{ keyword ? t('keywords.edit') : t('keywords.create') }}
				</UiDialogTitle>
			</UiDialogHeader>

			<form class="flex flex-col gap-4" @submit="save">
				<UiFormField v-slot="{ componentField }" name="text">
					<UiFormItem>
						<UiFormLabel class="flex gap-2">
							{{ t('keywords.triggerText') }}
						</UiFormLabel>
						<UiFormControl>
							<UiTextarea v-bind="componentField" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ field }" name="isRegularExpression">
					<UiFormItem>
						<UiFormLabel class="flex gap-2">
							{{ t('keywords.isRegular') }}
						</UiFormLabel>
						<UiFormControl>
							<UiSwitch
								:model-value="field.value"
								@update:model-value="field['onUpdate:modelValue']"
							/>
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
				<UiAlert>
					<InfoIcon class="h-4 w-4" />

					<UiAlertDescription>
						<i18n-t keypath="keywords.regularDescription">
							<a
								href="https://yourbasic.org/golang/regexp-cheat-sheet/#cheat-sheet"
								target="_blank"
								class="underline"
							>
								{{ t('keywords.regularDescriptionCheatSheet') }}
							</a>
						</i18n-t>
					</UiAlertDescription>
				</UiAlert>

				<UiFormField v-slot="{ componentField }" name="response">
					<UiFormItem>
						<UiFormLabel class="flex gap-2">
							{{ t('sharedTexts.response') }}
						</UiFormLabel>
						<UiFormControl>
							<div class="relative">
								<VariableInput
									:model-value="componentField.modelValue"
									input-type="textarea"
									@update:model-value="componentField['onUpdate:modelValue']"
								/>
							</div>
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ field }" name="isReply">
					<UiFormItem>
						<UiFormLabel class="flex gap-2">
							{{ t('sharedTexts.reply.label') }}
						</UiFormLabel>
						<UiFormControl>
							<UiSwitch
								:model-value="field.value"
								@update:model-value="field['onUpdate:modelValue']"
							/>
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<div class="flex flex-col gap-2">
					<UiLabel class="flex gap-2"> Roles </UiLabel>
					<FormRolesSelector class="xl:w-full xl:max-w-full" field-name="rolesIds" />
				</div>

				<UiFormField v-slot="{ componentField }" name="cooldown">
					<UiFormItem>
						<UiFormLabel class="flex gap-2">
							{{ t('keywords.cooldown') }}
						</UiFormLabel>
						<UiFormControl>
							<UiInput v-bind="componentField" type="number" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ componentField }" name="usageCount">
					<UiFormItem>
						<UiFormLabel class="flex gap-2">
							{{ t('keywords.usages') }}
						</UiFormLabel>
						<UiFormControl>
							<UiInput v-bind="componentField" type="number" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
				<UiButton type="submit">
					{{ t('sharedButtons.save') }}
				</UiButton>
			</form>
		</DialogOrSheet>
	</UiDialog>
</template>
