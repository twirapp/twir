<script setup lang="ts">
import { PencilIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useCommandsApi } from '@/api/commands/commands.js'
import { Button } from '@/components/ui/button'

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
				<Button as="a" :href="href" variant="default" @click="navigate">
					<div class="flex items-center min-w-20 justify-between gap-2">
						<span>!{{ command.name }}</span>
						<PencilIcon class="h-4 w-4" />
					</div>
				</Button>
			</RouterLink>
		</div>
	</div>
</template>
