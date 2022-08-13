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
import * as icons from '@tsuwari/ui-icons/icons';
import { computed, FunctionalComponent, SVGAttributes } from 'vue';

const props = withDefaults(
  defineProps<{
    name: icons.IconName;
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

type Icons = {
  [K in icons.IconName]: FunctionalComponent<SVGAttributes>;
};

const icon = computed(() => (icons as Icons)[props.name]);
</script>
