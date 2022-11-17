<template>
  <div :class="classes">
    <input
      v-model="inputValue"
      :type="inputType"
      :placeholder="placeholder"
      class="tsw-input"
      :disabled="disabled"
    />
    <button type="button" class="tsw-password-input-button" @click="changeVisibilify">
      <TswIcon :name="iconName" size="20px" />
    </button>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';

import TswIcon from '@/components/TswIcon/TswIcon.vue';
import { PswInputIconName, PswInputTypes } from '@/components/TswPasswordInput/types';

const props = withDefaults(
  defineProps<{
    name?: string;
    placeholder?: string;
    disabled?: boolean;
    id?: string;
    value: string;
    isError?: boolean;
  }>(),
  {
    id: undefined,
    name: undefined,
    placeholder: undefined,
    disabled: false,
    isError: false,
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
  set(value: string) {
    emit('update:value', value);
  },
});

const inputType = ref<PswInputTypes>('password');
const iconName = ref<PswInputIconName>('Eye');

const changeVisibilify = () => {
  if (inputType.value === 'password') {
    inputType.value = 'text';
    iconName.value = 'EyeOff';
  } else {
    inputType.value = 'password';
    iconName.value = 'Eye';
  }
};

const classes = computed(() =>
  [
    'tsw-password-input',
    props.isError ? 'tsw-text-input-error' : '',
    props.disabled ? 'tsw-input-disabled' : '',
    props.disabled || props.isError ? '' : 'tsw-password-input-actions',
  ].join(' '),
);
</script>

<style lang="postcss">
.tsw-password-input {
  @apply inline-flex items-center bg-black-20 rounded overflow-hidden border-b w-full;

  transition: border-bottom-color 150ms theme('transitionTimingFunction.DEFAULT');

  & > input {
    @apply px-[14px] py-[10px] bg-black-20 leading-5 text-sm text-white-95 flex-1;
  }
}

.tsw-password-input-button {
  @apply px-[12px] py-[10px];

  &:hover svg {
    stroke: rgba(theme('colors.white.95'), 0.65);
  }

  & > svg {
    @apply stroke-gray-60;
    transition: stroke 150ms theme('transitionTimingFunction.DEFAULT');
  }
}

.tsw-password-input-actions {
  @apply border-b-black-20;

  &:hover {
    @apply border-b-black-25;
  }

  &:focus-within {
    @apply border-b-purple-70;
  }
}
</style>
