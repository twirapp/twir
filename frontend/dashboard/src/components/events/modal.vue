<script setup lang='ts'>
import { IconTrash, IconGripVertical, IconPlus } from '@tabler/icons-vue';
import {
	type SelectOption,
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
	NAvatar,
	NButton,
} from 'naive-ui';
import { h, computed, onMounted, ref, watch, nextTick, type VNodeChild } from 'vue';
import { VueDraggableNext } from 'vue-draggable-next';
import { useI18n } from 'vue-i18n';

import { eventTypeSelectOptions, operationTypeSelectOptions, getOperation, flatEvents } from './helpers.js';
import type { EditableEvent, EventOperation } from './types.js';

import { useCommandsManager, useKeywordsManager, useObsOverlayManager, useTwitchRewards, useVariablesManager } from '@/api';

const props = defineProps<{
	event: EditableEvent | null
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableEvent>({
	description: '',
	enabled: true,
	onlineOnly: false,
	operations: [],
	type: '',
});

const selectedOperationsTab = ref(0);
const currentOperation = ref<EventOperation | null>(null);

watch(selectedOperationsTab, (v) => {
	currentOperation.value = formValue.value.operations[v];
}, { immediate: true });

onMounted(() => {
	if (props.event) {
		formValue.value = props.event;

		if (props.event.operations.length) {
			currentOperation.value = props.event.operations.at(0)!;
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

const { data: rewardsData, isLoading: isRewardsLoading, isError: isRewardsError } = useTwitchRewards();
const rewardsSelectOptions = computed(() => {
	return rewardsData.value?.rewards.map(r => ({
		value: r.id,
		label: r.title,
		image: r.image?.url4X,
	})) ?? [];
});
const renderRewardTag = (option: SelectOption & { image?: string }): VNodeChild => {
	return h(NSpace, { align: 'center' }, {
		default: () => [
			h(NAvatar, { src: option.image, round: true, size: 'small', style: 'display: flex;' }),
			h(NText, { }, { default: () =>  option.label }),
		],
	});
};

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
	formValue.value.operations.push({
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
};

const removeOperation = (index: number) => {
	formValue.value.operations = formValue.value.operations.filter((_, i) => i != index);
};
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
						<n-input v-model:value="formValue.description" type="textarea" />
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
						/>
					</n-form-item>

					<n-form-item
						v-if="formValue.type === 'REDEMPTION_CREATED'"
						:label="t('events.targetTwitchReward')"
						required
						path="rewardId"
					>
						<n-select
							v-model:value="formValue.rewardId"
							size="large"
							:options="rewardsSelectOptions"
							:placeholder="t('events.targetTwitchReward')"
							:loading="isRewardsLoading"
							:render-label="renderRewardTag"
							:disabled="isRewardsError"
						/>
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
						/>
					</n-form-item>
				</n-space>

				<n-space vertical>
					<n-text
						v-for="(variable, variableIndex) of flatEvents[formValue.type]?.variables"
						:key="variableIndex"
					>
						{{ variable }}
					</n-text>
				</n-space>
			</n-space>

			<n-divider title-placement="center">
				{{ t('events.operations.divider') }}
			</n-divider>
		</n-space>


		<n-space :wrap="false">
			<n-space vertical style="height:100%" :x-gap="5">
				<VueDraggableNext v-model="formValue.operations">
					<div
						v-for="(operation, operationIndex) of formValue.operations"
						:key="operationIndex"
						style="display:flex; gap: 5px; margin-top: 5px; width: 100%"
					>
						<n-button text>
							<IconGripVertical style="width: 18px" />
						</n-button>

						<n-button
							secondary
							size="small"
							style="flex-grow: 1;"
							@click="() => selectedOperationsTab = operationIndex"
						>
							{{ getOperation(operation.type)?.name.slice(0, 15) ?? '' }}
						</n-button>

						<n-button text>
							<IconTrash style="width: 18px; display: flex" @click="removeOperation(operationIndex)" />
						</n-button>
					</div>
				</VueDraggableNext>
				<n-button
					block
					size="small"
					secondary
					@click="addOperation"
				>
					<IconPlus />
				</n-button>
			</n-space>

			<n-divider vertical style="height:100%" />

			<div v-if="currentOperation">
				<n-space vertical style="gap: 0">
					<n-grid cols="3 s:1 m:3" :x-gap="5" :y-gap="5" responsive="screen">
						<n-grid-item :span="2">
							<n-form-item :label="t('events.operations.name')" required>
								<n-select v-model:value="currentOperation.type" :options="operationTypeSelectOptions" />
							</n-form-item>
						</n-grid-item>
						<n-grid-item :span="1">
							<n-form-item :label="t('events.delay')">
								<n-input-number v-model:value="currentOperation.delay" :min="0" :max="10" />
							</n-form-item>
						</n-grid-item>
						<n-grid-item :span="1">
							<n-form-item :label="t('events.repeat')">
								<n-input-number v-model:value="currentOperation.repeat" :min="0" :max="10" />
							</n-form-item>
						</n-grid-item>
					</n-grid>

					<n-divider title-placement="left" style="margin-top: 0px">
						{{ t('events.operations.values') }}
					</n-divider>

					<n-form-item
						v-if="getOperation(currentOperation.type)?.haveInput"
						:label="t('events.operations.input')"
						:path="`operations[${selectedOperationsTab}].input`"
						:rule="rules.input"
					>
						<n-input v-model:value="currentOperation.input" />
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
								<n-input v-model:value="currentOperation.timeoutMessage" />
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
				</n-space>
			</div>
		</n-space>
	</n-form>
</template>
