<script setup lang="ts">
import { IconDeviceFloppy, IconX } from '@tabler/icons-vue';
import { NTag, NButton, NInput, useNotification, NText, NPopconfirm } from 'naive-ui';
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
	<div class="command-name">
		<n-tag v-if="row.isGroup" :style="{ backgroundColor: row.groupColor }">
			<p :style="{ color: computeGroupTextColor(row.groupColor) }">
				{{ row.name }}
			</p>
		</n-tag>

		<n-popconfirm
			v-else
			:show-icon="false"
			style="width: 200px"
			:show="isEdit"
			trigger="click"
			placement="bottom-start"
			@clickoutside="reset"
		>
			<template #trigger>
				<div style="width: 100%" @click="setEdit">
					<n-text class="name">
						{{ row.name }}
					</n-text>
				</div>
			</template>

			<n-input v-model:value="name" size="small" :maxlength="50" type="textarea" autosize />

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
.name {
	text-decoration: underline dotted;
}

.command-name {
	display: inline-flex;
	align-items: flex-start;
	gap: 4px;
	justify-content: space-between;
}

.actions {
	display: flex;
	gap: 2px;
}
</style>
