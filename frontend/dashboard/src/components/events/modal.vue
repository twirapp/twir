<script setup lang='ts'>
import {
	NSpace,
	NSelect,
	SelectOption,
	NForm,
	NFormItem,
	FormInst,
	FormItemRule,
	FormRules,
	NInput,
	NText,
	NTimeline,
	NTimelineItem,
	NMention,
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
	operations: [],
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

const typeSelectOptions: SelectOption[] = Object.entries(EVENTS).map(([key, value]) => {
	return {
		value: key,
		label: value.name,
	};
});

const availableEventVariables = computed(() => {
	const evt = EVENTS[formValue.value.type];

	return evt?.variables?.map(v => ({
		label: `{${v}}`,
		value: `${v}}`,
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
					<n-mention :options="availableEventVariables" :prefix="['{', '}']" />
				</n-timeline-item>
				<n-timeline-item>
					qwe
				</n-timeline-item>
			</n-timeline>
		</n-space>
	</n-form>
</template>
