<script setup lang='ts'>
import {
	type FormInst,
	NInput,
	NForm,
	NFormItem,
	NDivider,
	NGrid,
	NGridItem,
	NInputNumber,
	NCheckbox,
	NCheckboxGroup,
} from 'naive-ui';
import { ref, onMounted, toRaw } from 'vue';

import { type EditableRole, permissions } from '@/components/roles/types.js';

const props = defineProps<{
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

onMounted(() => {
	if (!props.role) return;
	formValue.value = structuredClone(toRaw(props.role));
});
</script>

<template>
  <n-form ref="formRef">
    <n-form-item label="Name">
      <n-input v-model:value="formValue.name" />
    </n-form-item>

    <n-divider>Access by stats</n-divider>

    <n-grid :cols="12" :x-gap="5">
      <n-grid-item :span="6">
        <n-form-item label="Required watch time">
          <n-input-number
            v-model:value="formValue.settings!.requiredWatchTime"
            :min="0" :max="99999999"
          />
        </n-form-item>
      </n-grid-item>

      <n-grid-item :span="6">
        <n-form-item label="Required messages">
          <n-input-number
            v-model:value="formValue.settings!.requiredMessages"
            :min="0"
            :max="99999999"
          />
        </n-form-item>
      </n-grid-item>

      <n-grid-item :span="6">
        <n-form-item label="Required used channels points">
          <n-input-number
            v-model:value="formValue.settings!.requiredUserChannelPoints"
            :min="0"
            :max="999999999999"
          />
        </n-form-item>
      </n-grid-item>
    </n-grid>

    <n-divider>Permissions</n-divider>

    <n-checkbox-group v-model:value="formValue.permissions">
      <n-grid :cols="12" :x-gap="5">
        <n-grid-item v-for="(permission) of Object.entries(permissions)" :key="permission[0]" :span="6">
          <n-checkbox :value="permission[0]" :label="permission[1]" />
        </n-grid-item>
      </n-grid>
    </n-checkbox-group>
  </n-form>
</template>
