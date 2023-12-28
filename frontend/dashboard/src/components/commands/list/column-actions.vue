<script setup lang="ts">
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import { NButton, NPopconfirm } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { useCommandsManager, useUserAccessFlagChecker } from '@/api';
import { ListRowData } from '@/components/commands/types';


const emits = defineEmits<{ edit: [ListRowData] }>();
const props = defineProps<{ row: ListRowData }>();
const userCanManageCommands = useUserAccessFlagChecker('MANAGE_COMMANDS');

const manager = useCommandsManager();
const deleter = manager.deleteOne;

const { t } = useI18n();

function edit() {
	emits('edit', props.row);
}
</script>

<template>
	<template v-if="row.isGroup"></template>

	<div v-else class="actions">
		<n-button
			quaternary
			size="small"
			type="primary"
			:disabled="!userCanManageCommands"
			@click="edit"
		>
			<IconPencil />
		</n-button>
		<n-popconfirm
			:positive-text="t('deleteConfirmation.confirm')"
			:negative-text="t('deleteConfirmation.cancel')"
			:on-positive-click="() => deleter.mutate({ commandId: row.id })"
		>
			<template #trigger>
				<n-button quaternary size="small" type="error" :disabled="!userCanManageCommands">
					<IconTrash />
				</n-button>
			</template>
		</n-popconfirm>
	</div>
</template>

<style scoped>
.actions {
	display: flex;
	gap: 4px;
}
</style>
