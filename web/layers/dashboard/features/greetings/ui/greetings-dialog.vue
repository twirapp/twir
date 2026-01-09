<script setup lang="ts">
import { ref } from 'vue'

import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'

import {
	type GreetingsCreateInputInput,
	GreetingsCreateInputSchema,
	type GreetingsUpdateInputInput,
	GreetingsUpdateInputSchema,
} from '~/gql/validation-schemas.js'
import { type Greetings, useGreetingsApi } from '#layers/dashboard/api/greetings'
import DialogOrSheet from '#layers/dashboard/components/dialog-or-sheet.vue'
import TwitchUserSelect from '#layers/dashboard/components/twitchUsers/twitch-user-select.vue'


import VariableInput from '#layers/dashboard/components/variable-input.vue'


import { toast } from 'vue-sonner'

const props = defineProps<{
	greeting?: Greetings | null
}>()

const emits = defineEmits<{
	close: []
}>()

const greetingForm = useForm({
	validationSchema: toTypedSchema(
		props.greeting ? GreetingsUpdateInputSchema : GreetingsCreateInputSchema
	),
	keepValuesOnUnmount: true,
	validateOnMount: false,
	initialValues: {
		userId: props.greeting?.userId ?? '',
		text: props.greeting?.text,
		enabled: props.greeting?.enabled ?? true,
		isReply: props.greeting?.isReply ?? true,
		withShoutOut: props.greeting?.withShoutOut ?? false,
	},
})

const open = ref(false)

const greetingsApi = useGreetingsApi()
const greetingsUpdate = greetingsApi.useMutationUpdateGreetings()
const greetingsCreate = greetingsApi.useMutationCreateGreetings()

function isUpdate(
	values: GreetingsCreateInputInput | GreetingsUpdateInputInput
): values is GreetingsUpdateInputInput {
	return !props.greeting?.id && Object.values(values).some((v) => v === undefined)
}

const onSubmit = greetingForm.handleSubmit(async (values) => {
	try {
		if (isUpdate(values)) {
			await greetingsUpdate.executeMutation({
				id: props.greeting!.id,
				opts: values,
			})
		} else {
			await greetingsCreate.executeMutation({ opts: values })
		}
		emits('close')
		open.value = false

		toast.success('Saved', {
			duration: 2500,
		})
	} catch (e) {
		console.error(e)

		if ('message' in (e as Error)) {
			toast.error(`Error ${(e as Error).message}`)
		}
	}
})

const { t } = useI18n()
</script>

<template>
	<UiDialog v-model:open="open">
		<UiDialogTrigger as-child>
			<slot name="dialog-trigger" />
		</UiDialogTrigger>
		<DialogOrSheet class="sm:max-w-[424px]">
			<UiDialogHeader>
				<UiDialogTitle>
					{{ greeting ? t('greetings.edit') : t('greetings.create') }}
				</UiDialogTitle>
			</UiDialogHeader>
			<form @submit="onSubmit" class="grid gap-4 py-4">
				<UiFormField v-slot="{ componentField }" name="userId">
					<UiFormItem>
						<UiFormLabel>{{ t('sharedTexts.userName') }}</UiFormLabel>
						<UiFormControl>
							<TwitchUserSelect
								v-model="componentField.modelValue"
								@update:model-value="componentField['onUpdate:modelValue']"
								twir-only
							/>
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<div class="relative">
					<UiFormField v-slot="{ componentField }" name="text">
						<UiFormItem>
							<UiFormLabel>{{ t('sharedTexts.response') }}</UiFormLabel>
							<UiFormControl>
								<VariableInput v-bind="componentField" input-type="textarea" />
							</UiFormControl>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>
				</div>

				<UiFormField v-slot="{ value, handleChange }" name="withShoutOut">
					<UiFormItem class="flex justify-between items-center flex-wrap">
						<UiFormLabel>Send shoutout with greeting</UiFormLabel>
						<UiFormControl>
							<UiSwitch :model-value="value" @update:model-value="handleChange" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ value, handleChange }" name="isReply">
					<UiFormItem class="flex justify-between items-center flex-wrap">
						<UiFormLabel>{{ t('sharedTexts.reply.text') }}</UiFormLabel>
						<UiFormControl>
							<UiSwitch :model-value="value" @update:model-value="handleChange" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
				<UiDialogFooter>
					<UiButton type="submit">
						{{ t('sharedButtons.save') }}
					</UiButton>
				</UiDialogFooter>
			</form>
		</DialogOrSheet>
	</UiDialog>
</template>
