<script setup lang="ts">
import { useField, useForm, useSetFormValues, useSubmitForm } from 'vee-validate'
import { nextTick, onMounted, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import ModalCaps from './ui/form/moderation-form-caps.vue'
import ModalDenylist from './ui/form/moderation-form-denylist.vue'
import ModalEmotes from './ui/form/moderation-form-emotes.vue'
import ModalLanguage from './ui/form/moderation-form-language.vue'
import ModalLinks from './ui/form/moderation-form-links.vue'
import ModalLongMessage from './ui/form/moderation-form-longmessage.vue'
import ModalSymbols from './ui/form/moderation-form-symbols.vue'

import type {
	EditableItem,
} from '@/features/moderation/composables/use-moderation-form.ts'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
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
import {
	moderationValidationRules,
} from '@/features/moderation/composables/use-moderation-form.ts'
import { ModerationSettingsType , RoleTypeEnum } from '@/gql/graphql.ts'
import PageLayout from '@/layout/page-layout.vue'

const route = useRoute()
const router = useRouter()

const { t } = useI18n()

const manager = useModerationApi()

const { resetForm: formReset, setFieldValue, values } = useForm({
	validationSchema: moderationValidationRules,
	keepValuesOnUnmount: false,
	validateOnMount: false,
})

const { value: currentItemId, setValue: setCurrentItemId } = useField<string | undefined>('id')
const setFormValues = useSetFormValues<EditableItem>()
const { value: currentEditType } = useField<ModerationSettingsType>('type')

onMounted(async () => {
	formReset()

	const id = route.params.id as string | 'new'
	const type = route.query.ruleType as ModerationSettingsType

	if (id !== 'new') {
		// wait items to be loaded
		await manager.fetchItems()
		await nextTick()

		const item = manager.items.value.find(i => i.id === id)
		if (!item) return

		const values = structuredClone(toRaw(item))

		setFormValues(values)
	} else if (id === 'new' && type) {
		setFieldValue('type', type)
	} else {
		router.push({ name: 'Moderation' })
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
	<form @submit.prevent="handleSubmit">
		<PageLayout back-redirect-to="/dashboard/moderation?tab=rules" show-back sticky-header>
			<template #title>
				<span v-if="route.params.id === 'new'">
					{{ t('sharedTexts.create') }} {{ t(`moderation.types.${values.type}.name`) }}
				</span>
				<span v-else-if="values.id">
					Edit {{ values.name ?? t(`moderation.types.${values.type}.name`) }}
				</span>
			</template>

			<template #action>
				<Button type="submit">
					{{ t('sharedButtons.save') }}
				</Button>
			</template>

			<template #content>
				<Card>
					<CardHeader>
						<CardTitle class="flex flex-row justify-between">
							{{ t(`moderation.types.${values.type}.name`) }}
							<FormField v-slot="{ field }" name="enabled">
								<FormItem>
									<FormControl>
										<Switch
											:checked="field.value"
											default-checked
											@update:checked="field['onUpdate:modelValue']"
										/>
									</FormControl>
								</FormItem>
							</FormField>
						</CardTitle>
					</CardHeader>
					<CardContent class="flex flex-col gap-3">
						<FormField v-slot="{ field }" name="name">
							<FormItem>
								<FormLabel>{{ t('sharedTexts.name') }}</FormLabel>
								<FormControl>
									<Input type="text" v-bind="field" :placeholder="t('moderation.name')" />
								</FormControl>
								<FormMessage />
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
					</CardContent>
				</Card>
			</template>
		</PageLayout>
	</form>
</template>
