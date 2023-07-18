<script setup lang='ts'>
import {
	IconTrash,
	IconPlus,
	IconArrowNarrowUp,
	IconArrowNarrowDown,
} from '@tabler/icons-vue';
import {
	type FormInst,
	type FormRules,
	type FormItemRule,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NButton,
  NSlider,
  NGrid,
  NGridItem,
  NDynamicInput,
  NDivider,
  NSpace,
  NTimeline,
  NTimelineItem,
  NCheckbox,
} from 'naive-ui';
import { ref, onMounted, toRaw } from 'vue';

import { useTimersManager } from '@/api/index.js';
import type { EditableTimer, EditableTimerResponse } from '@/components/timers/types.js';

const props = defineProps<{
	timer?: EditableTimer | null
}>();
const emits = defineEmits<{
	close: []
}>();

const nowTime = new Date();

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
	responses: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error('Text is required');
			}
			if (value.length > 500) {
				return new Error('Text is too long');
			}
			return true;
		},
	},
};

onMounted(() => {
	if (!props.timer) return;
	formValue.value = structuredClone(toRaw(props.timer));
});

const timersManager = useTimersManager();
const timersUpdater = timersManager.update;
const timersCreator = timersManager.create;

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const data = formValue.value;

	if (data.id) {
		await timersUpdater.mutateAsync({
			id: data.id,
			timer: data,
		});
	} else {
		await timersCreator.mutateAsync({ data });
	}

	emits('close');
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
    ref="formRef"
    :model="formValue"
    :rules="rules"
  >
    <n-space vertical style="width: 100%">
      <n-form-item label="Name" path="name" show-require-mark>
        <n-input v-model:value="formValue.name" />
      </n-form-item>
      <n-form-item label="Time interval in minutes" path="timeInterval" show-require-mark>
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
      <!--      <n-form-item label="Interval in messages" path="messageInterval">-->
      <!--        <n-input-number v-model:value="formValue.messageInterval" :min="0" :max="5000" />-->
      <!--      </n-form-item>-->

      <n-divider>
        Responses
      </n-divider>

      <n-dynamic-input
        v-model:value="formValue.responses"
        placeholder="Timer response"
        class="groups"
        :create-button-props="{ class: 'create-button' } as any"
      >
        <template #default="{ value, index }: { value: EditableTimerResponse }">
          <n-space vertical style="width: 100%">
            <n-form-item :path="`responses[${index}].text`" :rule="rules.responses" show-require-mark>
              <n-input
                v-model:value="value.text"
                type="textarea"
                placeholder="Response text"
                :autosize="{
                  minRows: 1,
                  maxRows: 5
                }"
              />
            </n-form-item>
            <n-checkbox v-model:checked="value.isAnnounce">
              Use twitch announce
            </n-checkbox>
          </n-space>
        </template>

        <template #action="{ index, remove, move }">
          <div class="group-actions">
            <n-button size="small" type="error" quaternary @click="() => remove(index)">
              <IconTrash />
            </n-button>
            <n-button
              size="small"
              type="info"
              quaternary
              :disabled="index == 0"
              @click="() => move('up', index)"
            >
              <IconArrowNarrowUp />
            </n-button>
            <n-button
              size="small"
              type="info"
              quaternary
              :disabled="!!formValue.responses.length && index === formValue.responses.length-1"
              @click="() => move('down', index)"
            >
              <IconArrowNarrowDown />
            </n-button>
          </div>
        </template>
      </n-dynamic-input>

      <n-button
        dashed
        block
        @click="() => formValue.responses.push({ text: '', isAnnounce: false })"
      >
        <IconPlus /> Create
      </n-button>

      <n-divider />

      One response will be sent once every N intervals, example:
      <n-timeline>
        <n-timeline-item
          v-for="(response, index) in formValue.responses"
          :key="index"
          type="info"
          :title="`Response #${index+1}`"
          :content="response.text.slice(0, 50) + '...'"
          :time="new Date(nowTime.getTime() + (formValue.timeInterval * index * 60000)).toLocaleTimeString()"
          line-type="dashed"
        />
      </n-timeline>

      <n-divider />

      <n-button secondary type="success" block style="margin-top: 10px" @click="save">
        Save
      </n-button>
    </n-space>
  </n-form>
</template>

<style scoped>
.groups :deep(.create-button) {
	display: none;
}

.group-actions {
	display: flex;
	column-gap: 5px;
	align-items: center;
	margin-left: 10px;
}
</style>
