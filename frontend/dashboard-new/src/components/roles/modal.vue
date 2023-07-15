<script setup lang='ts'>
import { FormInst } from 'naive-ui';
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
});
</script>

<template>
  <n-form ref="formRef">
    <n-form-item>
    </n-form-item>
  </n-form>
</template>
