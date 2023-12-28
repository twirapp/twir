<script setup lang="ts">
import { useNotification, NSwitch } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { useCommandsManager, useUserAccessFlagChecker } from '@/api';
import { ListRowData } from '@/components/commands/types';


const props = defineProps<{ row: ListRowData }>();
const userCanManageCommands = useUserAccessFlagChecker('MANAGE_COMMANDS');

const manager = useCommandsManager();
const patcher = manager.patch!;

const { t } = useI18n();
const message = useNotification();

async function save(newValue: boolean) {
	await patcher?.mutateAsync({ commandId: props.row.id, enabled: newValue });

	message.success({
		title: t('sharedTexts.saved'),
		duration: 2500,
	});
}
</script>

<template>
	<n-switch
		v-if="!row.isGroup"
		:disabled="!userCanManageCommands"
		:value="row.enabled"
		@update-value="(newValue) => save(newValue)"
	></n-switch>
</template>

<style scoped>

</style>
