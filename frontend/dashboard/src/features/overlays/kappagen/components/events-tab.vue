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
import { Check, ChevronsUpDown, SettingsIcon } from 'lucide-vue-next'
import { Switch } from '@/components/ui/switch'
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
	<div class="h-[40dvh]">
		<div class="grid grid-cols-1 gap-2 max-h-[40dvh] overflow-y-auto">
			<div
				v-for="event in events"
				:key="event.key"
				class="flex gap-2 justify-between items-center bg-stone-700/40 p-2 rounded-md"
			>
				<span>{{ flatEvents[event.value.event]?.name ?? event.value.event }}</span>

				<div class="flex items-center gap-2">
					<button
						class="p-1 border-border border rounded-md bg-zinc-600/50 hover:bg-zinc-600/30 transition-colors"
					>
						<SettingsIcon class="size-4" />
					</button>
					<Switch />
				</div>
			</div>
		</div>
	</div>

	<!--	<div class="flex gap-2 flex-col" v-if="selectedEventIndex >= 0">-->
	<!--		<h3>Animations</h3>-->

	<!--		<div class="grid grid-cols-1 xl:grid-cols-2 gap-2">-->
	<!--			<div-->
	<!--				class="flex items-center gap-2 rounded-md bg-background/60 p-2"-->
	<!--				v-for="animation of animations"-->
	<!--			>-->
	<!--				<Checkbox-->
	<!--					:checked="-->
	<!--						!events.at(selectedEventIndex).value.disabledAnimations.includes(animation.value.style)-->
	<!--					"-->
	<!--					@update:checked="(v) => handleCheckboxChange(animation, v)"-->
	<!--				/>-->
	<!--				{{ animation.value.style }}-->
	<!--			</div>-->
	<!--		</div>-->
	<!--	</div>-->
</template>
