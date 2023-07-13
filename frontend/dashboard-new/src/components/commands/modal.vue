<script setup lang='ts'>
import { IconShieldHalfFilled } from '@tabler/icons-vue';
import type { Command, Command_Response } from '@twir/grpc/generated/api/api/commands';
import {
	NForm,
	NFormItem,
	FormInst,
	FormRules,
	FormItemRule,
	NInput,
	NDynamicTags,
	NGrid,
	NGridItem,
	NDivider,
	NInputGroup,
	NInputGroupLabel,
	NDynamicInput,
	NSelect,
} from 'naive-ui';
import { ref, reactive, onMounted, computed } from 'vue';

import { useRolesManager } from '@/api/index.js';
import TextWithVariables from '@/components/textWithVariables.vue';
import TwitchUsersMultiple from '@/components/twitchUsers/multiple.vue';

const props = defineProps<{
	command: Command | null
}>();

type FormCommand = Omit<Command, 'responses'> & {
	responses: Array<Omit<Command_Response, 'id' | 'commandId' | 'order'> & { id?: string }>
};

const formRef = ref<FormInst | null>(null);
const formValue = reactive<FormCommand>({
	name:'',
	aliases: [],
	responses: [],
	description: '',
	rolesIds: [],
	deniedUsersIds: [],
	allowedUsersIds: [],
});

onMounted(() => {
	if (props.command) {
		formValue.name = props.command.name;
		formValue.aliases = props.command.aliases;
		formValue.responses = props.command.responses;
		formValue.description = props.command.description;
		formValue.rolesIds = props.command.rolesIds;
		formValue.deniedUsersIds = props.command.deniedUsersIds;
		formValue.allowedUsersIds = props.command.allowedUsersIds;
	}
});

const rolesManager = useRolesManager();
const roles = rolesManager.getAll({});
const rolesSelectOptions = computed(() => {
	if (!roles.data?.value) return [];
	return roles.data.value.roles.map((role) => ({
		label: role.name,
		value: role.id,
	}));
});

const nameValidator = (rule: FormItemRule, value: string) => {
	if (!value) {
		return new Error('Please input a name');
	}
	if (value.startsWith('!')) {
		return new Error('Name cannot start with !');
	}
	return true;
};
const rules: FormRules = {
	name: [{
		trigger: ['input', 'blur'],
		validator: nameValidator,
	}],
	aliases: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: string[]) => {
			value.forEach((alias) => nameValidator(rule, alias));

			return true;
		},
	},
};
</script>

<template>
  <n-form ref="formRef" :model="formValue" :rules="rules">
    <n-grid :cols="12" :x-gap="10">
      <n-grid-item :span="6">
        <n-form-item label="Name" path="name">
          <n-input-group>
            <n-input-group-label>!</n-input-group-label>
            <n-input v-model:value="formValue.name" placeholder="Input Name" />
          </n-input-group>
        </n-form-item>
      </n-grid-item>
      <n-grid-item :span="6">
        <n-form-item label="Aliases" path="aliases">
          <n-dynamic-tags v-model:value="formValue.aliases" />
        </n-form-item>
      </n-grid-item>
    </n-grid>
    <n-form-item label="Description" path="description">
      <n-input v-model:value="formValue.description" placeholder="Description" type="textarea" autosize />
    </n-form-item>

    <n-divider>
      Responses
    </n-divider>

    <n-dynamic-input
      v-model:value="formValue.responses"
      :on-create="() => ({ text: '' })"
      placeholder="Response"
      show-sort-button
    >
      <template #default="{ value }: { value: Command_Response }">
        <text-with-variables
          v-model="value.text"
          inputType="textarea"
          :minRows="3"
          :maxRows="6"
        >
        </text-with-variables>
      </template>
    </n-dynamic-input>

    <n-divider>
      Permissions
    </n-divider>

    <n-form-item label="Roles" path="rolesIds">
      <n-select
        v-model:value="formValue.rolesIds"
        multiple
        :options="rolesSelectOptions"
        :loading="roles.isLoading.value"
      >
        <template #arrow>
          <IconShieldHalfFilled />
        </template>
      </n-select>
    </n-form-item>

    <n-grid :cols="12">
      <n-grid-item :span="6">
        <twitch-users-multiple v-model="formValue.deniedUsersIds" />
      </n-grid-item>

      <n-grid-item :span="6">
      </n-grid-item>
    </n-grid>
    {{ JSON.stringify(formValue) }}
  </n-form>
</template>
