<script setup lang='ts'>
import {
	type FormInst,
	type FormRules,
	type FormItemRule,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NSpace,
  NCheckbox,
  NDivider,
  NGrid,
  NGridItem,
  NSwitch,
  NCard,
  NText,
  NButton,
} from 'naive-ui';
import { ref, onMounted, toRaw } from 'vue';

import { useKeywordsManager } from '@/api/index.js';
import type { EditableKeyword } from '@/components/keywords/types.js';
import TextWithVariables from '@/components/textWithVariables.vue';

const props = defineProps<{
	keyword?: EditableKeyword | null
}>();
const emits = defineEmits<{
	close: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableKeyword>({
	text: '',
	response: '',
	cooldown: 0,
	enabled: true,
	isReply: true,
	usages: 0,
	isRegular: false,
});
onMounted(() => {
	if (!props.keyword) return;
	formValue.value = structuredClone(toRaw(props.keyword));
});

const keywordsManager = useKeywordsManager();
const keywordsUpdater = keywordsManager.update;
const keywordsCreator = keywordsManager.create;

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const data = formValue.value;

	if (data.id) {
		await keywordsUpdater.mutateAsync({
			id: data.id,
			keyword: data,
		});
	} else {
		await keywordsCreator.mutateAsync(data);
	}

	emits('close');
}

const rules: FormRules = {
	text: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error('Text is required');
			}

			return true;
		},
	},
	usages: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: number) => {
			if (value < 0) {
				return new Error('Usages are too short');
			}

			return true;
		},
	},
	response: {
		trigger: ['input', 'blue'],
		validator: (rule: FormItemRule, value: string) => {
			if (value?.length > 500) {
				return new Error('Response is too long');
			}
		},
	},
};
</script>

<template>
  <n-form
    ref="formRef"
    :model="formValue"
    :rules="rules"
  >
    <n-space vertical style="width: 100%">
      <n-form-item label="Text" path="text">
        <n-input v-model:value="formValue.text" />
      </n-form-item>
      <n-checkbox v-model:checked="formValue.isRegular">
        Is regular expression
      </n-checkbox>

      <n-form-item label="Response" path="response">
        <text-with-variables
          v-model="formValue.response"
          :min-rows="1"
          :max-rows="6"
          inputType="textarea"
        />
      </n-form-item>

      <n-divider>Settings</n-divider>

      <n-grid :cols="12" :x-gap="5" responsive="self">
        <n-grid-item :span="6">
          <n-form-item label="Cooldown" path="cooldown">
            <n-input-number v-model:value="formValue.cooldown" />
          </n-form-item>
        </n-grid-item>

        <n-grid-item :span="6">
          <n-form-item label="Usages counter" path="usages">
            <n-input-number v-model:value="formValue.usages" />
          </n-form-item>
        </n-grid-item>

        <n-grid-item :span="6">
          <n-card>
            <div class="settings-switch">
              <n-space vertical>
                <n-text>Reply</n-text>
                <n-text>Bot will send message as reply</n-text>
              </n-space>
              <n-switch v-model:value="formValue.isReply" />
            </div>
          </n-card>
        </n-grid-item>
      </n-grid>

      <n-button secondary type="success" block @click="save">
        Save
      </n-button>
    </n-space>
  </n-form>
</template>

<style scoped>
.settings-switch {
	display: flex;
	flex-direction: row;
	justify-content: space-between;
}
</style>
