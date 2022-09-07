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
  return { animation: 'fadeInBlob 1s', opacity: '0.6' };
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
    opacity: 0.6;
  }
}

.purple-blob {
  background-color: theme('colors.purple.60');
}

.blue-blob {
  background-color: #5767f8;
}

.cyan-blob {
  background-color: #24777c;
}
</style>
