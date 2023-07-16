<script setup lang='ts'>
import { IconPlus, IconTrash, IconDeviceFloppy } from '@tabler/icons-vue';
import type { Group } from '@twir/grpc/generated/api/api/commands_group';
import { NDynamicInput, NInput, NColorPicker, NFormItem, NGrid, NGridItem, NButton } from 'naive-ui';
import { toRaw, ref, watch, onMounted } from 'vue';

import { useCommandsGroupsManager } from '@/api/index.js';

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
    style="width: 100%"
    class="groups"
    :create-button-props="{ class: 'create-button' } as any"
  >
    <template #default="{ value }: { value: FormGroup }">
      <n-grid :cols="12" :x-gap="5">
        <n-grid-item :span="6">
          <n-form-item label="Name">
            <n-input v-model:value="value.name" type="text" />
          </n-form-item>
        </n-grid-item>
        <n-grid-item :span="5">
          <n-form-item label="Color">
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
    <IconPlus /> Create
  </n-button>
</template>

<style scoped>
.groups :deep(.create-button) {
	display: none;
}

.group-actions {
	display: flex;
	column-gap: 5px;
	align-items: center
}
</style>
