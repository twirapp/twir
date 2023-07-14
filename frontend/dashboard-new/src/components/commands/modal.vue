<script setup lang='ts'>
import { IconShieldHalfFilled } from '@tabler/icons-vue';
import type { Command, Command_Response } from '@twir/grpc/generated/api/api/commands';
import {
	NForm,
	NFormItem,
	FormInst,
	FormRules,
	FormItemRule,
	NInput,
	NDynamicTags,
	NGrid,
	NGridItem,
	NDivider,
	NInputGroup,
	NInputGroupLabel,
	NDynamicInput,
	NSelect,
	NInputNumber,
	NCard,
	NSwitch,
	NSpace,
	NAlert,
	NText,
} from 'naive-ui';
import { ref, reactive, onMounted, computed } from 'vue';

import { useRolesManager, useCommandsGroupsManager } from '@/api/index.js';
import TextWithVariables from '@/components/textWithVariables.vue';
import TwitchUsersMultiple from '@/components/twitchUsers/multiple.vue';

const props = defineProps<{
	command: Command | null
}>();

type FormCommand = Omit<
	Command,
	'responses' |
	'channelId' |
	'default' |
	'defaultName' |
	'id' |
	'group'
> & {
	responses: Array<Omit<Command_Response, 'id' | 'commandId' | 'order'>>,
	id?: string
};

const formRef = ref<FormInst | null>(null);
const formValue = reactive<FormCommand>({
	name:'',
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

const nameValidator = (rule: FormItemRule, value: string) => {
	if (!value) {
		return new Error('Please input a name');
	}
	if (value.startsWith('!')) {
		return new Error('Name cannot start with !');
	}
	return true;
};
const rules: FormRules = {
	name: [{
		trigger: ['input', 'blur'],
		validator: nameValidator,
	}],
	aliases: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: string[]) => {
			value.forEach((alias) => nameValidator(rule, alias));

			return true;
		},
	},
	cooldown: {
		trigger: ['input', 'blur'],
		validator: (rule: FormItemRule, value: number) => {
			if (value < 0) {
				return new Error('Cooldown cannot be negative');
			}
			return true;
		},
	},
};
</script>

<template>
  <n-form ref="formRef" :model="formValue" :rules="rules">
    <n-grid :cols="12" :x-gap="10">
      <n-grid-item :span="6">
        <n-form-item label="Name" path="name">
          <n-input-group>
            <n-input-group-label>!</n-input-group-label>
            <n-input v-model:value="formValue.name" placeholder="Input Name" />
          </n-input-group>
        </n-form-item>
      </n-grid-item>
      <n-grid-item :span="6">
        <n-form-item label="Aliases" path="aliases">
          <n-dynamic-tags v-model:value="formValue.aliases" />
        </n-form-item>
      </n-grid-item>
    </n-grid>
    <n-form-item label="Description" path="description">
      <n-input v-model:value="formValue.description" placeholder="Description" type="textarea" autosize />
    </n-form-item>

    <n-divider>
      Responses
    </n-divider>

    <n-dynamic-input
      v-if="formValue.module === 'CUSTOM'"
      v-model:value="formValue.responses"
      :on-create="() => ({ text: '' })"
      placeholder="Response"
      show-sort-button
    >
      <template #default="{ value }: { value: Command_Response }">
        <text-with-variables
          v-model="value.text"
          inputType="textarea"
          :minRows="3"
          :maxRows="6"
        >
        </text-with-variables>
      </template>
    </n-dynamic-input>

    <n-alert v-else type="info" :show-icon="false">
      Responses cannot be edited for built-in commands
    </n-alert>

    <n-divider>
      Permissions
    </n-divider>

    <n-form-item label="Roles" path="rolesIds">
      <n-select
        v-model:value="formValue.rolesIds"
        multiple
        :options="rolesSelectOptions"
        :loading="roles.isLoading.value"
      >
        <template #arrow>
          <IconShieldHalfFilled />
        </template>
      </n-select>
    </n-form-item>

    <n-grid :cols="12" :x-gap="5">
      <n-grid-item :span="6">
        <n-form-item label="Denied users" path="deniedUsersIds">
          <twitch-users-multiple
            v-model="formValue.deniedUsersIds"
            :initial-users-ids="formValue.deniedUsersIds"
          />
        </n-form-item>
      </n-grid-item>

      <n-grid-item :span="6">
        <n-form-item label="Allowed users" path="allowedUsersIds">
          <twitch-users-multiple
            v-model="formValue.allowedUsersIds"
            :initial-users-ids="formValue.allowedUsersIds"
          />
        </n-form-item>
      </n-grid-item>
    </n-grid>

    <n-divider>
      Restrictions by stats
    </n-divider>

    <n-grid :cols="12" :x-gap="5">
      <n-grid-item :span="6">
        <n-form-item label="Required watch time (hours)" path="requiredWatchTime">
          <n-input-number
            v-model:value="formValue.requiredWatchTime"
            :min="0"
            class="grid-stats-item"
          />
        </n-form-item>
      </n-grid-item>

      <n-grid-item :span="6">
        <n-form-item label="Required messages" path="requiredMessages">
          <n-input-number
            v-model:value="formValue.requiredMessages"
            :min="0"
            class="grid-stats-item"
          />
        </n-form-item>
      </n-grid-item>

      <n-grid-item :span="6">
        <n-form-item label="Required used channels points" path="requiredUsedChannelPoints">
          <n-input-number
            v-model:value="formValue.requiredUsedChannelPoints"
            :min="0"
            class="grid-stats-item"
          />
        </n-form-item>
      </n-grid-item>
    </n-grid>

    <n-divider>
      Cooldown
    </n-divider>

    <n-grid :cols="12" :x-gap="5">
      <n-grid-item :span="6">
        <n-form-item label="Cooldown (seconds)" path="cooldown">
          <n-input-number
            v-model:value="formValue.cooldown"
            :min="0"
            class="grid-stats-item"
          />
        </n-form-item>
      </n-grid-item>

      <n-grid-item :span="6">
        <n-form-item label="Cooldown type" path="cooldownType">
          <n-select
            v-model:value="formValue.cooldownType"
            :options="[
              { label: 'On command', value: 'GLOBAL' },
              { label: 'On user (not working on mods, subscribers)', value: 'PER_USER' },
            ]"
          />
        </n-form-item>
      </n-grid-item>
    </n-grid>

    <n-divider>
      Settings
    </n-divider>

    <n-grid :cols="12" :x-gap="5" :y-gap="5" item-responsive>
      <n-grid-item :span="6">
        <n-card>
          <div class="settings-card-body">
            <n-space vertical>
              <n-text>Reply</n-text>
              <n-text>Bot will send command response as reply</n-text>
            </n-space>
            <n-switch v-model:value="formValue.isReply" />
          </div>
        </n-card>
      </n-grid-item>

      <n-grid-item :span="6">
        <n-card>
          <div class="settings-card-body">
            <n-space vertical>
              <n-text>Visible</n-text>
              <n-text>Is command visible in commands list on public page and in chat commands variable</n-text>
            </n-space>
            <n-switch v-model:value="formValue.visible" />
          </div>
        </n-card>
      </n-grid-item>

      <n-grid-item :span="6">
        <n-card>
          <div class="settings-card-body">
            <n-space vertical>
              <n-text>Keep order</n-text>
              <n-text>Keep order of responses when sending them in chat</n-text>
            </n-space>
            <n-switch v-model:value="formValue.keepResponsesOrder" />
          </div>
        </n-card>
      </n-grid-item>

      <n-grid-item :span="6">
        <n-card>
          <div class="settings-card-body">
            <n-space vertical>
              <n-text>Online only</n-text>
              <n-text>Command will work only when stream online</n-text>
            </n-space>
            <n-switch v-model:value="formValue.onlineOnly" />
          </div>
        </n-card>
      </n-grid-item>
    </n-grid>

    <n-divider>
      Other
    </n-divider>

    <n-form-item label="Command group" path="groupId">
      <n-select
        v-model:value="formValue.groupId"
        :options="commandsGroupsOptions"
      />
    </n-form-item>
  </n-form>
</template>

<style scoped>
.grid-stats-item {
	width: 100%;
}

.settings-card-body {
	display: flex;
	flex-direction: row;
	justify-content: space-between;
}
</style>
