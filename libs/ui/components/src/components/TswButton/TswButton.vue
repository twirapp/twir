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
    <slot name="left" innerClass="btn-inner-el" />
    {{ text }}
    <slot name="right" innerClass="btn-inner-el" />
  </component>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

import { ButtonSize, ButtonVariant, ButtonType } from '@/components/TswButton/props.types';

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
  }>(),
  {
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
</script>

<style lang="postcss">
.btn {
  @apply inline-grid grid-flow-col items-center leading-tight cursor-pointer select-none;

  transition: box-shadow 150ms theme('transitionTimingFunction.in-out');

  &:focus {
    outline: none;
  }
}

.btn-lg {
  @apply text-lg px-5 py-3 rounded-md gap-x-2.5;

  &.btn-round {
    @apply px-6 rounded-full;
  }

  & > .btn-inner-el {
    @apply w-6 h-6;
  }
}

.btn-md {
  @apply px-3 py-2.5 rounded gap-x-2;

  &.btn-round {
    @apply px-4 rounded-full;
  }

  & > .btn-inner-el {
    @apply w-5 h-5;
  }
}

.btn-sm {
  @apply text-sm px-2 py-1.5 rounded gap-x-1.5;

  &.btn-round {
    @apply px-3 rounded-full;
  }

  & > .btn-inner-el {
    @apply w-4 h-4;
  }
}

.btn-outline-gray {
  @apply bg-black-90 border border-gray-70 text-white-20;

  &:hover {
    @apply bg-black-80;
  }

  &:focus {
    box-shadow: 0 0 0 3px theme('colors.black.80');
  }
}

.btn-solid-purple {
  @apply text-white-0 bg-purple-40;

  &:hover {
    @apply bg-purple-45;
  }

  &:focus {
    box-shadow: 0 0 0 3px rgba(theme('colors.purple.45'), 0.6);
  }
}

.btn-solid-gray {
  @apply bg-black-80 text-white-10;

  &:hover {
    @apply bg-black-75;
  }

  &:focus {
    box-shadow: 0 0 0 3px rgba(theme('colors.black.75'), 0.55);
  }
}

.btn-disabled {
  @apply opacity-50 pointer-events-none;
}
</style>
