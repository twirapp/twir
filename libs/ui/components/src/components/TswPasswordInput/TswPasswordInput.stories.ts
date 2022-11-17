import { Meta, Story } from '@storybook/vue3';

import TswPasswordInput from '@/components/TswPasswordInput/TswPasswordInput.vue';

export default {
  title: 'TswPasswordInput',
  component: TswPasswordInput,
} as Meta;

interface TswPasswordInputProps {
  name: string;
  placeholder: string;
  disabled: true;
  id: string;
  value: string;
  isError: boolean;
}

const Template: Story<TswPasswordInputProps> = (args) => ({
  components: { TswPasswordInput },
  setup() {
    return { args };
  },
  template: '<TswPasswordInput v-bind="args" />',
});

export const PasswordInput = Template.bind({});
PasswordInput.args = {
  placeholder: 'Please write your password',
  value: '123456789',
  disabled: undefined,
  id: 'password-input',
  isError: false,
  name: 'password-input',
};
