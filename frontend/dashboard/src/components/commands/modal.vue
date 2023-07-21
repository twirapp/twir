<script setup lang='ts'>
import { IconArrowNarrowDown, IconArrowNarrowUp, IconPlus, IconShieldHalfFilled, IconTrash } from '@tabler/icons-vue';
import {
    type FormInst,
    type FormItemRule,
    type FormRules,
    NAlert,
    NButton,
    NCard,
    NDivider,
    NDynamicInput,
    NDynamicTags,
    NForm,
    NFormItem,
    NGrid,
    NGridItem,
    NInput,
    NInputGroup,
    NInputGroupLabel,
    NInputNumber,
    NSelect,
    NSpace,
    NSwitch,
    NText,
} from 'naive-ui';
import { computed, onMounted, reactive, ref, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCommandsGroupsManager, useCommandsManager, useRolesManager } from '@/api/index.js';
import type { EditableCommand } from '@/components/commands/types.js';
import TextWithVariables from '@/components/textWithVariables.vue';
import TwitchUsersMultiple from '@/components/twitchUsers/multiple.vue';

const { t } = useI18n();

const props = defineProps<{
  command: EditableCommand | null
}>();

const emits = defineEmits<{
  close: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = reactive<EditableCommand>({
  name: '',
  aliases: [],
  responses: [],
  description: '',
  rolesIds: [],
  deniedUsersIds: [],
  allowedUsersIds: [],
  requiredMessages: 0,
  requiredUsedChannelPoints: 0,
  requiredWatchTime: 0,
  cooldown: 0,
  cooldownType: 'GLOBAL',
  isReply: true,
  visible: true,
  keepResponsesOrder: true,
  onlineOnly: false,
  enabled: true,
  groupId: undefined,
  module: 'CUSTOM',
});

onMounted(() => {
  if (props.command) {
    formValue.id = props.command.id;
    formValue.name = props.command.name;
    formValue.aliases = props.command.aliases;
    formValue.responses = props.command.responses;
    formValue.description = props.command.description;
    formValue.rolesIds = props.command.rolesIds;
    formValue.deniedUsersIds = props.command.deniedUsersIds;
    formValue.allowedUsersIds = props.command.allowedUsersIds;
    formValue.requiredMessages = props.command.requiredMessages;
    formValue.requiredUsedChannelPoints = props.command.requiredUsedChannelPoints;
    formValue.requiredWatchTime = props.command.requiredWatchTime;
    formValue.cooldown = props.command.cooldown;
    formValue.cooldownType = props.command.cooldownType;
    formValue.isReply = props.command.isReply;
    formValue.visible = props.command.visible;
    formValue.keepResponsesOrder = props.command.keepResponsesOrder;
    formValue.onlineOnly = props.command.onlineOnly;
    formValue.enabled = props.command.enabled;
    formValue.groupId = props.command.groupId;
    formValue.module = props.command.module;
  }
});

const rolesManager = useRolesManager();
const roles = rolesManager.getAll({});
const rolesSelectOptions = computed(() => {
  if (!roles.data?.value) return [];
  return roles.data.value.roles.map((role) => ({
    label: role.name,
    value: role.id,
  }));
});

const commandsGroupsManager = useCommandsGroupsManager();
const commandsGroups = commandsGroupsManager.getAll({});
const commandsGroupsOptions = computed(() => {
  if (!commandsGroups.data?.value) return [];
  return commandsGroups.data.value.groups.map((group) => ({
    label: group.name,
    value: group.id,
  }));
});

const nameValidator = (_: FormItemRule, value: string) => {
  if (!value) {
    return new Error(t('commands.modal.name.validations.empty'));
  }
  if (value.startsWith('!')) {
    return new Error(t('commands.modal.name.validations.startsWith'));
  }
  if (value.length > 25) {
    return new Error(t('commands.modal.name.validations.len'));
  }
  return true;
};
const rules: FormRules = {
  name: [{
    trigger: ['input', 'blur'],
    validator: nameValidator,
  }],
  description: {
    trigger: ['input', 'blur'],
    validator: (_: FormItemRule, value: string) => {
      if (value.length > 500) {
        return new Error('Description cannot be longer than 500 characters');
      }
      return true;
    },
  },
  responses: {
    trigger: ['input', 'blur', 'focus'],
    validator: (_: FormItemRule, value: string) => {
      if (value.length === 0) {
        return new Error(t('commands.modal.responses.validations.empty'));
      }
      if (value.length > 500) {
        return new Error(t('commands.modal.responses.validations.len'));
      }
      return true;
    },
  },
};

const commandsManager = useCommandsManager();
const commandsCreate = commandsManager.create;
const commandsUpdate = commandsManager.update;

async function save() {
  await formRef.value?.validate();

  const rawData = toRaw(formValue);
  const data = {
    ...rawData,
    responses: rawData.responses.map((r, i) => ({
      ...r,
      order: i,
    })),
    groupId: rawData.groupId === null ? undefined : rawData.groupId,
  };

  if (rawData.id) {
    await commandsUpdate.mutateAsync({
      id: rawData.id,
      command: data,
    });
  } else {
    await commandsCreate.mutateAsync(data);
  }

  emits('close');
}
</script>

<template>
	<n-form ref="formRef" :model="formValue" :rules="rules">
		<n-grid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="10">
			<n-grid-item :span="1">
				<n-form-item :label="t('commands.modal.name.label')" path="name" show-require-mark>
					<n-input-group>
						<n-input-group-label>!</n-input-group-label>
						<n-input v-model:value="formValue.name" placeholder="Input Name" />
					</n-input-group>
				</n-form-item>
			</n-grid-item>
			<n-grid-item :span="1">
				<n-form-item :label="t('commands.modal.aliases.label')" path="aliases">
					<n-dynamic-tags v-model:value="formValue.aliases" />
				</n-form-item>
			</n-grid-item>
		</n-grid>
		<n-form-item :label="t('commands.modal.description.label')" path="description">
			<n-input v-model:value="formValue.description" placeholder="Description" type="textarea" autosize />
		</n-form-item>

		<n-divider>
			{{ t('sharedTexts.responses') }}
		</n-divider>

		<div v-if="formValue.module === 'CUSTOM'">
			<n-dynamic-input
				v-if="formValue.module === 'CUSTOM'"
				v-model:value="formValue.responses"
				class="groups"
				placeholder="text"
				:create-button-props="{ class: 'create-button' } as any"
			>
				<template #default="{ value, index }">
					<n-form-item
						style="width: 100%"
						:path="`responses[${index}].text`"
						:rule="rules.responses"
					>
						<text-with-variables
							v-model="value.text"
							inputType="textarea"
							:minRows="3"
							:maxRows="6"
						>
						</text-with-variables>
					</n-form-item>
				</template>

				<template #action="{ index, remove, move }">
					<div class="group-actions">
						<n-button size="small" type="error" quaternary @click="() => remove(index)">
							<IconTrash />
						</n-button>
						<n-button
							size="small"
							type="info"
							quaternary
							:disabled="index == 0"
							@click="() => move('up', index)"
						>
							<IconArrowNarrowUp />
						</n-button>
						<n-button
							size="small"
							type="info"
							quaternary
							:disabled="!!formValue.responses.length && index === formValue.responses.length-1"
							@click="() => move('down', index)"
						>
							<IconArrowNarrowDown />
						</n-button>
					</div>
				</template>
			</n-dynamic-input>
			<n-button dashed block style="margin-top:10px" @click="() => formValue.responses.push({ text: '' })">
				<IconPlus />
				{{ t('sharedButtons.create') }}
			</n-button>
		</div>

		<n-alert v-else type="info" :show-icon="false">
			{{ t('commands.modal.responses.defaultWarning') }}
		</n-alert>

		<n-divider>
			{{ t('commands.modal.permissions.divider') }}
		</n-divider>

		<n-form-item :label="t('commands.modal.permissions.name')" path="rolesIds">
			<n-select
				v-model:value="formValue.rolesIds"
				multiple
				:options="rolesSelectOptions"
				:loading="roles.isLoading.value"
				:placeholder="t('commands.modal.permissions.placeholder')"
			>
				<template #arrow>
					<IconShieldHalfFilled />
				</template>
			</n-select>
		</n-form-item>

		<n-grid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5">
			<n-grid-item :span="1">
				<n-form-item :label="t('commands.modal.permissions.deniedUsers')" path="deniedUsersIds">
					<twitch-users-multiple
						v-model="formValue.deniedUsersIds"
						:initial-users-ids="formValue.deniedUsersIds"
					/>
				</n-form-item>
			</n-grid-item>

			<n-grid-item :span="1">
				<n-form-item :label="t('commands.modal.permissions.allowedUsers')" path="allowedUsersIds">
					<twitch-users-multiple
						v-model="formValue.allowedUsersIds"
						:initial-users-ids="formValue.allowedUsersIds"
					/>
				</n-form-item>
			</n-grid-item>
		</n-grid>

		<n-divider>
			{{ t('commands.modal.restrictions.name') }}
		</n-divider>

		<n-grid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5">
			<n-grid-item :span="1">
				<n-form-item
					:label="t('commands.modal.restrictions.watchTime')"
					path="requiredWatchTime"
				>
					<n-input-number
						v-model:value="formValue.requiredWatchTime"
						:min="0"
						class="grid-stats-item"
					/>
				</n-form-item>
			</n-grid-item>

			<n-grid-item :span="1">
				<n-form-item
					:label="t('commands.modal.restrictions.messages')"
					path="requiredMessages"
				>
					<n-input-number
						v-model:value="formValue.requiredMessages"
						:min="0"
						class="grid-stats-item"
					/>
				</n-form-item>
			</n-grid-item>

			<n-grid-item :span="1">
				<n-form-item
					:label="t('commands.modal.restrictions.channelsPoints')"
					path="requiredUsedChannelPoints"
				>
					<n-input-number
						v-model:value="formValue.requiredUsedChannelPoints"
						:min="0"
						class="grid-stats-item"
					/>
				</n-form-item>
			</n-grid-item>
		</n-grid>

		<n-divider>
			{{ t('commands.modal.cooldown.label') }}
		</n-divider>

		<n-grid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5">
			<n-grid-item :span="1">
				<n-form-item
					:label="t('commands.modal.cooldown.value')"
					path="cooldown"
				>
					<n-input-number
						v-model:value="formValue.cooldown"
						:min="0"
						class="grid-stats-item"
					/>
				</n-form-item>
			</n-grid-item>

			<n-grid-item :span="1">
				<n-form-item
					:label="t('commands.modal.cooldown.type.name')"
					path="cooldownType"
				>
					<n-select
						v-model:value="formValue.cooldownType"
						:options="[
							{
								label: t('commands.modal.cooldown.type.global'),
								value: 'GLOBAL',
							},
							{
								label: t('commands.modal.cooldown.type.user'),
								value: 'PER_USER'
							},
						]"
					/>
				</n-form-item>
			</n-grid-item>
		</n-grid>

		<n-divider>
			{{ t('commands.modal.settings.divider') }}
		</n-divider>

		<n-grid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5" :y-gap="5">
			<n-grid-item :span="1">
				<n-card>
					<div class="settings-card-body">
						<n-space vertical>
							<n-text>{{ t("commands.modal.settings.reply.label") }}</n-text>
							<n-text>{{ t("commands.modal.settings.reply.text") }}</n-text>
						</n-space>
						<n-switch v-model:value="formValue.isReply" />
					</div>
				</n-card>
			</n-grid-item>

			<n-grid-item :span="1">
				<n-card>
					<div class="settings-card-body">
						<n-space vertical>
							<n-text>{{ t('commands.modal.settings.visible.label') }}</n-text>
							<n-text>{{ t('commands.modal.settings.reply.text') }}</n-text>
						</n-space>
						<n-switch v-model:value="formValue.visible" />
					</div>
				</n-card>
			</n-grid-item>

			<n-grid-item :span="1">
				<n-card>
					<div class="settings-card-body">
						<n-space vertical>
							<n-text>{{ t('commands.modal.settings.keepOrder.label') }}</n-text>
							<n-text>{{ t('commands.modal.settings.keepOrder.text') }}</n-text>
						</n-space>
						<n-switch v-model:value="formValue.keepResponsesOrder" />
					</div>
				</n-card>
			</n-grid-item>

			<n-grid-item :span="1">
				<n-card>
					<div class="settings-card-body">
						<n-space vertical>
							<n-text>{{ t('commands.modal.settings.onlineOnly.label') }}</n-text>
							<n-text>{{ t('commands.modal.settings.onlineOnly.text') }}</n-text>
						</n-space>
						<n-switch v-model:value="formValue.onlineOnly" />
					</div>
				</n-card>
			</n-grid-item>
		</n-grid>

		<n-divider>
			{{ t('commands.modal.settings.other.divider') }}
		</n-divider>

		<n-form-item :label="t('commands.modal.settings.other.commandGroup')" path="groupId">
			<n-select
				v-model:value="formValue.groupId"
				:options="commandsGroupsOptions"
				clearable
				:fallback-option="undefined"
			/>
		</n-form-item>

		<n-button secondary type="success" block @click="save">
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-form>
</template>

<style scoped>
.groups :deep(.create-button) {
  display: none;
}

.group-actions {
  display: flex;
  margin-left: 5px;
  column-gap: 5px;
  align-items: center
}

.grid-stats-item {
  width: 100%;
}

.settings-card-body {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}
</style>
