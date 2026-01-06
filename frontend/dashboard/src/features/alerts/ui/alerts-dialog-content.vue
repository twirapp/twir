<script setup lang="ts">
import { computed, onMounted, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'

import AlertsDialogContentAudio from './alerts-dialog-content-audio.vue'

import { useTwitchGetUsers } from '@/api/twitch'
import { type Alert, useAlertsApi } from '@/api/alerts.js'
import { useCommandsApi } from '@/api/commands/commands.js'
import { useGreetingsApi } from '@/api/greetings.js'
import { useKeywordsApi } from '@/api/keywords.js'
import RewardsSelector from '@/components/rewardsSelector.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import FormMessage from '@/components/ui/form/FormMessage.vue'

const props = defineProps<{ alert?: Alert | null }>()
const emits = defineEmits<{ close: [] }>()

const { t } = useI18n()

const formSchema = z.object({
	name: z.string().min(1, t('alerts.validations.name')).max(30, t('alerts.validations.name')),
	audioId: z.string().optional(),
	audioVolume: z.number().min(0).max(100).default(75),
	commandIds: z.array(z.string()).default([]),
	rewardIds: z.array(z.string()).default([]),
	greetingsIds: z.array(z.string()).default([]),
	keywordsIds: z.array(z.string()).default([]),
})

const { handleSubmit, defineField, setValues } = useForm({
	validationSchema: toTypedSchema(formSchema),
	initialValues: {
		name: '',
		audioId: undefined,
		audioVolume: 75,
		commandIds: [],
		rewardIds: [],
		greetingsIds: [],
		keywordsIds: [],
	},
})

const [name, nameAttrs] = defineField('name')
const [commandIds] = defineField('commandIds')
const [rewardIds] = defineField('rewardIds')
const [keywordsIds] = defineField('keywordsIds')
const [greetingsIds] = defineField('greetingsIds')
const [audioId] = defineField('audioId')
const audioVolume = defineField('audioVolume')[0]

onMounted(() => {
	if (!props.alert) return
	const alertData = structuredClone(toRaw(props.alert))
	// Convert null to undefined for nullable fields to satisfy form schema
	setValues({
		name: alertData.name,
		audioId: alertData.audioId ?? undefined,
		audioVolume: alertData.audioVolume ?? 75,
		commandIds: alertData.commandIds ?? [],
		rewardIds: alertData.rewardIds ?? [],
		greetingsIds: alertData.greetingsIds ?? [],
		keywordsIds: alertData.keywordsIds ?? [],
	})
})

const manager = useAlertsApi()
const alertsCreateMutation = manager.useAlertsCreateMutation()
const alertsUpdateMutation = manager.useAlertsUpdateMutation()

const save = handleSubmit(async (formValues) => {
	const data = { ...formValues, id: undefined }

	if (props.alert?.id) {
		await alertsUpdateMutation.executeMutation({
			id: props.alert.id,
			opts: data,
		})
	} else {
		delete data.id
		await alertsCreateMutation.executeMutation({ opts: data })
	}

	emits('close')
})

const commandsApi = useCommandsApi()
const { data: commands } = commandsApi.useQueryCommands()
const commandsSelectOptions = computed(() =>
	commands.value?.commands.map((command) => ({ label: command.name, value: command.id }))
)

const greetingsApi = useGreetingsApi()
const { data: greetings } = greetingsApi.useQueryGreetings()
const greetingsUsersIds = computed(() => greetings.value?.greetings.map((g) => g.userId) ?? [])
const { data: twitchUsers } = useTwitchGetUsers({ ids: greetingsUsersIds })
const greetingsSelectOptions = computed(() => {
	if (!greetingsUsersIds.value.length || !twitchUsers.value?.length) return []
	return greetings.value?.greetings.map((g) => {
		const twitchUser = twitchUsers.value.find((u) => u.id === g.userId)
		return { label: twitchUser?.login ?? g.userId, value: g.id }
	})
})

const keywordsApi = useKeywordsApi()
const { data: keywords } = keywordsApi.useQueryKeywords()
const keywordsSelectOptions = computed(() =>
	keywords.value?.keywords.map((k) => ({ label: k.text, value: k.id }))
)
</script>

<template>
	<form @submit="save" class="p-6 pt-0">
		<div class="flex flex-col gap-6">
			<div class="flex flex-col gap-2">
				<Label for="name">Name <span class="text-destructive">*</span></Label>
				<Input id="name" v-model="name" v-bind="nameAttrs" :maxlength="30" />
				<FormMessage name="name" />
			</div>

			<div class="flex flex-col gap-2">
				<Label for="commandIds">{{ t('alerts.trigger.commands') }}</Label>
				<Select v-model="commandIds" multiple>
					<SelectTrigger id="commandIds">
						<SelectValue placeholder="Select commands" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem
							v-for="option in commandsSelectOptions"
							:key="option.value"
							:value="option.value"
						>
							{{ option.label }}
						</SelectItem>
					</SelectContent>
				</Select>
				<FormMessage name="commandIds" />
			</div>

			<div class="flex flex-col gap-2">
				<Label for="rewardIds">{{ t('alerts.trigger.rewards') }}</Label>
				<RewardsSelector v-model="rewardIds" multiple />
				<FormMessage name="rewardIds" />
			</div>

			<div class="flex flex-col gap-2">
				<Label for="keywordsIds">{{ t('alerts.trigger.keywords') }}</Label>
				<Select v-model="keywordsIds" multiple>
					<SelectTrigger id="keywordsIds">
						<SelectValue placeholder="Select keywords" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem
							v-for="option in keywordsSelectOptions"
							:key="option.value"
							:value="option.value"
						>
							{{ option.label }}
						</SelectItem>
					</SelectContent>
				</Select>
				<FormMessage name="keywordsIds" />
			</div>

			<div class="flex flex-col gap-2">
				<Label for="greetingsIds">{{ t('alerts.trigger.greetings') }}</Label>
				<Select v-model="greetingsIds" multiple>
					<SelectTrigger id="greetingsIds">
						<SelectValue placeholder="Select greetings" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem
							v-for="option in greetingsSelectOptions"
							:key="option.value"
							:value="option.value"
						>
							{{ option.label }}
						</SelectItem>
					</SelectContent>
				</Select>
				<FormMessage name="greetingsIds" />
			</div>
		</div>

		<Separator class="my-6" />

		<div class="flex flex-col gap-6">
			<AlertsDialogContentAudio
				v-model:audio-id="audioId"
				:initialVolume="audioVolume"
				@update:volume="(val) => audioVolume = val"
			/>

			<div class="flex justify-end">
				<Button type="submit">
					{{ t('sharedButtons.save') }}
				</Button>
			</div>
		</div>
	</form>
</template>
