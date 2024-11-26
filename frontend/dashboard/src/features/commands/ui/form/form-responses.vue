<script setup lang="ts">
import { BadgePlus, Ellipsis, GripVertical, Settings, Trash } from 'lucide-vue-next'
import { FieldArray, useField } from 'vee-validate'
import { computed, ref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { useI18n } from 'vue-i18n'

import { useCommandEditV2 } from '../../composables/use-command-edit-v2'

import TwitchCategorySearchShadcnMultiple from '@/components/twitch-category-search-shadcn-multiple.vue'
import {
	Alert,
	AlertDescription,
} from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'
import FormLabel from '@/components/ui/form/FormLabel.vue'
import VariableInput from '@/components/variable-input.vue'

const { t } = useI18n()

const { errors: responsesErrors, value, setValue } = useField<any[]>('responses')

function handlePush() {
	setValue([...value.value, { text: '', twitchCategoriesIds: [] }])
}

const responseDialogOpened = ref(false)

const { command } = useCommandEditV2()

const editable = computed(() => !command.value?.default)
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle :class="{ 'text-destructive': responsesErrors.length }">
				{{ t('sharedTexts.responses') }}
			</CardTitle>
		</CardHeader>
		<CardContent v-if="editable" class="flex flex-col gap-2">
			<FieldArray v-slot="{ fields, remove }" name="responses">
				<VueDraggable
					v-model="value"
					handle=".drag-handle"
					class="flex flex-col gap-2"
				>
					<div v-for="(field, index) in fields" :key="`responses-text-${field.key}`">
						<Dialog>
							<FormField v-slot="{ componentField }" :name="`responses[${index}].text`">
								<FormItem>
									<div class="relative flex items-center">
										<FormControl>
											<div class="w-full">
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
																<DialogTrigger as-child>
																	<DropdownMenuItem @click="responseDialogOpened = true">
																		<div class="flex items-center gap-2">
																			<Settings class="size-4" />
																			Settings
																		</div>
																	</DropdownMenuItem>
																</DialogTrigger>

																<DropdownMenuItem @click="remove(index)">
																	<div class="flex items-center gap-2">
																		<Trash class="size-4" />
																		Remove
																	</div>
																</DropdownMenuItem>
															</DropdownMenuContent>
														</DropdownMenu>
													</template>
												</VariableInput>
											</div>
										</FormControl>
									</div>
									<FormMessage />
								</FormItem>
							</FormField>

							<DialogContent>
								<DialogHeader>
									<DialogTitle>Edit response settings</DialogTitle>
								</DialogHeader>

								<FormField v-slot="{ componentField }" :name="`responses[${index}].twitchCategoriesIds`">
									<FormItem>
										<FormLabel>
											Category for response
										</FormLabel>
										<FormControl>
											<TwitchCategorySearchShadcnMultiple
												:id="componentField.name"
												:model-value="componentField.modelValue"
												@update:model-value="componentField['onUpdate:modelValue']"
											/>
										</FormControl>
									</FormItem>
								</FormField>

								<DialogFooter>
									<DialogClose>
										<Button>Close</Button>
									</DialogClose>
								</DialogFooter>
							</DialogContent>
						</Dialog>
					</div>
				</VueDraggable>

				<Button
					type="button"
					variant="outline"
					size="sm"
					class="text-xs w-full flex gap-2 items-center mt-2"
					:disabled="(fields.length ?? 0) >= 3"
					@click="handlePush"
				>
					<BadgePlus class="size-4" />
					Add response {{ fields.length ?? 0 }} / 3
				</Button>
			</FieldArray>
		</CardContent>

		<CardContent v-else>
			<Alert>
				<AlertDescription>
					{{ t('commands.modal.responses.defaultWarning') }}
				</AlertDescription>
			</Alert>
		</CardContent>
	</Card>
</template>
