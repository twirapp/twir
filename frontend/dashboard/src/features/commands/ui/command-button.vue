<script setup lang="ts">
import { IconPencil } from '@tabler/icons-vue'
import { NButton } from 'naive-ui'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useCommandsApi } from '@/api/commands/commands.js'

const props = defineProps<{
	name: string
	title?: string
}>()

const commandsManager = useCommandsApi()

const { data: commands } = commandsManager.useQueryCommands()

const command = computed(() => commands.value?.commands.find((command) => command.defaultName === props.name))

const { t } = useI18n()
</script>

<template>
	<div class="flex flex-col">
		<span>{{ props.title ?? t('games.command') }}</span>
		<div v-if="command" class="flex gap-1">
			<RouterLink v-slot="{ href, navigate }" custom :to="`/dashboard/commands/${command.module.toLowerCase()}/${command.id}`">
				<NButton tag="a" :href="href" secondary type="success" @click="navigate">
					<div class="flex items-center min-w-20 justify-between">
						<span>!{{ command.name }}</span>
						<IconPencil />
					</div>
				</NButton>
			</RouterLink>
		</div>
	</div>
</template>
