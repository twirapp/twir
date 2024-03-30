<script setup lang="ts">
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import { type Command } from '@twir/api/messages/commands/commands';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCommandsManager, useUserAccessFlagChecker } from '@/api';
import { Button } from '@/components/ui/button';
import DeleteConfirmation from '@/components/ui/delete-confirm.vue';
import { Switch } from '@/components/ui/switch';
import { useToast } from '@/components/ui/toast/use-toast';

const emits = defineEmits<{ edit: [] }>();
const props = defineProps<{ row: Command }>();
const userCanManageCommands = useUserAccessFlagChecker('MANAGE_COMMANDS');

const manager = useCommandsManager();
const deleter = manager.deleteOne;
const patcher = manager.patch!;

const { t } = useI18n();
const { toast } = useToast();

function edit() {
	emits('edit');
}

const showDelete = ref(false);

async function switchEnabled(newValue: boolean) {
	await patcher?.mutateAsync({ commandId: props.row.id, enabled: newValue });

	toast({
		title: t('sharedTexts.saved'),
		variant: 'success',
		duration: 1500,
	});
}
</script>

<template>
	<div class="flex items-center gap-4">
		<Switch
			:disabled="!userCanManageCommands"
			:checked="row.enabled"
			class="data-[state=unchecked]:bg-zinc-400"
			@update:checked="switchEnabled"
		/>
		<Button :disabled="!userCanManageCommands" size="icon" @click="edit">
			<IconPencil class="h-5 w-5" />
		</Button>
		<Button
			v-if="row.module === 'CUSTOM'"
			:disabled="!userCanManageCommands"
			variant="destructive"
			size="icon"
			@click="showDelete = true"
		>
			<IconTrash class="h-5 w-5" />
		</Button>
	</div>

	<DeleteConfirmation
		v-model:open="showDelete"
		@confirm="deleter.mutateAsync({ commandId: row.id })"
	/>
</template>

<style scoped>
.actions {
	display: flex;
	gap: 4px;
}
</style>
