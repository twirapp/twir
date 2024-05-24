<script setup lang='ts'>
import { InfoIcon } from 'lucide-vue-next'
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NA,
	NForm,
	NFormItem,
	NSpace,
} from 'naive-ui'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import type { KeywordResponse } from '@/api/keywords'
import type { MakeOptional } from '@/gql/graphql'

import { useKeywordsApi } from '@/api/keywords'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import VariableInput from '@/components/variable-input.vue'

const props = defineProps<{
	keyword?: MakeOptional<KeywordResponse, 'id'> | null
}>()

const emits = defineEmits<{
	close: []
}>()

const open = ref(false)
const formRef = ref<FormInst | null>(null)

const defaultFormValue: MakeOptional<KeywordResponse, 'id'> = {
	text: '',
	usageCount: 0,
	cooldown: 0,
	enabled: true,
	isRegularExpression: false,
	isReply: true,
	response: null,
}

const formValue = ref<MakeOptional<KeywordResponse, 'id'>>(structuredClone(defaultFormValue))
function resetFormValue() {
	formValue.value = structuredClone(defaultFormValue)
}

watch(() => props.keyword, (keyword) => {
	if (!keyword) return
	formValue.value = structuredClone(toRaw(keyword))
}, { immediate: true })

const keywordsApi = useKeywordsApi()
const updateMutation = keywordsApi.useMutationUpdateKeyword()
const createMutation = keywordsApi.useMutationCreateKeyword()

async function save() {
	if (!formRef.value || !formValue.value) return
	await formRef.value.validate()

	const data = formValue.value
	delete data.id

	try {
		if (props.keyword?.id) {
			await updateMutation.executeMutation({
				id: props.keyword.id,
				opts: data,
			})
		} else {
			await createMutation.executeMutation({ opts: data })
		}
		emits('close')
		open.value = false
		resetFormValue()
	} catch (e) {
		console.error(e)
	}
}

const { t } = useI18n()

const rules: FormRules = {
	text: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error(t('keywords.validations.triggerRequired'))
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
			if (!state && !keyword) resetFormValue()
		}"
	>
		<DialogTrigger as-child>
			<slot name="dialog-trigger" />
		</DialogTrigger>
		<DialogOrSheet class="sm:max-w-[424px]">
			<DialogHeader>
				<DialogTitle>
					{{ keyword ? t('greetings.edit') : t('greetings.create') }}
				</DialogTitle>
			</DialogHeader>
			<NForm ref="formRef" :model="formValue" :rules="rules">
				<div class="grid gap-4 py-4">
					<NSpace vertical class="w-full">
						<NFormItem :label="t('keywords.triggerText')" path="text" show-require-mark>
							<Input v-model="formValue.text" />
						</NFormItem>

						<div class="flex flex-col gap-2 pb-4">
							<div class="flex justify-between">
								<span>{{ t('keywords.isRegular') }}</span>
								<Switch
									:checked="formValue.isRegularExpression"
									@update:checked="(v) => formValue.isRegularExpression = v"
								/>
							</div>
							<Alert>
								<InfoIcon class="h-4 w-4" />

								<AlertDescription>
									<i18n-t
										keypath="keywords.regularDescription"
									>
										<NA
											href="https://yourbasic.org/golang/regexp-cheat-sheet/#cheat-sheet"
											target="_blank"
										>
											{{ t('keywords.regularDescriptionCheatSheet') }}
										</NA>
									</i18n-t>
								</AlertDescription>
							</Alert>
						</div>

						<NFormItem :label="t('sharedTexts.response')" path="response">
							<VariableInput v-model="formValue.response" input-type="textarea" />
						</NFormItem>

						<div class="flex justify-between pb-4">
							<span>{{ t('sharedTexts.reply.text') }}</span>
							<Switch
								:checked="formValue.isReply"
								@update:checked="(v) => formValue.isReply = v"
							/>
						</div>

						<NFormItem :label="t('keywords.cooldown')">
							<Input v-model="formValue.cooldown" type="number" />
						</NFormItem>

						<NFormItem :label="t('keywords.usages')">
							<Input v-model="formValue.usageCount" type="number" />
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
