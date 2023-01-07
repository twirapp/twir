<template>
  <div class="inline-flex flex-col items-center justify-center">
    <span class="text-4xl font-medium leading-[130%]">
      {{ useStatsFormatter().format(item.count) }}
    </span>
    <span class="text-gray-60 leading-normal whitespace-nowrap">
      {{ item.name.at(0)?.toUpperCase() + item.name.slice(1) }}
    </span>
  </div>
</template>

<script lang="ts" setup>
import { useStatsFormatter, StatInfo } from '@/services/landing';

defineProps<{
  item: StatInfo;
}>();
</script>
