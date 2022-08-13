import type { Meta, Story } from '@storybook/vue3';
import type { IconName } from '@tsuwari/ui-icons/icons';

import type { ButtonSize, ButtonType, ButtonVariant } from '@/components/TswButton/props.types';
import TswButton from '@/components/TswButton/Tswbutton.vue';
import TswIcon from '@/components/TswIcon/TswIcon.vue';
import TswLoader from '@/components/TswLoader/TswLoader.vue';

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
        value: ['outline-gray', 'solid-gray', 'solid-purple'] as ButtonVariant[],
      },
      control: 'select',
    },
    isRounded: {
      type: 'boolean',
    },
    disabled: {
      type: 'boolean',
    },
    href: {
      type: 'string',
    },
    targetBlank: {
      type: 'boolean',
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

export const OutlineGray = Template.bind({});
OutlineGray.args = {
  text: 'Hello world',
  size: 'md',
  type: 'button',
  variant: 'outline-gray',
  disabled: false,
  href: undefined,
  isRounded: false,
  targetBlank: undefined,
};

const TemplateWithIcons: Story<
  TswButtonProps & { rightIcon: IconName | null; leftIcon: IconName | null }
> = (args) => ({
  components: { TswButton, TswIcon },
  setup() {
    return { args };
  },
  template: `
    <TswButton v-bind="args">
      <template #left="{innerClass}" v-if="args.leftIcon !== null">
        <TswIcon name="${args.leftIcon}" :class="innerClass" />
      </template>
      <template #right="{innerClass}"  v-if="args.rightIcon !== null">
        <TswIcon name="${args.rightIcon}" :class="innerClass" />
      </template>
    </TswButton>
  `,
});

export const OutlineRoundWithIcons = TemplateWithIcons.bind({});
OutlineRoundWithIcons.args = {
  text: 'Hello world',
  size: 'md',
  type: 'button',
  variant: 'outline-gray',
  leftIcon: 'CommandLine',
  rightIcon: 'Bell',
  disabled: false,
  href: undefined,
  isRounded: true,
  targetBlank: undefined,
};

const TemplateWithLoader: Story<TswButtonProps> = (args) => ({
  components: { TswButton, TswLoader },
  setup() {
    return { args };
  },
  template: `
    <TswButton v-bind="args">
      <template #right>
        <TswLoader size="md" />
      </template>
    </TswButton>
  `,
});

export const OutlineWithLoader = TemplateWithLoader.bind({});
OutlineWithLoader.args = {
  text: 'Hello world',
  size: 'md',
  type: 'button',
  variant: 'outline-gray',
  disabled: false,
  href: undefined,
  isRounded: true,
  targetBlank: undefined,
};
