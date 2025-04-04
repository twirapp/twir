<script setup lang="ts">
import { LanguagesIcon, SettingsIcon } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useUserAccessFlagChecker } from '@/api'
import { useLanguagesApi } from '@/api/languages'
import Card from '@/components/card/card.vue'
import TwitchUsersSelect from '@/components/twitchUsers/twitch-users-select.vue'
import { Button } from '@/components/ui/button'
import { CardContent, CardHeader, CardTitle, Card as UICard } from '@/components/ui/card'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { MultiSelect } from '@/components/ui/multi-select'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { useChatTranslations } from '@/features/modules/composables/use-chat-translations'
import { ChannelRolePermissionEnum } from '@/gql/graphql.ts'

const { t } = useI18n()

const showSettings = ref(false)

const {
	handleSubmit,
	isLoading,
	fetching,
	exists,
} = useChatTranslations()

const languagesApi = useLanguagesApi()
const { data: languagesData } = languagesApi.useAvailableLanguages()

const availableLanguages = computed(() =>
	languagesData.value?.moderationLanguagesAvailableLanguages.languages.map(lang => ({
		label: lang.name,
		value: lang.iso_639_1,
	})) || [],
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
			<Button class="flex gap-2 items-center" variant="secondary" :disabled="!canManageModules" @click="showSettings = !showSettings">
				{{ t('sharedTexts.settings') }}
				<SettingsIcon class="size-4" />
			</Button>
		</template>
	</Card>

	<Dialog v-model:open="showSettings">
		<DialogContent class="sm:max-w-[600px]">
			<DialogHeader>
				<DialogTitle>{{ t('modules.chatTranslations.settings.title') }}</DialogTitle>
			</DialogHeader>

			<form class="space-y-6" @submit.prevent="handleSubmit">
				<FormField v-slot="{ componentField }" name="enabled">
					<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
						<div class="space-y-0.5">
							<FormLabel class="text-base">
								{{ t('modules.chatTranslations.settings.enabled.label') }}
							</FormLabel>
							<FormDescription>{{ t('modules.chatTranslations.settings.enabled.description') }}</FormDescription>
						</div>
						<FormControl>
							<Switch
								:checked="componentField.modelValue"
								@update:checked="componentField['onUpdate:modelValue']"
							/>
						</FormControl>
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="targetLanguage">
					<FormItem>
						<FormLabel>{{ t('modules.chatTranslations.settings.targetLanguage.label') }}</FormLabel>
						<FormControl>
							<Select
								:model-value="componentField.modelValue"
								@update:model-value="componentField['onUpdate:modelValue']"
							>
								<SelectTrigger>
									<SelectValue :placeholder="t('modules.chatTranslations.settings.targetLanguage.placeholder')" />
								</SelectTrigger>
								<SelectContent>
									<SelectItem
										v-for="lang in availableLanguages"
										:key="lang.value"
										:value="lang.value"
									>
										{{ lang.label }}
									</SelectItem>
								</SelectContent>
							</Select>
						</FormControl>
						<FormDescription>{{ t('modules.chatTranslations.settings.targetLanguage.description') }}</FormDescription>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="excludedLanguages">
					<FormItem>
						<FormLabel>{{ t('modules.chatTranslations.settings.excludedLanguages.label') }}</FormLabel>
						<FormControl>
							<MultiSelect
								v-model:model-value="componentField.modelValue"
								:options="availableLanguages"
								:placeholder="t('modules.chatTranslations.settings.excludedLanguages.placeholder')"
								@update:model-value="componentField['onUpdate:modelValue']"
							/>
						</FormControl>
						<FormDescription>{{ t('modules.chatTranslations.settings.excludedLanguages.description') }}</FormDescription>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="useItalic">
					<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
						<div class="space-y-0.5">
							<FormLabel class="text-base">
								{{ t('modules.chatTranslations.settings.useItalic.label') }}
							</FormLabel>
							<FormDescription>{{ t('modules.chatTranslations.settings.useItalic.description') }}</FormDescription>
						</div>
						<FormControl>
							<Switch
								:checked="componentField.modelValue"
								@update:checked="componentField['onUpdate:modelValue']"
							/>
						</FormControl>
					</FormItem>
				</FormField>

				<UICard>
					<CardHeader>
						<CardTitle>{{ t('modules.chatTranslations.settings.excludedUsers.title') }}</CardTitle>
					</CardHeader>
					<CardContent class="space-y-4">
						<FormField v-slot="{ field }" name="excludedUsersIDs">
							<FormItem>
								<TwitchUsersSelect
									v-model:model-value="field.value"
									class="flex-1"
									@update:model-value="field['onUpdate:modelValue']"
								/>
								<FormDescription>{{ t('modules.chatTranslations.settings.excludedUsers.description') }}</FormDescription>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</UICard>

				<div class="flex justify-end">
					<Button type="submit" :disabled="isLoading">
						{{ exists ? t('modules.chatTranslations.settings.buttons.update') : t('modules.chatTranslations.settings.buttons.create') }}
					</Button>
				</div>
			</form>
		</DialogContent>
	</Dialog>
</template>
