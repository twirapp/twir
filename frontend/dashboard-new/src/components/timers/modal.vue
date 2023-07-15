<script setup lang='ts'>
import {
	type FormInst,
	type FormRules,
	NForm,
	NFormItem,
	NInput,
	NInputNumber,
	NButton,
	NSlider,
	NGrid,
	NGridItem, FormItemRule,
} from 'naive-ui';
import { ref, onMounted } from 'vue';

import { EditableTimer } from '@/components/timers/types.js';

const props = defineProps<{
	timer?: EditableTimer | null
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableTimer>({
	name: '',
	enabled: true,
	messageInterval: 0,
	timeInterval: 5,
	responses: [],
});
const rules: FormRules = {
	name: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error('Name is required');
			}
			if (value.length > 50) {
				return new Error('Name is too long');
			}
			return true;
		},
	},
	timeInterval: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: number) => {
			if (value < 1) {
				return new Error('Time interval is too short');
			}
			if (value > 100) {
				return new Error('Time interval is too long');
			}
			return true;
		},
	},
	messageInterval: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: number) => {
			if (value < 0) {
				return new Error('Message interval is too short');
			}
			if (value > 5000) {
				return new Error('Message interval is too long');
			}
			return true;
		},
	},
};

onMounted(() => {
	if (!props.timer) return;
	formValue.value = props.timer;
});

async function save() {
return;
}

const sliderMarks = {
	5: '5',
	15: '15',
	30: '30',
	45: '45',
	60: '60',
	75: '75',
	90: '90',
	100: '100',
};
</script>

<template>
  <n-form
    :ref="formRef"
    :model="formValue"
    :rules="rules"
  >
    <n-space vertical style="width: 100%">
      <n-form-item label="Name" path="name">
        <n-input v-model:value="formValue.name" />
      </n-form-item>
      <n-form-item label="Time interval in minutes" path="timeInterval">
        <n-grid :cols="12" :x-gap="5">
          <n-grid-item :span="10">
            <n-slider
              v-model:value="formValue.timeInterval"
              :step="1"
              :marks="sliderMarks"
              :min="1"
            />
          </n-grid-item>
          <n-grid-item :span="2">
            <n-input-number v-model:value="formValue.timeInterval" size="small" :min="1" :max="100" />
          </n-grid-item>
        </n-grid>
      </n-form-item>
      <n-form-item label="Interval in messages" path="messageInterval">
        <n-input-number v-model:value="formValue.messageInterval" :min="0" :max="5000" />
      </n-form-item>
      <n-button secondary type="success" block @click="save">
        Save
      </n-button>
    </n-space>
  </n-form>
</template>

