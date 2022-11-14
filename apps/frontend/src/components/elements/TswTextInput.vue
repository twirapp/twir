<template>
  <div
    class="flex"
    :class="{
      'flex-col items-start': direction === 'col',
      'items-center': direction === 'row'
    }"
  >
    <label
      :for="id"
      :class="{
        'mr-4': direction === 'row',
        'mb-2': direction === 'col'
      }"
      class="inline-block label-text leading-tight text-[#AFAFAF] text-sm"
    >{{ label ?? fieldLabel }}</label>
    <input
      :id="id"
      v-model="value"
      :name="nameRef"
      type="text"
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
  direction?: 'row' | 'col',
  type?: 'text' | 'email'
}>(), { direction: 'row', type: 'text' });

const nameRef = toRef(props, 'name');
const { value, label: fieldLabel } = useField<string>(nameRef, {});
</script>