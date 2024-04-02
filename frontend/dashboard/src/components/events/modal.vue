<script setup lang="ts">
import { IconTrash, IconGripVertical, IconPlus } from '@tabler/icons-vue';
import { EventOperationFilterType } from '@twir/types/events';
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NSpace,
	NSelect,
	NForm,
	NFormItem,
	NInput,
	NText,
	NGrid,
	NGridItem,
	NInputNumber,
	NDivider,
	NSwitch,
	NAlert,
	NButton,
	useThemeVars,
	NModal,
} from 'naive-ui';
import { computed, onMounted, ref, watch, nextTick } from 'vue';
import { VueDraggableNext } from 'vue-draggable-next';
import { useI18n } from 'vue-i18n';

import {
	eventTypeSelectOptions,
	operationTypeSelectOptions,
	getOperation,
	flatEvents,
} from './helpers.js';
import type { EditableEvent, EventOperation } from './types.js';

import {
	useAlertsManager,
	useCommandsManager,
	useEventsManager,
	useKeywordsManager,
	useObsOverlayManager,
	useProfile,
	useVariablesManager,
} from '@/api/index.js';
import AlertModal from '@/components/alerts/list.vue';
import rewardsSelector from '@/components/rewardsSelector.vue';

const themeVars = useThemeVars();
const selectedTabBackground = computed(() => themeVars.value.cardColor);

const props = defineProps<{
	event: EditableEvent | null
}>();
const emits = defineEmits<{
	saved: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableEvent>({
	description: '',
	enabled: true,
	onlineOnly: false,
	operations: [],
	type: '',
});

const selectedOperationsTab = ref(-1);
const currentOperation = ref<EventOperation | null>(null);

watch(selectedOperationsTab, (v) => {
	currentOperation.value = formValue.value.operations[v];
}, { immediate: true });

onMounted(() => {
	if (props.event) {
		formValue.value = props.event;

		if (props.event.operations.length) {
			selectedOperationsTab.value = 0;
		}
	}
});

watch(() => formValue.value.type, () => {
	nextTick(formRef.value?.validate);
});

const rules: FormRules = {
	type: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (!v) return new Error('Type required');

			return true;
		},
	},
	description: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (!v) return new Error('Description required');

			return true;
		},
	},
	input: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (!v) return new Error('Please type something');
			if (v?.length > 100) return new Error('Too long input');

			return true;
		},
	},
	timeoutMessage: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (v?.length > 100) return new Error('Too long message');

			return true;
		},
	},
	commandId: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (formValue.value.type !== 'COMMAND_USED') return true;
			if (!v) return new Error('Please select command');

			return true;
		},
	},
	rewardId: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (formValue.value.type !== 'REDEMPTION_CREATED') return true;
			if (!v) return new Error('Please select reward');

			return true;
		},
	},
	keywordId: {
		trigger: ['input', 'blur', 'focus'],
		validator: (_: FormItemRule, v: string) => {
			if (formValue.value.type !== 'KEYWORD_MATCHED') return true;
			if (!v) return new Error('Please select keyword');

			return true;
		},
	},
};

const obsManager = useObsOverlayManager();
const obsSettings = obsManager.getSettings();

const obsScenes = computed(() => {
	return obsSettings.data.value?.scenes.map(s => ({
		value: s,
		label: s,
	})) ?? [];
});
const obsSources = computed(() => {
	return obsSettings.data.value?.sources.map(s => ({
		value: s,
		label: s,
	})) ?? [];
});
const obsAudioSources = computed(() => {
	return obsSettings.data.value?.audioSources.map(s => ({
		value: s,
		label: s,
	})) ?? [];
});

const variablesManager = useVariablesManager();
const { data: variablesData, isLoading: isVariablesLoading } = variablesManager.getAll({});
const variablesSelectOptions = computed(() => {
	return variablesData.value?.variables.map((v) => ({
		label: v.name,
		value: v.id,
	})) ?? [];
});

const commandsManager = useCommandsManager();
const { data: commandsData, isLoading: isCommandsLoading } = commandsManager.getAll({});
const commandsSelectOptions = computed(() => {
	return commandsData.value?.commands.map(c => ({
		label: c.name,
		value: c.id,
	})) ?? [];
});


const keywordsManager = useKeywordsManager();
const { data: keywordsData, isLoading: isKeywordsLoading } = keywordsManager.getAll({});
const keywordsSelectOptions = computed(() => {
	return keywordsData.value?.keywords.map(k => ({
		label: k.text,
		value: k.id,
	}));
});

const { t } = useI18n();

const addOperation = () => {
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
	});
	selectedOperationsTab.value = newLength - 1;
};

const removeOperation = (index: number) => {
	if (index === selectedOperationsTab.value) selectedOperationsTab.value = 0;
	formValue.value.operations = formValue.value.operations.filter((_, i) => i != index);
};

const eventsManager = useEventsManager();
const eventsUpdater = eventsManager.update;
const eventsCreator = eventsManager.create;

const { data: profile } = useProfile();

async function save() {
	if (!formRef.value || !profile.value) return;
	await formRef.value.validate();

	const event = {
		...formValue.value,
		channelId: profile.value.selectedDashboardId,
		id: formValue.value.id ?? '',
	};

	if (!formValue.value.id) {
		await eventsCreator.mutateAsync({
			event,
		});
	} else {
		await eventsUpdater.mutateAsync({
			id: formValue.value.id,
			event,
		});
	}

	emits('saved');
}

const manager = useAlertsManager();
const { data: alerts } = manager.getAll({});

const showAlertModal = ref(false);

function getOperationLabel(type: string): string {
	switch (type) {
		case 'SEND_MESSAGE':
			return t('events.operations.inputs.message');
		case 'UNVIP_RANDOM_IF_NO_SLOTS':
			return t('events.operations.inputs.vipSlots');
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
			return t('events.operations.inputs.username');
		case 'CHANGE_VARIABLE':
			return t('events.operations.inputs.variableValue');
		case 'INCREMENT_VARIABLE':
		case 'DECREMENT_VARIABLE':
			return t('events.operations.inputs.variableIncrementDecrement');
		default:
			return t('events.operations.inputs.default');
	}
}

const eventsOperationsFiltersTypes = Object.values(EventOperationFilterType).map((item) => ({
	label: item.toLowerCase().split('_').join(' '),
	value: item,
}));
</script>

<template>
	<n-form ref="formRef" :model="formValue" :rules="rules">
		<n-space vertical>
			<n-space justify="space-between" item-style="width: 49%">
				<n-space vertical item-style="width: 100%">
					<n-form-item :label="t('events.operations.name')" path="type" show-require-mark>
						<n-select v-model:value="formValue.type" filterable :options="eventTypeSelectOptions" />
					</n-form-item>

					<n-form-item :label="t('events.description')" path="description" show-require-mark>
						<n-input
							v-model:value="formValue.description"
							type="textarea"
							:autosize="{
								minRows: 1,
								maxRows: 5
							}"
							:maxlength="500"
						/>
					</n-form-item>

					<n-form-item
						v-if="formValue.type === 'COMMAND_USED'"
						:label="t('events.targetCommand')"
						required
						path="commandId"
					>
						<n-select
							v-model:value="formValue.commandId"
							:options="commandsSelectOptions"
							:placeholder="t('events.targetCommand')"
							:loading="isCommandsLoading"
							filterable
						/>
					</n-form-item>

					<n-form-item
						v-if="formValue.type === 'REDEMPTION_CREATED'"
						:label="t('events.targetTwitchReward')"
						required
						path="rewardId"
					>
						<rewards-selector v-model="formValue.rewardId" />
					</n-form-item>

					<n-form-item
						v-if="formValue.type === 'KEYWORD_MATCHED'"
						:label="t('events.targetKeyword')"
						required
						path="keywordId"
					>
						<n-select
							v-model:value="formValue.keywordId"
							:options="keywordsSelectOptions"
							:placeholder="t('events.targetKeyword')"
							:loading="isKeywordsLoading"
							filterable
						/>
					</n-form-item>

					<n-form-item :label="t('events.onlineOnly')">
						<n-switch v-model:value="formValue.onlineOnly" />
					</n-form-item>
				</n-space>

				<n-space vertical>
					<n-text
						v-for="(variable, variableIndex) of flatEvents[formValue.type]?.variables"
						:key="variableIndex"
					>
						{{ '{' + variable + '}' }} - {{ t(`events.variables.${variable}`) }}
					</n-text>
				</n-space>
			</n-space>

			<n-divider title-placement="center">
				{{ t('events.operations.divider') }}
			</n-divider>
		</n-space>


		<n-space :wrap="false">
			<n-space vertical class="h-full" :x-gap="5">
				<VueDraggableNext v-model="formValue.operations">
					<div
						v-for="(operation, operationIndex) of formValue.operations"
						:key="operationIndex"
						style="display:flex; gap: 5px; margin-top: 5px; width: 100%; padding: 5px; border-radius: 11px;"
						:style="{
							'background-color': selectedOperationsTab === operationIndex ? selectedTabBackground : undefined,
						}"
					>
						<n-button text>
							<IconGripVertical class="w-4" />
						</n-button>

						<n-button
							secondary
							size="small"
							style="flex-grow: 1;"
							:type="getOperation(operation.type)?.color ?? 'default'"
							@click="() => selectedOperationsTab = operationIndex"
						>
							{{ getOperation(operation.type)?.name.slice(0, 15) ?? '' }}
						</n-button>

						<n-button text>
							<IconTrash
								class="w-[18px] flex"
								@click="removeOperation(operationIndex)"
							/>
						</n-button>
					</div>
				</VueDraggableNext>
				<n-button
					block
					size="small"
					secondary
					:disabled="formValue.operations.length >= 10"
					@click="addOperation"
				>
					<IconPlus />
				</n-button>
			</n-space>

			<n-divider vertical class="h-full" />

			<div v-if="currentOperation">
				<n-space vertical class="gap-0">
					<n-grid cols="3 s:1 m:3" :x-gap="5" :y-gap="5" responsive="screen">
						<n-grid-item :span="2">
							<n-form-item :label="t('events.operations.name')" required>
								<n-select
									v-model:value="currentOperation.type" filterable
									:options="operationTypeSelectOptions"
								/>
							</n-form-item>
						</n-grid-item>
						<n-grid-item :span="1">
							<n-form-item :label="t('sharedTexts.status')">
								<n-switch v-model:value="currentOperation.enabled" />
							</n-form-item>
						</n-grid-item>
						<n-grid-item :span="1">
							<n-form-item :label="t('events.delay')">
								<n-input-number v-model:value="currentOperation.delay" :min="0" :max="1800" />
							</n-form-item>
						</n-grid-item>
						<n-grid-item :span="1">
							<n-form-item :label="t('events.repeat')">
								<n-input-number v-model:value="currentOperation.repeat" :min="0" :max="10" />
							</n-form-item>
						</n-grid-item>
					</n-grid>

					<n-divider title-placement="left" class="mt-0">
						{{ t('events.operations.values') }}
					</n-divider>

					<n-form-item
						v-if="getOperation(currentOperation.type)?.haveInput"
						:label="getOperationLabel(currentOperation.type)"
					>
						<n-input v-model:value="currentOperation.input" :maxlength="500" />
					</n-form-item>

					<n-form-item v-if="currentOperation.type === 'SEND_MESSAGE'" label="Use announce">
						<n-switch v-model:value="currentOperation.useAnnounce" />
					</n-form-item>

					<n-grid cols="4 s:1 m:4" :x-gap="5" :y-gap="5" responsive="screen">
						<n-grid-item :span="3">
							<n-form-item
								v-if="['TIMEOUT', 'TIMEOUT_RANDOM', 'BAN', 'BAN_RANDOM'].some(v => currentOperation!.type === v)"
								:label="t('events.operations.banMessage')"
								:path="`operations[${selectedOperationsTab}].timeoutMessage`"
								:rule="rules.timeoutMessage"
							>
								<n-input v-model:value="currentOperation.timeoutMessage" :maxlength="500" />
							</n-form-item>
						</n-grid-item>

						<n-grid-item :span="1">
							<n-form-item
								v-if="['TIMEOUT', 'TIMEOUT_RANDOM'].some(v => currentOperation!.type === v)"
								:label="t('events.operations.banTime')"
							>
								<n-input-number v-model:value="currentOperation.timeoutTime" />
							</n-form-item>
						</n-grid-item>

						<n-grid-item
							v-if="currentOperation.type.startsWith('OBS')
								&& (!obsSettings.data.value?.isConnected || !obsSettings.data.value?.serverPassword)
							"
							:span="4"
						>
							<n-alert :title="t('events.operations.obs.warningTitle')" type="error">
								{{ t('events.operations.obs.warningText') }}
							</n-alert>
						</n-grid-item>

						<n-grid-item v-if="currentOperation.type === 'OBS_SET_SCENE'" :span="2">
							<n-form-item :label="t('events.operations.obs.scene')">
								<n-select
									v-model:value="currentOperation.target"
									:options="obsScenes"
									:placeholder="t('events.operations.obs.scene')"
									:disabled="!obsSettings.data.value?.isConnected"
								/>
							</n-form-item>
						</n-grid-item>

						<n-grid-item v-if="currentOperation.type === 'OBS_TOGGLE_SOURCE'" :span="2">
							<n-form-item :label="t('events.operations.obs.source')">
								<n-select
									v-model:value="currentOperation.target"
									:options="obsSources"
									:placeholder="t('events.operations.obs.source')"
									:disabled="!obsSettings.data.value?.isConnected"
								/>
							</n-form-item>
						</n-grid-item>

						<n-grid-item
							v-if="[
								'OBS_TOGGLE_AUDIO',
								'OBS_AUDIO_SET_VOLUME',
								'OBS_AUDIO_DECREASE_VOLUME',
								'OBS_AUDIO_INCREASE_VOLUME',
								'OBS_ENABLE_AUDIO',
								'OBS_DISABLE_AUDIO'
							].some(v => v === currentOperation!.type)"
							:span="2"
						>
							<n-form-item :label="t('events.operations.obs.audioSource')">
								<n-select
									v-model:value="currentOperation.target"
									:options="obsAudioSources"
									:placeholder="t('events.operations.obs.audioSource')"
									:disabled="!obsSettings.data.value?.isConnected"
								/>
							</n-form-item>
						</n-grid-item>

						<n-grid-item
							v-if="currentOperation.type === 'TRIGGER_ALERT'"
							:span="2"
						>
							<n-form-item :label="t('events.operations.triggerAlert')">
								<div class="flex gap-2.5 w-[90%]">
									<n-button block type="info" @click="showAlertModal = true">
										{{
											alerts?.alerts.find(a => a.id === currentOperation!.target)?.name ?? t('sharedButtons.select')
										}}
									</n-button>
									<n-button
										:disabled="!currentOperation!.target"
										text
										type="error"
										@click="currentOperation!.target = undefined"
									>
										<IconTrash />
									</n-button>
								</div>
							</n-form-item>
						</n-grid-item>

						<n-grid-item
							v-if="currentOperation.type.endsWith('VARIABLE')"
							:span="2"
						>
							<n-form-item :label="t('events.targetVariable')">
								<n-select
									v-model:value="currentOperation.target"
									:options="variablesSelectOptions"
									:placeholder="t('events.targetVariable')"
									:loading="isVariablesLoading"
								/>
							</n-form-item>
						</n-grid-item>
					</n-grid>

					<n-divider title-placement="left" class="mt-0">
						{{ t('events.operations.filters.label') }}
					</n-divider>

					<div class="flex flex-col gap-2">
						<n-text :depth="3">
							{{ t('events.operations.filters.description') }}
						</n-text>

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
							<n-input
								v-model:value="currentOperation.filters[index].left"
								:placeholder="t('events.operations.filters.placeholderLeft')"
								:maxlength="50"
							/>
							<n-select
								v-model:value="currentOperation.filters[index].type"
								:options="eventsOperationsFiltersTypes"
							/>
							<n-input
								v-model:value="currentOperation.filters[index].right"
								:placeholder="t('events.operations.filters.placeholderRight')"
								:maxlength="50"
							/>
							<div class="flex justify-end mt-2">
								<n-button
									type="error"
									secondary
									@click="currentOperation.filters.splice(index, 1)"
								>
									<IconTrash />
									{{ t('sharedButtons.delete') }}
								</n-button>
							</div>
						</div>
					</div>
					<n-button
						secondary
						type="info"
						block
						style="margin-top: 6px;"
						:disabled="currentOperation.filters.length >= 5"
						@click="currentOperation.filters.push({
							left: '',
							right: '',
							type: EventOperationFilterType.EQUALS,
						})"
					>
						{{ t('sharedButtons.create') }}
					</n-button>
				</n-space>
			</div>
		</n-space>

		<n-button block secondary type="success" class="mt-4" @click="save">
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-form>

	<n-modal
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
		<alert-modal
			:with-select="true"
			@select="(id) => {
				if (!currentOperation) return;
				currentOperation.target = id
				showAlertModal = false
			}"
			@delete="(id) => {
				if (currentOperation && id === currentOperation.target) {
					currentOperation.target = undefined
				}
			}"
		/>
	</n-modal>
</template>
