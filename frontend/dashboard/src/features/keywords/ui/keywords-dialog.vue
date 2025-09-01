<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { InfoIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import type { KeywordResponse } from '@/api/keywords'

import { useKeywordsApi } from '@/api/keywords'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { useToast } from '@/components/ui/toast'
import VariableInput from '@/components/variable-input.vue'
import FormRolesSelector from '@/features/commands/ui/form-roles-selector.vue'

const props = defineProps<{
	keyword?: Omit<KeywordResponse, 'id'> & { id?: string }
}>()

const emits = defineEmits<{
	close: []
}>()

const { t } = useI18n()
const { toast } = useToast()

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
		}),
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
	{ immediate: true },
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
		toast({
			title: 'Error occured while saving keyword',
			variant: 'default',
		})
		console.error(e)
	}
})
</script>

<template>
	<Dialog
		v-model:open="open"
		@update:open="
			(state) => {
				if (!state && !keyword) resetFormValue()
			}
		"
	>
		<DialogTrigger as-child>
			<slot name="dialog-trigger" />
		</DialogTrigger>
		<DialogOrSheet class="sm:max-w-[424px]">
			<DialogHeader>
				<DialogTitle>
					{{ keyword ? t('keywords.edit') : t('keywords.create') }}
				</DialogTitle>
			</DialogHeader>

			<form class="flex flex-col gap-4" @submit="save">
				<FormField v-slot="{ componentField }" name="text">
					<FormItem>
						<FormLabel class="flex gap-2">
							{{ t('keywords.triggerText') }}
						</FormLabel>
						<FormControl>
							<Input v-bind="componentField" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="isRegularExpressionsRegular">
					<FormItem>
						<FormLabel class="flex gap-2">
							{{ t('keywords.isRegular') }}
						</FormLabel>
						<FormControl>
							<Switch :checked="componentField.modelValue" :update:checked="componentField['onUpdate:modelValue']" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
				<Alert>
					<InfoIcon class="h-4 w-4" />

					<AlertDescription>
						<i18n-t keypath="keywords.regularDescription">
							<a
								href="https://yourbasic.org/golang/regexp-cheat-sheet/#cheat-sheet"
								target="_blank"
								class="underline"
							>
								{{ t('keywords.regularDescriptionCheatSheet') }}
							</a>
						</i18n-t>
					</AlertDescription>
				</Alert>

				<FormField v-slot="{ componentField }" name="response">
					<FormItem>
						<FormLabel class="flex gap-2">
							{{ t('sharedTexts.response') }}
						</FormLabel>
						<FormControl>
							<div class="relative">
								<VariableInput :model-value="componentField.modelValue" input-type="textarea" @update:model-value="componentField['onUpdate:modelValue']" />
							</div>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="isReply">
					<FormItem>
						<FormLabel class="flex gap-2">
							{{ t('sharedTexts.reply.label') }}
						</FormLabel>
						<FormControl>
							<Switch :checked="componentField.modelValue" @update:checked="componentField['onUpdate:modelValue']" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<div class="flex flex-col gap-2">
					<Label class="flex gap-2">
						Roles
					</Label>
					<FormRolesSelector field-name="rolesIds" />
				</div>

				<FormField v-slot="{ componentField }" name="cooldown">
					<FormItem>
						<FormLabel class="flex gap-2">
							{{ t('keywords.cooldown') }}
						</FormLabel>
						<FormControl>
							<Input v-bind="componentField" type="number" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="usageCount">
					<FormItem>
						<FormLabel class="flex gap-2">
							{{ t('keywords.usages') }}
						</FormLabel>
						<FormControl>
							<Input v-bind="componentField" type="number" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
				<Button type="submit">
					{{ t('sharedButtons.save') }}
				</Button>
			</form>
		</DialogOrSheet>
	</Dialog>
</template>
