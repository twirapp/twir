<script setup lang="ts">
import {
	BadgePlus,
	Ellipsis,
	GripVertical,
	MessageCircleReplyIcon,
	Settings,
	Trash,
} from 'lucide-vue-next'
import { FieldArray, useField } from 'vee-validate'
import { computed, ref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'


import { useCommandEditV2 } from '../../composables/use-command-edit-v2'

import type { FormSchema } from '../../composables/use-command-edit-v2'

// ...existing code...

import TwitchCategorySearchShadcnMultiple from '#layers/dashboard/components/twitch-category-search-shadcn-multiple.vue'

import VariableInput from '#layers/dashboard/components/variable-input.vue'

const { t } = useI18n()
const { user: profile } = storeToRefs(useDashboardAuth())

const { errors: responsesErrors, value, setValue } = useField<FormSchema['responses']>('responses')

const maxCommandResponses = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan?.maxCommandsResponses ?? 3
})

function handlePush() {
	setValue([
		...value.value,
		{ text: '', twitchCategoriesIds: [], onlineOnly: false, offlineOnly: false },
	])
}

const responseDialogOpened = ref(false)

const { command } = useCommandEditV2()

const editable = computed(() => !command.value?.default)
</script>

<template>
	<UiCard>
		<UiCardHeader class="flex flex-row place-content-center flex-wrap">
			<UiCardTitle
				:class="{ 'text-destructive': responsesErrors.length }"
				class="flex items-center gap-2"
			>
				<MessageCircleReplyIcon />
				{{ t('sharedTexts.responses') }}
			</UiCardTitle>
		</UiCardHeader>
		<UiCardContent v-if="editable" class="flex flex-col gap-2 pt-4">
			<FieldArray v-slot="{ fields, remove }" name="responses">
				<VueDraggable v-model="value" handle=".drag-handle" class="flex flex-col gap-2">
					<div v-for="(field, index) in fields" :key="`responses-text-${field.key}`">
						<UiDialog>
							<UiFormField v-slot="{ componentField }" :name="`responses[${index}].text`">
								<UiFormItem>
									<div class="relative flex items-center">
										<UiFormControl>
											<div class="w-full">
												<div
													class="absolute flex left-0 rounded-l-md h-full bg-accent w-4 cursor-move drag-handle"
												>
													<GripVertical class="my-auto size-6" />
												</div>
												<VariableInput
													input-type="textarea"
													class="pl-6 pr-14!"
													:model-value="componentField.modelValue"
													:min-rows="1"
													:rows="1"
													popoverAlign="end"
													popoverSide="bottom"
													@update:model-value="componentField.onChange"
												>
													<template #additional-buttons>
														<UiDropdownMenu>
															<UiDropdownMenuTrigger as-child>
																<button class="hover:bg-accent p-1 rounded-md">
																	<Ellipsis class="size-4 opacity-50" />
																</button>
															</UiDropdownMenuTrigger>

															<UiDropdownMenuContent :hideWhenDetached="false">
																<UiDialogTrigger as-child>
																	<UiDropdownMenuItem @click="responseDialogOpened = true">
																		<div class="flex items-center gap-2">
																			<Settings class="size-4" />
																			Settings
																		</div>
																	</UiDropdownMenuItem>
																</UiDialogTrigger>

																<UiDropdownMenuItem @click="remove(index)">
																	<div class="flex items-center gap-2">
																		<Trash class="size-4" />
																		Remove
																	</div>
																</UiDropdownMenuItem>
															</UiDropdownMenuContent>
														</UiDropdownMenu>
													</template>
												</VariableInput>
											</div>
										</UiFormControl>
									</div>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiDialogContent>
								<UiDialogHeader>
									<UiDialogTitle>Edit response settings</UiDialogTitle>
								</UiDialogHeader>

								<UiFormField
									v-slot="{ componentField }"
									:name="`responses[${index}].twitchCategoriesIds`"
								>
									<UiFormItem>
										<UiFormLabel> Category for response </UiFormLabel>
										<UiFormControl>
											<TwitchCategorySearchShadcnMultiple
												:id="componentField.name"
												:model-value="componentField.modelValue"
												@update:model-value="componentField['onUpdate:modelValue']"
											/>
										</UiFormControl>
									</UiFormItem>
								</UiFormField>

								<div class="grid grid-cols-1 md:grid-cols-2 gap-2 w-full">
									<UiFormField
										v-slot="{ value, handleChange }"
										type="checkbox"
										:name="`responses[${index}].onlineOnly`"
									>
										<UiFormItem
											class="flex flex-row items-start gap-x-3 space-y-0 rounded-md border p-4"
										>
											<UiFormControl>
												<UiCheckbox :model-value="value" @update:model-value="handleChange" />
											</UiFormControl>
											<div class="space-y-1 leading-none">
												<UiFormLabel>{{ t('commands.modal.settings.onlineOnly.label') }}</UiFormLabel>
												<UiFormMessage />
											</div>
										</UiFormItem>
									</UiFormField>

									<UiFormField
										v-slot="{ value, handleChange }"
										type="checkbox"
										:name="`responses[${index}].offlineOnly`"
									>
										<UiFormItem
											class="flex flex-row items-start gap-x-3 space-y-0 rounded-md border p-4"
										>
											<UiFormControl>
												<UiCheckbox :model-value="value" @update:model-value="handleChange" />
											</UiFormControl>
											<div class="space-y-1 leading-none">
												<UiFormLabel>{{ t('commands.modal.settings.offlineOnly.label') }}</UiFormLabel>
												<UiFormMessage />
											</div>
										</UiFormItem>
									</UiFormField>
								</div>

								<UiDialogFooter>
									<UiDialogClose>
										<UiButton>Close</UiButton>
									</UiDialogClose>
								</UiDialogFooter>
							</UiDialogContent>
						</UiDialog>
					</div>
				</VueDraggable>

				<UiButton
					type="button"
					variant="outline"
					size="sm"
					class="text-xs w-full flex gap-2 items-center mt-2"
					:disabled="(fields.length ?? 0) >= maxCommandResponses"
					@click="handlePush"
				>
					<BadgePlus class="size-4" />
					Add response {{ fields.length ?? 0 }} / {{ maxCommandResponses }}
				</UiButton>
			</FieldArray>
		</UiCardContent>

		<UiCardContent v-else>
			<UiAlert>
				<UiAlertDescription>
					{{ t('commands.modal.responses.defaultWarning') }}
				</UiAlertDescription>
			</UiAlert>
		</UiCardContent>
	</UiCard>
</template>
