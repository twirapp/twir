import { Meta, Story } from '@storybook/vue3';

import { TextInputType } from '@/components/TswTextInput/props.types.js';
import TswTextInput from '@/components/TswTextInput/TswTextInput.vue';

export default {
  title: 'TswTextInput',
  component: TswTextInput,
  argTypes: {
    type: {
      type: {
        name: 'enum',
        value: ['email', 'password', 'tel', 'text'] as TextInputType[],
      },
      control: 'select',
    },
  },
} as Meta;

interface TswTextInputProps {
  name: string;
  placeholder: string;
  type: TextInputType;
  disabled: boolean;
  id: string;
  value: string;
  isError: boolean;
}

const Template: Story<TswTextInputProps> = (args) => ({
  components: { TswTextInput },
  setup() {
    return { args };
  },
  template: '<TswTextInput v-bind="args" v-model:value="args.value" />',
});

export const TextInput = Template.bind({});
TextInput.args = {
  placeholder: 'Write some text',
  disabled: false,
  id: 'text-input',
  isError: false,
  name: 'text-input',
  type: 'text',
  value: 'Hello world',
};
