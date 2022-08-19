<!-- xmlns attribute not parsed -->
<!-- eslint-disable vue/no-parsing-error -->
<template>
  <component
    :is="icon"
    xmlns="http://www.w3.org/2000/svg"
    :width="size"
    :height="size"
    :stroke="stroke"
    :fill="fill"
    :stroke-width="strokeWidth"
    stroke-linecap="round"
    stroke-linejoin="round"
  />
</template>

<script lang="ts" setup>
import type { IconName } from '@tsuwari/ui-icons/icons';
import { computed, defineAsyncComponent } from 'vue';

const props = withDefaults(
  defineProps<{
    name: IconName;
    size?: string;
    fill?: string;
    stroke?: string;
    strokeWidth?: number;
  }>(),
  {
    fill: 'none',
    size: undefined,
    stroke: 'white',
    strokeWidth: 1.5,
  },
);

const icon = computed(() =>
  defineAsyncComponent(async () => {
    return (await import('@tsuwari/ui-icons/icons')).default[props.name];
  }),
);
</script>
