<script setup lang="ts">
import { SettingsIcon } from 'lucide-vue-next'
import { useFieldArray } from 'vee-validate'

import type {
	KappagenOverlayAnimationStyle,
	KappagenOverlayAnimationsSettings,
	KappagenOverlayEvent,
} from '@/gql/graphql'

import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Switch } from '@/components/ui/switch'
import { flatEvents } from '@/features/events/constants/helpers.ts'

const { fields: events, update: updateEvent } = useFieldArray<KappagenOverlayEvent>('events')

const { fields: animations } = useFieldArray<KappagenOverlayAnimationsSettings>('animations')

function handleDisabledAnimationCheckboxChange(
	eventIndex: number,
	animationStyle: KappagenOverlayAnimationStyle
) {
	const ignored = events.value[eventIndex].value.disabledAnimations.includes(animationStyle)

	if (ignored) {
		updateEvent(eventIndex, {
			...events.value[eventIndex].value,
			disabledAnimations: events.value[eventIndex].value.disabledAnimations.filter(
				(anim) => anim !== animationStyle
			),
		})
	} else {
		updateEvent(eventIndex, {
			...events.value[eventIndex].value,
			disabledAnimations: [...events.value[eventIndex].value.disabledAnimations, animationStyle],
		})
	}
}

function handleEventEnabledChange(index: number, value: boolean) {
	const event = events.value[index]
	if (event) {
		updateEvent(index, { ...event.value, enabled: value })
	}
}
</script>

<template>
	<div class="h-[40dvh]">
		<div class="grid grid-cols-1 gap-2 max-h-[40dvh] overflow-y-auto">
			<div
				v-for="(event, eventIndex) in events"
				:key="event.key"
				class="flex gap-2 justify-between items-center bg-stone-700/40 p-2 rounded-md"
			>
				<span>{{ flatEvents[event.value.event]?.name ?? event.value.event }}</span>

				<div class="flex items-center gap-2">
					<Popover>
						<PopoverTrigger as-child>
							<button
								class="p-1 border-border border rounded-md bg-zinc-600/50 hover:bg-zinc-600/30 transition-colors"
								type="button"
							>
								<SettingsIcon class="size-4" />
							</button>
						</PopoverTrigger>
						<PopoverContent class="w-96 bg-zinc-800/60 backdrop-blur-md">
							<div class="grid grid-cols-2 gap-2">
								<div v-for="animation of animations" :key="animation.value.style">
									<div
										:key="animation.key"
										class="flex items-center gap-2 rounded-md bg-background/60 p-2"
									>
										<Checkbox
											:id="`animation-${animation.value.style}`"
											:checked="!event.value.disabledAnimations.includes(animation.value.style)"
											@update:checked="
												handleDisabledAnimationCheckboxChange(eventIndex, animation.value.style)
											"
										/>
										<Label :for="`animation-${animation.value.style}`" class="cursor-pointer">
											{{ animation.value.style }}
										</Label>
									</div>
								</div>
							</div>
						</PopoverContent>
					</Popover>
					<Switch
						:checked="event.value.enabled"
						@update:checked="(v: boolean) => handleEventEnabledChange(eventIndex, v)"
					/>
				</div>
			</div>
		</div>
	</div>
</template>
