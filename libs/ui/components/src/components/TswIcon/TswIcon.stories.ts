import { Meta, Story } from '@storybook/vue3';
import icons, { IconName } from '@tsuwari/ui-icons';

import TswIcon from '@/components/TswIcon/TswIcon.vue';

export default {
  title: 'TswIcon',
  component: TswIcon,
  argTypes: {
    name: {
      type: {
        name: 'enum',
        value: Object.keys(icons) as IconName[],
      },
      control: 'select',
    },
    fill: {
      type: 'string',
      control: 'color',
    },
    stroke: {
      type: 'string',
      control: 'color',
    },
  },
} as Meta;

const Template: Story<{
  name: IconName;
  size?: string;
  fill?: string;
  stroke?: string;
  strokeWidth?: number;
}> = (args) => ({
  components: { TswIcon },
  setup() {
    return { args };
  },
  template: '<TswIcon v-bind="args" />',
});

export const Icon = Template.bind({});
Icon.args = {
  name: 'Home',
  size: '36px',
  stroke: 'white',
  fill: 'none',
  strokeWidth: 1.5,
};
