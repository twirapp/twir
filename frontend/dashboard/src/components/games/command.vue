<script setup lang="ts">
import { computed, ref } from 'vue';
import { NButton, NModal } from 'naive-ui'
import { IconPencil } from '@tabler/icons-vue';

import { useCommandsManager } from '@/api';
import CommandModal from '../commands/modal.vue';

const props = defineProps<{
	name: string
}>();

const commandsManager = useCommandsManager();
const { data: commands } = commandsManager.getAll({});

const command = computed(() => commands.value?.commands.find((command) => command.defaultName === props.name));
const showCommandEditModal = ref(false)
</script>

<template>
	<h3>Command</h3>
	<div v-if="command" style="display: flex; gap: 5px;">
		<n-button secondary type="success" @click="() => showCommandEditModal = true">
			<div style="display: flex; align-items: center; min-width: 80px; justify-content: space-between;">
				<span>{{ command.name }}</span>
				<IconPencil />
			</div>
		</n-button>
	</div>

	<n-modal
		v-model:show="showCommandEditModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="command?.name"
		class="modal"
		:style="{
			width: '800px',
			top: '50px',
		}"
		:on-close="() => showCommandEditModal = false"
		v-if="command"
	>
		<command-modal
			:command="command"
			@close="() => {
				showCommandEditModal = false;
			}"
		/>
	</n-modal>
</template>
