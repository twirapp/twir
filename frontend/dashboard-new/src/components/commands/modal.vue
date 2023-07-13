<script setup lang='ts'>
import {
	IconTrash,
	IconPlus,
	IconArrowBigUp,
	IconArrowBigDown,
} from '@tabler/icons-vue';
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
	NButton,
	NSpace,
	NInputGroup,
	NDynamicInput,
} from 'naive-ui';
import { ref } from 'vue';

defineProps<{
	command: Command | null
}>();

type FormCommand = Omit<Command, 'responses'> & {
	responses: Omit<Command_Response, 'id' | 'commandId' | 'order'>[]
};

const formRef = ref<FormInst | null>(null);
const formValue = ref<FormCommand>({
	name:'',
	aliases: [],
	responses: [],
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
          <n-input v-model:value="formValue.name" placeholder="Input Name" />
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

    <!--    <n-button size="small" block @click="formValue.responses.push({ text: '' })">-->
    <!--      +-->
    <!--    </n-button>-->
    <!--    <div v-for="(response, responseIndex) in formValue.responses" :key="responseIndex">-->
    <!--      <n-form-item :path="'aliases'[responseIndex]">-->
    <!--        <n-input-group>-->
    <!--          <n-input v-model:value="response.text" placeholder="Response" />-->
    <!--          <n-button type="primary">-->
    <!--            Search-->
    <!--          </n-button>-->
    <!--          <n-button secondary type="error">-->
    <!--            <IconTrash />-->
    <!--          </n-button>-->
    <!--        </n-input-group>-->
    <!--      </n-form-item>-->
    <!--    </div>-->
    <n-dynamic-input
      v-model:value="formValue.responses"
      :on-create="() => ({ text: '' })"
      placeholder="Come on"
      show-sort-button
    >
      <template #default="{ value }">
        <n-input v-model:value="value.text" type="text" />
      </template>
      <!--      <template #action="{ index, create, remove, move }">-->
      <!--        <n-space style="margin-left: 20px">-->
      <!--          <n-button size="small" @click="() => create(index)">-->
      <!--            <IconPlus />-->
      <!--          </n-button>-->
      <!--          <n-button size="small" @click="() => remove(index)">-->
      <!--            <IconTrash />-->
      <!--          </n-button>-->
      <!--          <n-button size="small" @click="() => move('up', index)">-->
      <!--            <IconArrowBigUp />-->
      <!--          </n-button>-->
      <!--          <n-button size="small" @click="() => move('down', index)">-->
      <!--            <IconArrowBigDown />-->
      <!--          </n-button>-->
      <!--        </n-space>-->
      <!--      </template>-->
    </n-dynamic-input>
    {{ JSON.stringify(formValue) }}
  </n-form>
</template>
