<script setup lang="ts">
import { LanguagesIcon, SettingsIcon } from 'lucide-vue-next'
import { computed, ref } from 'vue'


import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useLanguagesApi } from '#layers/dashboard/api/languages'
import Card from '#layers/dashboard/components/card/card.vue'
import TwitchUsersSelect from '#layers/dashboard/components/twitchUsers/twitch-users-select.vue'







import { useChatTranslations } from '~/features/modules/composables/use-chat-translations'
import { ChannelRolePermissionEnum } from '~/gql/graphql.ts'

const { t } = useI18n()

const showSettings = ref(false)

const { handleSubmit, isLoading, fetching, exists } = useChatTranslations()

const languagesApi = useLanguagesApi()
const { data: languagesData } = languagesApi.useAvailableLanguages()

const availableLanguages = computed(
	() =>
		languagesData.value?.moderationLanguagesAvailableLanguages.languages.map((lang) => ({
			label: lang.name,
			value: lang.iso_639_1,
		})) || []
)

const canManageModules = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageModules)
</script>

<template>
	<Card
		:title="t('modules.chatTranslations.title')"
		:is-loading="fetching"
		:icon="LanguagesIcon"
		icon-height="30px"
		icon-width="30px"
		:description="t('modules.chatTranslations.description')"
	>
		<template #footer>
			<UiButton
				class="flex gap-2 items-center"
				variant="secondary"
				:disabled="!canManageModules"
				@click="showSettings = !showSettings"
			>
				{{ t('sharedTexts.settings') }}
				<SettingsIcon class="size-4" />
			</UiButton>
		</template>
	</Card>

	<UiDialog v-model:open="showSettings">
		<UiDialogContent class="sm:max-w-[600px]">
			<UiDialogHeader>
				<UiDialogTitle>{{ t('modules.chatTranslations.settings.title') }}</UiDialogTitle>
			</UiDialogHeader>

			<form class="space-y-6" @submit.prevent="handleSubmit">
				<UiFormField v-slot="{ componentField }" name="enabled">
					<UiFormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
						<div class="space-y-0.5">
							<UiFormLabel class="text-base">
								{{ t('modules.chatTranslations.settings.enabled.label') }}
							</UiFormLabel>
							<UiFormDescription>{{
								t('modules.chatTranslations.settings.enabled.description')
							}}</UiFormDescription>
						</div>
						<UiFormControl>
							<UiSwitch
								:model-value="componentField.modelValue"
								@update:model-value="componentField['onUpdate:modelValue']"
							/>
						</UiFormControl>
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ componentField }" name="targetLanguage">
					<UiFormItem>
						<UiFormLabel>{{ t('modules.chatTranslations.settings.targetLanguage.label') }}</UiFormLabel>
						<UiFormControl>
							<UiSelect
								:model-value="componentField.modelValue"
								@update:model-value="componentField['onUpdate:modelValue']"
							>
								<UiSelectTrigger>
									<UiSelectValue
										:placeholder="t('modules.chatTranslations.settings.targetLanguage.placeholder')"
									/>
								</UiSelectTrigger>
								<UiSelectContent>
									<UiSelectItem
										v-for="lang in availableLanguages"
										:key="lang.value"
										:value="lang.value"
									>
										{{ lang.label }}
									</UiSelectItem>
								</UiSelectContent>
							</UiSelect>
						</UiFormControl>
						<UiFormDescription>{{
							t('modules.chatTranslations.settings.targetLanguage.description')
						}}</UiFormDescription>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ componentField }" name="excludedLanguages">
					<UiFormItem>
						<UiFormLabel>{{
							t('modules.chatTranslations.settings.excludedLanguages.label')
						}}</UiFormLabel>
						<UiFormControl>
							<UiMultiSelect
								v-model:model-value="componentField.modelValue"
								:options="availableLanguages"
								:placeholder="t('modules.chatTranslations.settings.excludedLanguages.placeholder')"
								@update:model-value="componentField['onUpdate:modelValue']"
							/>
						</UiFormControl>
						<UiFormDescription>{{
							t('modules.chatTranslations.settings.excludedLanguages.description')
						}}</UiFormDescription>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ componentField }" name="useItalic">
					<UiFormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
						<div class="space-y-0.5">
							<UiFormLabel class="text-base">
								{{ t('modules.chatTranslations.settings.useItalic.label') }}
							</UiFormLabel>
							<UiFormDescription>{{
								t('modules.chatTranslations.settings.useItalic.description')
							}}</UiFormDescription>
						</div>
						<UiFormControl>
							<UiSwitch
								:model-value="componentField.modelValue"
								@update:model-value="componentField['onUpdate:modelValue']"
							/>
						</UiFormControl>
					</UiFormItem>
				</UiFormField>

				<UiCard>
					<UiCardHeader>
						<UiCardTitle>{{ t('modules.chatTranslations.settings.excludedUsers.title') }}</UiCardTitle>
					</UiCardHeader>
					<UiCardContent class="space-y-4">
						<UiFormField v-slot="{ field }" name="excludedUsersIDs">
							<UiFormItem>
								<TwitchUsersSelect
									v-model:model-value="field.value"
									class="flex-1"
									@update:model-value="field['onUpdate:modelValue']"
								/>
								<UiFormDescription>{{
									t('modules.chatTranslations.settings.excludedUsers.description')
								}}</UiFormDescription>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>
					</UiCardContent>
				</UiCard>

				<div class="flex justify-end">
					<UiButton type="submit" :disabled="isLoading">
						{{
							exists
								? t('modules.chatTranslations.settings.buttons.update')
								: t('modules.chatTranslations.settings.buttons.create')
						}}
					</UiButton>
				</div>
			</form>
		</UiDialogContent>
	</UiDialog>
</template>
