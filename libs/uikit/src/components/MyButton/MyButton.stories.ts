import { Meta, Story } from '@storybook/vue3';

import { Size, Type } from './props.types.js';

import MyButton from '@/components/MyButton/MyButton.vue';

export default {
  title: 'MyButton',
  component: MyButton,
  template: '<MyButton v-bind="args" />',
  argTypes: {
    text: { type: 'string', required: true },
    type: {
      type: {
        name: 'enum',
        value: ['button', 'submit', 'reset'] as Type[],
      },
      defaultValue: 'button' as Type,
      control: 'select',
    },
    size: {
      type: {
        name: 'enum',
        value: ['lg', 'md', 'sm'] as Size[],
      },
      control: 'select',
      defaultValue: 'md' as Size,
    },
  },
} as Meta;

const Template: Story<{ text: string }> = (args) => ({
  components: { MyButton },
  setup() {
    return { args };
  },
  template: '<MyButton v-bind="args" :type="args.type" />',
});

export const Primary = Template.bind({});
Primary.args = {
  text: 'Hello',
};
