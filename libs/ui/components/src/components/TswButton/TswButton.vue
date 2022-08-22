<template>
  <component
    :is="href ? 'a' : 'button'"
    :type="href ? undefined : type"
    :class="classes"
    :disabled="disabled ?? undefined"
    :href="href"
    :target="href && targetBlank ? '_black' : undefined"
    :tabindex="href ? 1 : undefined"
    @click.prevent="emitClickEvent"
  >
    <slot name="left" :innerClass="btnInnerClass" />
    {{ text }}
    <slot name="right" :innerClass="btnInnerClass" />
  </component>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

import { ButtonSize, ButtonVariant, ButtonType } from './props.types.js';

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

const emit = defineEmits({
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  click: (e: MouseEvent) => true,
});

const classes = computed(() =>
  [
    'btn',
    `btn-${props.size}`,
    props.isRounded ? 'btn-round' : '',
    `btn-${props.variant}`,
    props.disabled ? 'btn-disabled' : '',
  ].join(' '),
);

const emitClickEvent = (e: MouseEvent) => {
  emit('click', e);
};

const btnInnerClass = 'btn-inner-el';
</script>

<style lang="postcss">
button.btn {
  @apply inline-grid grid-flow-col items-center leading-tight cursor-pointer select-none;

  transition: box-shadow 150ms theme('transitionTimingFunction.in-out');

  &:focus {
    outline: none;
  }
}

button.btn-lg {
  @apply text-lg px-5 py-3 rounded-md gap-x-2.5;

  &.btn-round {
    @apply px-6 rounded-full;
  }

  & > .btn-inner-el {
    @apply w-6 h-6;
  }
}

button.btn-md {
  @apply px-3 py-2.5 rounded gap-x-2;

  &.btn-round {
    @apply px-4 rounded-full;
  }

  & > .btn-inner-el {
    @apply w-5 h-5;
  }
}

button.btn-sm {
  @apply text-sm px-2 py-1.5 rounded gap-x-1.5;

  &.btn-round {
    @apply px-3 rounded-full;
  }

  & > .btn-inner-el {
    @apply w-4 h-4;
  }
}

button.btn-outline-gray {
  @apply border border-gray-35 text-white-95;

  &:hover {
    @apply bg-black-15;
  }

  &:focus {
    box-shadow: 0 0 0 3px theme('colors.black.15');
  }
}

button.btn-solid-purple {
  @apply text-white-100 bg-purple-60;

  &:hover {
    @apply bg-purple-55;
  }

  &:focus {
    box-shadow: 0 0 0 3px rgba(theme('colors.purple.55'), 0.5);
  }
}

button.btn-solid-gray {
  @apply bg-black-15 text-white-95;

  &:hover {
    @apply bg-black-20;
  }

  &:focus {
    box-shadow: 0 0 0 3px rgba(theme('colors.black.20'), 0.55);
  }
}

button.btn-disabled {
  @apply opacity-50 pointer-events-none;
}
</style>
