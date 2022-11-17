<template>
  <component
    :is="tag"
    v-if="lazy ? isImageReady : true"
    class="inline-flex bg-contain rounded-full flex-shrink-0"
    :href="href"
    :style="{ backgroundImage: computedBgImage, width: computedSize, height: computedSize }"
  />
  <component
    :is="tag"
    v-else
    ref="placeholder"
    :href="href"
    class="inline-flex flex-shrink-0 bg-gray-30 rounded-full"
    :style="{ width: computedSize, height: computedSize }"
  />
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';

import useLazyImage from '@/hooks/useLazyImage.js';
import { cssPX, cssURL } from '@/utils/css.js';

const props = withDefaults(
  defineProps<{
    src: string;
    href?: string;
    size?: number;
    lazy?: boolean;
  }>(),
  {
    lazy: false,
    href: undefined,
    size: 36,
  },
);

const tag = computed(() => (props.href ? 'a' : 'div'));
const computedSize = computed(() => cssPX(props.size));
const computedBgImage = computed(() => cssURL(props.src));

const placeholder = ref<HTMLElement | undefined>();
const { isImageReady, execute } = useLazyImage(props.src, placeholder, false);

if (props.lazy) {
  if (!execute) {
    throw new Error('Cannot find execute function');
  }
  execute();
}
</script>
