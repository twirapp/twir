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
} from 'naive-ui';
import { computed, onMounted, ref } from 'vue';

import { EVENTS } from './events.js';
import { eventTypeSelectOptions, operationTypeSelectOptions, getOperation } from './helpers.js';
import { EditableEvent } from './types.js';

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
};

const availableEventVariables = computed(() => {
	const evt = EVENTS[formValue.value.type];

	return evt?.variables?.map(v => ({
		label: `{${v}}`,
		value: v,
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

			<n-timeline>
				<n-timeline-item
					v-for="(operation, operationIndex) of formValue.operations"
					:key="operationIndex"
					:type="getOperation(operation.type)?.color ?? 'default'"
				>
					<n-space vertical>
						<n-grid cols="3 s:1 m:3" :x-gap="5" :y-gap="5">
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

						<n-form-item v-if="getOperation(operation.type)?.haveInput" label="Operation input">
							<n-input v-model:value="operation.input" />
						</n-form-item>
					</n-space>
				</n-timeline-item>
			</n-timeline>
		</n-space>
	</n-form>
</template>
