<template>
  <TswIcon
    :name="arrowName"
    :size="size"
    :class="classes"
    :strokeWidth="strokeWidth"
    :stroke="stroke"
  />
</template>

<script lang="ts" setup>
import { computed } from 'vue';

import type { ArrowDirection, ArrowIconName } from '@/components/TswArrowIcon/props.types';
import TswIcon from '@/components/TswIcon/TswIcon.vue';

const props = withDefaults(
  defineProps<{
    direction: ArrowDirection;
    arrowName?: ArrowIconName;
    stroke?: string;
    strokeWidth?: number;
    size?: string;
  }>(),
  {
    direction: 'right',
    stroke: 'white',
    arrowName: 'ArrowTriangleMedium',
    strokeWidth: 1.5,
    size: undefined,
  },
);

const classes = computed(() => [`arrow-${props.direction}`].join(' '));
</script>

<style lang="postcss">

.arrow-bottom-right {
  transform: rotate(45deg);
}

.arrow-bottom {
  transform: rotate(90deg);
}

.arrow-bottom-left {
  transform: rotate(135deg);
}

.arrow-left {
  transform: rotate(180deg);
}

.arrow-top-left {
  transform: rotate(225deg);
}

.arrow-top {
  transform: rotate(-90deg);
}

.arrow-top-right {
  transform: rotate(-45deg);
}

.arrow-right {
  transform: rotate(0deg);
}

</style>
