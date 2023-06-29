<template>
  <button
    :class="`flex text-sm justify-between py-[6px] px-[7px] rounded-md items-center ${
      isActive ? 'bg-purple-60 text-white-100' : 'text-white-95 hover:bg-black-20'
    }`"
    @click="emit('change', locale)"
  >
    <TswIcon
      v-if="isActive"
      name="Check"
      :height="16"
      :width="16"
      class="stroke-white-100 absolute"
    />
    <span class="pl-[22px] pr-2 text-sm">{{ name }}</span>
  </button>
</template>

<script lang="ts" setup>
import { TswIcon } from '@twir/ui-components';

import type { Locale } from '@/locales';

defineProps<{
  name: string;
  locale: Locale;
  isActive: boolean;
}>();

const emit =
  defineEmits<{
    (event: 'change', locale: Locale): void;
  }>();
</script>
