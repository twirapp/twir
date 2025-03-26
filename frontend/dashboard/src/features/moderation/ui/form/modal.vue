<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { onMounted, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { z } from 'zod'

import { availableSettings, useEditableItem } from './helpers.ts'
import ModalCaps from './modal-caps.vue'
import ModalDenylist from './modal-denylist.vue'
import ModalEmotes from './modal-emotes.vue'
import ModalLanguage from './modal-language.vue'
import ModalLinks from './modal-links.vue'
import ModalLongMessage from './modal-longmessage.vue'
import ModalSymbols from './modal-symbols.vue'

import { useModerationManager } from '@/api'
import { Button } from '@/components/ui/button'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import { useToast } from '@/components/ui/toast'
import FormRolesSelector from '@/features/commands/ui/form-roles-selector.vue'
import { RoleTypeEnum } from '@/gql/graphql.ts'

const { t } = useI18n()

const manager = useModerationManager()
const updater = manager.update
const creator = manager.create

const { editableItem } = useEditableItem()

const validationSchema = toTypedSchema(z.object({
	id: z.string().optional(),
	banMessage: z.string().max(500),
	banTime: z.number().min(0).max(86400),
	warningMessage: z.string().max(500),
	maxWarnings: z.number().min(0).max(10),
	excludedRoles: z.array(z.string()),
	enabled: z.boolean(),
	checkClips: z.boolean(),
	triggerLength: z.number().min(0).max(10000),
	maxPercentage: z.number().min(0).max(100),
	denyList: z.array(z.string()),
	deniedChatLanguages: z.array(z.string()),
	type: z.string(),
}))

const moderationForm = useForm({
	validationSchema,
	keepValuesOnUnmount: false,
	validateOnMount: false,
})

onMounted(() => {
	if (editableItem.value?.id && editableItem.value?.data) {
		moderationForm.setFieldValue('id', editableItem.value.id)
		moderationForm.setValues(structuredClone(toRaw(editableItem.value.data)))
		return
	}

	const currentType = availableSettings.find(s => s.type === editableItem.value?.data?.type)

	moderationForm.setValues({
		type: currentType?.type,
		banMessage: currentType?.banMessage,
		banTime: currentType?.banTime,
		warningMessage: currentType?.warningMessage,
		maxWarnings: currentType?.maxWarnings,
		excludedRoles: currentType?.excludedRoles ?? [],
		enabled: currentType?.enabled,
		checkClips: currentType?.checkClips,
		triggerLength: currentType?.triggerLength,
		maxPercentage: currentType?.maxPercentage,
		denyList: currentType?.denyList ?? [],
		deniedChatLanguages: currentType?.deniedChatLanguages ?? [],
		id: editableItem.value?.id,
	})
})

const { toast } = useToast()

const handleSubmit = moderationForm.handleSubmit(async (values) => {
	if (!values?.id) {
		await creator.mutateAsync({
			data: values,
		})
	} else {
		await updater.mutateAsync({
			id: values.id,
			data: values,
		})
	}

	toast({
		title: t('sharedTexts.saved'),
		duration: 2000,
	})
})
</script>

<template>
	<form class="flex flex-col gap-3" @submit.prevent="handleSubmit">
		<FormField v-slot="{ field }" name="enabled">
			<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
				<div class="space-y-0.5">
					<FormLabel class="text-base">
						{{ t('sharedTexts.enabled') }}
					</FormLabel>
				</div>
				<FormControl>
					<Switch
						:checked="field.value"
						default-checked
						@update:checked="field['onUpdate:modelValue']"
					/>
				</FormControl>
			</FormItem>
		</FormField>

		<ModalSymbols
			v-if="editableItem?.data?.type === 'symbols'"
		/>

		<ModalLanguage
			v-if="editableItem?.data?.type === 'language'"
		/>

		<ModalLongMessage
			v-if="editableItem?.data?.type === 'long_message'"
		/>

		<ModalCaps
			v-if="editableItem?.data?.type === 'caps'"
		/>

		<ModalEmotes
			v-if="editableItem?.data?.type === 'emotes'"
		/>

		<ModalLinks
			v-if="editableItem?.data?.type === 'links'"
		/>

		<Separator label="Timeouts" />

		<FormField v-slot="{ componentField }" name="banMessage">
			<FormItem>
				<FormLabel>Timeout message</FormLabel>
				<FormControl>
					<Input type="text" v-bind="componentField" />
				</FormControl>
				<FormMessage />
			</FormItem>
		</FormField>

		<FormField v-slot="{ componentField }" name="banTime">
			<FormItem>
				<FormLabel>{{ t('moderation.banTime') }}</FormLabel>
				<FormControl>
					<Input type="number" v-bind="componentField" />
				</FormControl>
				<FormMessage />
				<FormDescription>
					{{ t('moderation.banDescription') }}
				</FormDescription>
			</FormItem>
		</FormField>

		<Separator label="Warnings" />

		<FormField v-slot="{ componentField }" name="warningMessage">
			<FormItem>
				<FormLabel>{{ t('moderation.warningMessage') }}</FormLabel>
				<FormControl>
					<Input type="text" v-bind="componentField" />
				</FormControl>
				<FormMessage />
			</FormItem>
		</FormField>

		<FormField v-slot="{ componentField }" name="maxWarnings">
			<FormItem>
				<FormLabel>{{ t('moderation.warningMaxCount') }}</FormLabel>
				<FormControl>
					<Input type="number" v-bind="componentField" />
				</FormControl>
				<FormMessage />
			</FormItem>
		</FormField>

		<Separator />

		<div class="flex flex-col gap-2">
			<Label>Select excluded from moderation roles:</Label>
			<FormRolesSelector
				field-name="excludedRoles"
				hide-everyone
				:excluded-types="[RoleTypeEnum.Broadcaster, RoleTypeEnum.Moderator]"
			/>
		</div>

		<ModalDenylist
			v-if="editableItem?.data?.type === 'deny_list'"
		/>

		<Button type="submit">
			{{ t('sharedButtons.save') }}
		</Button>
	</form>
</template>
