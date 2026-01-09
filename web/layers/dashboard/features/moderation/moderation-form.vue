<script setup lang="ts">
import { useField, useForm, useSetFormValues, useSubmitForm } from 'vee-validate'
import { nextTick, onMounted, toRaw } from 'vue'

import { useRoute, useRouter } from 'vue-router'

import ModalCaps from './ui/form/moderation-form-caps.vue'
import ModalDenylist from './ui/form/moderation-form-denylist.vue'
import ModalEmotes from './ui/form/moderation-form-emotes.vue'
import ModalLanguage from './ui/form/moderation-form-language.vue'
import ModalLinks from './ui/form/moderation-form-links.vue'
import ModalLongMessage from './ui/form/moderation-form-longmessage.vue'
import ModalOneManSpam from './ui/form/moderation-form-one-man-spam.vue'
import ModalSymbols from './ui/form/moderation-form-symbols.vue'

import type { EditableItem } from '~/features/moderation/composables/use-moderation-form.ts'








import FormRolesSelector from '~/features/commands/ui/form-roles-selector.vue'
import { useModerationApi } from '~/features/moderation/composables/use-moderation-api.ts'
import { moderationValidationRules } from '~/features/moderation/composables/use-moderation-form.ts'
// oxlint-disable-next-line consistent-type-imports
import { ModerationSettingsType, RoleTypeEnum } from '~/gql/graphql.ts'
import PageLayout from '~/layout/page-layout.vue'

const route = useRoute()
const router = useRouter()

const { t } = useI18n()

const manager = useModerationApi()

const {
	resetForm: formReset,
	setFieldValue,
	values,
} = useForm({
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

		const item = manager.items.value.find((i) => i.id === id)
		if (!item) return

		const values = structuredClone(toRaw(item))

		setFormValues(values)
	} else if (id === 'new' && type) {
		setFieldValue('type', type)
		setFieldValue('name', t(`moderation.types.${type}.name`))
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
				<UiButton type="submit">
					{{ t('sharedButtons.save') }}
				</UiButton>
			</template>

			<template #content>
				<UiCard>
					<UiCardHeader>
						<UiCardTitle class="flex flex-row justify-between">
							{{ t(`moderation.types.${values.type}.name`) }}
							<UiFormField v-slot="{ field }" name="enabled">
								<UiFormItem>
									<UiFormControl>
										<UiSwitch
											:model-value="field.value"
											default-checked
											@update:model-value="field['onUpdate:modelValue']"
										/>
									</UiFormControl>
								</UiFormItem>
							</UiFormField>
						</UiCardTitle>
					</UiCardHeader>
					<UiCardContent class="flex flex-col gap-3">
						<UiFormField v-slot="{ field }" name="name">
							<UiFormItem>
								<UiFormLabel>{{ t('sharedTexts.name') }}</UiFormLabel>
								<UiFormControl>
									<UiInput type="text" v-bind="field" :placeholder="t('moderation.name')" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<ModalSymbols v-if="currentEditType === ModerationSettingsType.Symbols" />

						<ModalLanguage v-if="currentEditType === ModerationSettingsType.Language" />

						<ModalLongMessage v-if="currentEditType === ModerationSettingsType.LongMessage" />

						<ModalCaps v-if="currentEditType === ModerationSettingsType.Caps" />

						<ModalEmotes v-if="currentEditType === ModerationSettingsType.Emotes" />

						<ModalLinks v-if="currentEditType === ModerationSettingsType.Links" />

						<ModalOneManSpam v-if="currentEditType === ModerationSettingsType.OneManSpam" />

						<UiSeparator label="Timeouts" />

						<UiFormField v-slot="{ componentField }" name="banMessage">
							<UiFormItem>
								<UiFormLabel>Timeout message</UiFormLabel>
								<UiFormControl>
									<UiInput type="text" v-bind="componentField" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ componentField }" name="banTime">
							<UiFormItem>
								<UiFormLabel>{{ t('moderation.banTime') }}</UiFormLabel>
								<UiFormControl>
									<UiInput type="number" v-bind="componentField" />
								</UiFormControl>
								<UiFormMessage />
								<UiFormDescription>
									{{ t('moderation.banDescription') }}
								</UiFormDescription>
							</UiFormItem>
						</UiFormField>

						<UiSeparator label="Warnings" />

						<UiFormField v-slot="{ componentField }" name="warningMessage">
							<UiFormItem>
								<UiFormLabel>{{ t('moderation.warningMessage') }}</UiFormLabel>
								<UiFormControl>
									<UiInput type="text" v-bind="componentField" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ componentField }" name="maxWarnings">
							<UiFormItem>
								<UiFormLabel>{{ t('moderation.warningMaxCount') }}</UiFormLabel>
								<UiFormControl>
									<UiInput type="number" v-bind="componentField" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiSeparator />

						<div class="flex flex-col gap-2">
							<UiLabel>Select excluded from moderation roles:</UiLabel>
							<FormRolesSelector
								field-name="excludedRoles"
								hide-everyone
								hide-broadcaster
								:excluded-types="[RoleTypeEnum.Broadcaster, RoleTypeEnum.Moderator]"
							/>
						</div>

						<ModalDenylist v-if="currentEditType === ModerationSettingsType.DenyList" />
					</UiCardContent>
				</UiCard>
			</template>
		</PageLayout>
	</form>
</template>
