<script setup lang="ts">
import { computed, onMounted, toRaw } from 'vue'

import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'

import AlertsDialogContentAudio from './alerts-dialog-content-audio.vue'

import { useTwitchGetUsers } from '#layers/dashboard/api/twitch'
import { type Alert, useAlertsApi } from '#layers/dashboard/api/alerts'
import { useCommandsApi } from '#layers/dashboard/api/commands/commands'
import { useGreetingsApi } from '#layers/dashboard/api/greetings'
import { useKeywordsApi } from '#layers/dashboard/api/keywords'
import RewardsSelector from '#layers/dashboard/components/rewardsSelector.vue'






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

const { handleSubmit, setValues, useFieldModel } = useForm({
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

const audioId = useFieldModel<'audioId'>('audioId')
const audioVolume = useFieldModel<'audioVolume'>('audioVolume')

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
				<UiFormField name="name" v-slot="{componentField}">
		      <UiFormItem>
						<UiFormLabel>Name <span class="text-destructive">*</span></UiFormLabel>
						<UiInput v-bind="componentField" :maxlength="30" />
						<UiFormMessage />
		      </UiFormItem>
				</UiFormField>
			</div>

			<div class="flex flex-col gap-2">
				<UiFormField name="commandIds" v-slot="{componentField}">
		      <UiFormItem>
						<UiFormLabel>Commands</UiFormLabel>
						<UiSelect v-bind="componentField" multiple>
							<UiSelectTrigger>
								<UiSelectValue placeholder="Select commands" />
							</UiSelectTrigger>
							<UiSelectContent class="max-h-60 overflow-y-auto">
								<UiSelectItem
									v-for="option in commandsSelectOptions"
									:key="option.value"
									:value="option.value"
								>
									{{ option.label }}
								</UiSelectItem>
							</UiSelectContent>
						</UiSelect>
						<UiFormMessage />
		      </UiFormItem>
				</UiFormField>
			</div>

			<div class="flex flex-col gap-2">
				<UiFormField name="rewardIds" v-slot="{componentField}">
		      <UiFormItem>
						<UiFormLabel>{{ t('alerts.trigger.rewards') }}</UiFormLabel>
						<RewardsSelector v-bind="componentField" multiple />
						<UiFormMessage />
		      </UiFormItem>
				</UiFormField>
			</div>

			<div class="flex flex-col gap-2">
				<UiFormField name="keywordsIds" v-slot="{componentField}">
		      <UiFormItem>
						<UiFormLabel>{{ t('alerts.trigger.keywords') }}</UiFormLabel>
						<UiSelect v-bind="componentField" multiple>
							<UiSelectTrigger>
								<UiSelectValue placeholder="Select keywords" />
							</UiSelectTrigger>
							<UiSelectContent>
								<UiSelectItem
									v-for="option in keywordsSelectOptions"
									:key="option.value"
									:value="option.value"
								>
									{{ option.label }}
								</UiSelectItem>
							</UiSelectContent>
						</UiSelect>
						<UiFormMessage />
		      </UiFormItem>
				</UiFormField>
			</div>

			<div class="flex flex-col gap-2">
				<UiFormField name="greetingsIds" v-slot="{componentField}">
		      <UiFormItem>
						<UiFormLabel>{{ t('alerts.trigger.greetings') }}</UiFormLabel>
						<UiSelect v-bind="componentField" multiple>
							<UiSelectTrigger>
								<UiSelectValue placeholder="Select greetings" />
							</UiSelectTrigger>
							<UiSelectContent>
								<UiSelectItem
									v-for="option in greetingsSelectOptions"
									:key="option.value"
									:value="option.value"
								>
									{{ option.label }}
								</UiSelectItem>
							</UiSelectContent>
						</UiSelect>
						<UiFormMessage />
		      </UiFormItem>
				</UiFormField>
			</div>
		</div>

		<UiSeparator class="my-6" />

		<div class="flex flex-col gap-6">
			<AlertsDialogContentAudio
				v-model:audio-id="audioId"
				v-model:volume="audioVolume"
			/>

			<div class="flex justify-end">
				<UiButton type="submit">
					{{ t('sharedButtons.save') }}
				</UiButton>
			</div>
		</div>
	</form>
</template>
