import { Meta, Story } from '@storybook/vue3';

import * as icons from './icons';

import TswIcon from '@/components/TswIcon/TswIcon.vue';

export default {
  title: 'TswIcon',
  component: TswIcon,
  argTypes: {
    name: {
      type: {
        name: 'enum',
        value: Object.keys(icons),
      },
      control: 'select',
    },
  },
} as Meta;

const Template: Story<{
  name: keyof typeof icons;
  width?: number;
  height?: number;
  title?: string;
  label?: string;
  strokeWidth?: number;
  stroke?: string;
  fill?: string;
  strokeStyle?: 'round' | 'butt' | 'square';
}> = (args) => ({
  components: { TswIcon },
  setup() {
    return { args };
  },
  template: '<TswIcon v-bind="args" />',
});

export const Icon = Template.bind({});
Icon.args = {
  name: 'ArrowInCircle',
  stroke: 'white',
};
