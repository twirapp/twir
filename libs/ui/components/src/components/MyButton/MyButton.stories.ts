import { Meta, Story } from '@storybook/vue3';
import { IconName } from '@tsuwari/ui-icons/icons';

import Icon from '../Icon/Icon.vue';
import MyButton from './MyButton.vue';
import { ButtonType, Size, Variant } from './props.types';

export default {
  title: 'MyButton',
  component: MyButton,
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
        value: ['lg', 'md', 'sm'] as Size[],
      },
      control: 'select',
    },
    variant: {
      type: {
        name: 'enum',
        value: ['outline-gray', 'solid-gray', 'solid-purple'] as Variant[],
      },
      control: 'select',
    },
    isRounded: {
      type: 'boolean',
    },
  },
} as Meta;

interface MyButtonProps {
  text: string;
  size: Size;
  variant: Variant;
  type: ButtonType;
  isRounded: boolean;
}

const Template: Story<MyButtonProps> = (args) => ({
  components: { MyButton },
  setup() {
    return { args };
  },
  template: '<MyButton v-bind="args" :type="args.type" />',
});

export const SolidPurple = Template.bind({});
SolidPurple.args = {
  text: 'Hello world',
  size: 'md',
  type: 'button',
  variant: 'solid-purple',
};

export const OutlineGray = Template.bind({});
OutlineGray.args = {
  text: 'Hello world',
  size: 'md',
  type: 'button',
  variant: 'outline-gray',
};

const TemplateWithRightIcon: Story<
  MyButtonProps & { rightIcon: IconName | null; leftIcon: IconName | null }
> = (args) => ({
  components: { MyButton, Icon },
  setup() {
    return { args };
  },
  template: `
    <MyButton v-bind="args">
      <template #leftIcon="{classes: leftIconClasses}">
        <Icon name="${args.leftIcon}" :class="leftIconClasses" />
      </template>
      <template #rightIcon="{classes: rightIconClasses}">
        <Icon name="${args.rightIcon}" :class="rightIconClasses" />
      </template>
    </MyButton>
  `,
});

export const ButtonWithIcon = TemplateWithRightIcon.bind({});
ButtonWithIcon.args = {
  text: 'Hello world',
  size: 'md',
  type: 'button',
  variant: 'outline-gray',
  leftIcon: 'Bell',
  rightIcon: 'ArrowMedium',
};
