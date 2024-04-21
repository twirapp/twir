<script setup lang="ts">
import { IconSearch } from '@tabler/icons-vue';
import { NModal, NButton, NInput, NIcon  } from 'naive-ui';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';

import List from './components/list.vue';
import { useCommandEdit } from './composables/use-command-edit';

import { useUserAccessFlagChecker } from '@/api';
import { useCommandsApi } from '@/api/commands/commands.js';
import ManageGroups from '@/components/commands/manageGroups.vue';

const route = useRoute();
const { t } = useI18n();
const userCanManageCommands = useUserAccessFlagChecker('MANAGE_COMMANDS');
const editCommand = useCommandEdit();

const commandsManager = useCommandsApi();
const { data: commandsResponse } = commandsManager.useQueryCommands();

const commandsFilter = ref('');
const commands = computed(() => {
	if (!commandsResponse.value?.commands) return [];

	const system = Array.isArray(route.params.system) ? route.params.system[0] : route.params.system;

	return commandsResponse.value.commands
		.filter(c => {
			if (system.toUpperCase() === 'CUSTOM') {
				return c.module === 'CUSTOM';
			}

			return c.module != 'CUSTOM';
		})
		.filter(c => {
			return c.name.includes(commandsFilter.value) ||
				c.aliases.some(a => a.includes(commandsFilter.value));
		});
});

const showManageGroupsModal = ref(false);
</script>

<template>
	<div class="flex flex-col gap-4">
		<div class="flex justify-between items-center flex-wrap gap-2">
			<div>
				<n-input
					v-model:value="commandsFilter"
					:placeholder="t('commands.searchPlaceholder')"
				>
					<template #prefix>
						<n-icon><IconSearch /></n-icon>
					</template>
				</n-input>
			</div>
			<div>
				<div class="flex gap-2">
					<n-button
						:disabled="!userCanManageCommands" secondary type="info"
						@click="showManageGroupsModal = true"
					>
						{{ t('commands.groups.manageButton') }}
					</n-button>

					<n-button
						secondary
						type="success"
						:disabled="!userCanManageCommands"
						@click="editCommand.createCommand"
					>
						{{ t('sharedButtons.create') }}
					</n-button>
				</div>
			</div>
		</div>

		<list
			:commands="commands"
			show-background
			enable-groups
		/>
	</div>

	<n-modal
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
		<manage-groups />
	</n-modal>
</template>
