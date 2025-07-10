<script setup lang="ts">
import { useFieldArray } from 'vee-validate'
import { computed } from 'vue'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
	EventType,
	KappagenOverlayAnimationSettings,
	KappagenOverlayAnimationsSettings,
	KappagenOverlayEvent,
} from '@/gql/graphql'
import { Plus, Trash2 } from 'lucide-vue-next'
import { flatEvents } from '@/features/events/constants/helpers.ts'
import { useKappagenApi } from '@/api'
import { Checkbox } from '@/components/ui/checkbox'

const {
	fields: events,
	push: addEvent,
	remove: removeEvent,
	update: updateEvent,
} = useFieldArray<KappagenOverlayEvent>('events')

const { fields: animations } = useFieldArray<KappagenOverlayAnimationsSettings>('animations')

const eventTypes = computed(() => {
	return Object.values(EventType).map((type) => ({
		value: type,
		label: type
			.replace(/_/g, ' ')
			.toLowerCase()
			.replace(/\b\w/g, (l) => l.toUpperCase()),
	}))
})

function createNewEvent() {
	addEvent({
		event: EventType.Follow,
		disabledAnimations: [],
		enabled: true,
	})
}

function handleCheckboxChange(
	eventIndex: number,
	animation: KappagenOverlayAnimationSettings,
	value: boolean
) {
	if (value) {
		updateEvent(eventIndex, {
			...events.value[eventIndex].value,
			disabledAnimations: events.value[eventIndex].value.disabledAnimations.filter(
				(anim) => anim !== animation.value.style
			),
		})
	} else {
		updateEvent(eventIndex, {
			...events.value[eventIndex].value,
			disabledAnimations: [
				...events.value[eventIndex].value.disabledAnimations,
				animation.value.style,
			],
		})
	}
}
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
				<Tabs
					:default-value="EventType.Follow"
					class="w-full flex flex-col lg:flex-row gap-4"
					orientation="vertical"
				>
					<TabsList class="w-full lg:w-[70%] lg:grid lg:grid-cols-1 flex flex-row flex-wrap">
						<TabsTrigger v-for="event of events" :key="event.key" :value="event.value.event">
							{{ flatEvents[event.value.event]?.name ?? event.value.event }}
						</TabsTrigger>
					</TabsList>

					<TabsContent
						v-for="(event, index) of events"
						:key="event.key"
						:value="event.value.event"
						class="w-auto"
					>
						<div class="grid grid-cols-1 gap-1">
							<div class="flex items-center gap-1" v-for="animation of animations">
								<Checkbox
									:checked="!event.value.disabledAnimations.includes(animation.value.style)"
									@update:checked="(v) => handleCheckboxChange(index, animation, v)"
								/>
								{{ animation.value.style }}
							</div>
						</div>
					</TabsContent>
				</Tabs>
			</CardContent>
		</Card>
	</div>
</template>
