<script setup lang="ts">
import { XIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { useCommandEditV2 } from '../../composables/use-command-edit-v2'

import { useCommandsGroupsApi } from '@/api/commands/commands-groups'
import { Alert, AlertDescription } from '@/components/ui/alert'
import Button from '@/components/ui/button/Button.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
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
			<CardTitle>General</CardTitle>

			<FormField v-slot="{ field }" name="enabled">
				<FormItem class="space-y-0 flex items-center gap-4">
					<FormLabel class="text-base">
						{{ t('sharedTexts.enabled') }}
					</FormLabel>
					<FormControl>
						<Switch
							:checked="field.value"
							default-checked
							@update:checked="field['onUpdate:modelValue']"
						/>
					</FormControl>
				</FormItem>
			</FormField>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
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
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="description">
				<FormItem>
					<FormLabel>{{ t('commands.modal.description.label') }}</FormLabel>
					<FormControl>
						<Input v-bind="componentField" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="groupId">
				<FormItem>
					<FormLabel>{{ t('commands.modal.settings.other.commandGroup') }}</FormLabel>

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
									<SelectGroup>
										<SelectItem
											v-for="(group) in groups?.commandsGroups"
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
						<AlertDescription>
							Group cannot be set for default command
						</AlertDescription>
					</Alert>
					<FormMessage />
				</FormItem>
			</FormField>
		</CardContent>
	</Card>
</template>
