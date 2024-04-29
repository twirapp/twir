<script setup lang="ts">
import { IconPencil } from '@tabler/icons-vue'
import { NButton } from 'naive-ui'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import EditModal from './edit-modal.vue'

import { useCommandsApi } from '@/api/commands/commands'
import { useCommandEdit } from '@/features/commands/composables/use-command-edit'

const props = defineProps<{
	name: string
	title?: string
}>()

const commandsManager = useCommandsApi()
const commandEdit = useCommandEdit()

const { data: commands } = commandsManager.useQueryCommands()

const command = computed(() => commands.value?.commands.find((command) => command.defaultName === props.name))

const { t } = useI18n()
</script>

<template>
	<div class="flex flex-col">
		<span>{{ props.title ?? t('games.command') }}</span>
		<div v-if="command" class="flex gap-1">
			<NButton secondary type="success" @click="() => commandEdit.editCommand(command!.id)">
				<div class="flex items-center min-w-20 justify-between">
					<span>!{{ command.name }}</span>
					<IconPencil />
				</div>
			</NButton>
		</div>
	</div>

	<EditModal />
</template>
