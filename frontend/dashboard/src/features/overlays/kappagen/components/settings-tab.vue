<script setup lang="ts">
import { Button } from '@/components/ui/button'

import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Plus, X } from 'lucide-vue-next'
import { useFieldArray } from 'vee-validate'
import CommandButton from '@/features/commands/ui/command-button.vue'
import { SliderRange, SliderRoot, SliderThumb, SliderTrack } from 'radix-vue'
import { Checkbox } from '@/components/ui/checkbox'
import { KappagenEmojiStyle } from '@/gql/graphql.ts'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'

const {
	fields: excludedEmotes,
	push: addExcludedEmote,
	remove: removeExcludedEmote,
} = useFieldArray('excludedEmotes')

const emotesCheckboxes = [
	{
		key: 'ffz',
		label: 'FFZ Emotes',
	},
	{
		key: 'bttv',
		label: 'BetterTTV Emotes',
	},
	{
		key: 'sevenTv',
		label: '7tv Emotes',
	},
]

const availableEmojiStyles = Object.entries(KappagenEmojiStyle).map(([key, value]) => {
	return {
		value: value,
		label: key.charAt(0).toUpperCase() + key.slice(1).toLowerCase(),
	}
})
</script>

<template>
	<div class="bg-stone-700/40 rounded-md w-full max-h-[40dvh] overflow-y-auto">
		<div class="flex flex-col gap-4 p-2">
			<div class="flex justify-between items-center">
				<span>Command for manual spawn</span>

				<CommandButton name="kappagen" title="" />
			</div>

			<FormField name="enableSpawn" v-slot="{ value, handleChange }">
				<FormItem class="flex flex-row items-center justify-between">
					<div class="flex flex-col w-full">
						<FormLabel>Enable Spawn</FormLabel>
						<div class="text-[0.8rem] text-muted-foreground">Spawn emote on each chat message</div>
					</div>

					<Switch :checked="value" @update:checked="handleChange" />
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField name="enableRave" v-slot="{ value, handleChange }">
				<FormItem class="flex flex-row items-center justify-between">
					<div class="flex flex-col w-full">
						<FormLabel>Enable Rave Mode</FormLabel>
						<div class="text-[0.8rem] text-muted-foreground">
							Enable special rave animations and effects
						</div>
					</div>
					<Switch :checked="value" @update:checked="handleChange" />
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField name="size.rationSmall" v-slot="{ value, handleChange, field }">
				<FormItem class="flex flex-col">
					<FormLabel>Small Emote Ratio</FormLabel>
					<div class="flex flex-row items-center gap-2">
						<div class="bg-stone-700/40 p-1 w-[80px] rounded-md text-center">{{ value }}</div>
						<SliderRoot
							:model-value="[value]"
							@update:model-value="([v]) => handleChange(v)"
							class="relative flex items-center select-none touch-none w-full h-5"
							:min="0.02"
							:max="0.07"
							:step="0.01"
						>
							<SliderTrack class="bg-gray-500 relative grow rounded-full h-[3px]">
								<SliderRange class="absolute bg-indigo-400 rounded-full h-full" />
							</SliderTrack>
							<SliderThumb class="block w-5 h-5 bg-white cursor-pointer rounded-[10px]" />
						</SliderRoot>
					</div>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField name="size.rationNormal" v-slot="{ value, handleChange }">
				<FormItem class="flex flex-col">
					<FormLabel>Normal emote ratio</FormLabel>
					<div class="flex flex-row items-center gap-2">
						<div class="bg-stone-700/40 p-1 w-[80px] rounded-md text-center">{{ value }}</div>
						<SliderRoot
							:model-value="[value]"
							@update:model-value="([v]) => handleChange(v)"
							class="relative flex items-center select-none touch-none w-full h-5"
							:min="0.05"
							:max="0.15"
							:step="0.01"
						>
							<SliderTrack class="bg-gray-500 relative grow rounded-full h-[3px]">
								<SliderRange class="absolute bg-indigo-400 rounded-full h-full" />
							</SliderTrack>
							<SliderThumb class="block w-5 h-5 bg-white cursor-pointer rounded-[10px]" />
						</SliderRoot>
					</div>
					<FormMessage />
				</FormItem>
			</FormField>

			<div class="space-y-2">
				<span>External emotes</span>
				<div class="grid grid-cols-2 gap-2">
					<FormField
						v-for="service of emotesCheckboxes"
						:name="`emotes.${service.key}Enabled`"
						v-slot="{ value, handleChange }"
						:key="service.key"
					>
						<FormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<FormControl class="flex items-center !m-0">
								<Checkbox :checked="value" @update:checked="handleChange" />
							</FormControl>
							<FormLabel class="select-none cursor-pointer !m-0">{{ service.label }}</FormLabel>
						</FormItem>
					</FormField>
				</div>
			</div>

			<FormField v-slot="{ componentField }" name="emotes.emojiStyle">
				<FormItem>
					<FormLabel>Emoji style</FormLabel>

					<Select v-bind="componentField">
						<FormControl>
							<SelectTrigger>
								<SelectValue placeholder="Select emoji style" />
							</SelectTrigger>
						</FormControl>
						<SelectContent>
							<SelectGroup>
								<SelectItem
									v-for="item of availableEmojiStyles"
									:key="item.value"
									:value="item.value"
								>
									{{ item.label }}
								</SelectItem>
							</SelectGroup>
						</SelectContent>
					</Select>
					<FormMessage />
				</FormItem>
			</FormField>

			<div class="flex flex-col gap-1">
				<FormField v-slot="{ value, handleChange }" name="emotes.time">
					<FormItem class="flex flex-col">
						<FormLabel>Max emote time</FormLabel>
						<div class="flex flex-row items-center gap-2">
							<div class="bg-stone-700/40 p-1 w-[80px] rounded-md text-center">{{ value }}s</div>
							<SliderRoot
								:model-value="[value]"
								@update:model-value="([v]) => handleChange(v)"
								class="relative flex items-center select-none touch-none w-full h-5"
								:min="1"
								:max="15"
								:step="1"
							>
								<SliderTrack class="bg-gray-500 relative grow rounded-full h-[3px]">
									<SliderRange class="absolute bg-indigo-400 rounded-full h-full" />
								</SliderTrack>
								<SliderThumb class="block w-5 h-5 bg-white cursor-pointer rounded-[10px]" />
							</SliderRoot>
						</div>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<div class="flex flex-col gap-1">
				<FormField v-slot="{ value, handleChange }" name="emotes.max">
					<FormItem class="flex flex-col">
						<FormLabel>Max emotes count</FormLabel>
						<div class="flex flex-row items-center gap-2">
							<div class="bg-stone-700/40 p-1 w-[80px] rounded-md text-center">{{ value }}</div>
							<SliderRoot
								:model-value="[value]"
								@update:model-value="([v]) => handleChange(v)"
								class="relative flex items-center select-none touch-none w-full h-5"
								:min="1"
								:max="500"
								:step="1"
							>
								<SliderTrack class="bg-gray-500 relative grow rounded-full h-[3px]">
									<SliderRange class="absolute bg-indigo-400 rounded-full h-full" />
								</SliderTrack>
								<SliderThumb class="block w-5 h-5 bg-white cursor-pointer rounded-[10px]" />
							</SliderRoot>
						</div>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<div class="space-y-2">
				<div
					v-for="(field, index) in excludedEmotes"
					:key="field.key"
					class="flex items-center gap-2"
				>
					<FormField :name="`excludedEmotes[${index}]`" v-slot="{ componentField }">
						<FormItem class="flex-1">
							<Input v-bind="componentField" placeholder="Enter emote name" />
						</FormItem>
					</FormField>
					<Button type="button" variant="outline" size="sm" @click="removeExcludedEmote(index)">
						<X class="h-4 w-4" />
					</Button>
				</div>
			</div>

			<div class="space-y-2">
				<span>Appearing emotes animations</span>

				<div class="grid grid-cols-1 lg:grid-cols-2 gap-2">
					<FormField name="animation.fadeIn" v-slot="{ value, handleChange }">
						<FormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<FormControl class="flex items-center !m-0">
								<Checkbox :checked="value" @update:checked="handleChange" />
							</FormControl>
							<FormLabel class="select-none cursor-pointer !m-0">Fade</FormLabel>
						</FormItem>
					</FormField>

					<FormField name="animation.zoomIn" v-slot="{ value, handleChange }">
						<FormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<FormControl class="flex items-center !m-0">
								<Checkbox :checked="value" @update:checked="handleChange" />
							</FormControl>
							<FormLabel class="select-none cursor-pointer !m-0">Zoom</FormLabel>
						</FormItem>
					</FormField>
				</div>
			</div>

			<div class="space-y-2">
				<span>Disappearing emotes animations</span>

				<div class="grid grid-cols-1 lg:grid-cols-2 gap-2">
					<FormField name="animation.fadeOut" v-slot="{ value, handleChange }">
						<FormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<FormControl class="flex items-center !m-0">
								<Checkbox :checked="value" @update:checked="handleChange" />
							</FormControl>
							<FormLabel class="select-none cursor-pointer !m-0">Fade</FormLabel>
						</FormItem>
					</FormField>

					<FormField name="animation.zoomOut" v-slot="{ value, handleChange }">
						<FormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<FormControl class="flex items-center !m-0">
								<Checkbox :checked="value" @update:checked="handleChange" />
							</FormControl>
							<FormLabel class="select-none cursor-pointer !m-0">Zoom</FormLabel>
						</FormItem>
					</FormField>
				</div>
			</div>
		</div>
	</div>
</template>
