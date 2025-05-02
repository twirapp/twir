<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { PlusIcon, Trash2 } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import * as z from 'zod'

import { useCommandsApi } from '@/api/commands/commands'
import { useEventsApi } from '@/api/events'
import { useKeywordsApi } from '@/api/keywords'
import { eventTypeSelectOptions, flatOperations, operationTypeSelectOptions } from '@/components/events/helpers'
import TwitchRewardsSelector from '@/components/rewardsSelector.vue'
import { Button } from '@/components/ui/button'
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card'
import {
	Form,
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
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import {
	Tabs,
	TabsContent,
	TabsList,
	TabsTrigger,
} from '@/components/ui/tabs'
import { useToast } from '@/components/ui/toast/use-toast'
import VariableInput from '@/components/variable-input.vue'
import PageLayout from '@/layout/page-layout.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const { toast } = useToast()
const eventsApi = useEventsApi()
const isNewEvent = computed(() => route.params.id === 'new')
const eventId = computed(() => isNewEvent.value ? '' : String(route.params.id))

// Fetch event data if editing
const { data: eventData, fetching: isLoadingEvent } = eventsApi.useQueryEventById(eventId.value)
const event = computed(() => eventData.value?.eventById)

// Fetch commands and keywords for selectors
const commandsApi = useCommandsApi()
const { data: commandsData } = commandsApi.useQueryCommands()
const commands = computed(() => commandsData.value?.commands || [])

const keywordsApi = useKeywordsApi()
const { data: keywordsData } = keywordsApi.useQueryKeywords()
const keywords = computed(() => keywordsData.value?.keywords || [])

// Form validation schema
const formSchema = toTypedSchema(z.object({
	type: z.string().min(1, t('events.validation.typeRequired')),
	description: z.string().min(1, t('events.validation.descriptionRequired')),
	enabled: z.boolean().default(true),
	onlineOnly: z.boolean().default(false),
	rewardId: z.string().optional(),
	commandId: z.string().optional(),
	keywordId: z.string().optional(),
	operations: z.array(z.object({
		type: z.string().min(1, t('events.validation.operationTypeRequired')),
		input: z.string().optional(),
		delay: z.number().min(0).default(0),
		repeat: z.number().min(0).default(0),
		useAnnounce: z.boolean().default(false),
		timeoutTime: z.number().min(0).default(0),
		timeoutMessage: z.string().optional(),
		target: z.string().optional(),
		enabled: z.boolean().default(true),
		filters: z.array(z.object({
			type: z.string().min(1, t('events.validation.filterTypeRequired')),
			left: z.string().min(1, t('events.validation.leftRequired')),
			right: z.string().min(1, t('events.validation.rightRequired')),
		})).default([]),
	})).default([]),
}))

// Initialize form
const form = useForm({
	validationSchema: formSchema,
	initialValues: {
		type: '',
		description: '',
		enabled: true,
		onlineOnly: false,
		operations: [],
	},
})

// Update form values when event data is loaded
watch(() => event.value, (newEvent) => {
	if (newEvent) {
		form.setValues({
			type: newEvent.type,
			description: newEvent.description,
			enabled: newEvent.enabled,
			onlineOnly: newEvent.onlineOnly,
			rewardId: newEvent.rewardId || undefined,
			commandId: newEvent.commandId || undefined,
			keywordId: newEvent.keywordId || undefined,
			operations: newEvent.operations.map(op => ({
				type: op.type,
				input: op.input || undefined,
				delay: op.delay,
				repeat: op.repeat,
				useAnnounce: op.useAnnounce,
				timeoutTime: op.timeoutTime,
				timeoutMessage: op.timeoutMessage || undefined,
				target: op.target || undefined,
				enabled: op.enabled,
				filters: op.filters.map(filter => ({
					type: filter.type,
					left: filter.left,
					right: filter.right,
				})),
			})),
		})
	}
}, { immediate: true })

// Operations management
const selectedOperationTab = ref(0)
const currentOperation = computed(() => {
	const operations = form.values.operations
	return operations.length > 0 ? operations[selectedOperationTab.value] : null
})

function addOperation() {
	const operations = [...form.values.operations]
	operations.push({
		type: '',
		input: '',
		delay: 0,
		repeat: 0,
		useAnnounce: false,
		timeoutTime: 0,
		timeoutMessage: '',
		target: '',
		enabled: true,
		filters: [],
	})
	form.setFieldValue('operations', operations)
	selectedOperationTab.value = operations.length - 1
}

function removeOperation(index: number) {
	const operations = [...form.values.operations]
	operations.splice(index, 1)
	form.setFieldValue('operations', operations)

	if (selectedOperationTab.value >= operations.length) {
		selectedOperationTab.value = Math.max(0, operations.length - 1)
	}
}

function addFilter(operationIndex: number) {
	const operations = [...form.values.operations]
	operations[operationIndex].filters.push({
		type: 'EQUALS',
		left: '',
		right: '',
	})
	form.setFieldValue('operations', operations)
}

function removeFilter(operationIndex: number, filterIndex: number) {
	const operations = [...form.values.operations]
	operations[operationIndex].filters.splice(filterIndex, 1)
	form.setFieldValue('operations', operations)
}

// Form submission
const createEventMutation = eventsApi.useMutationCreateEvent()
const updateEventMutation = eventsApi.useMutationUpdateEvent()
const isSubmitting = ref(false)

async function onSubmit(values: z.infer<typeof formSchema>) {
	isSubmitting.value = true

	try {
		if (isNewEvent.value) {
			await createEventMutation.executeMutation({
				input: {
					type: values.type,
					description: values.description,
					enabled: values.enabled,
					onlineOnly: values.onlineOnly,
					rewardId: values.rewardId,
					commandId: values.commandId,
					keywordId: values.keywordId,
					operations: values.operations.map(op => ({
						type: op.type,
						input: op.input,
						delay: op.delay,
						repeat: op.repeat,
						useAnnounce: op.useAnnounce,
						timeoutTime: op.timeoutTime,
						timeoutMessage: op.timeoutMessage,
						target: op.target,
						enabled: op.enabled,
						filters: op.filters,
					})),
				},
			})

			toast({
				title: t('events.createSuccess'),
				description: t('events.createSuccessDescription'),
			})
		} else {
			await updateEventMutation.executeMutation({
				id: eventId.value,
				input: {
					type: values.type,
					description: values.description,
					enabled: values.enabled,
					onlineOnly: values.onlineOnly,
					rewardId: values.rewardId,
					commandId: values.commandId,
					keywordId: values.keywordId,
					operations: values.operations.map(op => ({
						type: op.type,
						input: op.input,
						delay: op.delay,
						repeat: op.repeat,
						useAnnounce: op.useAnnounce,
						timeoutTime: op.timeoutTime,
						timeoutMessage: op.timeoutMessage,
						target: op.target,
						enabled: op.enabled,
						filters: op.filters,
					})),
				},
			})

			toast({
				title: t('events.updateSuccess'),
				description: t('events.updateSuccessDescription'),
			})
		}

		router.push('/dashboard/events')
	} catch (error) {
		console.error(error)
		toast({
			title: isNewEvent.value ? t('events.createError') : t('events.updateError'),
			description: isNewEvent.value ? t('events.createErrorDescription') : t('events.updateErrorDescription'),
			variant: 'destructive',
		})
	} finally {
		isSubmitting.value = false
	}
}

function goBack() {
	router.push('/dashboard/events')
}
</script>

<template>
	<PageLayout>
		<template #title>
			{{ isNewEvent ? t('events.create') : t('events.edit') }}
		</template>

		<template #action>
			<Button variant="outline" @click="goBack">
				{{ t('sharedTexts.back') }}
			</Button>
		</template>

		<template #content>
			<div v-if="isLoadingEvent && !isNewEvent" class="flex justify-center items-center h-64">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
			</div>

			<Form v-else :form="form" class="space-y-6" @submit="onSubmit">
				<Card>
					<CardHeader>
						<CardTitle>{{ t('events.basicInfo') }}</CardTitle>
						<CardDescription>{{ t('events.basicInfoDescription') }}</CardDescription>
					</CardHeader>
					<CardContent class="space-y-4">
						<FormField
							v-slot="{ componentField }"
							name="type"
						>
							<FormItem>
								<FormLabel>{{ t('events.type') }}</FormLabel>
								<FormControl>
									<Select
										v-bind="componentField"
										:placeholder="t('events.selectType')"
									>
										<SelectTrigger>
											<SelectValue :placeholder="t('events.selectType')" />
										</SelectTrigger>
										<SelectContent>
											<SelectItem
												v-for="option in eventTypeSelectOptions"
												:key="option.value"
												:value="option.value"
											>
												{{ option.label }}
											</SelectItem>
										</SelectContent>
									</Select>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ componentField }"
							name="description"
						>
							<FormItem>
								<FormLabel>{{ t('events.description') }}</FormLabel>
								<FormControl>
									<Input v-bind="componentField" :placeholder="t('events.descriptionPlaceholder')" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<FormField
								v-slot="{ value, handleChange }"
								name="enabled"
							>
								<FormItem class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
									<div class="space-y-0.5">
										<FormLabel>{{ t('events.enabled') }}</FormLabel>
										<FormDescription>
											{{ t('events.enabledDescription') }}
										</FormDescription>
									</div>
									<FormControl>
										<Switch
											:checked="value"
											@update:checked="handleChange"
										/>
									</FormControl>
								</FormItem>
							</FormField>

							<FormField
								v-slot="{ value, handleChange }"
								name="onlineOnly"
							>
								<FormItem class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
									<div class="space-y-0.5">
										<FormLabel>{{ t('events.onlineOnly') }}</FormLabel>
										<FormDescription>
											{{ t('events.onlineOnlyDescription') }}
										</FormDescription>
									</div>
									<FormControl>
										<Switch
											:checked="value"
											@update:checked="handleChange"
										/>
									</FormControl>
								</FormItem>
							</FormField>
						</div>

						<div v-if="form.values.type === 'REDEMPTION_CREATED'">
							<FormField
								v-slot="{ componentField }"
								name="rewardId"
							>
								<FormItem>
									<FormLabel>{{ t('events.reward') }}</FormLabel>
									<FormControl>
										<TwitchRewardsSelector v-bind="componentField" />
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>

						<div v-if="form.values.type === 'COMMAND_USED'">
							<FormField
								v-slot="{ componentField }"
								name="commandId"
							>
								<FormItem>
									<FormLabel>{{ t('events.command') }}</FormLabel>
									<FormControl>
										<Select
											v-bind="componentField"
											:placeholder="t('events.selectCommand')"
										>
											<SelectTrigger>
												<SelectValue :placeholder="t('events.selectCommand')" />
											</SelectTrigger>
											<SelectContent>
												<SelectItem
													v-for="command in commands"
													:key="command.id"
													:value="command.id"
												>
													{{ command.name }}
												</SelectItem>
											</SelectContent>
										</Select>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>

						<div v-if="form.values.type === 'KEYWORD_USED'">
							<FormField
								v-slot="{ componentField }"
								name="keywordId"
							>
								<FormItem>
									<FormLabel>{{ t('events.keyword') }}</FormLabel>
									<FormControl>
										<Select
											v-bind="componentField"
											:placeholder="t('events.selectKeyword')"
										>
											<SelectTrigger>
												<SelectValue :placeholder="t('events.selectKeyword')" />
											</SelectTrigger>
											<SelectContent>
												<SelectItem
													v-for="keyword in keywords"
													:key="keyword.id"
													:value="keyword.id"
												>
													{{ keyword.text }}
												</SelectItem>
											</SelectContent>
										</Select>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>{{ t('events.operations') }}</CardTitle>
						<CardDescription>{{ t('events.operationsDescription') }}</CardDescription>
					</CardHeader>
					<CardContent>
						<div v-if="form.values.operations.length === 0" class="text-center py-8">
							<p class="text-muted-foreground mb-4">
								{{ t('events.noOperations') }}
							</p>
							<Button @click="addOperation">
								<PlusIcon class="mr-2 h-4 w-4" />
								{{ t('events.addOperation') }}
							</Button>
						</div>

						<div v-else>
							<Tabs v-model="selectedOperationTab">
								<div class="flex justify-between items-center mb-4">
									<TabsList>
										<TabsTrigger
											v-for="(operation, index) in form.values.operations"
											:key="index"
											:value="index"
										>
											{{ t('events.operation') }} {{ index + 1 }}
										</TabsTrigger>
									</TabsList>

									<div class="flex gap-2">
										<Button variant="outline" size="sm" @click="addOperation">
											<PlusIcon class="h-4 w-4" />
										</Button>
										<Button
											v-if="form.values.operations.length > 0"
											variant="destructive"
											size="sm"
											@click="removeOperation(selectedOperationTab)"
										>
											<Trash2 class="h-4 w-4" />
										</Button>
									</div>
								</div>

								<div v-for="(operation, opIndex) in form.values.operations" :key="opIndex">
									<TabsContent :value="opIndex" class="space-y-4">
										<FormField
											v-slot="{ componentField }"
											:name="`operations.${opIndex}.type`"
										>
											<FormItem>
												<FormLabel>{{ t('events.operationType') }}</FormLabel>
												<FormControl>
													<Select
														v-bind="componentField"
														:placeholder="t('events.selectOperationType')"
													>
														<SelectTrigger>
															<SelectValue :placeholder="t('events.selectOperationType')" />
														</SelectTrigger>
														<SelectContent>
															<SelectItem
																v-for="option in operationTypeSelectOptions"
																:key="option.value"
																:value="option.value"
															>
																{{ option.label }}
															</SelectItem>
														</SelectContent>
													</Select>
												</FormControl>
												<FormMessage />
											</FormItem>
										</FormField>

										<div v-if="operation.type && flatOperations[operation.type]?.haveInput">
											<FormField
												v-slot="{ componentField }"
												:name="`operations.${opIndex}.input`"
											>
												<FormItem>
													<FormLabel>{{ t('events.input') }}</FormLabel>
													<FormControl>
														<VariableInput v-bind="componentField" input-type="textarea" />
													</FormControl>
													<FormMessage />
												</FormItem>
											</FormField>
										</div>

										<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
											<FormField
												v-slot="{ componentField }"
												:name="`operations.${opIndex}.delay`"
											>
												<FormItem>
													<FormLabel>{{ t('events.delay') }}</FormLabel>
													<FormControl>
														<Input
															v-bind="componentField"
															type="number"
															min="0"
															:placeholder="t('events.delayPlaceholder')"
														/>
													</FormControl>
													<FormDescription>{{ t('events.delayDescription') }}</FormDescription>
													<FormMessage />
												</FormItem>
											</FormField>

											<FormField
												v-slot="{ componentField }"
												:name="`operations.${opIndex}.repeat`"
											>
												<FormItem>
													<FormLabel>{{ t('events.repeat') }}</FormLabel>
													<FormControl>
														<Input
															v-bind="componentField"
															type="number"
															min="0"
															:placeholder="t('events.repeatPlaceholder')"
														/>
													</FormControl>
													<FormDescription>{{ t('events.repeatDescription') }}</FormDescription>
													<FormMessage />
												</FormItem>
											</FormField>
										</div>

										<div v-if="operation.type && flatOperations[operation.type]?.additionalValues?.includes('useAnnounce')">
											<FormField
												v-slot="{ value, handleChange }"
												:name="`operations.${opIndex}.useAnnounce`"
											>
												<FormItem class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
													<div class="space-y-0.5">
														<FormLabel>{{ t('events.useAnnounce') }}</FormLabel>
														<FormDescription>
															{{ t('events.useAnnounceDescription') }}
														</FormDescription>
													</div>
													<FormControl>
														<Switch
															:checked="value"
															@update:checked="handleChange"
														/>
													</FormControl>
												</FormItem>
											</FormField>
										</div>

										<div v-if="operation.type && flatOperations[operation.type]?.additionalValues?.includes('timeoutTime')">
											<FormField
												v-slot="{ componentField }"
												:name="`operations.${opIndex}.timeoutTime`"
											>
												<FormItem>
													<FormLabel>{{ t('events.timeoutTime') }}</FormLabel>
													<FormControl>
														<Input
															v-bind="componentField"
															type="number"
															min="0"
															:placeholder="t('events.timeoutTimePlaceholder')"
														/>
													</FormControl>
													<FormDescription>{{ t('events.timeoutTimeDescription') }}</FormDescription>
													<FormMessage />
												</FormItem>
											</FormField>
										</div>

										<div v-if="operation.type && flatOperations[operation.type]?.additionalValues?.includes('timeoutMessage')">
											<FormField
												v-slot="{ componentField }"
												:name="`operations.${opIndex}.timeoutMessage`"
											>
												<FormItem>
													<FormLabel>{{ t('events.timeoutMessage') }}</FormLabel>
													<FormControl>
														<VariableInput v-bind="componentField" input-type="textarea" />
													</FormControl>
													<FormMessage />
												</FormItem>
											</FormField>
										</div>

										<div v-if="operation.type && flatOperations[operation.type]?.additionalValues?.includes('target')">
											<FormField
												v-slot="{ componentField }"
												:name="`operations.${opIndex}.target`"
											>
												<FormItem>
													<FormLabel>{{ t('events.target') }}</FormLabel>
													<FormControl>
														<Input
															v-bind="componentField"
															:placeholder="t('events.targetPlaceholder')"
														/>
													</FormControl>
													<FormDescription>{{ t('events.targetDescription') }}</FormDescription>
													<FormMessage />
												</FormItem>
											</FormField>
										</div>

										<FormField
											v-slot="{ value, handleChange }"
											:name="`operations.${opIndex}.enabled`"
										>
											<FormItem class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
												<div class="space-y-0.5">
													<FormLabel>{{ t('events.operationEnabled') }}</FormLabel>
													<FormDescription>
														{{ t('events.operationEnabledDescription') }}
													</FormDescription>
												</div>
												<FormControl>
													<Switch
														:checked="value"
														@update:checked="handleChange"
													/>
												</FormControl>
											</FormItem>
										</FormField>

										<!-- Filters section -->
										<div class="mt-6">
											<div class="flex justify-between items-center mb-4">
												<h3 class="text-lg font-medium">
													{{ t('events.filters') }}
												</h3>
												<Button variant="outline" size="sm" @click="addFilter(opIndex)">
													<PlusIcon class="h-4 w-4 mr-2" />
													{{ t('events.addFilter') }}
												</Button>
											</div>

											<div v-if="operation.filters.length === 0" class="text-center py-4 border rounded-md">
												<p class="text-muted-foreground">
													{{ t('events.noFilters') }}
												</p>
											</div>

											<div v-for="(filter, filterIndex) in operation.filters" :key="filterIndex" class="border rounded-md p-4 mb-4">
												<div class="flex justify-between items-center mb-4">
													<h4 class="font-medium">
														{{ t('events.filter') }} {{ filterIndex + 1 }}
													</h4>
													<Button
														variant="destructive"
														size="sm"
														@click="removeFilter(opIndex, filterIndex)"
													>
														<Trash2 class="h-4 w-4" />
													</Button>
												</div>

												<div class="space-y-4">
													<FormField
														v-slot="{ componentField }"
														:name="`operations.${opIndex}.filters.${filterIndex}.type`"
													>
														<FormItem>
															<FormLabel>{{ t('events.filterType') }}</FormLabel>
															<FormControl>
																<Select
																	v-bind="componentField"
																	:placeholder="t('events.selectFilterType')"
																>
																	<SelectTrigger>
																		<SelectValue :placeholder="t('events.selectFilterType')" />
																	</SelectTrigger>
																	<SelectContent>
																		<SelectItem value="EQUALS">
																			{{ t('events.filterEquals') }}
																		</SelectItem>
																		<SelectItem value="NOT_EQUALS">
																			{{ t('events.filterNotEquals') }}
																		</SelectItem>
																		<SelectItem value="CONTAINS">
																			{{ t('events.filterContains') }}
																		</SelectItem>
																		<SelectItem value="NOT_CONTAINS">
																			{{ t('events.filterNotContains') }}
																		</SelectItem>
																		<SelectItem value="STARTS_WITH">
																			{{ t('events.filterStartsWith') }}
																		</SelectItem>
																		<SelectItem value="ENDS_WITH">
																			{{ t('events.filterEndsWith') }}
																		</SelectItem>
																		<SelectItem value="GREATER_THAN">
																			{{ t('events.filterGreaterThan') }}
																		</SelectItem>
																		<SelectItem value="LESS_THAN">
																			{{ t('events.filterLessThan') }}
																		</SelectItem>
																		<SelectItem value="GREATER_THAN_OR_EQUALS">
																			{{ t('events.filterGreaterThanOrEquals') }}
																		</SelectItem>
																		<SelectItem value="LESS_THAN_OR_EQUALS">
																			{{ t('events.filterLessThanOrEquals') }}
																		</SelectItem>
																	</SelectContent>
																</Select>
															</FormControl>
															<FormMessage />
														</FormItem>
													</FormField>

													<FormField
														v-slot="{ componentField }"
														:name="`operations.${opIndex}.filters.${filterIndex}.left`"
													>
														<FormItem>
															<FormLabel>{{ t('events.filterLeft') }}</FormLabel>
															<FormControl>
																<VariableInput v-bind="componentField" />
															</FormControl>
															<FormDescription>{{ t('events.filterLeftDescription') }}</FormDescription>
															<FormMessage />
														</FormItem>
													</FormField>

													<FormField
														v-slot="{ componentField }"
														:name="`operations.${opIndex}.filters.${filterIndex}.right`"
													>
														<FormItem>
															<FormLabel>{{ t('events.filterRight') }}</FormLabel>
															<FormControl>
																<VariableInput v-bind="componentField" />
															</FormControl>
															<FormDescription>{{ t('events.filterRightDescription') }}</FormDescription>
															<FormMessage />
														</FormItem>
													</FormField>
												</div>
											</div>
										</div>
									</TabsContent>
								</div>
							</Tabs>
						</div>
					</CardContent>
				</Card>

				<div class="flex justify-end gap-4">
					<Button variant="outline" type="button" @click="goBack">
						{{ t('sharedTexts.cancel') }}
					</Button>
					<Button type="submit" :disabled="isSubmitting">
						{{ isSubmitting ? t('sharedTexts.saving') : (isNewEvent ? t('sharedTexts.create') : t('sharedTexts.save')) }}
					</Button>
				</div>
			</Form>
		</template>
	</PageLayout>
</template>
