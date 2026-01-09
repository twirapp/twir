<script setup lang="ts">
import { EditIcon, WrenchIcon, XIcon } from 'lucide-vue-next'


import { useCommandEditV2 } from '../../composables/use-command-edit-v2'

import { useCommandsGroupsApi } from '#layers/dashboard/api/commands/commands-groups'
import ManageGroups from '#layers/dashboard/components/commands/manageGroups.vue'
import DialogOrSheet from '#layers/dashboard/components/dialog-or-sheet.vue'










const { t } = useI18n()

const groupsApi = useCommandsGroupsApi()
const { data: groups } = groupsApi.useQueryGroups()
const { isCustom } = useCommandEditV2()

function computeSelectedGroupColor(id: string) {
	if (!groups?.value?.commandsGroups) {
		return ''
	}

	const group = groups.value.commandsGroups.find((group) => group.id === id)
	return group?.color || ''
}
</script>

<template>
	<UiCard>
		<UiCardHeader class="flex flex-row justify-between flex-wrap">
			<div></div>

			<UiCardTitle class="flex items-center gap-2">
				<WrenchIcon />
				General
			</UiCardTitle>

			<UiCardAction>
				<UiFormField v-slot="{ field }" name="enabled">
					<UiFormItem class="space-y-0 flex items-center gap-4">
						<UiFormControl>
							<UiSwitch
								:model-value="field.value"
								default-checked
								@update:model-value="field['onUpdate:modelValue']"
							/>
						</UiFormControl>
					</UiFormItem>
				</UiFormField>
			</UiCardAction>
		</UiCardHeader>
		<UiCardContent class="flex flex-col gap-4 pt-4">
			<UiFormField v-slot="{ componentField }" name="name">
				<UiFormItem>
					<UiFormLabel>{{ t('sharedTexts.name') }}</UiFormLabel>
					<UiFormControl>
						<UiInput type="text" v-bind="componentField" />
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField v-slot="{ field, errorMessage }" name="aliases">
				<UiFormItem>
					<UiFormLabel>{{ t('commands.modal.aliases.label') }}</UiFormLabel>
					<UiFormControl>
						<UiTagsInput
							:invalid="!!errorMessage"
							:model-value="field.value"
							@update:model-value="field['onUpdate:modelValue']"
						>
							<UiTagsInputItem v-for="item in field.value" :key="item" :value="item">
								<UiTagsInputItemText />
								<UiTagsInputItemDelete />
							</UiTagsInputItem>

							<UiTagsInputInput :placeholder="t('commands.modal.aliases.label')" class="h-7" />
						</UiTagsInput>
					</UiFormControl>
					<UiFormDescription>
						Alternative names for triggering this command. Press enter to submit.
					</UiFormDescription>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField v-slot="{ componentField }" name="description">
				<UiFormItem>
					<UiFormLabel>{{ t('commands.modal.description.label') }}</UiFormLabel>
					<UiFormControl>
						<UiInput v-bind="componentField" />
					</UiFormControl>
					<UiFormDescription>
						Description which is showed on public page, in command help, in FFZ addon, e.t.c.
					</UiFormDescription>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField v-slot="{ componentField }" name="groupId">
				<UiFormItem>
					<UiFormLabel class="flex gap-2">
						<span>{{ t('commands.modal.settings.other.commandGroup') }}</span>
						<UiDialog>
							<UiDialogTrigger as-child>
								<span class="flex flex-row gap-1 items-center cursor-pointer underline">
									{{ t('commands.groups.manageButton') }}
									<EditIcon class="size-4" />
								</span>

								<DialogOrSheet>
									<UiDialogTitle>{{ t('commands.groups.manageButton') }}</UiDialogTitle>
									<ManageGroups />
								</DialogOrSheet>
							</UiDialogTrigger>
						</UiDialog>
					</UiFormLabel>

					<div v-if="isCustom" class="flex flex-row gap-2">
						<UiFormControl>
							<UiSelect v-bind="componentField">
								<UiSelectTrigger>
									<UiSelectValue
										:style="{ color: computeSelectedGroupColor(componentField.modelValue) }"
										placeholder="No group"
									/>
								</UiSelectTrigger>
								<UiSelectContent>
									<div v-if="!groups?.commandsGroups.length" class="p-2">No groups created</div>
									<UiSelectGroup v-else>
										<UiSelectItem
											v-for="group in groups?.commandsGroups"
											:key="group.id"
											:value="group.id"
											:style="{ color: group.color }"
										>
											{{ group.name }}
										</UiSelectItem>
									</UiSelectGroup>
								</UiSelectContent>
							</UiSelect>
						</UiFormControl>
						<UiButton
							v-if="componentField['onUpdate:modelValue']"
							variant="outline"
							type="button"
							@click="componentField['onUpdate:modelValue'](null)"
						>
							<XIcon class="size-4" />
						</UiButton>
					</div>
					<UiAlert v-else class="py-2">
						<UiAlertDescription> Group cannot be set for default command </UiAlertDescription>
					</UiAlert>
					<UiFormMessage />
					<UiFormDescription>
						Groups used to create "folder" of commands in dashboard and public page, so you can
						stick related commands together.
					</UiFormDescription>
				</UiFormItem>
			</UiFormField>
		</UiCardContent>
	</UiCard>
</template>
