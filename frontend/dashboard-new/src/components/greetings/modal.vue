<script setup lang='ts'>
import { FormInst } from 'naive-ui';
import { ref, onMounted, toRaw } from 'vue';

import type { EditableGreeting } from '@/components/greetings/types.js';

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
</script>

<template>
  <div>
  </div>
</template>
