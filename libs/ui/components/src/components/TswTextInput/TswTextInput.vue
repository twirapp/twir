<template>
  <input
    :id="id"
    v-model="inputValue"
    :name="name"
    :type="type"
    :placeholder="placeholder"
    :class="classes"
    :disabled="disabled"
  />
</template>

<script lang="ts" setup>
import { computed } from 'vue';

import { TextInputType } from '@/components/TswTextInput/props.types';

const props = withDefaults(
  defineProps<{
    name?: string;
    placeholder?: string;
    type?: TextInputType;
    disabled?: boolean;
    id?: string;
    value: string;
    isError?: boolean;
  }>(),
  {
    isError: false,
    id: undefined,
    placeholder: undefined,
    disabled: false,
    type: 'text',
    name: undefined,
  },
);

const emit = defineEmits({
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  'update:value': (value: string) => true,
});

const inputValue = computed({
  get() {
    return props.value;
  },
  set(value) {
    emit('update:value', value);
  },
});

const classes = computed(() =>
  [
    'tsw-input',
    'tws-text-input',
    props.disabled ? 'tsw-input-disabled' : '',
    props.isError ? 'tsw-text-input-error' : '',
    props.disabled || props.isError ? '' : 'tsw-text-input-actions',
  ].join(' '),
);
</script>

<style lang="postcss">
.tws-text-input {
  @apply bg-black-20 rounded text-white-95 leading-5 px-[14px] py-[10px] border-b text-sm w-full;

  transition: border-bottom-color 150ms theme('transitionTimingFunction.DEFAULT');
}

.tsw-text-input-actions {
  border-bottom-color: rgba(theme('colors.gray.35'), 0.5);

  &:hover {
    border-bottom-color: theme('colors.gray.35');
  }

  &:focus {
    @apply border-b-purple-70;
  }
}
</style>
