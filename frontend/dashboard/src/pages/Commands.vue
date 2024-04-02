<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import { useCommandsManager } from '@/api/index.js';
import List from '@/components/commands/list.vue';

const route = useRoute();
const commandsManager = useCommandsManager();
const { data: commandsResponse } = commandsManager.getAll({});

const commands = computed(() => {
	const system = Array.isArray(route.params.system) ? route.params.system[0] : route.params.system;

	return commandsResponse.value?.commands.filter(c => {
		if (system.toUpperCase() === 'CUSTOM') {
			return c.module === 'CUSTOM';
		}

		return c.module != 'CUSTOM';
	}) ?? [];
});
</script>

<template>
	<list
		:commands="commands"
		show-header
		enable-groups
		show-background
		:showCreateButton="route.params.system === 'custom'"
	/>
</template>
