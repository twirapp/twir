<template>
  <button :type="type" :class="classes" class="py">
    {{ text }}
  </button>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

import { Size, Variant, Type } from './props.types.js';

const props = withDefaults(
  defineProps<{
    text: string;
    size?: Size;
    variant?: Variant;
    type?: Type;
  }>(),
  {
    size: 'md',
    type: 'button',
    variant: 'solid-purple',
  },
);

const base = 'inline-grid grid-flow-col leading-tight';

const sizeClasses: { [K in Size]: string } = {
  lg: 'text-lg px-5 py-3 rounded-md gap-x-2',
  md: 'px-3 py-2.5 rounded gap-x-2',
  sm: 'text-sm px-2 py-1.5 rounded gap-x-[1.5]',
};

const variantClasses: { [K in Variant]: string } = {
  'outline-gray': '',
  'solid-gray': '',
  'solid-purple': 'text-white-0 bg-purple-40 hover:bg-purple-45',
};

const classes = computed(
  () => `${base} ${sizeClasses[props.size]} ${variantClasses[props.variant]}`,
);
</script>
