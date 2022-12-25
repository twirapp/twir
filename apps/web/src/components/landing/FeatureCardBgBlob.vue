<template>
  <Transition>
    <div
      v-if="isActive"
      class="absolute -z-[1] select-none bg-contain bg-no-repeat"
      :style="styles"
    />
  </Transition>
</template>

<script lang="ts" setup>
import { toRef } from 'vue';

import PinkBlob from '@/assets/blob-pink.png';
import useFollowingBgImage from '@/utils/useFollowingBgImage.js';

const props =
  defineProps<{
    card: HTMLElement | null;
  }>();

const cardRef = toRef(props, 'card');

const { styles, isActive } = useFollowingBgImage(cardRef, PinkBlob, 0.5);
</script>

<style lang="postcss" scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.3s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
