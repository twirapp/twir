<script setup lang='ts'>
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NForm,
	NFormItem,
	NSpace,
	NSwitch,
} from 'naive-ui'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import type { GreetingsCreateInput } from '@/gql/graphql'

import { type Greetings, useGreetingsApi } from '@/api/greetings'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import TwitchUsersSelect from '@/components/twitchUsers/twitch-users-select.vue'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import VariableInput from '@/components/variable-input.vue'

const props = defineProps<{
	greeting?: Greetings | null
}>()

const emits = defineEmits<{
	close: []
}>()

const open = ref(false)
const formRef = ref<FormInst | null>(null)
const defaultFormValue: Omit<Greetings, 'twitchProfile'> = {
	id: '',
	text: '',
	userId: '',
	enabled: true,
	isReply: true,
}

const formValue = ref(structuredClone(defaultFormValue))
function resetFormValue() {
	formValue.value = structuredClone(defaultFormValue)
}

watch(() => props.greeting, (greeting) => {
	if (!greeting) return
	formValue.value = structuredClone(toRaw(greeting))
}, { immediate: true })

const greetingsApi = useGreetingsApi()
const greetingsUpdate = greetingsApi.useMutationUpdateGreetings()
const greetingsCreate = greetingsApi.useMutationCreateGreetings()

async function save() {
	if (!formRef.value || !formValue.value) return
	await formRef.value.validate()

	const data = formValue.value
	const opts: GreetingsCreateInput = {
		enabled: data.enabled,
		isReply: data.isReply,
		text: data.text,
		userId: data.userId,
	}

	if (data.id) {
		await greetingsUpdate.executeMutation({
			id: data.id,
			opts,
		})
	} else {
		await greetingsCreate.executeMutation({ opts })
	}

	emits('close')
	open.value = false
}

const { t } = useI18n()

const rules: FormRules = {
	userId: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error(t('greetings.validations.userName'))
			}

			return true
		},
	},
	text: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error(t('greetings.validations.textRequired'))
			}

			return true
		},
	},
}
</script>

<template>
	<Dialog
		v-model:open="open"
		@update:open="(state) => {
			if (!state && !greeting) resetFormValue()
		}"
	>
		<DialogTrigger as-child>
			<slot name="dialog-trigger" />
		</DialogTrigger>
		<DialogOrSheet class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>
					{{ greeting ? t('greetings.edit') : t('greetings.create') }}
				</DialogTitle>
			</DialogHeader>
			<NForm ref="formRef" :model="formValue" :rules="rules">
				<div class="grid gap-4 py-4">
					<NSpace vertical class="w-full">
						<NFormItem :label="t('sharedTexts.userName')" path="userId" show-require-mark>
							<TwitchUsersSelect v-model="formValue.userId" :initial="formValue.userId" twir-only />
						</NFormItem>
						<NFormItem :label="t('sharedTexts.response')" path="text" show-require-mark>
							<VariableInput v-model="formValue.text" input-type="textarea" />
						</NFormItem>

						<NFormItem :label="t('sharedTexts.reply.text')" path="text">
							<NSwitch v-model:value="formValue.isReply" />
						</NFormItem>
					</NSpace>
				</div>
				<DialogFooter>
					<Button type="submit" @click="save">
						{{ t('sharedButtons.save') }}
					</Button>
				</DialogFooter>
			</NForm>
		</DialogOrSheet>
	</Dialog>
</template>
