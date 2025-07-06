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
import { EventType } from '@/gql/graphql'
import { Plus, Trash2 } from 'lucide-vue-next'

const { fields: events, push: addEvent, remove: removeEvent } = useFieldArray('events')

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
			<CardContent class="space-y-4">
				<div class="space-y-4">
					<div
						v-for="(event, index) in events"
						:key="event.key"
						class="border rounded-lg p-4 space-y-4"
					>
						<div class="flex items-center justify-between">
							<h4 class="font-medium">Event {{ index + 1 }}</h4>
							<Button type="button" variant="destructive" size="sm" @click="removeEvent(index)">
								<Trash2 class="h-4 w-4" />
							</Button>
						</div>

						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<FormField :name="`events[${index}].event`">
								<FormItem>
									<FormLabel>Event Type</FormLabel>
									<Select v-model="event.value.event">
										<SelectTrigger>
											<SelectValue placeholder="Select event type" />
										</SelectTrigger>
										<SelectContent>
											<SelectItem
												v-for="eventType in eventTypes"
												:key="eventType.value"
												:value="eventType.value"
											>
												{{ eventType.label }}
											</SelectItem>
										</SelectContent>
									</Select>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField :name="`events[${index}].enabled`">
								<FormItem
									class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm"
								>
									<div class="space-y-0.5">
										<FormLabel>Enabled</FormLabel>
										<div class="text-[0.8rem] text-muted-foreground">Enable this event trigger</div>
									</div>
									<Switch v-model:checked="event.value.enabled" />
									<FormMessage />
								</FormItem>
							</FormField>
						</div>

						<FormField :name="`events[${index}].disabledAnimations`">
							<FormItem>
								<FormLabel>Disabled Animations</FormLabel>
								<Input
									v-model="event.value.disabledAnimations"
									placeholder="Comma-separated list of disabled animations"
								/>
								<div class="text-[0.8rem] text-muted-foreground">
									Enter animation names separated by commas to disable them for this event
								</div>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>
				</div>

				<Button type="button" variant="outline" @click="createNewEvent">
					<Plus class="h-4 w-4 mr-2" />
					Add Event
				</Button>
			</CardContent>
		</Card>
	</div>
</template>
