<script setup lang='ts'>
import 'plyr/dist/plyr.css';
import Plyr from 'plyr';
import { ref, computed, onMounted, onUnmounted } from 'vue';

const player = ref<HTMLVideoElement | null>(null);

const plyr = computed(() => {
	if (!player.value) return null;

	return new Plyr(player.value);
});

onMounted(() => {
	if (!plyr.value) {
		throw new Error('Plyr is not initialized');
	}
});

onUnmounted(() => {
	if (!plyr.value) return;

	plyr.value!.destroy();
});
</script>

<template>
  <video ref="player" />
</template>

<style scoped lang='postcss'>

</style>
