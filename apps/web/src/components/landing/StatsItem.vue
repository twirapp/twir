<template>
  <li class="inline-flex flex-col items-center justify-center">
    <span class="text-4xl font-medium leading-[130%]">
      {{ item.stat }}
    </span>
    <span class="text-gray-60 leading-normal whitespace-nowrap">
      {{ item.description }}
    </span>
  </li>
</template>

<script lang="ts" setup>
import type { StatInfo } from '@/data/landing/statInfo.js';

defineProps<{
  item: StatInfo;
}>();
</script>
