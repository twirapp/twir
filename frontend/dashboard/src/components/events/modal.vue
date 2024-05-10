<script setup lang="ts">
import { dragAndDrop } from '@formkit/drag-and-drop/vue'
import { IconGripVertical, IconPlus, IconTrash } from '@tabler/icons-vue'
import { EventOperationFilterType } from '@twir/types/events'
import { toRef } from '@vueuse/core'
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NAlert,
	NButton,
	NDivider,
	NForm,
	NFormItem,
	NGrid,
	NGridItem,
	NInput,
	NInputNumber,
	NModal,
	NSelect,
	NSpace,
	NSwitch,
	NText,
	useThemeVars,
} from 'naive-ui'
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import {
	eventTypeSelectOptions,
	flatEvents,
	getOperation,
	operationTypeSelectOptions,
} from './helpers.js'
import Table from '../table.vue'

import type { EditableEvent, EventOperation } from './types.js'

import { useAlertsQuery } from '@/api/alerts'
import { useCommandsApi } from '@/api/commands/commands'
import {
	useEventsManager,
	useObsOverlayManager,
	useProfile,
} from '@/api/index.js'
import { useKeywordsApi } from '@/api/keywords.js'
import { useVariablesApi } from '@/api/variables.js'
import { OPERATIONS } from '@/components/events/operations'
import rewardsSelector from '@/components/rewardsSelector.vue'
import { useAlertsTable } from '@/features/alerts/composables/use-alerts-table.js'

const props = defineProps<{
	event: EditableEvent | null
}>()
const emits = defineEmits<{
	saved: []
}>()
const themeVars = useThemeVars()
const selectedTabBackground = computed(() => themeVars.value.cardColor)

const formRef = ref<FormInst | null>(null)
const formValue = ref<EditableEvent>({
	description: '',
	enabled: true,
	onlineOnly: false,
	operations: [],
	type: '',
})

const selectedOperationsTab = ref(-1)
const currentOperation = ref<EventOperation | null>(null)

watch(selectedOperationsTab, (v) => {
	currentOperation.value = formValue.value.operations[v]
}, { immediate: true })

onMounted(() => {
	if (props.event) {
		formValue.value = props.event

		if (props.event.operations.length) {
			selectedOperationsTab.value = 0
		}
	}
})

watch(() => formValue.value.type, () => {
	nextTick(formRef.value?.validate)
})

const rules: FormRules = {
	type: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (!v) return new Error('Type required')

			return true
		},
	},
	description: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (!v) return new Error('Description required')

			return true
		},
	},
	input: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (!v) return new Error('Please type something')
			if (v?.length > 100) return new Error('Too long input')

			return true
		},
	},
	timeoutMessage: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (v?.length > 100) return new Error('Too long message')

			return true
		},
	},
	commandId: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (formValue.value.type !== 'COMMAND_USED') return true
			if (!v) return new Error('Please select command')

			return true
		},
	},
	rewardId: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (formValue.value.type !== 'REDEMPTION_CREATED') return true
			if (!v) return new Error('Please select reward')

			return true
		},
	},
	keywordId: {
		trigger: ['input', 'blur', 'focus'],
		validator: (_: FormItemRule, v: string) => {
			if (formValue.value.type !== 'KEYWORD_MATCHED') return true
			if (!v) return new Error('Please select keyword')

			return true
		},
	},
}

const obsManager = useObsOverlayManager()
const obsSettings = obsManager.getSettings()

const obsScenes = computed(() => {
	return obsSettings.data.value?.scenes.map(s => ({
		value: s,
		label: s,
	})) ?? []
})
const obsSources = computed(() => {
	return obsSettings.data.value?.sources.map(s => ({
		value: s,
		label: s,
	})) ?? []
})
const obsAudioSources = computed(() => {
	return obsSettings.data.value?.audioSources.map(s => ({
		value: s,
		label: s,
	})) ?? []
})

const { customVariables: variables } = useVariablesApi()
const variablesSelectOptions = computed(() => {
	return variables.value.map((v) => ({
		label: v.name,
		value: v.id,
	})) ?? []
})

const commandsManager = useCommandsApi()
const { data: commandsData, fetching: isCommandsLoading } = commandsManager.useQueryCommands()
const commandsSelectOptions = computed(() => {
	return commandsData.value?.commands.map(c => ({
		label: c.name,
		value: c.id,
	})) ?? []
})

const keywordsManager = useKeywordsApi()
const { data: keywordsData, fetching: isKeywordsLoading } = keywordsManager.useQueryKeywords()
const keywordsSelectOptions = computed(() => {
	return keywordsData.value?.keywords.map(k => ({
		label: k.text,
		value: k.id,
	}))
})

const { data: alerts } = useAlertsQuery()

const { t } = useI18n()

function addOperation() {
	const newLength = formValue.value.operations.push({
		delay: 0,
		enabled: true,
		filters: [],
		repeat: 0,
		timeoutTime: 0,
		type: 'SEND_MESSAGE',
		useAnnounce: false,
		input: '',
		target: '',
		timeoutMessage: '',
	})
	selectedOperationsTab.value = newLength - 1
}

function removeOperation(index: number) {
	if (index === selectedOperationsTab.value) {
		selectedOperationsTab.value = 0
	}
	formValue.value.operations = formValue.value.operations.filter((_, i) => i !== index)

	if (!formValue.value.operations.length) {
		currentOperation.value = null
	}
}

const eventsManager = useEventsManager()
const eventsUpdater = eventsManager.update
const eventsCreator = eventsManager.create

const { data: profile } = useProfile()

async function save() {
	if (!formRef.value || !profile.value) return
	await formRef.value.validate()

	const event = {
		...formValue.value,
		channelId: profile.value.selectedDashboardId,
		id: formValue.value.id ?? '',
	}

	if (!formValue.value.id) {
		await eventsCreator.mutateAsync({
			event,
		})
	} else {
		await eventsUpdater.mutateAsync({
			id: formValue.value.id,
			event,
		})
	}

	emits('saved')
}

const showAlertModal = ref(false)

function getOperationLabel(type: string): string {
	switch (type) {
		case 'SEND_MESSAGE':
			return t('events.operations.inputs.message')
		case 'UNVIP_RANDOM_IF_NO_SLOTS':
			return t('events.operations.inputs.vipSlots')
		case 'BAN':
		case 'UNBAN':
		case 'TIMEOUT':
		case 'VIP':
		case 'UNVIP':
		case 'MOD':
		case 'UNMOD':
		case 'ALLOW_COMMAND_TO_USER':
		case 'REMOVE_ALLOW_COMMAND_TO_USER':
		case 'DENY_COMMAND_TO_USER':
		case 'REMOVE_DENY_COMMAND_TO_USER':
			return t('events.operations.inputs.username')
		case 'CHANGE_VARIABLE':
			return t('events.operations.inputs.variableValue')
		case 'INCREMENT_VARIABLE':
		case 'DECREMENT_VARIABLE':
			return t('events.operations.inputs.variableIncrementDecrement')
		default:
			return t('events.operations.inputs.default')
	}
}

const eventsOperationsFiltersTypes = Object.values(EventOperationFilterType).map((item) => ({
	label: item.toLowerCase().split('_').join(' '),
	value: item,
}))

const operations = toRef(formValue.value, 'operations')

const dragParentRef = ref<HTMLElement>()
dragAndDrop({
	parent: dragParentRef,
	values: operations,
	dragHandle: '.drag-handle',
})

function variableText(variable: string) {
	return `{${variable}}`
}

const alertsTable = useAlertsTable({
	onSelect(alert) {
		if (!currentOperation.value) return
		currentOperation.value.target = alert.id
		showAlertModal.value = false
	},
})

const filteredOperationTypeSelectOptions = computed(() => {
	return operationTypeSelectOptions.filter(option => {
		if (OPERATIONS[option.value as string].dependsOnEvents) {
			return OPERATIONS[option.value as string].dependsOnEvents!.includes(formValue.value.type)
		}

		return true
	})
})
</script>

<template>
	<NForm ref="formRef" :model="formValue" :rules="rules">
		<NSpace vertical>
			<NSpace justify="space-between" item-style="width: 49%">
				<NSpace vertical item-style="width: 100%">
					<NFormItem :label="t('events.operations.name')" path="type" show-require-mark>
						<NSelect v-model:value="formValue.type" filterable :options="eventTypeSelectOptions" />
					</NFormItem>

					<NFormItem :label="t('events.description')" path="description" show-require-mark>
						<NInput
							v-model:value="formValue.description"
							type="textarea"
							:autosize="{
								minRows: 1,
								maxRows: 5,
							}"
							:maxlength="500"
						/>
					</NFormItem>

					<NFormItem
						v-if="formValue.type === 'COMMAND_USED'"
						:label="t('events.targetCommand')"
						required
						path="commandId"
					>
						<NSelect
							v-model:value="formValue.commandId"
							:options="commandsSelectOptions"
							:placeholder="t('events.targetCommand')"
							:loading="isCommandsLoading"
							filterable
						/>
					</NFormItem>

					<NFormItem
						v-if="formValue.type === 'REDEMPTION_CREATED'"
						:label="t('events.targetTwitchReward')"
						required
						path="rewardId"
					>
						<rewards-selector v-model="formValue.rewardId" />
					</NFormItem>

					<NFormItem
						v-if="formValue.type === 'KEYWORD_MATCHED'"
						:label="t('events.targetKeyword')"
						required
						path="keywordId"
					>
						<NSelect
							v-model:value="formValue.keywordId"
							:options="keywordsSelectOptions"
							:placeholder="t('events.targetKeyword')"
							:loading="isKeywordsLoading"
							filterable
						/>
					</NFormItem>

					<NFormItem :label="t('events.onlineOnly')">
						<NSwitch v-model:value="formValue.onlineOnly" />
					</NFormItem>
				</NSpace>

				<NSpace vertical>
					<NText
						v-for="(variable, variableIndex) of flatEvents[formValue.type]?.variables"
						:key="variableIndex"
					>
						{{ variableText(variable) }} - {{ t(`events.variables.${variable}`) }}
					</NText>
				</NSpace>
			</NSpace>

			<NDivider title-placement="center">
				{{ t('events.operations.divider') }}
			</NDivider>
		</NSpace>

		<NSpace :wrap="false">
			<NSpace vertical class="h-full" :x-gap="5">
				<div ref="dragParentRef">
					<div
						v-for="(operation, operationIndex) of formValue.operations"
						:key="operation.type + operationIndex"
						style="display:flex; gap: 5px; margin-top: 5px; width: 100%; padding: 5px; border-radius: 11px;"
						:style="{
							'background-color': selectedOperationsTab === operationIndex ? selectedTabBackground : undefined,
						}"
					>
						<NButton text class="drag-handle">
							<IconGripVertical class="w-4" />
						</NButton>

						<NButton
							secondary
							size="small"
							style="flex-grow: 1;"
							:type="getOperation(operation.type)?.color ?? 'default'"
							@click="() => selectedOperationsTab = operationIndex"
						>
							{{ getOperation(operation.type)?.name.slice(0, 15) ?? '' }}
						</NButton>

						<NButton text>
							<IconTrash
								class="w-[18px] flex"
								@click="removeOperation(operationIndex)"
							/>
						</NButton>
					</div>
				</div>

				<NButton
					block
					size="small"
					secondary
					:disabled="formValue.operations.length >= 10"
					@click="addOperation"
				>
					<IconPlus />
				</NButton>
			</NSpace>

			<NDivider vertical class="h-full" />

			<div v-if="currentOperation">
				<NSpace vertical class="gap-0">
					<NGrid cols="3 s:1 m:3" :x-gap="5" :y-gap="5" responsive="screen">
						<NGridItem :span="2">
							<NFormItem :label="t('events.operations.name')" required>
								<NSelect
									v-model:value="currentOperation.type"
									filterable
									:options="filteredOperationTypeSelectOptions"
								/>
							</NFormItem>
						</NGridItem>
						<NGridItem :span="1">
							<NFormItem :label="t('sharedTexts.status')">
								<NSwitch v-model:value="currentOperation.enabled" />
							</NFormItem>
						</NGridItem>
						<NGridItem :span="1">
							<NFormItem :label="t('events.delay')">
								<NInputNumber v-model:value="currentOperation.delay" :min="0" :max="1800" />
							</NFormItem>
						</NGridItem>
						<NGridItem :span="1">
							<NFormItem :label="t('events.repeat')">
								<NInputNumber v-model:value="currentOperation.repeat" :min="0" :max="10" />
							</NFormItem>
						</NGridItem>
					</NGrid>

					<NDivider title-placement="left" class="mt-0">
						{{ t('events.operations.values') }}
					</NDivider>

					<NFormItem
						v-if="getOperation(currentOperation.type)?.haveInput"
						:label="getOperationLabel(currentOperation.type)"
					>
						<NInput v-model:value="currentOperation.input" :maxlength="500" />
					</NFormItem>

					<NFormItem v-if="currentOperation.type === 'SEND_MESSAGE'" label="Use announce">
						<NSwitch v-model:value="currentOperation.useAnnounce" />
					</NFormItem>

					<NGrid cols="4 s:1 m:4" :x-gap="5" :y-gap="5" responsive="screen">
						<NGridItem :span="3">
							<NFormItem
								v-if="['TIMEOUT', 'TIMEOUT_RANDOM', 'BAN', 'BAN_RANDOM'].some(v => currentOperation!.type === v)"
								:label="t('events.operations.banMessage')"
								:path="`operations[${selectedOperationsTab}].timeoutMessage`"
								:rule="rules.timeoutMessage"
							>
								<NInput v-model:value="currentOperation.timeoutMessage" :maxlength="500" />
							</NFormItem>
						</NGridItem>

						<NGridItem :span="1">
							<NFormItem
								v-if="['TIMEOUT', 'TIMEOUT_RANDOM'].some(v => currentOperation!.type === v)"
								:label="t('events.operations.banTime')"
							>
								<NInputNumber v-model:value="currentOperation.timeoutTime" />
							</NFormItem>
						</NGridItem>

						<NGridItem
							v-if="currentOperation.type.startsWith('OBS')
								&& (!obsSettings.data.value?.isConnected || !obsSettings.data.value?.serverPassword)
							"
							:span="4"
						>
							<NAlert :title="t('events.operations.obs.warningTitle')" type="error">
								{{ t('events.operations.obs.warningText') }}
							</NAlert>
						</NGridItem>

						<NGridItem v-if="currentOperation.type === 'OBS_SET_SCENE'" :span="2">
							<NFormItem :label="t('events.operations.obs.scene')">
								<NSelect
									v-model:value="currentOperation.target"
									:options="obsScenes"
									:placeholder="t('events.operations.obs.scene')"
									:disabled="!obsSettings.data.value?.isConnected"
								/>
							</NFormItem>
						</NGridItem>

						<NGridItem v-if="currentOperation.type === 'OBS_TOGGLE_SOURCE'" :span="2">
							<NFormItem :label="t('events.operations.obs.source')">
								<NSelect
									v-model:value="currentOperation.target"
									:options="obsSources"
									:placeholder="t('events.operations.obs.source')"
									:disabled="!obsSettings.data.value?.isConnected"
								/>
							</NFormItem>
						</NGridItem>

						<NGridItem
							v-if="[
								'OBS_TOGGLE_AUDIO',
								'OBS_AUDIO_SET_VOLUME',
								'OBS_AUDIO_DECREASE_VOLUME',
								'OBS_AUDIO_INCREASE_VOLUME',
								'OBS_ENABLE_AUDIO',
								'OBS_DISABLE_AUDIO',
							].some(v => v === currentOperation!.type)"
							:span="2"
						>
							<NFormItem :label="t('events.operations.obs.audioSource')">
								<NSelect
									v-model:value="currentOperation.target"
									:options="obsAudioSources"
									:placeholder="t('events.operations.obs.audioSource')"
									:disabled="!obsSettings.data.value?.isConnected"
								/>
							</NFormItem>
						</NGridItem>

						<NGridItem
							v-if="currentOperation.type === 'TRIGGER_ALERT'"
							:span="2"
						>
							<NFormItem :label="t('events.operations.triggerAlert')">
								<div class="flex gap-2.5 w-[90%]">
									<NButton block type="info" @click="showAlertModal = true">
										{{
											alerts?.channelAlerts.find(a => a.id === currentOperation!.target)?.name ?? t('sharedButtons.select')
										}}
									</NButton>
									<NButton
										:disabled="!currentOperation!.target"
										text
										type="error"
										@click="currentOperation!.target = undefined"
									>
										<IconTrash />
									</NButton>
								</div>
							</NFormItem>
						</NGridItem>

						<NGridItem
							v-if="currentOperation.type.endsWith('VARIABLE')"
							:span="2"
						>
							<NFormItem :label="t('events.targetVariable')">
								<NSelect
									v-model:value="currentOperation.target"
									:options="variablesSelectOptions"
									:placeholder="t('events.targetVariable')"
								/>
							</NFormItem>
						</NGridItem>
					</NGrid>

					<NDivider title-placement="left" class="mt-0">
						{{ t('events.operations.filters.label') }}
					</NDivider>

					<div class="flex flex-col gap-2">
						<NText :depth="3">
							{{ t('events.operations.filters.description') }}
						</NText>

						<div
							v-if="!currentOperation.filters.length"
							class="flex p-4 border border-yellow-700 justify-center rounded items-center"
						>
							{{ t('events.operations.filters.empty') }}
						</div>

						<div
							v-for="(_, index) of currentOperation.filters"
							:key="index"
							class="flex flex-col gap-0.5 border border-zinc-600 p-2 rounded"
						>
							<NInput
								v-model:value="currentOperation.filters[index].left"
								:placeholder="t('events.operations.filters.placeholderLeft')"
								:maxlength="50"
							/>
							<NSelect
								v-model:value="currentOperation.filters[index].type"
								:options="eventsOperationsFiltersTypes"
							/>
							<NInput
								v-model:value="currentOperation.filters[index].right"
								:placeholder="t('events.operations.filters.placeholderRight')"
								:maxlength="50"
							/>
							<div class="flex justify-end mt-2">
								<NButton
									type="error"
									secondary
									@click="currentOperation?.filters.splice(index, 1)"
								>
									<IconTrash />
									{{ t('sharedButtons.delete') }}
								</NButton>
							</div>
						</div>
					</div>
					<NButton
						secondary
						type="info"
						block
						style="margin-top: 6px;"
						:disabled="currentOperation.filters.length >= 5"
						@click="currentOperation?.filters.push({
							left: '',
							right: '',
							type: EventOperationFilterType.EQUALS,
						})"
					>
						{{ t('sharedButtons.create') }}
					</NButton>
				</NSpace>
			</div>
		</NSpace>

		<NButton block secondary type="success" style="margin-top: 12px;" @click="save">
			{{ t('sharedButtons.save') }}
		</NButton>
	</NForm>

	<NModal
		v-model:show="showAlertModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Select alert"
		class="modal"
		:style="{
			width: '1000px',
			top: '50px',
		}"
		:on-close="() => showAlertModal = false"
	>
		<Table
			:table="alertsTable.table"
			:is-loading="alertsTable.isLoading.value"
		/>
	</NModal>
</template>
