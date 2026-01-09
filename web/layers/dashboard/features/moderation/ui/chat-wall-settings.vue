<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed, watch } from 'vue'

import * as z from 'zod'

import { useCommandsApi } from '#layers/dashboard/api/commands/commands.ts'
import { useModerationChatWall } from '#layers/dashboard/api/moderation-chat-wall.ts'





import CommandsList from '~/features/commands/ui/list.vue'
import { cn } from '~/lib/utils'

const { t } = useI18n()
const api = useModerationChatWall()
const update = api.useUpdateSettings()

const schema = z.object({
	muteSubscribers: z.boolean(),
	muteVips: z.boolean(),
})
const chatSettingsForm = useForm({
	validationSchema: toTypedSchema(schema),
})

const { data, fetching } = api.useSettings()
watch(
	data,
	(v) => {
		if (!v) return
		chatSettingsForm.setValues({
			muteSubscribers: v.chatWallSettings.muteSubscribers,
			muteVips: v.chatWallSettings.muteVips,
		})
	},
	{ immediate: true }
)

const onSubmit = chatSettingsForm.handleSubmit(async (values) => {
	await update.executeMutation({
		opts: values,
	})
})

const commandsApi = useCommandsApi()
const { data: commands } = commandsApi.useQueryCommands()

const chatWallCommandsNames = [
	'chat wall delete',
	'chat wall ban',
	'chat wall timeout',
	'chat wall stop',
]

const chatWallCommands = computed(() => {
	return commands.value?.commands?.filter((c) => {
		return c.defaultName && chatWallCommandsNames.includes(c.defaultName)
	})
})
</script>

<template>
	<UiCard>
		<UiCardHeader>
			<UiCardTitle>{{ t('chatWall.settings.title') }}</UiCardTitle>
		</UiCardHeader>
		<form @submit.prevent="onSubmit">
			<UiCardContent :class="cn('relative', { 'pointer-events-none': fetching })">
				<div
					v-if="fetching"
					class="absolute inset-0 z-50 flex items-center justify-center bg-background/80 backdrop-blur-xs"
				>
					<div
						class="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"
					/>
				</div>
				<div class="flex flex-col gap-4">
					<UiFormField v-slot="{ field }" name="muteSubscribers">
						<UiFormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
							<div class="space-y-0.5">
								<UiFormLabel class="text-base">
									{{ t('chatWall.settings.muteSubscribers') }}
								</UiFormLabel>
							</div>
							<UiFormControl>
								<UiSwitch
									:model-value="field.value"
									default-checked
									@update:model-value="field['onUpdate:modelValue']"
								/>
							</UiFormControl>
						</UiFormItem>
					</UiFormField>

					<UiFormField v-slot="{ field }" name="muteVips">
						<UiFormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
							<div class="space-y-0.5">
								<UiFormLabel class="text-base">
									{{ t('chatWall.settings.muteVips') }}
								</UiFormLabel>
							</div>
							<UiFormControl>
								<UiSwitch
									:model-value="field.value"
									default-checked
									@update:model-value="field['onUpdate:modelValue']"
								/>
							</UiFormControl>
						</UiFormItem>
					</UiFormField>
				</div>
			</UiCardContent>
			<UiCardFooter class="justify-end">
				<UiButton type="submit" :disabled="fetching">
					{{ t('sharedButtons.save') }}
				</UiButton>
			</UiCardFooter>
		</form>

		<UiSeparator />

		<UiCardHeader>
			<UiCardTitle>{{ t('chatWall.commands.title') }}</UiCardTitle>
		</UiCardHeader>
		<UiCardContent>
			<div class="flex flex-row flex-wrap gap-4">
				<CommandsList v-if="chatWallCommands" :commands="chatWallCommands" show-background />
			</div>
		</UiCardContent>
	</UiCard>
</template>

<style scoped>
.backdrop-blur-xs {
	backdrop-filter: blur(4px);
}
</style>
