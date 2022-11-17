<template>
  <component
    :is="href ? 'a' : 'button'"
    :type="href ? undefined : type"
    :class="classes"
    :disabled="disabled ?? undefined"
    :href="href"
    :target="href && targetBlank ? '_black' : undefined"
    :tabindex="href ? 1 : undefined"
  >
    <TswIcon v-if="leftIcon" :name="leftIcon" :class="btnInnerClass" />
    {{ text }}
    <TswIcon v-if="rightIcon" :name="rightIcon" :class="btnInnerClass" />
  </component>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

import { ButtonSize, ButtonVariant, ButtonType } from './props.types.js';

import TswIcon from '@/components/TswIcon/TswIcon.vue';
import { IconName } from '@/components/TswIcon/types.js';

const props = withDefaults(
  defineProps<{
    text: string;
    size?: ButtonSize;
    variant?: ButtonVariant;
    type?: ButtonType;
    isRounded?: boolean;
    href?: string;
    disabled?: boolean;
    targetBlank?: true;
    leftIcon?: IconName;
    rightIcon?: IconName;
    // align?: 'center' | 'right' | 'left';
  }>(),
  {
    // align: 'center',
    leftIcon: undefined,
    rightIcon: undefined,
    size: 'md',
    type: 'button',
    variant: 'solid-purple',
    isRounded: false,
    href: undefined,
    targetBlank: undefined,
    disabled: false,
  },
);

const classes = computed(() =>
  [
    'btn',
    `btn-${props.size}`,
    props.isRounded ? 'btn-round' : '',
    `btn-${props.variant}`,
    props.disabled ? 'btn-disabled' : '',
  ].join(' '),
);

const btnInnerClass = 'btn-inner-el';
</script>

<style lang="postcss">
.btn {
  @apply inline-grid grid-flow-col items-center leading-tight cursor-pointer select-none;

  transition: box-shadow 150ms theme('transitionTimingFunction.DEFAULT');
}

.btn.btn-lg {
  @apply px-[20px] py-[12px] rounded-md gap-x-2.5;

  &.btn-round {
    @apply px-6 rounded-full;
  }

  & > .btn-inner-el {
    @apply w-6 h-6;
  }
}

.btn.btn-md {
  @apply px-[14px] py-[10px] rounded-md gap-x-2;

  &.btn-round {
    @apply px-6 rounded-full;
  }

  & > .btn-inner-el {
    @apply w-5 h-5;
  }
}

.btn.btn-sm {
  @apply text-sm px-[10px] py-[7px] rounded gap-x-1.5;

  &.btn-round {
    @apply px-4 rounded-full;
  }

  & > .btn-inner-el {
    @apply w-4 h-4;
  }
}

.btn.btn-solid-purple {
  @apply text-white-100 bg-purple-60;

  &:hover {
    @apply bg-purple-50;
  }

  &:focus {
    box-shadow: 0 0 0 3px rgba(theme('colors.purple.55'), 0.5);
  }

  & > .btn-inner-el {
    @apply stroke-white-100;
  }
}

.btn.btn-solid-gray {
  @apply bg-black-10 text-white-95 outline outline-1 outline-offset-0 outline-black-25;

  &:hover {
    @apply bg-black-20;
  }

  & > .btn-inner-el {
    @apply stroke-gray-70;
  }
}

.btn.btn-disabled {
  @apply opacity-50 pointer-events-none;
}
</style>
