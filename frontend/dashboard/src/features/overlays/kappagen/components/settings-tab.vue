<script setup lang="ts">
import { X } from 'lucide-vue-next'
import { SliderRange, SliderRoot, SliderThumb, SliderTrack } from 'reka-ui'
import { useFieldArray } from 'vee-validate'

import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import CommandButton from '@/features/commands/ui/command-button.vue'
import { KappagenEmojiStyle } from '@/gql/graphql.ts'

const { fields: excludedEmotes, remove: removeExcludedEmote } = useFieldArray('excludedEmotes')

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
		value,
		label: key.charAt(0).toUpperCase() + key.slice(1).toLowerCase(),
	}
})

function getEmojiPreview(style: KappagenEmojiStyle) {
	if (style === KappagenEmojiStyle.Blobmoji) {
		return `https://cdn2.frankerfacez.com/static/emoji/images/blob/1f63a.png`
	}

	return `https://cdn2.frankerfacez.com/static/emoji/images/${style.toLowerCase()}/1f63a.png`
}
</script>

<template>
	<div class="bg-stone-700/40 rounded-md w-full max-h-[40dvh] overflow-y-auto">
		<div class="flex flex-col gap-4 p-2">
			<div class="flex justify-between items-center">
				<span>Command for manual spawn</span>

				<CommandButton name="kappagen" title="" />
			</div>

			<FormField v-slot="{ value, handleChange }" name="enableSpawn">
				<FormItem class="flex flex-row items-center justify-between">
					<div class="flex flex-col w-full">
						<FormLabel>Enable Spawn</FormLabel>
						<div class="text-[0.8rem] text-muted-foreground">Spawn emote on each chat message</div>
					</div>

					<Switch :model-value="value" @update:model-value="handleChange" />
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ value, handleChange }" name="enableRave">
				<FormItem class="flex flex-row items-center justify-between">
					<div class="flex flex-col w-full">
						<FormLabel>Enable Rave Mode</FormLabel>
						<div class="text-[0.8rem] text-muted-foreground">
							Enable special rave animations and effects
						</div>
					</div>
					<Switch :model-value="value" @update:model-value="handleChange" />
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ value, handleChange }" name="size.rationSmall">
				<FormItem class="flex flex-col">
					<FormLabel>Small Emote Ratio</FormLabel>
					<div class="flex flex-row items-center gap-2">
						<div class="bg-stone-700/40 p-1 w-[80px] rounded-md text-center">
							{{ Math.round((value - 0.02) / 0.01) + 1 }}
						</div>
						<SliderRoot
							:model-value="[value]"
							class="relative flex items-center select-none touch-none w-full h-5"
							:min="0.02"
							:max="0.07"
							:step="0.01"
							@update:model-value="(v) => handleChange(v?.[0] ?? 0.02)"
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

			<FormField v-slot="{ value, handleChange }" name="size.rationNormal">
				<FormItem class="flex flex-col">
					<FormLabel>Normal emote ratio</FormLabel>
					<div class="flex flex-row items-center gap-2">
						<div class="bg-stone-700/40 p-1 w-[80px] rounded-md text-center">
							{{ Math.round((value - 0.05) / 0.01) + 1 }}
						</div>
						<SliderRoot
							:model-value="[value]"
							class="relative flex items-center select-none touch-none w-full h-5"
							:min="0.05"
							:max="0.15"
							:step="0.01"
							@update:model-value="(v) => handleChange(v?.[0] ?? 0.05)"
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
						v-slot="{ value, handleChange }"
						:key="service.key"
						:name="`emotes.${service.key}Enabled`"
					>
						<FormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<FormControl class="flex items-center m-0!">
								<Checkbox :checked="value" @update:checked="handleChange" />
							</FormControl>
							<FormLabel class="select-none cursor-pointer m-0!">
								{{ service.label }}
							</FormLabel>
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
									<div class="flex flex-row gap-2 items-center">
										<img
											v-if="item.value !== KappagenEmojiStyle.None"
											class="size-4"
											:src="getEmojiPreview(item.value)"
										/>
										<span>
											{{ item.label }}
										</span>
									</div>
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
								class="relative flex items-center select-none touch-none w-full h-5"
								:min="1"
								:max="15"
								:step="1"
								@update:model-value="(v) => handleChange(v?.[0] ?? 1)"
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
							<div class="bg-stone-700/40 p-1 w-[80px] rounded-md text-center">
								{{ value }}
							</div>
							<SliderRoot
								:model-value="[value]"
								class="relative flex items-center select-none touch-none w-full h-5"
								:min="1"
								:max="500"
								:step="1"
								@update:model-value="(v) => handleChange(v?.[0] ?? 1)"
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
					<FormField v-slot="{ componentField }" :name="`excludedEmotes[${index}]`">
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
					<FormField v-slot="{ value, handleChange }" name="animation.fadeIn">
						<FormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<FormControl class="flex items-center m-0!">
								<Checkbox :checked="value" @update:checked="handleChange" />
							</FormControl>
							<FormLabel class="select-none cursor-pointer m-0!"> Fade </FormLabel>
						</FormItem>
					</FormField>

					<FormField v-slot="{ value, handleChange }" name="animation.zoomIn">
						<FormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<FormControl class="flex items-center m-0!">
								<Checkbox :checked="value" @update:checked="handleChange" />
							</FormControl>
							<FormLabel class="select-none cursor-pointer m-0!"> Zoom </FormLabel>
						</FormItem>
					</FormField>
				</div>
			</div>

			<div class="space-y-2">
				<span>Disappearing emotes animations</span>

				<div class="grid grid-cols-1 lg:grid-cols-2 gap-2">
					<FormField v-slot="{ value, handleChange }" name="animation.fadeOut">
						<FormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<FormControl class="flex items-center m-0!">
								<Checkbox :checked="value" @update:checked="handleChange" />
							</FormControl>
							<FormLabel class="select-none cursor-pointer m-0!"> Fade </FormLabel>
						</FormItem>
					</FormField>

					<FormField v-slot="{ value, handleChange }" name="animation.zoomOut">
						<FormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<FormControl class="flex items-center m-0!">
								<Checkbox :checked="value" @update:checked="handleChange" />
							</FormControl>
							<FormLabel class="select-none cursor-pointer m-0!"> Zoom </FormLabel>
						</FormItem>
					</FormField>
				</div>
			</div>
		</div>
	</div>
</template>
