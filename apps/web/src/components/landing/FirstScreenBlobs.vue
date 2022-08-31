<template>
  <div class="blob purple-blob" :style="blobStyles"></div>
  <div class="blob blue-blob" :style="blobStyles"></div>
</template>

<script lang="ts" setup>
import { computed, StyleValue } from 'vue';

const blobStyles = computed<StyleValue>(() => {
  // In Firefox the blur is not the same as everywhere else and the color transparency is too high
  if (window.navigator.userAgent.includes('Firefox')) {
    return { animation: 'fadeInFirefox 1s', opacity: '0.25' };
  }
  return { animation: 'fadeIn 1s', opacity: '0.6' };
});
</script>

<style lang="postcss">
.blob {
  z-index: -10;
  filter: blur(180px);
  position: absolute;
  border-radius: 50%;
  animation: fadeIn 1s;
}

@keyframes fadeInFirefox {
  from {
    opacity: 0;
  }
  to {
    opacity: 0.25;
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 0.6;
  }
}

.purple-blob {
  width: 324px;
  height: 410px;
  transform: rotate(28deg);
  right: -100px;
  bottom: -40px;
  background-color: theme('colors.purple.60');
}

.blue-blob {
  width: 308px;
  height: 420px;
  transform: rotate(-26deg);
  left: -100px;
  bottom: -80px;
  background-color: #5767f8;
}
</style>
