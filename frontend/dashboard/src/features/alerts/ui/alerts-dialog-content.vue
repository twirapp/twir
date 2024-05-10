<script setup lang="ts">
import { PlayIcon, TrashIcon } from 'lucide-vue-next'
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NDivider,
	NForm,
	NFormItem,
	NScrollbar,
	NSelect,
	NSlider,
	NSpace,
} from 'naive-ui'
import { computed, onMounted, ref, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'

import {
	useFiles,
	useProfile,
	useTwitchGetUsers,
} from '@/api'
import { type Alert, useAlertsCreateMutation, useAlertsUpdateMutation } from '@/api/alerts.js'
import { useCommandsApi } from '@/api/commands/commands.js'
import { useGreetingsApi } from '@/api/greetings.js'
import { useKeywordsApi } from '@/api/keywords.js'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import FilesPicker from '@/components/files/files.vue'
import RewardsSelector from '@/components/rewardsSelector.vue'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { playAudio } from '@/helpers/playAudio.js'

const props = defineProps<{ alert?: Alert | null }>()
const emits = defineEmits<{ close: [] }>()

const { data: profile } = useProfile()

const formRef = ref<FormInst | null>(null)
const formValue = ref<Alert>({
	id: '',
	name: '',
	channelId: profile.value?.selectedDashboardId,
	audioId: undefined,
	audioVolume: 100,
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

const { data: files } = useFiles()
const selectedAudio = computed(() => {
	return files.value?.files
		.find((file) => file.id === formValue.value.audioId)
})
const showAudioModal = ref(false)

async function testAudio() {
	if (!selectedAudio.value?.id || !profile.value) return

	const query = new URLSearchParams({
		channel_id: profile.value.selectedDashboardId,
		file_id: selectedAudio.value.id,
	})

	const req = await fetch(`${window.location.origin}/api-old/files/?${query}`)
	if (!req.ok) {
		console.error(await req.text())
		return
	}

	await playAudio(await req.arrayBuffer(), formValue.value.audioVolume ?? 50)
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
		<NSpace vertical class="w-full">
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

			<NDivider />

			<NFormItem :label="t('alerts.select.audio')">
				<div class="flex gap-2 w-full">
					<Dialog v-model:open="showAudioModal" @update:open="showAudioModal = false">
						<DialogTrigger as-child>
							<Button class="w-full" @click="showAudioModal = true">
								{{ selectedAudio?.name ?? t('sharedButtons.select') }}
							</Button>
						</DialogTrigger>

						<DialogOrSheet class="p-0">
							<DialogHeader class="p-6 border-b-[1px]">
								<DialogTitle>
									{{ t('alerts.select.audio') }}
								</DialogTitle>
							</DialogHeader>

							<NScrollbar class="p-6 max-h-[85vh]" trigger="none">
								<FilesPicker
									mode="picker"
									tab="audios"
									@select="(id) => {
										formValue.audioId = id
										showAudioModal = false
									}"
									@delete="(id) => {
										if (id === formValue.audioId) {
											formValue.audioId = undefined
										}
									}"
								/>
							</NScrollbar>
						</DialogOrSheet>
					</Dialog>

					<Button
						class="min-w-10"
						size="icon"
						variant="secondary"
						:disabled="!formValue.audioId"
						@click="testAudio"
					>
						<PlayIcon class="size-4" />
					</Button>

					<Button
						class="min-w-10"
						size="icon"
						variant="destructive"
						:disabled="!formValue.audioId"
						@click="formValue.audioId = undefined"
					>
						<TrashIcon class="size-4" />
					</Button>
				</div>
			</NFormItem>

			<NFormItem :label="t('alerts.audioVolume', { volume: formValue.audioVolume })">
				<NSlider
					v-model:value="formValue.audioVolume!"
					:step="1"
					:min="1"
					:max="100"
					:marks="{ 1: '1', 100: '100' }"
					:show-tooltip="false"
					:tooltip="false"
				/>
			</NFormItem>
		</NSpace>

		<div class="flex justify-end">
			<Button @click="save">
				{{ t('sharedButtons.save') }}
			</Button>
		</div>
	</NForm>
</template>
