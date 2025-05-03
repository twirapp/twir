<script setup lang="ts">
import { CheckIcon, ChevronsUpDownIcon } from 'lucide-vue-next'
import { useField } from 'vee-validate'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useCommandsApi } from '@/api/commands/commands'
import { EventType } from '@/api/events.ts'
import { useKeywordsApi } from '@/api/keywords'
import { EventsOptions } from '@/components/events/events.ts'
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
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import {
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { cn } from '@/lib/utils'

const { t } = useI18n()
// Fetch commands and keywords for selectors
const commandsApi = useCommandsApi()
const { data: commandsData } = commandsApi.useQueryCommands()
const commands = computed(() => commandsData.value?.commands || [])

const keywordsApi = useKeywordsApi()
const { data: keywordsData } = keywordsApi.useQueryKeywords()
const keywords = computed(() => keywordsData.value?.keywords || [])
const { value: currentEventType, setValue: setCurrentEventType } = useField<EventType>('type')

const typeSelectOptions = Object.values(EventsOptions).map<{
	isGroup: boolean
	name: string
	value?: EventType
	childrens: Array<{ name: string, value: EventType }>
}>((value) => {
	if (value.childrens) {
		return {
			isGroup: true,
			name: value.name,
			childrens: Object.values(value.childrens!).map(c => ({ name: c.name, value: c.enumValue! })),
		}
	}

	return {
		isGroup: false,
		name: value.name,
		value: value.enumValue,
		childrens: [],
	}
})
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle>{{ t('events.basicInfo') }}</CardTitle>
			<CardDescription>{{ t('events.basicInfoDescription') }}</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			<FormField
				name="type"
			>
				<FormItem class="flex flex-col">
					<FormLabel>{{ t('events.type') }}</FormLabel>
					<FormControl>
						<Popover>
							<PopoverTrigger as-child>
								<FormControl>
									<Button
										variant="outline"
										role="combobox"
										:class="cn('w-[200px] justify-between', !currentEventType && 'text-muted-foreground')"
									>
										{{ currentEventType ? getEventName(currentEventType) ? 'Select language...' }}
										<ChevronsUpDownIcon class="ml-2 h-4 w-4 shrink-0 opacity-50" />
									</Button>
								</FormControl>
							</PopoverTrigger>
							<PopoverContent class="w-[200px] p-0">
								<Command>
									<CommandInput placeholder="Search language..." />
									<CommandEmpty>Nothing found.</CommandEmpty>
									<CommandList>
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
													@select="() => {
														setCurrentEventType(event.value)
													}"
												>
													{{ event.name }}
													<CheckIcon
														:class="cn('ml-auto h-4 w-4', currentEventType === event.value ? 'opacity-100' : 'opacity-0')"
													/>
												</CommandItem>
											</CommandGroup>

											<CommandItem
												v-else
												:key="selectOption.value!"
												:value="selectOption.value!"
												@select="() => {
													setCurrentEventType(selectOption.value!)
												}"
											>
												{{ selectOption.name }}
												<CheckIcon :class="cn('ml-auto h-4 w-4', currentEventType === selectOption.value ? 'opacity-100' : 'opacity-0')" />
											</CommandItem>
										</template>
										<CommandGroup>
										</CommandGroup>
									</CommandList>
								</Command>
							</PopoverContent>
						</Popover>
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
							<FormLabel>{{ t('sharedTexts.enabled') }}</FormLabel>
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

			<div v-if="currentEventType === EventType.RedemptionCreated">
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

			<div v-if="currentEventType === EventType.CommandUsed">
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

			<div v-if="currentEventType === EventType.KeywordMatched">
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
</template>
