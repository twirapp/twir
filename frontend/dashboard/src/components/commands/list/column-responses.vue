<script setup lang="ts">
import { IconDeviceFloppy, IconX } from '@tabler/icons-vue';
import type { Command_Response } from '@twir/grpc/generated/api/api/commands';
import { useNotification, NText, NInput, NButton, NPopconfirm } from 'naive-ui';
import { ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCommandsManager, useUserAccessFlagChecker } from '@/api';
import { ListRowData } from '@/components/commands/types';

const props = defineProps<{ row: ListRowData }>();
const userCanManageCommands = useUserAccessFlagChecker('MANAGE_COMMANDS');
const responses = ref<Array<Command_Response>>([]);

watch(() => props.row.responses, (v) => {
	responses.value = structuredClone(toRaw(v));
}, { immediate: true });

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
	<div v-else>
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
				@keydown.enter="save"
			/>
		</template>

		<n-popconfirm
			v-else
			:show-icon="false"
			style="width: 400px"
			:show="isEdit"
			trigger="click"
			placement="bottom-start"
			@clickoutside="reset"
		>
			<template #trigger>
				<div class="responses-list" @click="setEdit">
					<n-text
						v-for="(response) in row.responses"
						:key="response.id"
						class="response"
					>
						{{ response.text }}
					</n-text>
				</div>
			</template>

			<div class="responses-list">
				<n-input
					v-for="(response) in responses"
					:key="response.id"
					v-model:value="response.text"
					size="tiny"
					type="textarea"
					autosize
					:maxlength="500"
				/>
			</div>

			<template #action>
				<div class="actions">
					<n-button secondary type="error" size="tiny" @click="reset">
						<IconX />
					</n-button>
					<n-button secondary type="success" size="tiny" @click="save">
						<IconDeviceFloppy />
					</n-button>
				</div>
			</template>
		</n-popconfirm>
	</div>
</template>

<style scoped>
.editable-responses {
	display: flex;
	gap: 4px;
	align-items: center;
}

.responses-list {
	display: flex;
	flex-direction: column;
	gap: 4px;
	width: 100%;
}

.response {
	text-decoration: underline dotted;
}

.actions {
	display: flex;
	gap: 2px
}
</style>
