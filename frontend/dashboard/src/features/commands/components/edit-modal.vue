<script setup lang='ts'>
import {
	IconArrowNarrowDown,
	IconArrowNarrowUp,
	IconPlus,
	IconSquare,
	IconSquareCheck,
	IconTrash
} from '@tabler/icons-vue'
import chunk from 'lodash.chunk'
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NAlert,
	NButton,
	NButtonGroup,
	NCard,
	NDivider,
	NDynamicInput,
	NForm,
	NFormItem,
	NGrid,
	NGridItem,
	NInput,
	NInputGroup,
	NInputGroupLabel,
	NInputNumber,
	NModal,
	NSelect,
	NSpace,
	NSwitch,
	NTabPane,
	NTabs,
	NTag,
	NText
} from 'naive-ui'
import { storeToRefs } from 'pinia'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useCommandEdit } from '../composables/use-command-edit.js'

import { useCommandsGroupsApi } from '@/api/commands/commands-groups'
import { useRoles } from '@/api/roles'
import TwitchCategorySearch from '@/components/twitch-category-search.vue'
import TwitchUsersMultiple from '@/components/twitchUsers/multiple.vue'
import { TagsInput, TagsInputInput, TagsInputItem, TagsInputItemDelete, TagsInputItemText } from '@/components/ui/tags-input'
import VariableInput from '@/components/variable-input.vue'

const { t } = useI18n()
const commandEdit = useCommandEdit()
const { formValue, isOpened: isEditOpened } = storeToRefs(commandEdit)

const formRef = ref<FormInst | null>(null)

const rolesManager = useRoles()
const { data: roles } = rolesManager.useRolesQuery()

const rolesSelectOptions = computed(() => {
	if (!roles.value?.roles) return []
	return roles.value.roles.map((role) => ({
		label: role.name,
		value: role.id
	}))
})

const commandsGroupsManager = useCommandsGroupsApi()
const commandsGroups = commandsGroupsManager.useQueryGroups()
const commandsGroupsOptions = computed(() => {
	if (!commandsGroups.data?.value) return []
	return commandsGroups.data.value.commandsGroups.map((group) => ({
		label: group.name,
		value: group.id
	}))
})

function nameValidator(_: FormItemRule, value: string) {
	if (!value) {
		return new Error(t('commands.modal.name.validations.empty'))
	}
	if (value.startsWith('!')) {
		return new Error(t('commands.modal.name.validations.startsWith'))
	}
	if (value.length > 25) {
		return new Error(t('commands.modal.name.validations.len'))
	}
	return true
}
const rules: FormRules = {
	name: [{
		trigger: ['input', 'blur'],
		validator: nameValidator
	}],
	description: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (value.length > 500) {
				return new Error('Description cannot be longer than 500 characters')
			}
			return true
		}
	},
	responses: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (value.length === 0) {
				return new Error(t('commands.modal.responses.validations.empty'))
			}
			if (value.length > 500) {
				return new Error(t('commands.modal.responses.validations.len'))
			}
			return true
		}
	}
}

const createButtonProps = { class: 'create-button' } as any
</script>

<template>
	<NModal
		:show="isEditOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="formValue?.name ?? t('commands.newCommandTitle')"
		class="modal"
		:style="{
			width: '800px',
			height: '90dvh',
		}"
		:on-close="commandEdit.close"
		content-style="padding: 5px;"
	>
		<div v-if="formValue" class="flex flex-col justify-between h-full">
			<NForm ref="formRef" :model="formValue" :rules="rules" class="flex flex-col h-[95%] flex-grow">
				<NTabs class="h-full" type="line" animated placement="left" default-value="general">
					<NTabPane name="general" tab="General">
						<div class="flex gap-2">
							<NFormItem :label="t('commands.modal.name.label')" path="name" show-require-mark class="w-[90%]">
								<NInputGroup>
									<NInputGroupLabel>!</NInputGroupLabel>
									<NInput v-model:value="formValue.name" placeholder="Name of command" :maxlength="25" :on-input="() => formValue!.name.startsWith('!') && (formValue!.name = formValue!.name.slice(1))" />
								</NInputGroup>
							</NFormItem>
							<NFormItem :label="t('sharedTexts.enabled')" path="enabled">
								<NSwitch v-model:value="formValue.enabled" />
							</NFormItem>
						</div>

						<NFormItem :label="t('commands.modal.aliases.label')" path="aliases">
							<TagsInput v-model="formValue.aliases" :max="25" class="bg-zinc-700 w-full">
								<TagsInputItem v-for="item in formValue.aliases" :key="item" :value="item">
									<TagsInputItemText />
									<TagsInputItemDelete />
								</TagsInputItem>

								<TagsInputInput placeholder="Write new aliase here..." />
							</TagsInput>
						</NFormItem>

						<NDivider>
							{{ t('sharedTexts.responses') }}
						</NDivider>

						<NText>
							<i18n-t keypath="commands.modal.responses.description">
								<NTag>$(sender)</NTag>
							</i18n-t>
						</NText>

						<div v-if="formValue.module === 'CUSTOM'">
							<NDynamicInput
								v-if="formValue.module === 'CUSTOM'"
								v-model:value="formValue.responses"
								class="groups"
								placeholder="text"
								:create-button-props="createButtonProps"
							>
								<template #default="{ value, index }">
									<NFormItem
										class="w-full"
										:path="`responses[${index}].text`"
										:rule="rules.responses"
									>
										<VariableInput
											v-model="value.text"
											inputType="textarea"
											:minRows="3"
											:maxRows="6"
										/>
									</NFormItem>
								</template>

								<template #action="{ index, remove, move }">
									<div class="flex items-center ml-1 gap-x-1">
										<NButton size="small" type="error" quaternary @click="() => remove(index)">
											<IconTrash />
										</NButton>
										<NButton
											size="small"
											type="info"
											quaternary
											:disabled="index === 0"
											@click="() => move('up', index)"
										>
											<IconArrowNarrowUp />
										</NButton>
										<NButton
											size="small"
											type="info"
											quaternary
											:disabled="!!formValue.responses.length && index === formValue.responses.length - 1"
											@click="() => move('down', index)"
										>
											<IconArrowNarrowDown />
										</NButton>
									</div>
								</template>
							</NDynamicInput>
							<NButton
								dashed
								block
								style="margin-top:10px"
								:disabled="formValue.responses.length >= 3"
								@click="() => formValue!.responses.push({ text: '', order: formValue!.responses.length })"
							>
								<IconPlus />
								{{ t('commands.modal.responses.add') }}
							</NButton>
						</div>

						<NAlert v-else type="info" :show-icon="false">
							{{ t('commands.modal.responses.defaultWarning') }}
						</NAlert>
					</NTabPane>
					<NTabPane name="permissions" :tab="t('commands.modal.permissions.divider')">
						<NFormItem :label="t('commands.modal.permissions.name')" path="rolesIds">
							<div class="flex gap-1 flex-col">
								<NButtonGroup
									v-for="(group, index) of chunk(rolesSelectOptions.sort(), 5)"
									:key="index"
								>
									<NButton
										v-for="option of group"
										:key="option.value"
										:type="formValue.rolesIds.includes(option.value) ? 'success' : 'default'"
										secondary
										@click="() => {
											if (formValue!.rolesIds.includes(option.value)) {
												formValue!.rolesIds = formValue!.rolesIds.filter(r => r !== option.value)
											}
											else {
												formValue!.rolesIds.push(option.value)
											}
										}"
									>
										<template #icon>
											<IconSquareCheck v-if="formValue.rolesIds.includes(option.value)" />
											<IconSquare v-else />
										</template>
										{{ option.label }}
									</NButton>
								</NButtonGroup>
							</div>
						</NFormItem>

						<NFormItem :label="t('commands.modal.permissions.deniedUsers')" path="deniedUsersIds">
							<TwitchUsersMultiple
								v-model="formValue.deniedUsersIds"
								:initial-users-ids="formValue.deniedUsersIds"
							/>
						</NFormItem>

						<NFormItem :label="t('commands.modal.permissions.allowedUsers')" path="allowedUsersIds">
							<TwitchUsersMultiple
								v-model="formValue.allowedUsersIds"
								:initial-users-ids="formValue.allowedUsersIds"
							/>
						</NFormItem>

						<NDivider>
							{{ t('commands.modal.restrictions.name') }}
						</NDivider>

						<NGrid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5">
							<NGridItem :span="1">
								<NFormItem
									:label="t('commands.modal.restrictions.watchTime')"
									path="requiredWatchTime"
								>
									<NInputNumber
										v-model:value="formValue.requiredWatchTime"
										:min="0"
										class="grid-stats-item"
									/>
								</NFormItem>
							</NGridItem>

							<NGridItem :span="1">
								<NFormItem
									:label="t('commands.modal.restrictions.messages')"
									path="requiredMessages"
								>
									<NInputNumber
										v-model:value="formValue.requiredMessages"
										:min="0"
										class="grid-stats-item"
									/>
								</NFormItem>
							</NGridItem>

							<NGridItem :span="1">
								<NFormItem
									:label="t('commands.modal.restrictions.channelsPoints')"
									path="requiredUsedChannelPoints"
								>
									<NInputNumber
										v-model:value="formValue.requiredUsedChannelPoints"
										:min="0"
										class="grid-stats-item"
									/>
								</NFormItem>
							</NGridItem>
						</NGrid>
					</NTabPane>
					<NTabPane name="cooldown" :tab="t('commands.modal.cooldown.label')">
						<NGrid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5">
							<NGridItem :span="1">
								<NFormItem
									:label="t('commands.modal.cooldown.value')"
									path="cooldown"
								>
									<NInputNumber
										v-model:value="formValue.cooldown"
										:min="0"
										class="grid-stats-item"
									/>
								</NFormItem>
							</NGridItem>

							<NGridItem :span="1">
								<NFormItem
									:label="t('commands.modal.cooldown.type.name')"
									path="cooldownType"
								>
									<NSelect
										v-model:value="formValue.cooldownType"
										:options="[
											{
												label: t('commands.modal.cooldown.type.global'),
												value: 'GLOBAL',
											},
											{
												label: t('commands.modal.cooldown.type.user'),
												value: 'PER_USER',
											},
										]"
									/>
								</NFormItem>
							</NGridItem>
						</NGrid>

						<div class="flex flex-col gap-1">
							<NButtonGroup
								v-for="(group, index) of chunk(rolesSelectOptions.sort(), 5)"
								:key="index"
							>
								<NButton
									v-for="option of group"
									:key="option.value"
									:type="formValue.cooldownRolesIds.includes(option.value) ? 'success' : 'default'"
									secondary
									@click="() => {
										if (formValue!.cooldownRolesIds.includes(option.value)) {
											formValue!.cooldownRolesIds = formValue!.cooldownRolesIds.filter(r => r !== option.value)
										}
										else {
											formValue!.cooldownRolesIds.push(option.value)
										}
									}"
								>
									<template #icon>
										<IconSquareCheck v-if="formValue.cooldownRolesIds.includes(option.value)" />
										<IconSquare v-else />
									</template>
									{{ option.label }}
								</NButton>
							</NButtonGroup>
						</div>
					</NTabPane>
					<NTabPane name="settings" :tab="t('commands.modal.settings.divider')">
						<NGrid cols="1 s:2 m:2 l:2" responsive="screen" :x-gap="5" :y-gap="5">
							<NGridItem :span="1">
								<NCard class="h-full">
									<div class="settings-card-body">
										<NSpace vertical>
											<NText>{{ t("sharedTexts.reply.label") }}</NText>
											<NText>{{ t("sharedTexts.reply.text") }}</NText>
										</NSpace>
										<NSwitch v-model:value="formValue.isReply" />
									</div>
								</NCard>
							</NGridItem>

							<NGridItem :span="1">
								<NCard class="h-full">
									<div class="settings-card-body">
										<NSpace vertical>
											<NText>{{ t('commands.modal.settings.visible.label') }}</NText>
											<NText>{{ t('commands.modal.settings.visible.text') }}</NText>
										</NSpace>
										<NSwitch v-model:value="formValue.visible" />
									</div>
								</NCard>
							</NGridItem>

							<NGridItem :span="1">
								<NCard class="h-full">
									<div class="settings-card-body">
										<NSpace vertical>
											<NText>{{ t('commands.modal.settings.keepOrder.label') }}</NText>
											<NText>{{ t('commands.modal.settings.keepOrder.text') }}</NText>
										</NSpace>
										<NSwitch v-model:value="formValue.keepResponsesOrder" />
									</div>
								</NCard>
							</NGridItem>

							<NGridItem :span="1">
								<NCard class="h-full">
									<div class="settings-card-body">
										<NSpace vertical>
											<NText>{{ t('commands.modal.settings.onlineOnly.label') }}</NText>
											<NText>{{ t('commands.modal.settings.onlineOnly.text') }}</NText>
										</NSpace>
										<NSwitch v-model:value="formValue.onlineOnly" />
									</div>
								</NCard>
							</NGridItem>
						</NGrid>

						<NDivider>
							{{ t('commands.modal.settings.other.divider') }}
						</NDivider>

						<NFormItem :label="t('commands.modal.gameCategories.label')" path="enabledGameCategories">
							<TwitchCategorySearch v-model="formValue.enabledCategories" multiple />
						</NFormItem>

						<NFormItem :label="t('commands.modal.description.label')" path="description">
							<NInput
								v-model:value="formValue.description" placeholder="Description" type="textarea"
								autosize
							/>
						</NFormItem>

						<NFormItem :label="t('commands.modal.settings.other.commandGroup')" path="groupId">
							<NButton v-if="!commandsGroupsOptions.length" secondary disabled>
								No groups created
							</NButton>
							<div
								v-else
								class="flex flex-col gap-1"
							>
								<NButtonGroup
									v-for="(group, index) of chunk(commandsGroupsOptions.sort(), 4)"
									:key="index"
								>
									<NButton
										v-for="option of group"
										:key="option.value"
										:type="formValue.groupId === option.value ? 'success' : 'default'"
										secondary
										@click="() => {
											if (formValue!.groupId === option.value) {
												formValue!.groupId = undefined
											}
											else {
												formValue!.groupId = option.value
											}
										}"
									>
										<template #icon>
											<IconSquareCheck v-if="formValue.groupId === option.value" />
											<IconSquare v-else />
										</template>
										{{ option.label }}
									</NButton>
								</NButtonGroup>
							</div>
						</NFormItem>
					</NTabPane>
				</NTabs>
			</NForm>

			<NButton class="mt-2" secondary type="success" block @click="commandEdit.save">
				{{ t('sharedButtons.save') }}
			</NButton>
		</div>
	</NModal>
</template>

<style scoped>
.groups :deep(.create-button) {
	@apply hidden;
}

.grid-stats-item {
	@apply w-full;
}

.settings-card-body {
	@apply flex flex-row justify-between;
}
</style>
