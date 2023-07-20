<script setup lang='ts'>
import {
	NSpace,
	NSelect,
	NForm,
	NFormItem,
	FormInst,
	FormItemRule,
	FormRules,
	NInput,
	NText,
	NTimeline,
	NTimelineItem,
	NGrid,
	NGridItem,
	NInputNumber,
	NDivider,
	NSwitch,
	NAlert,
} from 'naive-ui';
import { computed, onMounted, ref } from 'vue';

import { EVENTS } from './events.js';
import { eventTypeSelectOptions, operationTypeSelectOptions, getOperation } from './helpers.js';
import { EditableEvent } from './types.js';

import { useObsOverlayManager } from '@/api';

const props = defineProps<{
	event: EditableEvent | null
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableEvent>({
	description: '',
	enabled: true,
	onlineOnly: false,
	operations: [
		{
			delay: 0,
			enabled: true,
			filters: [],
			repeat: 0,
			timeoutTime: 0,
			timeoutMessage: '',
			type: 'SEND_MESSAGE',
			useAnnounce: false,
			input: '',
			target: '',
		},
	],
	type: '',
});

onMounted(() => {
	if (props.event) {
		formValue.value = props.event;
	}
});

const rules: FormRules = {
	type: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (!v) throw new Error('Type required');

			return true;
		},
	},
	description: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (!v) throw new Error('Description required');

			return true;
		},
	},
	input: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (v?.length > 1) throw new Error('Too long input');

			return true;
		},
	},
	timeoutMessage: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, v: string) => {
			if (v?.length > 1) throw new Error('Too long message');

			return true;
		},
	},
};

const availableEventVariables = computed(() => {
	const evt = EVENTS[formValue.value.type];

	return evt?.variables?.map(v => ({
		label: `{${v}}`,
		value: v,
	})) ?? [];
});

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
</script>

<template>
	<n-form ref="formRef" :model="formValue" :rules="rules">
		<n-space vertical>
			<n-space justify="space-between" item-style="width: 49%">
				<n-space vertical item-style="width: 100%">
					<n-form-item label="Type" path="type" show-require-mark>
						<n-select v-model:value="formValue.type" filterable :options="eventTypeSelectOptions" />
					</n-form-item>

					<n-form-item label="Description" path="description" show-require-mark>
						<n-input v-model:value="formValue.description" type="textarea" />
					</n-form-item>
				</n-space>

				<n-space vertical>
					<n-text v-for="variable of availableEventVariables" :key="variable.value">
						{{ variable }}
					</n-text>
				</n-space>
			</n-space>

			<n-timeline size="large">
				<n-timeline-item
					v-for="(operation, operationIndex) of formValue.operations"
					:key="operationIndex"
					:type="getOperation(operation.type)?.color ?? 'default'"
				>
					<n-space vertical style="gap: 0">
						<n-grid cols="3 s:1 m:3" :x-gap="5" :y-gap="5" responsive="screen">
							<n-grid-item>
								<n-form-item label="Operation">
									<n-select v-model:value="operation.type" :options="operationTypeSelectOptions" />
								</n-form-item>
							</n-grid-item>
							<n-grid-item>
								<n-form-item label="Delay">
									<n-input-number v-model:value="operation.delay" />
								</n-form-item>
							</n-grid-item>
							<n-grid-item>
								<n-form-item label="Repeat">
									<n-input-number v-model:value="operation.repeat" />
								</n-form-item>
							</n-grid-item>
						</n-grid>

						<n-form-item
							v-if="getOperation(operation.type)?.haveInput"
							label="Operation input"
							:path="`operations[${operationIndex}].input`"
							:rule="rules.input"
						>
							<n-input v-model:value="operation.input" />
						</n-form-item>

						<n-divider title-placement="left">
							Settings
						</n-divider>

						<n-form-item v-if="operation.type === 'SEND_MESSAGE'" label="Use announce">
							<n-switch v-model:value="operation.useAnnounce" />
						</n-form-item>

						<n-grid cols="4 s:1 m:4" :x-gap="5" :y-gap="5" responsive="screen">
							<n-grid-item :span="3">
								<n-form-item
									v-if="['TIMEOUT', 'TIMEOUT_RANDOM', 'BAN', 'BAN_RANDOM'].some(v => operation.type === v)"
									label="Ban message"
									:path="`operations[${operationIndex}].timeoutMessage`"
									:rule="rules.timeoutMessage"
								>
									<n-input v-model:value="operation.timeoutMessage" />
								</n-form-item>
							</n-grid-item>

							<n-grid-item :span="1">
								<n-form-item
									v-if="['TIMEOUT', 'TIMEOUT_RANDOM'].some(v => operation.type === v)"
									label="Ban time"
								>
									<n-input-number v-model:value="operation.timeoutTime" />
								</n-form-item>
							</n-grid-item>

							<n-grid-item
								v-if="operation.type.startsWith('OBS')
									&& (!obsSettings.data.value?.isConnected || !obsSettings.data.value?.serverPassword)
								"
								:span="4"
							>
								<n-alert title="You have to configure obs first" type="error">
									Seems like you not connected Twir with obs, please do it on overlays page.
								</n-alert>
							</n-grid-item>

							<n-grid-item v-if="operation.type === 'OBS_SET_SCENE'" :span="1">
								<n-form-item label="Obs scene">
									<n-select
										v-model:value="operation.target"
										:options="obsScenes"
										placeholder="Select obs scene"
										:disabled="!obsSettings.data.value?.isConnected"
									/>
								</n-form-item>
							</n-grid-item>

							<n-grid-item v-if="operation.type === 'OBS_TOGGLE_SOURCE'" :span="1">
								<n-form-item label="Obs source">
									<n-select
										v-model:value="operation.target"
										:options="obsSources"
										placeholder="Select obs source"
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
								].some(v => v === operation.type)"
								:span="1"
							>
								<n-form-item label="Obs audio source">
									<n-select
										v-model:value="operation.target"
										:options="obsAudioSources"
										placeholder="Select obs audio source"
										:disabled="!obsSettings.data.value?.isConnected"
									/>
								</n-form-item>
							</n-grid-item>
						</n-grid>
					</n-space>
				</n-timeline-item>
			</n-timeline>
		</n-space>
	</n-form>
</template>
