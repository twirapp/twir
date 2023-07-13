<script setup lang='ts'>
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
	NSpace,
} from 'naive-ui';
import { ref, reactive, onMounted } from 'vue';

import TextWithVariables from '@/components/textWithVariables.vue';

const props = defineProps<{
	command: Command | null
}>();

type FormCommand = Omit<Command, 'responses'> & {
	responses: Omit<Command_Response, 'id' | 'commandId' | 'order'>[]
};

const formRef = ref<FormInst | null>(null);
const formValue = reactive<FormCommand>({
	name:'',
	aliases: [],
	responses: [],
});

onMounted(() => {
	if (props.command) {
		formValue.name = props.command.name;
		formValue.aliases = props.command.aliases;
		formValue.responses = props.command.responses;
	}
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
    {{ JSON.stringify(formValue) }}
  </n-form>
</template>
