<template>
  <component
    :is="isBgImage ? 'div' : 'img'"
    v-if="lazy ? isImageReady : true"
    :style="imgStyles"
    :src="isBgImage ? undefined : src"
  ></component>
  <div v-else ref="placeholder" :style="{ width: w, height: h }"></div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';

import useLazyImage from '@/hooks/useLazyImage.js';
import { cssPX, cssURL } from '@/utils/css.js';

const props = withDefaults(
  defineProps<{
    src: string;
    width?: number;
    height?: number;
    lazy?: boolean;
    renderType?: 'bg-image' | 'image';
  }>(),
  { lazy: false, renderType: 'image', width: undefined, height: undefined, alt: undefined },
);
const placeholder = ref<HTMLElement | undefined>();
const { isImageReady, execute } = useLazyImage(props.src, placeholder, false);

if (props.lazy) {
  if (!execute) {
    throw new Error('Cannot find execute function');
  }
  execute();
}

const w = computed(() => cssPX(props.width));
const h = computed(() => cssPX(props.height));
const backgroundImage = computed(() => cssURL(props.src));

const isBgImage = computed(() => props.renderType === 'bg-image');

const imgStyles = computed(() =>
  isBgImage.value
    ? {
        backgroundImage: backgroundImage.value,
        width: w.value,
        height: h.value,
      }
    : { width: w.value, height: h.value },
);
</script>
