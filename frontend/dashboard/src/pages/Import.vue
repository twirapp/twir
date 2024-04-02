<script setup lang="ts">
import Nightbot from './import/nightbot.vue';
import StreamElements from './import/streamelements.vue';
import Streamlabs from './import/streamlabs.vue';
</script>

<template>
	<div class="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
		<nightbot />
		<stream-elements />
		<streamlabs />
	</div>
</template>
