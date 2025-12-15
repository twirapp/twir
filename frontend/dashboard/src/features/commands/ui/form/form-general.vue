<script setup lang="ts">
import { EditIcon, WrenchIcon, XIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { useCommandEditV2 } from '../../composables/use-command-edit-v2'

import { useCommandsGroupsApi } from '@/api/commands/commands-groups'
import ManageGroups from '@/components/commands/manageGroups.vue'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Alert, AlertDescription } from '@/components/ui/alert'
import Button from '@/components/ui/button/Button.vue'
import { Card, CardAction, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import {
	TagsInput,
	TagsInputInput,
	TagsInputItem,
	TagsInputItemDelete,
	TagsInputItemText,
} from '@/components/ui/tags-input'

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
	<Card>
		<CardHeader class="flex flex-row justify-between flex-wrap">
			<div></div>

			<CardTitle class="flex items-center gap-2">
				<WrenchIcon />
				General
			</CardTitle>

			<CardAction>
				<FormField v-slot="{ field }" name="enabled">
					<FormItem class="space-y-0 flex items-center gap-4">
						<FormControl>
							<Switch
								:model-value="field.value"
								default-checked
								@update:model-value="field['onUpdate:modelValue']"
							/>
						</FormControl>
					</FormItem>
				</FormField>
			</CardAction>
		</CardHeader>
		<CardContent class="flex flex-col gap-4 pt-4">
			<FormField v-slot="{ componentField }" name="name">
				<FormItem>
					<FormLabel>{{ t('sharedTexts.name') }}</FormLabel>
					<FormControl>
						<Input type="text" v-bind="componentField" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ field, errorMessage }" name="aliases">
				<FormItem>
					<FormLabel>{{ t('commands.modal.aliases.label') }}</FormLabel>
					<FormControl>
						<TagsInput
							:invalid="!!errorMessage"
							:model-value="field.value"
							@update:model-value="field['onUpdate:modelValue']"
						>
							<TagsInputItem v-for="item in field.value" :key="item" :value="item">
								<TagsInputItemText />
								<TagsInputItemDelete />
							</TagsInputItem>

							<TagsInputInput :placeholder="t('commands.modal.aliases.label')" />
						</TagsInput>
					</FormControl>
					<FormDescription>
						Alternative names for triggering this command. Press enter to submit.
					</FormDescription>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="description">
				<FormItem>
					<FormLabel>{{ t('commands.modal.description.label') }}</FormLabel>
					<FormControl>
						<Input v-bind="componentField" />
					</FormControl>
					<FormDescription>
						Description which is showed on public page, in command help, in FFZ addon, e.t.c.
					</FormDescription>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="groupId">
				<FormItem>
					<FormLabel class="flex gap-2">
						<span>{{ t('commands.modal.settings.other.commandGroup') }}</span>
						<Dialog>
							<DialogTrigger as-child>
								<span class="flex flex-row gap-1 items-center cursor-pointer underline">
									{{ t('commands.groups.manageButton') }}
									<EditIcon class="size-4" />
								</span>

								<DialogOrSheet>
									<DialogTitle>{{ t('commands.groups.manageButton') }}</DialogTitle>
									<ManageGroups />
								</DialogOrSheet>
							</DialogTrigger>
						</Dialog>
					</FormLabel>

					<div v-if="isCustom" class="flex flex-row gap-2">
						<FormControl>
							<Select v-bind="componentField">
								<SelectTrigger>
									<SelectValue
										:style="{ color: computeSelectedGroupColor(componentField.modelValue) }"
										placeholder="No group"
									/>
								</SelectTrigger>
								<SelectContent>
									<div v-if="!groups?.commandsGroups.length" class="p-2">No groups created</div>
									<SelectGroup v-else>
										<SelectItem
											v-for="group in groups?.commandsGroups"
											:key="group.id"
											:value="group.id"
											:style="{ color: group.color }"
										>
											{{ group.name }}
										</SelectItem>
									</SelectGroup>
								</SelectContent>
							</Select>
						</FormControl>
						<Button
							v-if="componentField['onUpdate:modelValue']"
							variant="outline"
							type="button"
							@click="componentField['onUpdate:modelValue'](null)"
						>
							<XIcon class="size-4" />
						</Button>
					</div>
					<Alert v-else class="py-2">
						<AlertDescription> Group cannot be set for default command </AlertDescription>
					</Alert>
					<FormMessage />
					<FormDescription>
						Groups used to create "folder" of commands in dashboard and public page, so you can
						stick related commands together.
					</FormDescription>
				</FormItem>
			</FormField>
		</CardContent>
	</Card>
</template>
