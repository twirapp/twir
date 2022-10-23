<template>
  <div v-if="!isTargetVisible || !isReady" ref="placeholder" :style="style" :class="classes"></div>
  <component
    :is="renderType === 'bg-image' ? 'div' : image"
    v-else
    :class="classes"
    :style="styles"
  ></component>
</template>

<script lang="ts" setup>
import { useImage, useIntersectionObserver } from '@vueuse/core';
import { computed, CSSProperties, ref } from 'vue';

import { cssURL } from '@/utils/css.js';

const props =
  defineProps<{
    src: string;
    style?: CSSProperties;
    class?: string;
    renderType: 'bg-image' | 'image';
  }>();

const placeholder = ref<HTMLElement | undefined>();
const isTargetVisible = ref<boolean>(false);
const { isReady, execute, state: image } = useImage({ src: props.src }, { immediate: false });

const { stop } = useIntersectionObserver(placeholder, async ([{ isIntersecting }]) => {
  if (isIntersecting) {
    isTargetVisible.value = true;
    stop();
    execute();
  }
});

const styles = computed(() => {
  if (props.renderType === 'bg-image') {
    return { backgroundImage: cssURL(props.src), ...props.style };
  }
  return props.style;
});

const classes = computed(() => props.class);
</script>
