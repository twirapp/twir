<script setup lang='ts'>
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NForm,
	NFormItem,
	NSpace,
} from 'naive-ui'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import type { CustomVariable } from '@/api/variables'
import type { VariableCreateInput } from '@/gql/graphql'

import { useVariablesApi } from '@/api/variables'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { VariableType } from '@/gql/graphql'

const props = defineProps<{
	variable?: CustomVariable | null
}>()

const emits = defineEmits<{
	close: []
}>()

const open = ref(false)
const formRef = ref<FormInst | null>(null)
const defaultFormValue: CustomVariable = {
	id: '',
	evalValue: `// semicolons (;) matters, do not forget put them on end of statements.
const request = await fetch('https://jsonplaceholder.typicode.com/todos/1');
const response = await request.json();
// you should return value from your script
return response.title;`,
	name: '',
	description: null,
	type: VariableType.Text,
	response: '',
}

const formValue = ref(structuredClone(defaultFormValue))
function resetFormValue() {
	formValue.value = structuredClone(defaultFormValue)
}

watch(() => props.variable, (variable) => {
	if (!variable) return
	formValue.value = structuredClone(toRaw(variable))
}, { immediate: true })

const variablesApi = useVariablesApi()
const variablesUpdate = variablesApi.useMutationUpdateVariable()
const variablesCreate = variablesApi.useMutationCreateVariable()

async function save() {
	if (!formRef.value || !formValue.value) return
	await formRef.value.validate()

	const data = formValue.value
	const opts: VariableCreateInput = {
		description: data.description,
		evalValue: data.evalValue,
		name: data.name,
		type: data.type,
		response: data.response,
	}

	try {
		if (data.id) {
			await variablesUpdate.executeMutation({
				id: data.id,
				opts,
			})
		} else {
			await variablesCreate.executeMutation({ opts })
		}
		emits('close')
		open.value = false
	} catch (e) {
		console.error(e)
	}
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
			if (!state && !variable) resetFormValue()
		}"
	>
		<DialogTrigger as-child>
			<slot name="dialog-trigger" />
		</DialogTrigger>
		<DialogOrSheet class="min-w-[40%]">
			<DialogHeader>
				<DialogTitle>
					{{ variable ? t('greetings.edit') : t('greetings.create') }}
				</DialogTitle>
			</DialogHeader>
			<NForm ref="formRef" :model="formValue" :rules="rules">
				<div class="grid gap-4 py-4">
					<NSpace vertical class="w-full">
						<NFormItem :label="t('sharedTexts.name')" path="name" show-require-mark>
							<Input v-model="formValue.name" />
						</NFormItem>

						<NFormItem :label="t('variables.type')" path="type" show-require-mark>
							<Select v-model="formValue.type">
								<SelectTrigger>
									<SelectValue placeholder="Select a type" />
								</SelectTrigger>
								<SelectContent>
									<SelectItem :value="VariableType.Text">
										Text
									</SelectItem>
									<SelectItem :value="VariableType.Number">
										Number
									</SelectItem>
									<SelectItem :value="VariableType.Script">
										Script
									</SelectItem>
								</SelectContent>
							</Select>
						</NFormItem>

						<NFormItem v-if="formValue.type !== VariableType.Script" :label="t('sharedTexts.response')" path="response" show-require-mark>
							<Textarea v-model="formValue.response" />
						</NFormItem>

						<vue-monaco-editor
							v-else
							v-model:value="formValue.evalValue"
							theme="vs-dark"
							height="500px"
							language="javascript"
						/>
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
