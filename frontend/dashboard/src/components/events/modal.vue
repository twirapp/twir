<script setup lang='ts'>
import {
	type SelectGroupOption,
	type SelectOption,
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
	NAutoComplete,
} from 'naive-ui';
import { computed, onMounted, ref } from 'vue';

import { EVENTS } from './events.js';
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

const typeSelectOptions: (SelectOption | SelectGroupOption)[] = Object.entries(EVENTS).map(([key, value]) => {
	return {
		value: key,
		label: value.name,
		type: value.type ? 'group' : undefined,
		children: value.childrens
		? Object.entries(value.childrens).map(([childKey, childValue]) => ({
			value: childKey,
			label: childValue.name,
		}))
		: [],
	};
});

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
						<n-select v-model:value="formValue.type" :options="typeSelectOptions" />
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
				<n-timeline-item>
					<n-auto-complete
						:options="availableEventVariables.map(v => ({
							...v,
							value: `${formValue.operations[0].input} {${v.value}}`
						}))"
					/>
				</n-timeline-item>
				<n-timeline-item>
					qwe
				</n-timeline-item>
			</n-timeline>
		</n-space>
	</n-form>
</template>
