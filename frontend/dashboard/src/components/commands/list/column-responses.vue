<script setup lang="ts">
import { IconDeviceFloppy, IconX } from '@tabler/icons-vue';
import { useNotification, NText, NInput, NButton } from 'naive-ui';
import { ref, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCommandsManager, useUserAccessFlagChecker } from '@/api';
import { ListRowData } from '@/components/commands/types';

const props = defineProps<{ row: ListRowData }>();
const userCanManageCommands = useUserAccessFlagChecker('MANAGE_COMMANDS');
const responses = ref(structuredClone(toRaw(props.row.responses)));
const description = ref(toRaw(props.row.description));
const isEdit = ref(false);

const manager = useCommandsManager();
const updater = manager.update;

const { t } = useI18n();
const message = useNotification();

async function save() {
	await updater.mutateAsync({
		id: props.row.id,
		command: {
			...props.row,
			responses: responses.value,
			description: description.value,
		},
	});

	message.success({
		title: t('sharedTexts.saved'),
		duration: 2500,
	});

	isEdit.value = false;
}

function reset() {
	responses.value = structuredClone(toRaw(props.row.responses));
	description.value = toRaw(props.row.description);
	isEdit.value = false;
}

function setEdit() {
	if (!userCanManageCommands.value) return;

	isEdit.value = true;
}
</script>

<template>
	<div v-if="row.isGroup"></div>
	<div v-else class="editable-responses">
		<div v-if="isEdit" class="actions">
			<n-button text type="error" size="tiny" @click="reset">
				<IconX />
			</n-button>
			<n-button text type="success" size="tiny" @click="save">
				<IconDeviceFloppy />
			</n-button>
		</div>

		<template v-if="row.module !== 'CUSTOM'">
			<n-text v-if="!isEdit" class="response" @click="setEdit">
				{{ row.description }}
			</n-text>
			<n-input
				v-else
				v-model:value="description"
				size="tiny"
				type="textarea"
				autosize
			/>
		</template>

		<template v-else>
			<div v-if="!isEdit" class="responses-list">
				<n-text
					v-for="response in responses"
					:key="response.id"
					class="response"
					@click="setEdit"
				>
					{{ response.text }}
				</n-text>
			</div>
			<div v-else class="responses-list">
				<n-input
					v-for="response in responses"
					:key="response.id"
					v-model:value="response.text"
					:placeholder="t('commands.customResponses.placeholder')"
					size="tiny"
					type="textarea"
					autosize
				/>
			</div>
		</template>
	</div>
</template>

<style scoped>
.editable-responses {
	display: flex;
	gap: 4px;
}

.responses-list {
	display: flex;
	flex-direction: column;
	gap: 4px;
}

.response {
	text-decoration: underline dotted;
}

.actions {
	display: flex;
	flex-direction: column;
	gap: 4px
}
</style>
