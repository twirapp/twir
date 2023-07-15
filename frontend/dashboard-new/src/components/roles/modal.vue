<script setup lang='ts'>
import {
	type FormInst,
	NInput,
	NForm,
	NFormItem,
} from 'naive-ui';
import { ref } from 'vue';

import type { EditableRole } from '@/components/roles/types.js';

defineProps<{
	role?: EditableRole | null
}>();
defineEmits<{
	close: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableRole>({
	name: '',
	type: '',
	permissions: [],
	users: [],
	settings: {
		requiredMessages: 0,
		requiredUserChannelPoints: 0,
		requiredWatchTime: 0,
	},
});
</script>

<template>
  <n-form ref="formRef">
    <n-form-item label="Name">
      <n-input v-model:value="formValue.name" />
    </n-form-item>
  </n-form>
</template>
