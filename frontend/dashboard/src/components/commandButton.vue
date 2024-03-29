<script setup lang="ts">
import { IconPencil } from '@tabler/icons-vue';
import { NButton, NModal } from 'naive-ui';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCommandsManager } from '@/api';
import CommandModal from '@/components/commands/modal.vue';

const props = defineProps<{
	name: string
	title?: string
}>();

const commandsManager = useCommandsManager();
const { data: commands } = commandsManager.getAll({});

const command = computed(() => commands.value?.commands.find((command) => command.defaultName === props.name));
const showCommandEditModal = ref(false);

const { t } = useI18n();
</script>

<template>
	<div style="display: flex; flex-direction: column;">
		<span>{{ props.title ?? t('games.command') }}</span>
		<div v-if="command" style="display: flex; gap: 5px;">
			<n-button secondary type="success" @click="() => showCommandEditModal = true">
				<div
					style="display: flex; align-items: center; min-width: 80px; justify-content: space-between;"
				>
					<span>!{{ command.name }}</span>
					<IconPencil />
				</div>
			</n-button>
		</div>
	</div>

	<n-modal
		v-if="command"
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
	>
		<command-modal
			:command="command"
			@close="() => {
				showCommandEditModal = false;
			}"
		/>
	</n-modal>
</template>
