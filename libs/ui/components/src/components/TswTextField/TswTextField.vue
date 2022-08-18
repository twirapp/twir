<template>
  <div :class="classes">
    <label v-if="label" :for="id">{{ label }}</label>
    <component
      :is="inputComponents[inputVariant]"
      :id="id"
      v-model:value="inputValue"
      :isError="isError"
      :disabled="disabled"
      :name="name"
      :placeholder="placeholder"
      :type="inputVariant === 'text' ? type : undefined"
    />
    <span v-if="isError" class="error-message">
      {{ errorMessage }}
    </span>
    <span v-if="infoMessage !== undefined && !isError" class="info-message">
      {{ infoMessage }}
    </span>
  </div>
</template>

<script lang="ts" setup>
import { useField } from 'vee-validate';
import { computed, toRef } from 'vue';

import TswPasswordInput from '@/components/TswPasswordInput/TswPasswordInput.vue';
import { InputVariantType } from '@/components/TswTextField/types';
import { TextInputType } from '@/components/TswTextInput/props.types';
import TswTextInput from '@/components/TswTextInput/TswTextInput.vue';

type InputComponents = { [K in InputVariantType]: typeof TswPasswordInput | typeof TswTextInput };

const inputComponents: InputComponents = {
  text: TswTextInput,
  password: TswPasswordInput,
};

const props = withDefaults(
  defineProps<{
    inputVariant?: InputVariantType;
    value?: string;
    initialErrors?: string[];
    label?: string;
    name: string;
    disabled?: boolean;
    id?: string;
    placeholder?: string;
    type?: TextInputType;
    infoMessage?: string;
  }>(),
  {
    infoMessage: undefined,
    initialErrors: undefined,
    placeholder: undefined,
    id: undefined,
    disabled: false,
    label: undefined,
    inputVariant: 'text',
    value: '',
    type: 'text',
  },
);

if (props.type !== 'text' && props.inputVariant !== 'text') {
  console.warn('Props type has no affect because inputVariant is not text');
}

const {
  errorMessage,
  value: inputValue,
  setErrors,
} = useField<string>(toRef(props, 'name'), undefined, {
  validateOnValueUpdate: false,
  initialValue: props.value,
});

setErrors(props.initialErrors);

const isError = computed(() => errorMessage.value !== undefined);

const classes = computed(() =>
  [
    'tsw-text-field',
    props.disabled ? 'tsw-disabled' : '',
    errorMessage.value !== undefined ? 'tsw-text-input-error' : '',
    props.disabled || errorMessage.value !== undefined ? '' : 'tsw-text-input-actions',
  ].join(' '),
);
</script>

<style lang="postcss">
.tsw-text-field {
  @apply inline-grid gap-y-2 w-full;

  label {
    @apply text-sm text-white-95;
  }

  .error-message {
    @apply text-xs text-red-200;
  }

  .info-message {
    @apply text-xs text-gray-60;
  }
}
</style>
