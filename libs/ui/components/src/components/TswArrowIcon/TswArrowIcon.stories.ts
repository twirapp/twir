import { Meta, Story } from '@storybook/vue3';

import { ArrowDirection, ArrowSize } from '@/components/TswArrowIcon/props.types';
import TswArrowIcon from '@/components/TswArrowIcon/TswArrowIcon.vue';

export default {
  title: 'TswArrowIcon',
  component: TswArrowIcon,
  argTypes: {
    arrowType: {
      type: {
        name: 'enum',
        value: ['in-circle', 'narrow', 'triangle-lg', 'triangle-md'] as ArrowSize[],
      },
      control: 'select',
    },
    direction: {
      type: {
        name: 'enum',
        value: [
          'bottom',
          'left',
          'right',
          'top',
          'bottom-left',
          'bottom-right',
          'top-left',
          'top-right',
        ] as ArrowDirection[],
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
  arrowType?: ArrowSize;
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
  arrowType: 'triangle-md',
  direction: 'right',
  stroke: 'white',
  strokeWidth: 1.5,
};
