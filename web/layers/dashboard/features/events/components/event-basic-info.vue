<script setup lang="ts">
import { useField } from 'vee-validate'
import { computed, ref } from 'vue'

import { useCommandsApi } from '~~/layers/dashboard/api/commands/commands'

// oxlint-disable-next-line consistent-type-imports
import { EventType } from '~~/layers/dashboard/api/events.js'
import { useKeywordsApi } from '~~/layers/dashboard/api/keywords'
import TwitchRewardsSelector from '~~/layers/dashboard/components/rewardsSelector.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Switch } from '@/components/ui/switch'
import { EventsOptions } from '~~/layers/dashboard/features/events/constants/events.js'
import { getEventName } from '~~/layers/dashboard/features/events/constants/helpers.js'
import PlatformSelector from '~~/layers/dashboard/components/platform-selector.vue'
import { cn } from '~~/layers/dashboard/lib/utils'

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
	<Card>
		<CardHeader>
			<CardTitle>General</CardTitle>
		</CardHeader>
		<CardContent class="space-y-4">
			<FormField name="type">
				<FormItem class="flex flex-col">
					<FormLabel>{{ t('events.type') }}</FormLabel>
					<FormControl>
						<Popover v-model:open="opened">
							<PopoverTrigger as-child>
								<FormControl>
									<Button
										type="button"
										variant="outline"
										role="combobox"
										:class="
											cn('w-100 justify-between', !currentEventType && 'text-muted-foreground')
										"
									>
										{{ currentEventType ? getEventName(currentEventType) : 'Select...' }}
										<Icon name="lucide:chevrons-up-down" class="ml-2 h-4 w-4 shrink-0 opacity-50" />
									</Button>
								</FormControl>
							</PopoverTrigger>
							<PopoverContent class="w-50 p-0">
								<Command>
									<CommandInput placeholder="Search trigger..." />
									<CommandList>
										<CommandEmpty>Nothing found.</CommandEmpty>
										<template v-for="selectOption of typeSelectOptions">
											<CommandGroup
												v-if="selectOption.isGroup"
												:key="selectOption.name"
												:heading="selectOption.name"
											>
												<CommandItem
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
													<Icon name="lucide:check"
														:class="
															cn(
																'ml-auto h-4 w-4',
																currentEventType === event.value ? 'opacity-100' : 'opacity-0'
															)
														" />
												</CommandItem>
											</CommandGroup>

											<CommandGroup v-else>
												<CommandItem
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
													<Icon name="lucide:check"
														:class="
															cn(
																'ml-auto h-4 w-4',
																currentEventType === selectOption.value
																	? 'opacity-100'
																	: 'opacity-0'
															)
														" />
												</CommandItem>
											</CommandGroup>
										</template>
									</CommandList>
								</Command>
							</PopoverContent>
						</Popover>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="description">
				<FormItem>
					<FormLabel>{{ t('events.description') }}</FormLabel>
					<FormControl>
						<Input v-bind="componentField" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ value, handleChange }" name="platforms">
				<FormItem>
					<FormLabel>{{ t('sharedTexts.platforms') }}</FormLabel>
					<FormControl>
						<PlatformSelector :model-value="value" @update:model-value="handleChange" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<FormField v-slot="{ value, handleChange }" name="enabled">
					<FormItem
						class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-xs"
					>
						<div class="space-y-0.5">
							<FormLabel>{{ t('sharedTexts.enabled') }}</FormLabel>
						</div>
						<FormControl>
							<Switch :model-value="value" @update:model-value="handleChange" />
						</FormControl>
					</FormItem>
				</FormField>

				<FormField v-slot="{ value, handleChange }" name="onlineOnly">
					<FormItem
						class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-xs"
					>
						<div class="space-y-0.5">
							<FormLabel>{{ t('events.onlineOnly') }}</FormLabel>
						</div>
						<FormControl>
							<Switch :model-value="value" @update:model-value="handleChange" />
						</FormControl>
					</FormItem>
				</FormField>
			</div>

			<div v-if="currentEventType === EventType.RedemptionCreated">
				<FormField v-slot="{ componentField }" name="rewardId">
					<FormItem>
						<FormLabel>{{ t('events.reward') }}</FormLabel>
						<FormControl>
							<TwitchRewardsSelector v-bind="componentField" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<div v-if="currentEventType === EventType.CommandUsed">
				<FormField name="commandId">
					<FormItem class="flex flex-col gap-2">
						<FormLabel>Command</FormLabel>
						<FormControl>
							<Popover>
								<PopoverTrigger as-child>
									<FormControl>
										<Button
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
											<Icon name="lucide:chevrons-up-down" class="ml-2 h-4 w-4 shrink-0 opacity-50" />
										</Button>
									</FormControl>
								</PopoverTrigger>
								<PopoverContent class="w-[200px] p-0">
									<Command>
										<CommandInput placeholder="Search command..." />
										<CommandEmpty>Nothing found.</CommandEmpty>
										<CommandList>
											<CommandGroup>
												<CommandItem
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
													<Icon name="lucide:check"
														:class="
															cn(
																'ml-auto h-4 w-4',
																currentCommandId === command.id ? 'opacity-100' : 'opacity-0'
															)
														" />
												</CommandItem>
											</CommandGroup>
										</CommandList>
									</Command>
								</PopoverContent>
							</Popover>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<div v-if="currentEventType === EventType.KeywordMatched">
				<FormField name="keywordId">
					<FormItem>
						<FormLabel>{{ t('events.keyword') }}</FormLabel>
						<FormControl class="flex flex-col gap-2">
							<Popover>
								<PopoverTrigger as-child>
									<FormControl>
										<Button
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
											<Icon name="lucide:chevrons-up-down" class="ml-2 h-4 w-4 shrink-0 opacity-50" />
										</Button>
									</FormControl>
								</PopoverTrigger>
								<PopoverContent class="w-[200px] p-0">
									<Command>
										<CommandInput placeholder="Search Keyword..." />
										<CommandEmpty>Nothing found.</CommandEmpty>
										<CommandList>
											<CommandGroup>
												<CommandItem
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
													<Icon name="lucide:check"
														:class="
															cn(
																'ml-auto h-4 w-4',
																currentKeywordId === keyword.id ? 'opacity-100' : 'opacity-0'
															)
														" />
												</CommandItem>
											</CommandGroup>
										</CommandList>
									</Command>
								</PopoverContent>
							</Popover>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>
		</CardContent>
	</Card>
</template>
