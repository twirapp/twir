<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { InfoIcon } from 'lucide-vue-next'
import { useField, useForm } from 'vee-validate'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import type { KeywordResponse } from '@/api/keywords'

import { useKeywordsApi } from '@/api/keywords'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Dialog, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { FormControl, FormField, FormItem, FormLabel, FormMessage, FormDescription } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Checkbox } from '@/components/ui/checkbox'
import { Textarea } from '@/components/ui/textarea'
import { toast } from 'vue-sonner'
import VariableInput from '@/components/variable-input.vue'
import FormRolesSelector from '@/features/commands/ui/form-roles-selector.vue'

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
			platforms: z.array(z.string()).default([]),
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
		platforms: [],
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
			platforms: [],
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

const { value: platforms, setValue: setPlatforms } = useField<string[]>('platforms')

function togglePlatform(platform: string, checked: boolean) {
	const current = platforms.value ?? []
	if (checked) {
		setPlatforms([...current, platform])
	} else {
		setPlatforms(current.filter((p) => p !== platform))
	}
}

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
							<Textarea v-bind="componentField" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ field }" name="isRegularExpression">
					<FormItem>
						<FormLabel class="flex gap-2">
							{{ t('keywords.isRegular') }}
						</FormLabel>
						<FormControl>
							<Switch
								:model-value="field.value"
								@update:model-value="field['onUpdate:modelValue']"
							/>
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
								<VariableInput
									:model-value="componentField.modelValue"
									input-type="textarea"
									@update:model-value="componentField['onUpdate:modelValue']"
								/>
							</div>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ field }" name="isReply">
					<FormItem>
						<FormLabel class="flex gap-2">
							{{ t('sharedTexts.reply.label') }}
						</FormLabel>
						<FormControl>
							<Switch
								:model-value="field.value"
								@update:model-value="field['onUpdate:modelValue']"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<div class="flex flex-col gap-2">
					<Label class="flex gap-2"> Roles </Label>
					<FormRolesSelector class="xl:w-full xl:max-w-full" field-name="rolesIds" />
				</div>

				<div class="flex flex-col gap-2 mt-2">
					<FormLabel>Platforms</FormLabel>
					<FormDescription class="mb-2">
						Select which platforms this keyword runs on. If none selected, it runs on all platforms.
					</FormDescription>
					<div class="flex gap-4">
						<FormItem class="flex items-center gap-2 space-y-0">
							<Checkbox
								:model-value="platforms?.includes('twitch')"
								@update:model-value="(checked) => togglePlatform('twitch', !!checked)"
							/>
							<FormLabel class="font-normal">Twitch</FormLabel>
						</FormItem>
						<FormItem class="flex items-center gap-2 space-y-0">
							<Checkbox
								:model-value="platforms?.includes('kick')"
								@update:model-value="(checked) => togglePlatform('kick', !!checked)"
							/>
							<FormLabel class="font-normal">Kick</FormLabel>
						</FormItem>
					</div>
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
