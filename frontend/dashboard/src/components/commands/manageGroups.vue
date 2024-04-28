<script setup lang="ts">
import { IconPlus, IconTrash, IconDeviceFloppy } from '@tabler/icons-vue';
import {
	NDynamicInput,
	NInput,
	NColorPicker,
	NFormItem,
	NGrid,
	NGridItem,
	NButton,
} from 'naive-ui';
import { toRaw, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCommandsGroupsApi } from '@/api/commands/commands-groups';
import type { CommandGroup } from '@/gql/graphql';

const { t } = useI18n();

const groupsManager = useCommandsGroupsApi();
const { data } = groupsManager.useQueryGroups();
const groupsCreator = groupsManager.useMutationCreateGroup();
const groupsDeleter = groupsManager.useMutationDeleteGroup();
const groupsUpdater = groupsManager.useMutationUpdateGroup();

type FormGroup = Omit<CommandGroup, 'id'> & { id?: string }

const groups = ref<FormGroup[]>([]);

watch(data, (data) => {
	groups.value = data?.commandsGroups ? toRaw(data?.commandsGroups) : [];
}, { immediate: true });

async function create(name: string, color: string) {
	await groupsCreator.executeMutation({ opts: { color, name } });
}

async function deleteGroup(index: number) {
	const group = groups.value[index];
	if (!group?.id) return;
	await groupsDeleter.executeMutation({ id: group.id });
}

async function update(index: number) {
	const group = groups.value[index];
	if (!group?.id) return;
	await groupsUpdater.executeMutation({ id: group.id, opts: { name: group.name, color: group.color } });
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
