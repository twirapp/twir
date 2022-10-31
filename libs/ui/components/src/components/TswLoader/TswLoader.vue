<template>
  <span :class="classes" />
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { LoaderSize } from '@/components/TswLoader/props.types';

const props = withDefaults(defineProps<{ size?: LoaderSize }>(), { size: 'md' });

const classes = computed(() => ['tsw-loader', `tsw-loader-${props.size}`].join(' '));
</script>

<style lang="postcss">
.tsw-loader {
  @apply inline-block relative;

  &::before {
    content: '';

    @apply border-white-95 border-b-purple-70 rounded-full animate-spin absolute m-auto top-0 bottom-0 left-0 right-0;
  }
}

.tsw-loader-sm {
  @apply w-6 h-6;

  &::before {
    @apply w-5 h-5 border-2;
  }
}

.tsw-loader-md {
  @apply w-9 h-9;

  &::before {
    @apply w-6 h-6 border-[2.5px];
  }
}

.tsw-loader-lg {
  @apply w-12 h-12;

  &::before {
    @apply w-8 h-8 border-[2.5px];
  }
}
</style>
