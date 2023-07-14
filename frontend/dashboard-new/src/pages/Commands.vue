<script setup lang='ts'>
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import { useCommandsManager } from '@/api/index.js';
import List from '@/components/commands/list.vue';

const route = useRoute();
const commandsManager = useCommandsManager();
const { data: commandsResponse } = commandsManager.getAll({});

const commands = computed(() => {
	const system = Array.isArray(route.params.system) ? route.params.system[0] : route.params.system;
	return commandsResponse.value?.commands.filter(c => c.module.toLowerCase() === system.toLowerCase()) ?? [];
});
</script>

<template>
  <list :commands="commands" :showHeader="true" />
</template>

<style scoped>
.title {
	display: flex;
	justify-content: space-between;
	align-items: center;
}
</style>

