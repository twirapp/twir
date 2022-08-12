<template>
  <button :type="type" :class="classes" :as="as">
    <slot name="leftIcon" :classes="iconClasses" />
    {{ text }}
    <slot name="rightIcon" :classes="iconClasses" />
  </button>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

import { Size, Variant, ButtonType, ButtonTag } from '@/components/MyButton/props.types';

const props = withDefaults(
  defineProps<{
    text: string;
    size?: Size;
    variant?: Variant;
    type?: ButtonType;
    isRounded?: boolean;
    as?: ButtonTag;
  }>(),
  {
    as: 'button',
    isRounded: false,
    size: 'md',
    type: 'button',
    variant: 'solid-purple',
  },
);

type SizeClasses = { [K in Size]: string };
type VariantClasses = { [K in Variant]: string };

const base = 'inline-grid grid-flow-col items-center leading-tight';

const sizes: SizeClasses = {
  lg: 'text-lg px-5 py-3 rounded-md gap-x-2',
  md: 'px-3 py-2.5 rounded gap-x-1.5',
  sm: 'text-sm px-2 py-1.5 rounded gap-x-1',
};

const roundedSizes: SizeClasses = {
  lg: 'text-lg px-6 py-3 rounded-full gap-x-2',
  md: 'px-4 py-2.5 rounded-full gap-x-1.5',
  sm: 'text-sm px-3 py-1.5 rounded-full gap-x-1',
};

const variants: VariantClasses = {
  'outline-gray': 'bg-black-90 border border-gray-70 text-white-20 hover:bg-black-80',
  'solid-gray': '',
  'solid-purple': 'text-white-0 bg-purple-40 hover:bg-purple-45',
};

const iconSizes: SizeClasses = {
  lg: 'w-6 h-6',
  md: 'w-5 h-5',
  sm: 'w-4 h-4',
};

const classes = computed(() =>
  [
    base,
    props.isRounded ? roundedSizes[props.size] : sizes[props.size],
    variants[props.variant],
  ].join(' '),
);

const iconClasses = computed(() => iconSizes[props.size]);
</script>
