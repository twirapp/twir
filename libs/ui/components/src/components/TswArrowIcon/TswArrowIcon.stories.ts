import { Meta, Story } from '@storybook/vue3';

import { ArrowDirection, ArrowSize } from '@/components/TswArrowIcon/props.types';
import TswArrowIcon from '@/components/TswArrowIcon/TswArrowIcon.vue';

export default {
  title: 'TswArrowIcon',
  component: TswArrowIcon,
  argTypes: {
    arrowSize: {
      type: {
        name: 'enum',
        value: ['lg', 'md', 'in-circle'] as ArrowSize[],
      },
      control: 'select',
    },
    direction: {
      type: {
        name: 'enum',
        value: ['bottom', 'left', 'right', 'top'] as ArrowDirection[],
      },
      control: 'select',
    },
    stroke: {
      control: {
        type: 'color',
      },
    },
  },
} as Meta;

const Template: Story<{
  direction: ArrowDirection;
  arrowSize?: ArrowSize;
  stroke?: string;
  strokeWidth?: number;
  size?: string;
}> = (args) => ({
  components: { TswArrowIcon },
  setup() {
    return { args };
  },
  template: '<TswArrowIcon v-bind="args" />',
});

export const RightArrow = Template.bind({});
RightArrow.args = {
  size: '36px',
  arrowSize: 'md',
  direction: 'right',
  stroke: 'white',
  strokeWidth: 1.5,
};
