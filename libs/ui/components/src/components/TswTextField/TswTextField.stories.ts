import { Meta, Story } from '@storybook/vue3';

import TswTextField from '@/components/TswTextField/TswTextField.vue';
import { InputVariantType } from '@/components/TswTextField/types';
import { TextInputType } from '@/components/TswTextInput/props.types';

export default {
  title: 'TswTextField',
  component: TswTextField,
} as Meta;

interface TswTextFieldProps {
  inputVariant?: InputVariantType;
  value?: string;
  initialErrors?: string[];
  label?: string;
  name: string;
  disabled?: boolean;
  id?: string;
  placeholder?: string;
  type?: TextInputType;
  infoMessage?: string;
}

const Template: Story<TswTextFieldProps> = (args) => ({
  components: { TswTextField },
  setup() {
    return { args };
  },
  template: '<TswTextField v-bind="args" v-model:modelValue="args.modelValue" />',
});

export const TextField = Template.bind({});
TextField.args = {
  placeholder: 'Write some text',
  disabled: false,
  id: 'text-field',
  infoMessage: 'Must be at least 5 characters',
  initialErrors: undefined,
  inputVariant: 'text',
  label: 'Text field label',
  name: 'text-field',
  type: 'text',
  value: '',
};

export const ErrorTextField = Template.bind({});
ErrorTextField.args = {
  placeholder: 'Write some text',
  disabled: false,
  id: 'text-field',
  infoMessage: 'Must be at least 5 characters',
  initialErrors: ['Something went wrong'],
  inputVariant: 'text',
  label: 'Text field label',
  name: 'text-field',
  type: 'text',
  value: 'Random wrong value',
};
