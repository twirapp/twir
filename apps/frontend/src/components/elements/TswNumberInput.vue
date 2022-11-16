<template>
  <div
    class="flex"
    :class="{
      'flex-col items-start': direction === 'col'
    }"
  >
    <label
      :for="id"
      class="inline-block label-text leading-tight text-[#AFAFAF] text-sm"
      :class="{
        'mr-5': direction === 'row',
        'mb-2': direction === 'col'
      }"
    >{{ label ?? fieldLabel }}</label>
    <input
      :id="id"
      :name="nameRef"
      type="number"
      class="bg-[#202020] form-control input input-sm px-3 py-1.5 rounded text-[#F5F5F5] w-full"
      :class="{
        'border border-[#C83B2B]': isError
      }"
      :value="value"
      @change="change"
      @blur="handleBlur"
      @reset="handleReset"
    >
    <span
      v-if="isError"
      class="mt-2 text-red-500 text-xs"
    >{{ errorMessage }}</span>
  </div>
</template>

<script lang="ts" setup>
import { useField } from 'vee-validate';
import { computed, toRef } from 'vue';

const props = withDefaults(defineProps<{
  id: string,
  name: string,
  label: string,
  direction?: 'row' | 'col'
}>(), { direction: 'row' });

const change = (e: Event) => {
  handleChange(+(e.target as HTMLInputElement).value);
};

const nameRef = toRef(props, 'name');
const {  label: fieldLabel, errorMessage, errors, handleChange, handleBlur, handleReset, value } = useField<number>(nameRef, {});

const isError = computed(() => errors.value.length > 0);
</script>