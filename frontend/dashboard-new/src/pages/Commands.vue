<script setup lang='ts'>
import { Command } from '@twir/grpc/generated/api/api/commands';
import { NButton } from 'naive-ui';
import { useRoute } from 'vue-router';

import { useCommandsManager } from '@/api/index.js';
import List from '@/components/commands/list.vue';

const route = useRoute();
const commandsManager = useCommandsManager();
const { data: commandsResponse } = commandsManager.getAll({});

const commands: Command = [{
	name: 'test2',
	responses: [
		{ text: 'qwe' },
	],
}];
</script>

<template>
  <div>
    <div class="title">
      <div>
        <h2>Commands</h2>
      </div>
      <div>
        <n-button type="primary">
          Create
        </n-button>
      </div>
    </div>

    <list :commands="commands" />
  </div>
</template>

<style scoped>
.title {
	display: flex;
	justify-content: space-between;
	align-items: center;
}
</style>

