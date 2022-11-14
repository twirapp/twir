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
      v-model.number="value"
      :name="nameRef"
      type="number"
      class="bg-[#202020] form-control input input-bordered input-sm px-3 py-1.5 rounded text-[#F5F5F5] w-full"
    >
  </div>
</template>

<script lang="ts" setup>
import { useField } from 'vee-validate';
import { toRef } from 'vue';

const props = withDefaults(defineProps<{
  id: string,
  name: string,
  label: string,
  direction?: 'row' | 'col'
}>(), { direction: 'row' });

const nameRef = toRef(props, 'name');
const { value, label: fieldLabel } = useField<number>(nameRef, {});
</script>