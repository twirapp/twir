<template>
  <TswIcon
    :name="iconName"
    :size="size"
    :class="classes"
    :strokeWidth="strokeWidth"
    :stroke="stroke"
  />
</template>

<script lang="ts" setup>
import type { IconName } from '@tsuwari/ui-icons/icons';
import { computed } from 'vue';

import type { ArrowDirection, ArrowSize } from '@/components/TswArrowIcon/props.types';
import TswIcon from '@/components/TswIcon/TswIcon.vue';

const props = withDefaults(
  defineProps<{
    direction: ArrowDirection;
    arrowSize?: ArrowSize;
    stroke?: string;
    strokeWidth?: number;
    size?: string;
  }>(),
  {
    direction: 'right',
    stroke: 'white',
    arrowSize: 'md',
    strokeWidth: 1.5,
    size: undefined,
  },
);

const iconName = computed<IconName>(() => {
  switch (props.arrowSize) {
    case 'lg':
      return 'ArrowLarge';
    case 'md':
      return 'ArrowMedium';
    case 'in-circle':
      return 'ArrowInCircle';
    default:
      return 'ArrowMedium';
  }
});

const classes = computed(() => [`arrow-${props.direction}`].join(' '));
</script>

<style lang="postcss">
.arrow-right {
  transform: rotate(0deg);
}

.arrow-bottom {
  transform: rotate(90deg);
}

.arrow-left {
  transform: rotate(180deg);
}

.arrow-top {
  transform: rotate(-90deg);
}
</style>
