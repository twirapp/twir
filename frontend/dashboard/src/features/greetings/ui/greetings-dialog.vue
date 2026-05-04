<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';

import {
	type GreetingsCreateInputInput,
	GreetingsCreateInputSchema,
	type GreetingsUpdateInputInput,
	GreetingsUpdateInputSchema,
} from '@/gql/validation-schemas.js';
import { type Greetings, useGreetingsApi } from '@/api/greetings';
import DialogOrSheet from '@/components/dialog-or-sheet.vue';
import TwitchUserSelect from '@/components/twitchUsers/twitch-user-select.vue';
import { Button } from '@/components/ui/button';
import {
	Dialog,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog';
import VariableInput from '@/components/variable-input.vue';
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { Switch } from '@/components/ui/switch';
import { toast } from 'vue-sonner';

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
		userId: props.greeting?.twitchProfile?.id ?? props.greeting?.userId ?? '',
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
	_values: GreetingsCreateInputInput | GreetingsUpdateInputInput
): _values is GreetingsUpdateInputInput {
	return !!props.greeting?.id
}

const onSubmit = greetingForm.handleSubmit(async (values) => {
	try {
		let e: unknown
		if (isUpdate(values)) {
			const { error } = await greetingsUpdate.executeMutation({
				id: props.greeting!.id,
				opts: values,
			})
			e = error
		} else {
			const { error } = await greetingsCreate.executeMutation({ opts: values })
			e = error
		}

		if (!e) {
			toast.success('Saved', {
				duration: 2500,
			})
			emits('close')
			open.value = false
		}
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
	<Dialog v-model:open="open">
		<DialogTrigger as-child>
			<slot name="dialog-trigger" />
		</DialogTrigger>
		<DialogOrSheet class="sm:max-w-[424px]">
			<DialogHeader>
				<DialogTitle>
					{{ greeting ? t('greetings.edit') : t('greetings.create') }}
				</DialogTitle>
			</DialogHeader>
			<form @submit="onSubmit" class="grid gap-4 py-4">
				<FormField v-slot="{ componentField }" name="userId">
					<FormItem>
						<FormLabel>{{ t('sharedTexts.userName') }}</FormLabel>
						<FormControl>
							<TwitchUserSelect
								v-model="componentField.modelValue"
								@update:model-value="componentField['onUpdate:modelValue']"
								twir-only
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<div class="relative">
					<FormField v-slot="{ componentField }" name="text">
						<FormItem>
							<FormLabel>{{ t('sharedTexts.response') }}</FormLabel>
							<FormControl>
								<VariableInput v-bind="componentField" input-type="textarea" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>

				<FormField v-slot="{ value, handleChange }" name="withShoutOut">
					<FormItem class="flex justify-between items-center flex-wrap">
						<FormLabel>Send shoutout with greeting</FormLabel>
						<FormControl>
							<Switch :model-value="value" @update:model-value="handleChange" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ value, handleChange }" name="isReply">
					<FormItem class="flex justify-between items-center flex-wrap">
						<FormLabel>{{ t('sharedTexts.reply.text') }}</FormLabel>
						<FormControl>
							<Switch :model-value="value" @update:model-value="handleChange" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
				<DialogFooter>
					<Button type="submit">
						{{ t('sharedButtons.save') }}
					</Button>
				</DialogFooter>
			</form>
		</DialogOrSheet>
	</Dialog>
</template>
