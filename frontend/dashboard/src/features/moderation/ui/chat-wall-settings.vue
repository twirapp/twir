<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import { useCommandsApi } from '@/api/commands/commands.ts'
import { useModerationChatWall } from '@/api/moderation-chat-wall.ts'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import CommandsList from '@/features/commands/ui/list.vue'
import { cn } from '@/lib/utils'

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
watch(data, (v) => {
	if (!v) return
	chatSettingsForm.setValues({
		muteSubscribers: v.chatWallSettings.muteSubscribers,
		muteVips: v.chatWallSettings.muteVips,
	})
}, { immediate: true })

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
	return commands.value?.commands?.filter(c => chatWallCommandsNames.includes(c.defaultName))
})
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle>{{ t('chatWall.settings.title') }}</CardTitle>
		</CardHeader>
		<form @submit.prevent="onSubmit">
			<CardContent :class="cn('relative', { 'pointer-events-none': fetching })">
				<div
					v-if="fetching"
					class="absolute inset-0 z-50 flex items-center justify-center bg-background/80 backdrop-blur-sm"
				>
					<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
				</div>
				<div class="flex flex-col gap-4">
					<FormField v-slot="{ field }" name="muteSubscribers">
						<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
							<div class="space-y-0.5">
								<FormLabel class="text-base">
									{{ t('chatWall.settings.muteSubscribers') }}
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

					<FormField v-slot="{ field }" name="muteVips">
						<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
							<div class="space-y-0.5">
								<FormLabel class="text-base">
									{{ t('chatWall.settings.muteVips') }}
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
				</div>
			</CardContent>
			<CardFooter class="justify-end">
				<Button type="submit" :disabled="fetching">
					{{ t('sharedButtons.save') }}
				</Button>
			</CardFooter>
		</form>

		<Separator />

		<CardHeader>
			<CardTitle>{{ t('chatWall.commands.title') }}</CardTitle>
		</CardHeader>
		<CardContent>
			<div class="flex flex-row flex-wrap gap-4">
				<CommandsList
					v-if="chatWallCommands"
					:commands="chatWallCommands"
					show-background
				/>
			</div>
		</CardContent>
	</Card>
</template>

<style scoped>
.backdrop-blur-sm {
  backdrop-filter: blur(4px);
}
</style>
