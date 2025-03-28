<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { watch } from 'vue'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import { useModerationChatWall } from '@/api/moderation-chat-wall.ts'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import CommandButton from '@/features/commands/ui/command-button.vue'
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
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle>Setup</CardTitle>
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
									Mute subscribers
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
									Mute VIPs
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
			<CardTitle>Commands</CardTitle>
		</CardHeader>
		<CardContent>
			<div class="flex flex-col gap-4">
				<CommandButton
					name="chat wall delete"
					title="Chat wall delete"
				/>
				<CommandButton
					name="chat wall ban"
					title="Chat wall ban"
				/>
				<CommandButton
					name="chat wall timeout"
					title="Chat wall timeout"
				/>
				<CommandButton
					name="chat wall stop"
					title="Chat wall stop"
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
