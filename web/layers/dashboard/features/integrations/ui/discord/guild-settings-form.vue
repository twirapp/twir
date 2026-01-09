<script setup lang="ts">
import '@discord-message-components/vue/dist/style.css'

import {
	DiscordEmbed,
	DiscordEmbedField,
	DiscordEmbedFields,
	DiscordMention,
	DiscordMessage,
	DiscordMessages,
	// eslint-disable-next-line ts/ban-ts-comment
	// @ts-expect-error
} from '@discord-message-components/vue'
import { toTypedSchema } from '@vee-validate/zod'
import { Info } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { computed, watch } from 'vue'


import TwitchMultipleUsersSelector from '#layers/dashboard/components/twitchUsers/twitch-users-select.vue'







import { toast } from 'vue-sonner'
import {
	useDiscordGuildInfo,
	useDiscordIntegration,
} from '~/features/integrations/composables/discord/use-discord-integration.js'


import StreamStarting from '~/assets/images/streamStarting.jpeg?url'
// ...existing code...
import {
	type DiscordGuildUpdateInputInput,
	DiscordGuildUpdateInputSchema,
} from '~/gql/validation-schemas.ts'
import { DiscordChannelType } from '~/gql/graphql.ts'

const props = defineProps<{
	guildId: string
}>()

const { profile: currentUser } = storeToRefs(useDashboardAuth())

const { t } = useI18n()

const { updateGuild, guilds } = useDiscordIntegration()
const { channels, roles, isLoading: isGuildInfoLoading } = useDiscordGuildInfo(() => props.guildId)

const form = useForm<DiscordGuildUpdateInputInput>({
	validationSchema: toTypedSchema(DiscordGuildUpdateInputSchema),
})

// Update form when initialValues change
watch(
	guilds,
	(newGuilds) => {
		const g = newGuilds.find((guild) => guild.id === props.guildId)
		if (!g) return

		form.setValues({
			liveNotificationEnabled: g.liveNotificationEnabled,
			liveNotificationChannelsIds: g.liveNotificationChannelsIds,
			liveNotificationShowTitle: g.liveNotificationShowTitle,
			liveNotificationShowCategory: g.liveNotificationShowCategory,
			liveNotificationShowViewers: g.liveNotificationShowViewers,
			liveNotificationMessage: g.liveNotificationMessage,
			liveNotificationShowPreview: g.liveNotificationShowPreview,
			liveNotificationShowProfileImage: g.liveNotificationShowProfileImage,
			offlineNotificationMessage: g.offlineNotificationMessage,
			shouldDeleteMessageOnOffline: g.shouldDeleteMessageOnOffline,
			additionalUsersIdsForLiveCheck: g.additionalUsersIdsForLiveCheck,
		})
	},
	{ immediate: true }
)

const textChannelOptions = computed(() => {
	return channels.value
		.filter((c) =>
			[DiscordChannelType.ChannelTypeText, DiscordChannelType.ChannelTypeAnnouncement].includes(
				c.type
			)
		)
		.map((c) => ({
			label: `#${c.name}`,
			value: c.id,
			disabled: !c.canSendMessages,
		}))
})

function getRoleColor(roleName: string) {
	const role = roles.value.find((r) => r.name === roleName.replace('@', ''))
	if (!role || role.color === '0') return null

	const hexColor = Number(role.color).toString(16)
	return `#${hexColor.padStart(6, '0')}`
}

const onSubmit = form.handleSubmit(async (values) => {
	const result = await updateGuild(props.guildId, values)

	if (result.error) {
		toast.error(t('sharedTexts.error'), {
			description: result.error.message,
			duration: 5000,
		})
	} else {
		toast.success(t('sharedTexts.saved'), {
			duration: 2500,
		})
	}
})

// Computed values for preview
const formValues = computed(() => form.values)

const messageWithMentions = computed(() => {
	return formValues.value.liveNotificationMessage?.split(' ')
})
</script>

<template>
	<form class="flex flex-col gap-6" @submit="onSubmit">
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<!-- Settings Column -->
			<UiCard>
				<UiCardContent class="pt-6">
					<div class="flex flex-col gap-4">
						<!-- Toggle Switches -->
						<div class="flex flex-col gap-3">
							<UiFormField v-slot="{ value, handleChange }" name="liveNotificationEnabled">
								<UiFormItem class="flex flex-row items-center gap-3">
									<UiFormControl>
										<UiCheckbox :model-value="value" @update:model-value="handleChange" />
									</UiFormControl>
									<UiFormLabel class="mt-0!">{{ t('sharedTexts.enabled') }}</UiFormLabel>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ value, handleChange }" name="liveNotificationShowTitle">
								<UiFormItem class="flex flex-row items-center gap-3">
									<UiFormControl>
										<UiCheckbox :model-value="value" @update:model-value="handleChange" />
									</UiFormControl>
									<UiFormLabel class="mt-0!">{{
										t('integrations.discord.alerts.showTitle')
									}}</UiFormLabel>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ value, handleChange }" name="liveNotificationShowCategory">
								<UiFormItem class="flex flex-row items-center gap-3">
									<UiFormControl>
										<UiCheckbox :model-value="value" @update:model-value="handleChange" />
									</UiFormControl>
									<UiFormLabel class="mt-0!">{{
										t('integrations.discord.alerts.showCategory')
									}}</UiFormLabel>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ value, handleChange }" name="liveNotificationShowPreview">
								<UiFormItem class="flex flex-row items-center gap-3">
									<UiFormControl>
										<UiCheckbox :model-value="value" @update:model-value="handleChange" />
									</UiFormControl>
									<UiFormLabel class="mt-0!">{{
										t('integrations.discord.alerts.showPreview')
									}}</UiFormLabel>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ value, handleChange }" name="liveNotificationShowProfileImage">
								<UiFormItem class="flex flex-row items-center gap-3">
									<UiFormControl>
										<UiCheckbox :model-value="value" @update:model-value="handleChange" />
									</UiFormControl>
									<UiFormLabel class="mt-0!">{{
										t('integrations.discord.alerts.showProfileImage')
									}}</UiFormLabel>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ value, handleChange }" name="liveNotificationShowViewers">
								<UiFormItem class="flex flex-row items-center gap-3">
									<UiFormControl>
										<UiCheckbox :model-value="value" @update:model-value="handleChange" />
									</UiFormControl>
									<UiFormLabel class="mt-0!">{{
										t('integrations.discord.alerts.showViewers')
									}}</UiFormLabel>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>
						</div>

						<UiSeparator />

						<!-- Channel Selector -->
						<UiFormField v-slot="{ value, handleChange }" name="liveNotificationChannelsIds">
							<UiFormItem>
								<UiFormLabel>{{ t('integrations.discord.alerts.channelsSelect') }}</UiFormLabel>
								<UiFormControl>
									<UiMultiSelect
										:model-value="value"
										:options="textChannelOptions"
										:disabled="isGuildInfoLoading"
										placeholder="Select channels..."
										@update:model-value="handleChange"
									/>
								</UiFormControl>
								<UiFormDescription />
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiSeparator />

						<!-- Additional Users -->
						<UiFormField v-slot="{ value, handleChange }" name="additionalUsersIdsForLiveCheck">
							<UiFormItem>
								<UiFormLabel>{{
									t('integrations.discord.alerts.additionalUsersIdsForLiveCheck')
								}}</UiFormLabel>
								<UiFormControl>
									<TwitchMultipleUsersSelector
										:model-value="value"
										:max="10"
										@update:model-value="handleChange"
									/>
								</UiFormControl>
								<UiFormDescription>
									<span v-if="value.length >= 10" class="text-orange-500">
										Maximum additional users selected
									</span>
								</UiFormDescription>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiSeparator />

						<!-- Online Message -->
						<UiFormField v-slot="{ componentField }" name="liveNotificationMessage">
							<UiFormItem>
								<UiFormLabel>{{ t('integrations.discord.alerts.streamOnlineLabel') }}</UiFormLabel>
								<UiFormControl>
									<UiTextarea
										v-bind="componentField"
										:placeholder="t('integrations.discord.alerts.streamOnlinePlaceholder')"
									/>
								</UiFormControl>
								<UiFormDescription>
									{userName}, {displayName}, {title}, {categoryName} – supported variables. Use
									@rolename to mention roles.
								</UiFormDescription>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<!-- Offline Message -->
						<UiFormField v-slot="{ componentField }" name="offlineNotificationMessage">
							<UiFormItem>
								<UiFormLabel>{{ t('integrations.discord.alerts.streamOfflineLabel') }}</UiFormLabel>
								<UiFormControl>
									<UiTextarea
										v-bind="componentField"
										:placeholder="t('integrations.discord.alerts.streamOfflinePlaceholder')"
										:disabled="formValues.shouldDeleteMessageOnOffline"
									/>
								</UiFormControl>
								<UiFormDescription> {userName}, {displayName} – supported variables </UiFormDescription>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ value, handleChange }" name="shouldDeleteMessageOnOffline">
							<UiFormItem class="flex flex-row items-center gap-3">
								<UiFormControl>
									<UiCheckbox :model-value="value" @update:model-value="handleChange" />
								</UiFormControl>
								<UiFormLabel class="mt-0!">{{
									t('integrations.discord.alerts.shouldDeleteMessageOnOffline')
								}}</UiFormLabel>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiAlert>
							<Info class="h-4 w-4" />
							<UiAlertTitle>Info</UiAlertTitle>
							<UiAlertDescription>
								{{ t('integrations.discord.alerts.updateAlert') }}
							</UiAlertDescription>
						</UiAlert>

						<UiButton type="submit" class="w-full sm:w-auto">
							{{ t('sharedButtons.save') }}
						</UiButton>
					</div>
				</UiCardContent>
			</UiCard>

			<div class="flex flex-col gap-4">
				<DiscordMessages class="rounded-md">
					<DiscordMessage :bot="true" author="TwirApp" avatar="/twir.svg">
						<template v-for="(m, idx) of messageWithMentions" :key="idx">
							<DiscordMention
								v-if="m.startsWith('@')"
								type="role"
								:highlight="true"
								:role-color="getRoleColor(m.trim())"
							>
								{{ m.replace('@', '') }}
							</DiscordMention>
							<template v-else>
								{{ m }}
							</template>
							{{ ' ' }}
						</template>
						<template #embeds>
							<DiscordEmbed
								embed-title="Today we are doing amazing things!"
								:url="`https://twitch.tv/${currentUser?.login}`"
								:timestamp="new Date()"
								footer-icon="/twir.svg"
								border-color="#6441a5"
								:thumbnail="
									formValues.liveNotificationShowProfileImage ? currentUser?.avatar : null
								"
								:image="formValues.liveNotificationShowPreview ? StreamStarting : null"
							>
								<template #fields>
									<DiscordEmbedFields>
										<DiscordEmbedField
											v-if="formValues.liveNotificationShowTitle"
											field-title="Title"
											inline
										>
											Today we are doing amazing things!
										</DiscordEmbedField>
										<DiscordEmbedField
											v-if="formValues.liveNotificationShowViewers"
											field-title="Viewers"
											inline
										>
											5
										</DiscordEmbedField>
										<DiscordEmbedField
											v-if="formValues.liveNotificationShowCategory"
											field-title="Category"
										>
											Software and game development
										</DiscordEmbedField>
									</DiscordEmbedFields>
								</template>

								<template #footer> TwirApp </template>
							</DiscordEmbed>
						</template>
					</DiscordMessage>
				</DiscordMessages>
			</div>
		</div>
	</form>
</template>
