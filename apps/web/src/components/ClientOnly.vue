<template>
  <slot v-if="isClient && renderClient" />
  <slot v-else name="server" />
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';

withDefaults(
  defineProps<{
    renderClient?: boolean;
  }>(),
  { renderClient: true },
);

const isClient = ref(false);

onMounted(() => {
  isClient.value = true;
});
</script>
