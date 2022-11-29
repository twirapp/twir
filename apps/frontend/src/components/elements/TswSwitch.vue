<template>
  <div
    class="flex form-check"
    :class="{
      'flex-col items-start': direction === 'col',
    }"
  >
    <label
      class="form-check-label inline-block leading-tight text-[#AFAFAF] text-sm"
      :for="id"
    >
      {{ label ?? fieldLabel }}
    </label>
    <div
      class="form-switch"
      :class="{
        'pl-6': direction === 'row',
        'pl-0 pt-2': direction === 'col',
      }"
    >
      <input
        :id="id"
        v-model="value"
        :name="nameRef"
        class="align-top appearance-none bg-[#595959] bg-contain bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
        type="checkbox"
        @change="() => setTouched(true)"
      >
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useField } from 'vee-validate';
import { toRef } from 'vue';

const props = withDefaults(
  defineProps<{
    id: string;
    name: string;
    label: string;
    direction?: 'row' | 'col';
  }>(),
  { direction: 'row' },
);

const nameRef = toRef(props, 'name');
const { value, label: fieldLabel, setTouched } = useField<boolean>(nameRef);
</script>
