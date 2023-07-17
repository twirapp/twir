<script setup lang='ts'>
import {
	type FormInst,
	type FormRules,
	type FormItemRule,
} from 'naive-ui';
import { ref, onMounted, toRaw } from 'vue';

import { useGreetingsManager } from '@/api/index.js';
import type { EditableGreeting } from '@/components/greetings/types.js';
import TextWithVariables from '@/components/textWithVariables.vue';
import TwitchUserSearch from '@/components/twitchUsers/single.vue';

const props = defineProps<{
	greeting?: EditableGreeting | null
}>();
const emits = defineEmits<{
	close: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableGreeting>({
	enabled: true,
	isReply: true,
	text: '',
	userId: '',
	id: '',
});

onMounted(() => {
	if (!props.greeting) return;
	formValue.value = structuredClone(toRaw(props.greeting));
});

const greetingsManager = useGreetingsManager();
const greetingsUpdater = greetingsManager.update;
const greetingsCreator = greetingsManager.create;

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const data = formValue.value;

	if (data.id) {
		await greetingsUpdater.mutateAsync({
			id: data.id,
			greeting: data,
		});
	} else {
		await greetingsCreator.mutateAsync(data);
	}

	emits('close');
}

const rules: FormRules = {
	userId: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error('User is required');
			}

			return true;
		},
	},
	text: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error('Text is required');
			}

			return true;
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
      <n-form-item label="User" path="userId">
        <twitch-user-search
          v-model="formValue.userId"
          :initial-user-id="props.greeting?.userId"
        />
      </n-form-item>
      <n-form-item label="Text" path="text">
        <text-with-variables v-model="formValue.text" />
      </n-form-item>

      <n-form-item label="Use twitch reply" path="text">
        <n-switch v-model:value="formValue.isReply" />
      </n-form-item>
    </n-space>

    <n-button secondary type="success" block @click="save">
      Save
    </n-button>
  </n-form>
</template>
