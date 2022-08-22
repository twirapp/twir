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
import type { IconName } from '@tsuwari/ui-icons';
import { computed } from 'vue';

import type { ArrowDirection, ArrowSize } from '@/components/TswArrowIcon/props.types';
import TswIcon from '@/components/TswIcon/TswIcon.vue';

const props = withDefaults(
  defineProps<{
    direction: ArrowDirection;
    arrowType?: ArrowSize;
    stroke?: string;
    strokeWidth?: number;
    size?: string;
  }>(),
  {
    direction: 'right',
    stroke: 'white',
    arrowType: 'triangle-md',
    strokeWidth: 1.5,
    size: undefined,
  },
);

const iconName = computed<IconName>(() => {
  switch (props.arrowType) {
    case 'triangle-lg':
      return 'ArrowLarge';
    case 'triangle-md':
      return 'ArrowMedium';
    case 'in-circle':
      return 'ArrowInCircle';
    case 'narrow':
      return 'ArrowNarrow';
    default:
      return 'ArrowMedium';
  }
});

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
