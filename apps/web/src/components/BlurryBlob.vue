<template>
  <div :class="`blob ${color}-blob`" :style="blobStyles"></div>
</template>

<script lang="ts" setup>
import { computed, StyleValue } from 'vue';

const blobStyles = computed<StyleValue>(() => {
  // In Firefox the blur is not the same as everywhere else and the color saturation is too high
  if (window.navigator.userAgent.includes('Firefox')) {
    return { animation: 'fadeInFirefox 1s', opacity: '0.25' };
  }
  return { animation: 'fadeInBlob 1s', opacity: '0.8' };
});

defineProps<{
  color: 'purple' | 'blue' | 'cyan';
}>();
</script>

<style lang="postcss">
.blob {
  z-index: -10;
  filter: blur(180px);
  position: absolute;
  border-radius: 50%;
}

@keyframes fadeInFirefox {
  from {
    opacity: 0;
  }
  to {
    opacity: 0.25;
  }
}

@keyframes fadeInBlob {
  from {
    opacity: 0;
  }
  to {
    opacity: 0.8;
  }
}

.purple-blob {
  /* background-color: theme('colors.purple.60'); */
  background: linear-gradient(180deg, #ff00e5 0%, #4617ff 100%);
}

.blue-blob {
  background: linear-gradient(180deg, #0094ff 0%, #801efd 100%);
  /* background-color: #5767f8; */
}

.cyan-blob {
  background: linear-gradient(180deg, #00e9f8 0%, #0c8aff 100%);
  /* background-color: #24777c; */
}
</style>
