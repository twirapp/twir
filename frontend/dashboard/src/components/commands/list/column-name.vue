<script setup lang="ts">
import { IconDeviceFloppy, IconX } from '@tabler/icons-vue';
import { NTag, NButton, NInput, useNotification, NText } from 'naive-ui';
import { ref, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCommandsManager, useUserAccessFlagChecker } from '@/api';
import { ListRowData } from '@/components/commands/types';

const props = defineProps<{ row: ListRowData }>();
const userCanManageCommands = useUserAccessFlagChecker('MANAGE_COMMANDS');

const rgbaPattern = /rgba?\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*(?:,\s*([\d.]+)\s*)?\)/;
const computeGroupTextColor = (color?: string) => {
	const result = rgbaPattern.exec(color ?? '');
	if (!result) return '#c2b7b7';
	const [r, g, b] = result.slice(1).map(i => parseInt(i, 10));

	const bright = (
		(((r * 299) + (g * 587) + (b * 114)) / 1000) - 128
	) * -1000;

	return `rgba(${bright},${bright},${bright})`;
};

const name = ref(toRaw(props.row.name));
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
			name: name.value,
		},
	});

	message.success({
		title: t('sharedTexts.saved'),
		duration: 2500,
	});

	isEdit.value = false;
}


function reset() {
	isEdit.value = false;
	name.value = props.row.name;
}

function setEdit() {
	if (!userCanManageCommands.value) return;

	isEdit.value = true;
}
</script>

<template>
	<n-tag v-if="row.isGroup" :style="{ backgroundColor: row.groupColor }">
		<p :style="{ color: computeGroupTextColor(row.groupColor) }">
			{{ row.name }}
		</p>
	</n-tag>
	<n-text v-else-if="!isEdit" class="name" @click="setEdit">
		{{ row.name }}
	</n-text>
	<div v-else class="editable-name">
		<n-input v-model:value="name" size="small" />
		<div class="actions">
			<n-button text type="error" size="tiny" @click="reset">
				<IconX />
			</n-button>
			<n-button text type="success" size="tiny" @click="save">
				<IconDeviceFloppy />
			</n-button>
		</div>
	</div>
</template>

<style scoped>
.name {
	text-decoration: underline dotted;
}

.editable-name {
	display: flex;
	align-items: center;
	gap: 4px;
	justify-content: space-between;
}

.actions {
	display: flex;
	gap: 2px;
}
</style>
