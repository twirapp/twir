<template>
  <button
    :class="`flex text-sm justify-between py-[6px] px-[8px] rounded items-center ${
      isActive ? 'bg-purple-60 text-white-100' : 'text-white-95 hover:bg-black-20'
    }`"
    @click="emit('change', locale)"
  >
    {{ name }}
    <TswIcon
      v-if="isActive"
      name="Check"
      size="16px"
      direction="bottom"
      :strokeWidth="1.5"
      :class="`${isActive ? 'stroke-white-100' : 'stroke-gray-60'}`"
    />
  </button>
</template>

<script lang="ts" setup>
import { TswIcon } from '@tsuwari/ui-components';

import type { Locale } from '@/types/locale';

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
