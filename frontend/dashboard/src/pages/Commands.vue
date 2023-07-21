<script setup lang='ts'>
import { NResult } from 'naive-ui';
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import { useCommandsManager, useUserAccessFlagChecker } from '@/api/index.js';
import List from '@/components/commands/list.vue';

const route = useRoute();
const commandsManager = useCommandsManager();
const { data: commandsResponse } = commandsManager.getAll({});

const userCanViewCommands = useUserAccessFlagChecker('VIEW_COMMANDS');

const commands = computed(() => {
	const system = Array.isArray(route.params.system) ? route.params.system[0] : route.params.system;
	return commandsResponse.value?.commands.filter(c => c.module.toLowerCase() === system.toLowerCase()) ?? [];
});
</script>

<template>
	<n-result v-if="!userCanViewCommands" status="403" title="You haven't access to view commands" />
	<list v-else :commands="commands" :showHeader="true" :showCreateButton="route.params.system === 'custom'" />
</template>

<style scoped>
.title {
	display: flex;
	justify-content: space-between;
	align-items: center;
}
</style>

