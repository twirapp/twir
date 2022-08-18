<template>
  <span :class="classes" />
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { LoaderSize } from '@/components/TswLoader/props.types';

const props = withDefaults(defineProps<{ size?: LoaderSize }>(), { size: 'md' });

const classes = computed(() => ['loader', `loader-${props.size}`].join(' '));
</script>

<style lang="postcss">
.loader {
  @apply inline-flex animate-spin border-b-purple-70 rounded-full border-white-95;
}

.loader-sm {
  @apply border w-[13px] h-[13px] m-[1px];
}

.loader-md {
  @apply border-[1.5px] w-[15px] h-[15px] m-[2px];
}

.loader-lg {
  @apply border-2 w-[21px] h-[21px] m-[2.5px];
}
</style>
