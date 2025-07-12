<script setup lang="ts">
import { useFieldArray } from 'vee-validate'
import { computed, ref } from 'vue'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
	EventType,
	KappagenOverlayAnimationSettings,
	KappagenOverlayAnimationsSettings,
	KappagenOverlayEvent,
} from '@/gql/graphql'

import { flatEvents } from '@/features/events/constants/helpers.ts'
import { Checkbox } from '@/components/ui/checkbox'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Check, ChevronsUpDown } from 'lucide-vue-next'
const { fields: events, update: updateEvent } = useFieldArray<KappagenOverlayEvent>('events')

const { fields: animations } = useFieldArray<KappagenOverlayAnimationsSettings>('animations')

function handleCheckboxChange(animation: KappagenOverlayAnimationSettings, value: boolean) {
	if (value) {
		updateEvent(selectedEventIndex.value, {
			...events.value[selectedEventIndex.value].value,
			disabledAnimations: events.value[selectedEventIndex.value].value.disabledAnimations.filter(
				(anim) => anim !== animation.value.style
			),
		})
	} else {
		updateEvent(selectedEventIndex.value, {
			...events.value[selectedEventIndex.value].value,
			disabledAnimations: [
				...events.value[selectedEventIndex.value].value.disabledAnimations,
				animation.value.style,
			],
		})
	}
}

const selectedEvent = ref<string>(EventType.Follow)
const selectedEventIndex = computed(() => {
	return events.value.findIndex((event) => event.value.event === selectedEvent.value)
})
</script>

<template>
	<div class="space-y-6">
		<Card>
			<CardHeader>
				<CardTitle>Event Configuration</CardTitle>
				<CardDescription>
					Configure which events trigger Kappagen animations and which animations to disable for
					specific events
				</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4 w-full">
				<div class="flex flex-col gap-2">
					<Popover>
						<PopoverTrigger as-child>
							<Button variant="outline" role="combobox" class="w-[400px] justify-between">
								{{ flatEvents[selectedEvent]?.name ?? selectedEvent }}
								<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
							</Button>
						</PopoverTrigger>
						<PopoverContent v-model="selectedEvent" class="w-[400px] p-0">
							<Command>
								<CommandInput placeholder="Search event..." />
								<CommandEmpty>Nothing found.</CommandEmpty>
								<CommandList>
									<CommandGroup>
										<CommandItem
											v-for="event in events"
											:key="event.key"
											:value="event.value.event"
											@select="selectedEvent = event.value.event"
										>
											{{ flatEvents[event.value.event]?.name ?? event.value.event }}
											<Check
												class="ml-auto h-4 w-4"
												:class="[selectedEvent === event.value.event ? 'opacity-100' : 'opacity-0']"
											/>
										</CommandItem>
									</CommandGroup>
								</CommandList>
							</Command>
						</PopoverContent>
					</Popover>
				</div>

				<div class="flex gap-2 flex-col" v-if="selectedEventIndex >= 0">
					<h3>Animations</h3>

					<div class="grid grid-cols-1 md:grid-cols-2 gap-2">
						<div
							class="flex items-center gap-2 rounded-md bg-background/60 p-2"
							v-for="animation of animations"
						>
							<Checkbox
								:checked="
									!events
										.at(selectedEventIndex)
										.value.disabledAnimations.includes(animation.value.style)
								"
								@update:checked="(v) => handleCheckboxChange(animation, v)"
							/>
							{{ animation.value.style }}
						</div>
					</div>
				</div>
			</CardContent>
		</Card>
	</div>
</template>
