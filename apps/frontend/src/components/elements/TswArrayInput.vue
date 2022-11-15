<template>
  <div>
    <span class="flex items-center label py-1">
      <span class="label-text text-[#AFAFAF] text-sm">{{ label ?? fieldLabel }}</span>
      <span
        class="bg-green-600 cursor-pointer duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-green-700 inline-block ml-2 p-1 rounded text-white text-xs transition"
        @click="value.push('')"
      ><Add /></span>
    </span>
    <div
      v-if="value.length"
      class="gap-x-3 gap-y-2 grid grid-cols-1 input-group lg:grid-cols-2 max-h-40 md:grid-cols-2 mt-1 overflow-x-hidden overflow-y-auto pt-1 scrollbar sm:grid-cols-2 xl:grid-cols-3"
    >
      <div
        v-for="(_, index) in value"
        :key="index"
        class="flex flex-wrap items-stretch relative"
      >
        <input
          v-model="value[index]"
          :name="`${nameRef}[${index}]`"
          type="text"
          class="bg-[#202020] border border-[#3E3E3E] flex-auto flex-grow flex-shrink leading-normal px-3 py-1.5 relative rounded rounded-r-none text-[#F5F5F5] w-px"
          @change="() => setTouched(true)"
        >
        <div
          class="-mr-px cursor-pointer flex"
          @click="value.splice(index, 1)"
        >
          <span
            class="bg-[#404040] flex hover:bg-red-700 items-center px-3 py-1.5 rounded rounded-l-none"
          >
            <Remove />
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useField } from 'vee-validate';
import { toRef } from 'vue';

import Add from '@/assets/buttons/add.svg?component';
import Remove from '@/assets/buttons/remove.svg?component';

const props = defineProps<{
  name: string,
  label?: string,
}>();

const nameRef = toRef(props, 'name');
const { value, label: fieldLabel, setTouched } = useField<string[]>(nameRef, {});
</script>