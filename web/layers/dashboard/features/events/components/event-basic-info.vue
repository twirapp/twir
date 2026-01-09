<script setup lang="ts">
import { CheckIcon, ChevronsUpDownIcon } from 'lucide-vue-next'
import { useField } from 'vee-validate'
import { computed, ref } from 'vue'


import { useCommandsApi } from '#layers/dashboard/api/commands/commands'

// oxlint-disable-next-line consistent-type-imports
import { EventType } from '#layers/dashboard/api/events.ts'
import { useKeywordsApi } from '#layers/dashboard/api/keywords'
import TwitchRewardsSelector from '#layers/dashboard/components/rewardsSelector.vue'







import { EventsOptions } from '~/features/events/constants/events.ts'
import { getEventName } from '~/features/events/constants/helpers.ts'
import { cn } from '~/lib/utils'

const { t } = useI18n()
// Fetch commands and keywords for selectors
const commandsApi = useCommandsApi()
const { data: commandsData } = commandsApi.useQueryCommands()
const commands = computed(() => commandsData.value?.commands || [])

const keywordsApi = useKeywordsApi()
const { data: keywordsData } = keywordsApi.useQueryKeywords()
const keywords = computed(() => keywordsData.value?.keywords || [])

const { value: currentEventType, setValue: setCurrentEventType } = useField<EventType>('type')
const { value: currentCommandId, setValue: setCurrentCommandId } = useField<string>('commandId')
const { value: currentKeywordId, setValue: setCurrentKeywordId } = useField<string>('keywordId')

const typeSelectOptions = Object.values(EventsOptions).map<{
	isGroup: boolean
	name: string
	value?: EventType
	childrens: Array<{ name: string; value: EventType }>
}>((value) => {
	if (value.childrens) {
		return {
			isGroup: true,
			name: value.name,
			childrens: Object.values(value.childrens!).map((c) => ({
				name: c.name,
				value: c.enumValue!,
			})),
		}
	}

	return {
		isGroup: false,
		name: value.name,
		value: value.enumValue,
		childrens: [],
	}
})

const opened = ref(false)
</script>

<template>
	<UiCard>
		<UiCardHeader>
			<UiCardTitle>General</UiCardTitle>
		</UiCardHeader>
		<UiCardContent class="space-y-4">
			<UiFormField name="type">
				<UiFormItem class="flex flex-col">
					<UiFormLabel>{{ t('events.type') }}</UiFormLabel>
					<UiFormControl>
						<UiPopover v-model:open="opened">
							<UiPopoverTrigger as-child>
								<UiFormControl>
									<UiButton
										type="button"
										variant="outline"
										role="combobox"
										:class="
											cn('w-100 justify-between', !currentEventType && 'text-muted-foreground')
										"
									>
										{{ currentEventType ? getEventName(currentEventType) : 'Select...' }}
										<ChevronsUpDownIcon class="ml-2 h-4 w-4 shrink-0 opacity-50" />
									</UiButton>
								</UiFormControl>
							</UiPopoverTrigger>
							<UiPopoverContent class="w-50 p-0">
								<UiCommand>
									<UiCommandInput placeholder="Search trigger..." />
									<UiCommandList>
										<UiCommandEmpty>Nothing found.</UiCommandEmpty>
										<template v-for="selectOption of typeSelectOptions">
											<UiCommandGroup
												v-if="selectOption.isGroup"
												:key="selectOption.name"
												:heading="selectOption.name"
											>
												<UiCommandItem
													v-for="event of selectOption.childrens"
													:key="event.value"
													:value="event.value"
													@select="
														() => {
															setCurrentEventType(event.value)
															opened = false
														}
													"
												>
													{{ event.name }}
													<CheckIcon
														:class="
															cn(
																'ml-auto h-4 w-4',
																currentEventType === event.value ? 'opacity-100' : 'opacity-0'
															)
														"
													/>
												</UiCommandItem>
											</UiCommandGroup>

											<UiCommandGroup v-else>
												<UiCommandItem
													:key="selectOption.value!"
													:value="selectOption.value!"
													@select="
														() => {
															setCurrentEventType(selectOption.value!)
															opened = false
														}
													"
												>
													{{ selectOption.name }}
													<CheckIcon
														:class="
															cn(
																'ml-auto h-4 w-4',
																currentEventType === selectOption.value
																	? 'opacity-100'
																	: 'opacity-0'
															)
														"
													/>
												</UiCommandItem>
											</UiCommandGroup>
										</template>
									</UiCommandList>
								</UiCommand>
							</UiPopoverContent>
						</UiPopover>
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField v-slot="{ componentField }" name="description">
				<UiFormItem>
					<UiFormLabel>{{ t('events.description') }}</UiFormLabel>
					<UiFormControl>
						<UiInput v-bind="componentField" />
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<UiFormField v-slot="{ value, handleChange }" name="enabled">
					<UiFormItem
						class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-xs"
					>
						<div class="space-y-0.5">
							<UiFormLabel>{{ t('sharedTexts.enabled') }}</UiFormLabel>
						</div>
						<UiFormControl>
							<UiSwitch :model-value="value" @update:model-value="handleChange" />
						</UiFormControl>
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ value, handleChange }" name="onlineOnly">
					<UiFormItem
						class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-xs"
					>
						<div class="space-y-0.5">
							<UiFormLabel>{{ t('events.onlineOnly') }}</UiFormLabel>
						</div>
						<UiFormControl>
							<UiSwitch :model-value="value" @update:model-value="handleChange" />
						</UiFormControl>
					</UiFormItem>
				</UiFormField>
			</div>

			<div v-if="currentEventType === EventType.RedemptionCreated">
				<UiFormField v-slot="{ componentField }" name="rewardId">
					<UiFormItem>
						<UiFormLabel>{{ t('events.reward') }}</UiFormLabel>
						<UiFormControl>
							<TwitchRewardsSelector v-bind="componentField" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>

			<div v-if="currentEventType === EventType.CommandUsed">
				<UiFormField name="commandId">
					<UiFormItem class="flex flex-col gap-2">
						<UiFormLabel>Command</UiFormLabel>
						<UiFormControl>
							<UiPopover>
								<UiPopoverTrigger as-child>
									<UiFormControl>
										<UiButton
											type="button"
											variant="outline"
											role="combobox"
											:class="
												cn(
													'w-[200px] justify-between',
													!currentCommandId && 'text-muted-foreground'
												)
											"
										>
											{{
												currentCommandId
													? commands.find((c) => c.id === currentCommandId)?.name
													: 'Select command'
											}}
											<ChevronsUpDownIcon class="ml-2 h-4 w-4 shrink-0 opacity-50" />
										</UiButton>
									</UiFormControl>
								</UiPopoverTrigger>
								<UiPopoverContent class="w-[200px] p-0">
									<UiCommand>
										<UiCommandInput placeholder="Search command..." />
										<UiCommandEmpty>Nothing found.</UiCommandEmpty>
										<UiCommandList>
											<UiCommandGroup>
												<UiCommandItem
													v-for="command in commands"
													:key="command.id"
													:value="command.name"
													@select="
														() => {
															setCurrentCommandId(command.id)
														}
													"
												>
													{{ command.name }}
													<CheckIcon
														:class="
															cn(
																'ml-auto h-4 w-4',
																currentCommandId === command.id ? 'opacity-100' : 'opacity-0'
															)
														"
													/>
												</UiCommandItem>
											</UiCommandGroup>
										</UiCommandList>
									</UiCommand>
								</UiPopoverContent>
							</UiPopover>
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>

			<div v-if="currentEventType === EventType.KeywordMatched">
				<UiFormField name="keywordId">
					<UiFormItem>
						<UiFormLabel>{{ t('events.keyword') }}</UiFormLabel>
						<UiFormControl class="flex flex-col gap-2">
							<UiPopover>
								<UiPopoverTrigger as-child>
									<UiFormControl>
										<UiButton
											type="button"
											variant="outline"
											role="combobox"
											:class="
												cn(
													'w-[200px] justify-between',
													!currentKeywordId && 'text-muted-foreground'
												)
											"
										>
											{{
												currentKeywordId
													? keywords.find((c) => c.id === currentKeywordId)?.text
													: 'Select keyword'
											}}
											<ChevronsUpDownIcon class="ml-2 h-4 w-4 shrink-0 opacity-50" />
										</UiButton>
									</UiFormControl>
								</UiPopoverTrigger>
								<UiPopoverContent class="w-[200px] p-0">
									<UiCommand>
										<UiCommandInput placeholder="Search Keyword..." />
										<UiCommandEmpty>Nothing found.</UiCommandEmpty>
										<UiCommandList>
											<UiCommandGroup>
												<UiCommandItem
													v-for="keyword in keywords"
													:key="keyword.id"
													:value="keyword.text"
													@select="
														() => {
															setCurrentKeywordId(keyword.id)
														}
													"
												>
													{{ keyword.text }}
													<CheckIcon
														:class="
															cn(
																'ml-auto h-4 w-4',
																currentKeywordId === keyword.id ? 'opacity-100' : 'opacity-0'
															)
														"
													/>
												</UiCommandItem>
											</UiCommandGroup>
										</UiCommandList>
									</UiCommand>
								</UiPopoverContent>
							</UiPopover>
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>
		</UiCardContent>
	</UiCard>
</template>
