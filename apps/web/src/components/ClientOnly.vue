<template>
  <slot v-if="show" />
  <slot v-else name="server" />
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';

const show = ref(false);

onMounted(() => {
  show.value = true;
});
</script>
