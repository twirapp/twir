<script setup lang="ts">
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NDivider,
	NForm,
	NFormItem,
	NSelect,
} from 'naive-ui'
import { computed, onMounted, ref, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'

import AlertsDialogContentAudio from './alerts-dialog-content-audio.vue'

import {
	useTwitchGetUsers,
} from '@/api'
import { type Alert, useAlertsCreateMutation, useAlertsUpdateMutation } from '@/api/alerts.js'
import { useCommandsApi } from '@/api/commands/commands.js'
import { useGreetingsApi } from '@/api/greetings.js'
import { useKeywordsApi } from '@/api/keywords.js'
import RewardsSelector from '@/components/rewardsSelector.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

const props = defineProps<{ alert?: Alert | null }>()
const emits = defineEmits<{ close: [] }>()

const formRef = ref<FormInst | null>(null)
const formValue = ref<Alert>({
	id: '',
	name: '',
	audioId: undefined,
	audioVolume: 75,
	commandIds: [],
	rewardIds: [],
	greetingsIds: [],
	keywordsIds: [],
})

onMounted(() => {
	if (!props.alert) return
	formValue.value = structuredClone(toRaw(props.alert))
})

const { t } = useI18n()

const rules: FormRules = {
	name: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length || value.length > 30) {
				return new Error(t('alerts.validations.name'))
			}

			return true
		},
	},
}

const alertsCreateMutation = useAlertsCreateMutation()
const alertsUpdateMutation = useAlertsUpdateMutation()

async function save() {
	if (!formRef.value || !formValue.value) return
	await formRef.value.validate()

	const data = { ...formValue.value, id: undefined }

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
}

const commandsApi = useCommandsApi()
const { data: commands } = commandsApi.useQueryCommands()
const commandsSelectOptions = computed(() => commands.value?.commands
	.map((command) => ({ label: command.name, value: command.id })),
)

const greetingsApi = useGreetingsApi()
const { data: greetings } = greetingsApi.useQueryGreetings()
const greetingsUsersIds = computed(() => greetings.value?.greetings.map(g => g.userId) ?? [])
const { data: twitchUsers } = useTwitchGetUsers({ ids: greetingsUsersIds })
const greetingsSelectOptions = computed(() => {
	if (!greetingsUsersIds.value.length || !twitchUsers.value?.users.length) return []
	return greetings.value?.greetings.map(g => {
		const twitchUser = twitchUsers.value.users.find(u => u.id === g.userId)
		return { label: twitchUser?.login ?? g.userId, value: g.id }
	})
})

const keywordsApi = useKeywordsApi()
const { data: keywords } = keywordsApi.useQueryKeywords()
const keywordsSelectOptions = computed(() => keywords.value?.keywords
	.map(k => ({ label: k.text, value: k.id })),
)
</script>

<template>
	<NForm
		ref="formRef"
		class="p-6 pt-0"
		:model="formValue"
		:rules="rules"
	>
		<div class="flex flex-col gap-6">
			<NFormItem label="Name" path="name" show-require-mark>
				<Input v-model="formValue.name" :maxlength="30" />
			</NFormItem>

			<NFormItem :label="t('alerts.trigger.commands')" path="commandIds">
				<NSelect
					v-model:value="formValue.commandIds"
					:fallback-option="false"
					filterable
					multiple
					:options="commandsSelectOptions"
				/>
			</NFormItem>

			<NFormItem :label="t('alerts.trigger.rewards')" path="rewardIds">
				<RewardsSelector v-model="formValue.rewardIds!" multiple />
			</NFormItem>

			<NFormItem :label="t('alerts.trigger.keywords')" path="rewardIds">
				<NSelect
					v-model:value="formValue.keywordsIds"
					:fallback-option="false"
					filterable
					multiple
					:options="keywordsSelectOptions"
				/>
			</NFormItem>

			<NFormItem :label="t('alerts.trigger.greetings')" path="rewardIds">
				<NSelect
					v-model:value="formValue.greetingsIds"
					:fallback-option="false"
					filterable
					multiple
					:options="greetingsSelectOptions"
				/>
			</NFormItem>
		</div>

		<NDivider />

		<div class="flex flex-col gap-6">
			<AlertsDialogContentAudio
				v-model:audio-id="formValue.audioId"
				:initialVolume="formValue.audioVolume"
				@update:volume="formValue.audioVolume = $event"
			/>

			<div class="flex justify-end">
				<Button @click="save">
					{{ t('sharedButtons.save') }}
				</Button>
			</div>
		</div>
	</NForm>
</template>
