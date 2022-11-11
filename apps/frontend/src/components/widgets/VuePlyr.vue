<template>
  <video ref="target" />
</template>

<script setup lang="ts">
import Plyr, { Options } from 'plyr';
import 'plyr/dist/plyr.css';
import { onBeforeUnmount, onMounted, ref } from 'vue';

const props = withDefaults(
  defineProps<{
    options?: Options;
  }>(),
  { options: () => ({}) },
);

const target = ref<HTMLElement | null>(null);
const player = ref<Plyr | null>(null);

const emit = defineEmits<{
  (event: 'init', player: Plyr): void;
}>();

onMounted(() => {
  if (target.value === null) {
    throw new Error('Cannot find video element to init a Flyr');
  }
  player.value = new Plyr(target.value, props.options);
  emit('init', player.value as Plyr);
});

onBeforeUnmount(() => {
  try {
    if (player.value === null) {
      return console.warn('Cannot find Plyr player to destroy it');
    }
    player.value.destroy();
  } catch (e) {
    console.error('Cannot destroy Plyr player');
  }
});
</script>
