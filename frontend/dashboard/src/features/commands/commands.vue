<script setup lang="ts">
import { SearchIcon } from 'lucide-vue-next'
import { NButton, NIcon, NInput, NModal } from 'naive-ui'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import { useCommandEdit } from './composables/use-command-edit'
import List from './ui/list.vue'

import { useUserAccessFlagChecker } from '@/api'
import { useCommandsApi } from '@/api/commands/commands.js'
import ManageGroups from '@/components/commands/manageGroups.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const route = useRoute()
const { t } = useI18n()
const userCanManageCommands = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageCommands)
const editCommand = useCommandEdit()

const commandsManager = useCommandsApi()
const { data: commandsResponse } = commandsManager.useQueryCommands()

const excludedModules = ['7tv']

const commandsFilter = ref('')
const commands = computed(() => {
	if (!commandsResponse.value?.commands) return []

	const system = Array.isArray(route.params.system) ? route.params.system[0] : route.params.system

	return commandsResponse.value.commands
		.filter(c => {
			if (system.toUpperCase() === 'CUSTOM') {
				return c.module === 'CUSTOM'
			}

			return c.module !== 'CUSTOM' && !excludedModules.includes(c.module)
		})
		.filter(c => {
			return c.name.includes(commandsFilter.value) || c.aliases.some(a => a.includes(commandsFilter.value))
		})
})

const showManageGroupsModal = ref(false)
</script>

<template>
	<div class="flex flex-col gap-4">
		<div class="flex justify-between items-center flex-wrap gap-2">
			<div>
				<NInput
					v-model:value="commandsFilter"
					:placeholder="t('commands.searchPlaceholder')"
				>
					<template #prefix>
						<NIcon><SearchIcon /></NIcon>
					</template>
				</NInput>
			</div>
			<div>
				<div class="flex gap-2">
					<NButton
						:disabled="!userCanManageCommands" secondary type="info"
						@click="showManageGroupsModal = true"
					>
						{{ t('commands.groups.manageButton') }}
					</NButton>

					<NButton
						secondary
						type="success"
						:disabled="!userCanManageCommands"
						@click="editCommand.createCommand"
					>
						{{ t('sharedButtons.create') }}
					</NButton>
				</div>
			</div>
		</div>

		<List
			:commands="commands"
			show-background
			enable-groups
		/>
	</div>

	<NModal
		v-model:show="showManageGroupsModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="t('commands.groups.manageButton')"
		class="modal"
		:style="{
			width: '600px',
		}"
		:on-close="() => showManageGroupsModal = false"
	>
		<ManageGroups />
	</NModal>
</template>
