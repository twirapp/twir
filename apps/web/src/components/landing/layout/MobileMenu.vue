<template>
  <div
    class="mobile-menu hidden"
    :style="{
      top: `${headerHeight}px`,
      height: `${windowHeight - headerHeight}px`,
    }"
  ></div>
</template>

<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { useWindowSize } from '@vueuse/core';

import { headerHeightStore } from '@/stores/landing/header.js';

const headerHeight = useStore(headerHeightStore);
const { height: windowHeight } = useWindowSize();
</script>
