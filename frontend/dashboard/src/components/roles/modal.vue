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
	NTabs,
	NTabPane,
	NButton,
} from 'naive-ui';
import { ref, onMounted, toRaw } from 'vue';

import { useRolesManager } from '@/api/index.js';
import { type EditableRole, permissions } from '@/components/roles/types.js';

const props = defineProps<{
	role?: EditableRole | null
}>();
const emits = defineEmits<{
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

const rolesManager =useRolesManager();
const rolesUpdater = rolesManager.update;
const rolesCreator = rolesManager.create;

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const data = formValue.value;

	if (data.id) {
		await rolesUpdater.mutateAsync({
			id: data.id,
			role: data,
		});
	} else {
		await rolesCreator.mutateAsync(data);
	}

	emits('close');
}
</script>

<template>
  <n-form ref="formRef">
    <n-tabs
      class="card-tabs"
      default-value="settings"
      size="large"
      animated
      pane-wrapper-style="margin: 0 -4px"
      pane-style="padding-left: 4px; padding-right: 4px; box-sizing: border-box;"
    >
      <n-tab-pane name="settings" tab="Settings">
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
            <n-grid-item
              v-for="(permission) of Object.entries(permissions)"
              :key="permission[0]"
              :span="6"
            >
              <n-checkbox :value="permission[0]" :label="permission[1]" :style="{ display: permission[1] == '' ? 'none' : undefined }" />
            </n-grid-item>
          </n-grid>
        </n-checkbox-group>
      </n-tab-pane>
      <n-tab-pane name="users" tab="Users">
        Users
      </n-tab-pane>
    </n-tabs>

    <n-divider />

    <n-button secondary type="success" block style="margin-top:15px" @click="save">
      Save
    </n-button>
  </n-form>
</template>

<style scoped>
.card-tabs .n-tabs-nav--bar-type {
	padding-left: 4px;
}
</style>
