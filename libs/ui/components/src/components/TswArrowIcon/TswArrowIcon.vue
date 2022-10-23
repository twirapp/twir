<template>
  <TswIcon :name="arrowName" :class="classes" :strokeWidth="strokeWidth" />
</template>

<script lang="ts" setup>
import { computed } from 'vue';

import type { ArrowDirection, ArrowIconName } from '@/components/TswArrowIcon/props.types';
import TswIcon from '@/components/TswIcon/TswIcon.vue';

const props = withDefaults(
  defineProps<{
    direction: ArrowDirection;
    arrowName?: ArrowIconName;
    strokeWidth?: number;
    size?: string;
    class?: string;
  }>(),
  {
    direction: 'right',
    arrowName: 'ArrowTriangleMedium',
    strokeWidth: undefined,
    size: undefined,
    class: undefined,
  },
);

const classes = computed(() => {
  return props.class + ` arrow-${props.direction}`;
});
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
