<script setup lang="ts">
import { useField, useResetForm, useSetFormValues, useSubmitForm } from 'vee-validate'
import { onMounted, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'

import ModalCaps from './modal-caps.vue'
import ModalDenylist from './modal-denylist.vue'
import ModalEmotes from './modal-emotes.vue'
import ModalLanguage from './modal-language.vue'
import ModalLinks from './modal-links.vue'
import ModalLongMessage from './modal-longmessage.vue'
import ModalSymbols from './modal-symbols.vue'

import type { EditableItem } from '@/features/moderation/composables/use-moderation-form.ts'

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
import FormRolesSelector from '@/features/commands/ui/form-roles-selector.vue'
import { useModerationApi } from '@/features/moderation/composables/use-moderation-api.ts'
import { ModerationSettingsType , RoleTypeEnum } from '@/gql/graphql.ts'

const { t } = useI18n()

const manager = useModerationApi()

const formReset = useResetForm()
const { value: currentItemId, setValue: setCurrentItemId } = useField<string | undefined>('id')
const setFormValues = useSetFormValues<EditableItem>()
const { value: currentEditType } = useField<ModerationSettingsType>('type')

onMounted(() => {
	formReset()

	if (currentItemId.value) {
		const item = manager.items.value.find(i => i.id === currentItemId.value)
		if (!item) return

		const values = structuredClone(toRaw(item))

		setFormValues(values)
	}
})

const handleSubmit = useSubmitForm<EditableItem>(async (values) => {
	if (!currentItemId.value) {
		const newRule = await manager.create(values)
		if (newRule?.id) {
			setCurrentItemId(newRule.id)
		}
	} else {
		await manager.update(currentItemId.value, values)
	}
})
</script>

<template>
	<form class="flex flex-col gap-3" @submit.prevent="handleSubmit">
		<FormField v-slot="{ field }" name="name">
			<FormItem>
				<FormLabel>{{ t('sharedTexts.name') }}</FormLabel>
				<FormControl>
					<Input type="text" v-bind="field" placeholder="Name of filter, just for your reference" />
				</FormControl>
				<FormMessage />
			</FormItem>
		</FormField>

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
			v-if="currentEditType === ModerationSettingsType.Symbols"
		/>

		<ModalLanguage
			v-if="currentEditType === ModerationSettingsType.Language"
		/>

		<ModalLongMessage
			v-if="currentEditType === ModerationSettingsType.LongMessage"
		/>

		<ModalCaps
			v-if="currentEditType === ModerationSettingsType.Caps"
		/>

		<ModalEmotes
			v-if="currentEditType === ModerationSettingsType.Emotes"
		/>

		<ModalLinks
			v-if="currentEditType === ModerationSettingsType.Links"
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
				hide-broadcaster
				:excluded-types="[RoleTypeEnum.Broadcaster, RoleTypeEnum.Moderator]"
			/>
		</div>

		<ModalDenylist
			v-if="currentEditType === ModerationSettingsType.DenyList"
		/>

		<Button type="submit">
			{{ t('sharedButtons.save') }}
		</Button>
	</form>
</template>
