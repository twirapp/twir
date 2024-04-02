<script setup lang="ts">
import { IconPlus, IconTrash, IconDeviceFloppy } from '@tabler/icons-vue';
import type { Group } from '@twir/api/messages/commands_group/commands_group';
import {
	NDynamicInput,
	NInput,
	NColorPicker,
	NFormItem,
	NGrid,
	NGridItem,
	NButton,
} from 'naive-ui';
import { toRaw, ref, watch, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCommandsGroupsManager } from '@/api/index.js';

const { t } = useI18n();

const groupsManager = useCommandsGroupsManager();
const groupsData = groupsManager.getAll({});
const groupsCreator = groupsManager.create;
const groupsDeleter = groupsManager.deleteOne;
const groupsUpdater = groupsManager.update;

type FormGroup = Omit<Group, 'id' | 'channelId'> & { id?: string }

const groups = ref<FormGroup[]>([]);

onMounted(() => {
	groupsData.refetch();
});

watch(groupsData.data, (data) => {
	groups.value = data?.groups ? toRaw(data.groups) : [];
}, { immediate: true });

async function create(name: string, color: string) {
	await groupsCreator.mutateAsync({ color, name });
}

async function deleteGroup(index: number) {
	const group = groups.value[index];
	if (!group?.id) return;
	await groupsDeleter.mutateAsync({ id: group.id });
}

async function update(index: number) {
	const group = groups.value[index];
	if (!group?.id) return;
	await groupsUpdater.mutateAsync({ id: group.id, name: group.name, color: group.color });
}

const swatches = [
	'rgba(116, 242, 202, 1)',
	'rgba(208, 48, 80, 1)',
];
</script>

<template>
	<n-dynamic-input
		v-model:value="groups"
		:on-remove="(a) => deleteGroup(a)"
		class="groups w-full"
		:create-button-props="{ class: 'create-button' } as any"
	>
		<template #default="{ value }: { value: FormGroup }">
			<n-grid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5">
				<n-grid-item :span="1">
					<n-form-item :label="t('commands.groups.name')">
						<n-input v-model:value="value.name" type="text" />
					</n-form-item>
				</n-grid-item>
				<n-grid-item :span="1">
					<n-form-item :label="t('commands.groups.color')">
						<n-color-picker
							v-model:value="value.color"
							:show-alpha="true"
							:swatches="swatches"
							:modes="['rgb']"
						/>
					</n-form-item>
				</n-grid-item>
			</n-grid>
		</template>

		<template #action="{ index, remove }">
			<div class="group-actions">
				<n-button size="small" type="success" quaternary @click="() => update(index)">
					<IconDeviceFloppy />
				</n-button>
				<n-button size="small" type="error" quaternary @click="() => remove(index)">
					<IconTrash />
				</n-button>
			</div>
		</template>
	</n-dynamic-input>
	<n-button dashed block @click="() => create(`New Group #${groups.length + 1}`, swatches[0])">
		<IconPlus />
		{{ t('sharedButtons.create') }}
	</n-button>
</template>

<style scoped>
.groups :deep(.create-button) {
	@apply hidden;
}

.group-actions {
	@apply flex gap-x-1 items-center;
}
</style>
