<script setup lang="ts">
import { IconPlayerPlay, IconTrash } from '@tabler/icons-vue'
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NButton,
	NDivider,
	NForm,
	NFormItem,
	NInput,
	NModal,
	NSelect,
	NSlider,
	NSpace
} from 'naive-ui'
import { storeToRefs } from 'pinia'
import { computed, onMounted, ref, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'

import rewardsSelector from '../rewardsSelector.vue'

import type { EditableAlert } from './types.js'

import {
	useAlertsManager,
	useFiles,
	useProfile,
	useTwitchGetUsers
} from '@/api'
import { useCommandsApi } from '@/api/commands/commands'
import { useGreetingsApi } from '@/api/greetings'
import { useKeywordsApi } from '@/api/keywords'
import FilesPicker from '@/components/files/files.vue'
import { playAudio } from '@/helpers/index.js'

const props = defineProps<{
	alert?: EditableAlert | null
}>()
const emits = defineEmits<{
	close: []
}>()

const formRef = ref<FormInst | null>(null)
const formValue = ref<EditableAlert>({
	id: '',
	name: '',
	audioId: undefined,
	audioVolume: 100,
	commandIds: [],
	rewardIds: [],
	greetingsIds: [],
	keywordsIds: []
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
		}
	}
}

const manager = useAlertsManager()
const creator = manager.create
const updater = manager.update

async function save() {
	if (!formRef.value || !formValue.value) return
	await formRef.value.validate()

	const data = formValue.value

	if (data.id) {
		await updater.mutateAsync({
			...data,
			id: data.id!
		})
	} else {
		await creator.mutateAsync(data)
	}

	emits('close')
}

const { data: files } = useFiles()
const selectedAudio = computed(() => files.value?.files.find(f => f.id === formValue.value.audioId))
const showAudioModal = ref(false)

const { data: profile } = storeToRefs(useProfile())

async function testAudio() {
	if (!selectedAudio.value?.id || !profile.value) return

	const query = new URLSearchParams({
		channel_id: profile.value.selectedDashboardId,
		file_id: selectedAudio.value.id
	})

	const req = await fetch(`${window.location.origin}/api-old/files/?${query}`)
	if (!req.ok) {
		console.error(await req.text())
		return
	}

	await playAudio(await req.arrayBuffer(), formValue.value.audioVolume)
}

const commandsManager = useCommandsApi()
const { data: commands } = commandsManager.useQueryCommands()
const commandsSelectOptions = computed(() => commands.value?.commands
	.map((command) => ({ label: command.name, value: command.id }))
)

const greetingsManager = useGreetingsApi()
const { data: greetings } = greetingsManager.useQueryGreetings()
const greetingsUsersIds = computed(() => greetings.value?.greetings.map(g => g.userId) ?? [])
const { data: twitchUsers } = useTwitchGetUsers({ ids: greetingsUsersIds })
const greetingsSelectOptions = computed(() => {
	if (!greetingsUsersIds.value.length || !twitchUsers.value?.users.length) return []
	return greetings.value?.greetings.map(g => {
		const twitchUser = twitchUsers.value.users.find(u => u.id === g.userId)
		return { label: twitchUser?.login ?? g.userId, value: g.id }
	})
})

const keywordsManager = useKeywordsApi()
const { data: keywords } = keywordsManager.useQueryKeywords()
const keywordsSelectOptions = computed(() => keywords.value?.keywords
	.map(k => ({ label: k.text, value: k.id }))
)
</script>

<template>
	<NForm
		ref="formRef"
		:model="formValue"
		:rules="rules"
	>
		<NSpace vertical class="w-full">
			<NFormItem label="Name" path="name" show-require-mark>
				<NInput v-model:value="formValue.name" :maxlength="30" />
			</NFormItem>

			<NDivider />

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
				<rewardsSelector v-model="formValue.rewardIds" multiple />
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
				<div class="flex gap-2.5 w-[85%]">
					<NButton class="overflow-hidden text-nowrap" block type="info" @click="showAudioModal = true">
						{{ selectedAudio?.name ?? t('sharedButtons.select') }}
					</NButton>
					<NButton
						:disabled="!formValue.audioId" text type="error"
						@click="formValue.audioId = undefined"
					>
						<IconTrash />
					</NButton>
					<NButton :disabled="!formValue.audioId" text type="info" @click="testAudio">
						<IconPlayerPlay />
					</NButton>
				</div>
			</NFormItem>

			<NFormItem :label="t('alerts.audioVolume', { volume: formValue.audioVolume })">
				<NSlider
					v-model:value="formValue.audioVolume"
					:step="1"
					:min="1"
					:max="100"
					:marks="{ 1: '1', 100: '100' }"
					:show-tooltip="false"
					:tooltip="false"
				/>
			</NFormItem>

			<NFormItem :label="t('alerts.select.image')">
				<NButton block type="info" disabled>
					Soon...
				</NButton>
			</NFormItem>

			<NFormItem :label="t('alerts.select.text')">
				<NButton block type="info" disabled>
					Soon...
				</NButton>
			</NFormItem>
		</NSpace>

		<NButton secondary type="success" block @click="save">
			{{ t('sharedButtons.save') }}
		</NButton>
	</NForm>

	<NModal
		v-model:show="showAudioModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="t('alerts.select.audio')"
		class="modal"
		:style="{
			width: '1000px',
			top: '50px',
		}"
		:on-close="() => showAudioModal = false"
	>
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
	</NModal>
</template>
