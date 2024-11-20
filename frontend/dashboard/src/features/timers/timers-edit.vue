<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { BadgePlus, Ellipsis, GripVertical, Trash } from 'lucide-vue-next'
import { FieldArray, useForm } from 'vee-validate'
import { computed, onMounted, ref, toRaw } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { DialogFooter } from '@/components/ui/dialog'
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
import { Label } from '@/components/ui/label'
import { Slider } from '@/components/ui/slider'
import VariableInput from '@/components/variable-input.vue'
import { formSchema, useTimersEdit } from '@/features/timers/composables/use-timers-edit.js'
import PageLayout from '@/layout/page-layout.vue'

const route = useRoute()
const { t } = useI18n()
const { findTimer, submit } = useTimersEdit()

const loading = ref(true)

const { resetForm, handleSubmit, controlledValues, errors, setValues } = useForm({
	validationSchema: toTypedSchema(formSchema),
	initialValues: {
		timeInterval: 1,
		messageInterval: 0,
		responses: [{ text: '', isAnnounce: false }],
	},
})

onMounted(async () => {
	resetForm()

	if (typeof route.params.id === 'string') {
		const timer = await findTimer(route.params.id)
		if (timer) {
			setValues(toRaw(timer))
		}
	}

	loading.value = false
})

const onSubmit = handleSubmit(submit)

const responsesHasError = computed(() => {
	return Object.keys(errors.value).some((key) => key.startsWith('responses'))
})
</script>

<template>
	<PageLayout>
		<template #title>
			{{ route.params.id === 'create' ? t('sharedTexts.create') : t('sharedTexts.edit') }}
		</template>

		<template #content>
			<form
				class="flex flex-col gap-4 max-w-4xl mx-auto"
				:class="{ 'blur-sm': loading }"
				@submit="onSubmit"
			>
				<FormField v-slot="{ componentField }" name="name">
					<FormItem>
						<FormLabel>{{ t('sharedTexts.name') }}</FormLabel>
						<FormControl>
							<Input type="text" v-bind="componentField" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<Card>
					<CardContent class="pt-6">
						<FormField v-slot="{ componentField }" name="timeInterval">
							<FormItem>
								<FormLabel>{{ t('timers.table.columns.intervalInMinutes') }}</FormLabel>
								<FormControl>
									<Slider
										:model-value="[componentField.modelValue]"
										:max="100"
										:default-value="[0, 100]"
										:min="0"
										:step="1"
										@update:model-value="(v) => {
											if (!v) return
											componentField.onChange(v[0])
										}"
									/>
								</FormControl>
								<FormMessage />
								<FormDescription class="flex justify-end">
									<span>{{ componentField.modelValue }} minutes</span>
								</FormDescription>
							</FormItem>
						</FormField>

						<FormField v-slot="{ componentField }" name="messageInterval">
							<FormItem>
								<FormLabel>{{ t('timers.table.columns.intervalInMessages') }}</FormLabel>
								<FormControl>
									<Input type="number" placeholder="0" v-bind="componentField" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</CardContent>
				</Card>

				<Label
					:class="{ 'text-destructive': responsesHasError }"
				>{{ t('sharedTexts.responses') }}</Label>

				<FieldArray v-slot="{ fields, push, remove }" name="responses">
					<VueDraggable
						v-model="controlledValues.responses!"
						handle=".drag-handle"
						class="flex flex-col gap-2"
					>
						<div v-for="(field, index) in fields" :key="`responses-text-${field.key}`">
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
											</div>
										</FormControl>
									</div>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</VueDraggable>
					<Button
						type="button"
						variant="outline"
						size="sm"
						class="text-xs w-full flex gap-2 items-center mt-2"
						:disabled="(controlledValues.responses?.length ?? 0) >= 10"
						@click="push({ text: '', isAnnounce: false })"
					>
						<BadgePlus class="size-4" />
						Add response {{ controlledValues.responses?.length ?? 0 }} / 10
					</Button>
				</FieldArray>

				<DialogFooter>
					<Button type="submit">
						{{ t('sharedButtons.save') }}
					</Button>
				</DialogFooter>
			</form>
		</template>
	</PageLayout>
</template>
