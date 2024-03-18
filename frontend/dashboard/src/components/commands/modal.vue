<script setup lang='ts'>
import {
	IconArrowNarrowDown,
	IconArrowNarrowUp,
	IconPlus,
	IconSquareCheck,
	IconSquare,
	IconTrash,
} from '@tabler/icons-vue';
import chunk from 'lodash.chunk';
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
	NButtonGroup,
	NTabs,
	NTabPane,
} from 'naive-ui';
import { computed, onMounted, ref, toRaw } from 'vue';
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
const formValue = ref<EditableCommand>({
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
	cooldownRolesIds: [],
});

onMounted(() => {
	if (!props.command) return;

	formValue.value = toRaw(props.command);
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

	const rawData = toRaw(formValue.value);
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

const createButtonProps = { class: 'create-button' } as any;
</script>

<template>
	<div class="flex flex-col justify-between h-full">
		<n-form ref="formRef" :model="formValue" :rules="rules" class="flex flex-col h-[95%] flex-grow">
			<n-tabs class="h-full" type="line" animated placement="left" default-value="general">
				<n-tab-pane name="general" tab="General">
					<div style="display: flex; gap: 8px">
						<n-form-item :label="t('commands.modal.name.label')" path="name" show-require-mark style="width: 90%">
							<n-input-group>
								<n-input-group-label>!</n-input-group-label>
								<n-input v-model:value="formValue.name" placeholder="Name of command" :maxlength="25" :on-input="() => formValue.name.startsWith('!') && (formValue.name = formValue.name.slice(1))" />
							</n-input-group>
						</n-form-item>
						<n-form-item :label="t('sharedTexts.enabled')" path="enabled">
							<n-switch v-model:value="formValue.enabled" />
						</n-form-item>
					</div>

					<n-form-item :label="t('commands.modal.aliases.label')" path="aliases">
						<n-dynamic-tags v-model:value="formValue.aliases" :input-props="{ maxlength: 25 }" />
					</n-form-item>

					<n-form-item :label="t('commands.modal.description.label')" path="description">
						<n-input
							v-model:value="formValue.description" placeholder="Description" type="textarea"
							autosize
						/>
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
							:create-button-props="createButtonProps"
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
						<n-button
							dashed
							block
							style="margin-top:10px"
							:disabled="formValue.responses.length >= 3"
							@click="() => formValue.responses.push({ text: '' })"
						>
							<IconPlus />
							{{ t('sharedButtons.create') }}
						</n-button>
					</div>

					<n-alert v-else type="info" :show-icon="false">
						{{ t('commands.modal.responses.defaultWarning') }}
					</n-alert>
				</n-tab-pane>
				<n-tab-pane name="permissions" :tab="t('commands.modal.permissions.divider')">
					<n-form-item :label="t('commands.modal.permissions.name')" path="rolesIds">
						<div style="display: flex; flex-direction: column; gap: 5px;">
							<n-button-group
								v-for="(group, index) of chunk(rolesSelectOptions.sort(), 5)"
								:key="index"
							>
								<n-button
									v-for="option of group"
									:key="option.value"
									:type="formValue.rolesIds.includes(option.value) ? 'success' : 'default'"
									secondary
									@click="() => {
										if (formValue.rolesIds.includes(option.value)) {
											formValue.rolesIds = formValue.rolesIds.filter(r => r !== option.value)
										} else {
											formValue.rolesIds.push(option.value)
										}
									}"
								>
									<template #icon>
										<IconSquareCheck v-if="formValue.rolesIds.includes(option.value)" />
										<IconSquare v-else />
									</template>
									{{ option.label }}
								</n-button>
							</n-button-group>
						</div>
					</n-form-item>


					<n-form-item :label="t('commands.modal.permissions.deniedUsers')" path="deniedUsersIds">
						<twitch-users-multiple
							v-model="formValue.deniedUsersIds"
							:initial-users-ids="formValue.deniedUsersIds"
						/>
					</n-form-item>

					<n-form-item :label="t('commands.modal.permissions.allowedUsers')" path="allowedUsersIds">
						<twitch-users-multiple
							v-model="formValue.allowedUsersIds"
							:initial-users-ids="formValue.allowedUsersIds"
						/>
					</n-form-item>

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
				</n-tab-pane>
				<n-tab-pane name="cooldown" :tab="t('commands.modal.cooldown.label')">
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

					<div style="display: flex; flex-direction: column; gap: 5px;">
						<n-button-group
							v-for="(group, index) of chunk(rolesSelectOptions.sort(), 5)"
							:key="index"
						>
							<n-button
								v-for="option of group"
								:key="option.value"
								:type="formValue.cooldownRolesIds.includes(option.value) ? 'success' : 'default'"
								secondary
								@click="() => {
									if (formValue.cooldownRolesIds.includes(option.value)) {
										formValue.cooldownRolesIds = formValue.cooldownRolesIds.filter(r => r !== option.value)
									} else {
										formValue.cooldownRolesIds.push(option.value)
									}
								}"
							>
								<template #icon>
									<IconSquareCheck v-if="formValue.cooldownRolesIds.includes(option.value)" />
									<IconSquare v-else />
								</template>
								{{ option.label }}
							</n-button>
						</n-button-group>
					</div>
				</n-tab-pane>
				<n-tab-pane name="settings" :tab="t('commands.modal.settings.divider')">
					<n-grid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5" :y-gap="5">
						<n-grid-item :span="1">
							<n-card style="height: 100%">
								<div class="settings-card-body">
									<n-space vertical>
										<n-text>{{ t("sharedTexts.reply.label") }}</n-text>
										<n-text>{{ t("sharedTexts.reply.text") }}</n-text>
									</n-space>
									<n-switch v-model:value="formValue.isReply" />
								</div>
							</n-card>
						</n-grid-item>

						<n-grid-item :span="1">
							<n-card style="height: 100%">
								<div class="settings-card-body">
									<n-space vertical>
										<n-text>{{ t('commands.modal.settings.visible.label') }}</n-text>
										<n-text>{{ t('commands.modal.settings.visible.text') }}</n-text>
									</n-space>
									<n-switch v-model:value="formValue.visible" />
								</div>
							</n-card>
						</n-grid-item>

						<n-grid-item :span="1">
							<n-card style="height: 100%">
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
							<n-card style="height: 100%">
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
						<n-button v-if="!commandsGroupsOptions.length" secondary disabled>
							No groups created
						</n-button>
						<div
							v-else
							style="display: flex; flex-direction: column; gap: 5px;"
						>
							<n-button-group
								v-for="(group, index) of chunk(commandsGroupsOptions.sort(), 4)"
								:key="index"
							>
								<n-button
									v-for="option of group"
									:key="option.value"
									:type="formValue.groupId === option.value ? 'success' : 'default'"
									secondary
									@click="() => {
										if (formValue.groupId === option.value) {
											formValue.groupId = undefined
										} else {
											formValue.groupId = option.value
										}
									}"
								>
									<template #icon>
										<IconSquareCheck v-if="formValue.groupId === option.value" />
										<IconSquare v-else />
									</template>
									{{ option.label }}
								</n-button>
							</n-button-group>
						</div>
					</n-form-item>
				</n-tab-pane>
			</n-tabs>
		</n-form>


		<n-button style="margin-top: 8px" secondary type="success" block @click="save">
			{{ t('sharedButtons.save') }}
		</n-button>
	</div>
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
