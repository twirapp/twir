import type { Meta, Story } from '@storybook/vue3';

import { IconName } from '../TswIcon/types.js';

import type { ButtonSize, ButtonType, ButtonVariant } from '@/components/TswButton/props.types';
import TswButton from '@/components/TswButton/TswButton.vue';
import * as icons from '@/components/TswIcon/icons';

export default {
  title: 'TswButton',
  component: TswButton,
  argTypes: {
    text: { type: 'string', required: true },
    type: {
      type: {
        name: 'enum',
        value: ['button', 'submit', 'reset'] as ButtonType[],
      },
      control: 'select',
    },
    size: {
      type: {
        name: 'enum',
        value: ['lg', 'md', 'sm'] as ButtonSize[],
      },
      control: 'select',
    },
    variant: {
      type: {
        name: 'enum',
        value: ['solid-gray', 'solid-purple'] as ButtonVariant[],
      },
      control: 'select',
    },
    onClick: {
      action: 'click',
      table: {
        disable: true,
      },
    },
    leftIcon: {
      type: {
        name: 'enum',
        value: Object.keys(icons) as IconName[],
      },
      control: 'select',
    },
    rightIcon: {
      type: {
        name: 'enum',
        value: Object.keys(icons) as IconName[],
      },
      control: 'select',
    },
  },
} as Meta;

interface TswButtonProps {
  text: string;
  size?: ButtonSize;
  variant?: ButtonVariant;
  type?: ButtonType;
  isRounded?: boolean;
  href?: string;
  disabled?: boolean;
  targetBlank?: true;
}

const Template: Story<TswButtonProps> = (args) => ({
  components: { TswButton },
  setup() {
    return { args };
  },
  template: '<TswButton v-bind="args" />',
});

export const SolidPurple = Template.bind({});
SolidPurple.args = {
  text: 'Hello world',
  size: 'md',
  type: 'button',
  variant: 'solid-purple',
  disabled: false,
  href: undefined,
  isRounded: false,
  targetBlank: undefined,
};
