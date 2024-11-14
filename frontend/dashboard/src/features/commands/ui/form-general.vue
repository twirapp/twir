<script setup lang="ts">
import { BadgePlus, Ellipsis, GripVertical, Trash } from 'lucide-vue-next'
import { FieldArray } from 'vee-validate'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import {
	DropdownMenu,
	DropdownMenuCheckboxItem,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { TagsInput, TagsInputInput, TagsInputItem, TagsInputItemDelete, TagsInputItemText } from '@/components/ui/tags-input'
import VariableInput from '@/components/variable-input.vue'

const { t } = useI18n()

const checkboxes = computed(() => {
	return [
		{
			name: 'isReply',
			label: t('sharedTexts.reply.label'),
			description: t('sharedTexts.reply.text'),
		},
		{
			name: 'visible',
			label: t('commands.modal.settings.visible.label'),
			description: t('commands.modal.settings.visible.text'),
		},
		{
			name: 'keepOrder',
			label: t('commands.modal.settings.keepOrder.label'),
			description: t('commands.modal.settings.keepOrder.text'),
		},
		{
			name: 'onlineOnly',
			label: t('commands.modal.settings.onlineOnly.label'),
			description: t('commands.modal.settings.onlineOnly.text'),
		},
	]
})
</script>

<template>
	<div class="flex flex-col gap-4">
		<Card>
			<CardHeader class="flex flex-row justify-between flex-wrap">
				<CardTitle>General</CardTitle>
				<div class="flex items-center gap-4">
					<FormField v-slot="{ value, handleChange }" name="enabled">
						<FormLabel class="text-base">
							Enabled
						</FormLabel>
						<FormControl>
							<Switch
								:checked="value"
								@update:checked="handleChange"
							/>
						</FormControl>
					</FormField>
				</div>
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
					<FormLabel>{{ t('commands.modal.aliases.label') }}</FormLabel>
					<FormItem>
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
						<FormMessage />
					</FormItem>
				</FormField>
			</CardContent>
		</Card>

		<Card>
			<CardHeader>
				<CardTitle>Responses</CardTitle>
			</CardHeader>
			<CardContent>
				<FieldArray v-slot="{ fields, push, remove }" name="responses">
					<div v-for="(field, index) in fields" :key="`responses-text-${field.key}`">
						<FormField v-slot="{ componentField }" :name="`responses[${index}].text`">
							<FormItem>
								<div class="relative flex items-center">
									<FormControl>
										<div class="absolute flex left-0 rounded-l-md h-full bg-accent w-4 cursor-move drag-handle">
											<GripVertical class="my-auto size-6" />
										</div>
										<VariableInput
											input-type="textarea"
											class="pl-6 !pr-14"
											:model-value="componentField.modelValue"
											:min-rows="1"
											:rows="1"
											popoverAlign="end"
											popoverSide="bottom"
											@update:model-value="componentField.onChange"
										>
											<template #additional-buttons>
												<DropdownMenu>
													<DropdownMenuTrigger as-child>
														<button class="hover:bg-accent p-1 rounded-md">
															<Ellipsis class="size-4 opacity-50" />
														</button>
													</DropdownMenuTrigger>

													<DropdownMenuContent :hideWhenDetached="false">
														<DropdownMenuCheckboxItem
															v-model:checked="(field.value as any).isAnnounce"
														>
															Send as announcement
														</DropdownMenuCheckboxItem>
														<DropdownMenuItem @click="remove">
															<div class="flex items-center gap-2">
																<Trash class="size-4" />
																Remove
															</div>
														</DropdownMenuItem>
													</DropdownMenuContent>
												</DropdownMenu>
											</template>
										</VariableInput>
									</FormControl>
								</div>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>

					<Button
						type="button"
						variant="outline"
						size="sm"
						class="text-xs w-full flex gap-2 items-center mt-2"
						:disabled="(fields.length ?? 0) >= 3"
						@click="push({ text: '' })"
					>
						<BadgePlus class="size-4" />
						Add response {{ fields.length ?? 0 }} / 3
					</Button>
				</FieldArray>
			</CardContent>
		</Card>

		<Card>
			<CardHeader>
				<CardTitle>Conditions</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="grid grid-cols-1 lg:grid-cols-2 gap-2">
					<FormField
						v-for="checkbox of checkboxes"
						:key="checkbox.name"
						v-slot="{ value, handleChange }"
						type="checkbox"
						:name="checkbox.name"
					>
						<FormItem class="flex flex-row items-start gap-x-3 space-y-0 rounded-md border p-4" :class="[`col-span-${checkbox.span ?? 1}`]">
							<FormControl>
								<Checkbox :checked="value" @update:checked="handleChange" />
							</FormControl>
							<div class="space-y-1 leading-none">
								<FormLabel>{{ checkbox.label }}</FormLabel>
								<FormDescription v-if="checkbox.description">
									{{ checkbox.description }}
								</FormDescription>
								<FormMessage />
							</div>
						</FormItem>
					</FormField>
				</div>
			</CardContent>
		</Card>
	</div>
</template>
