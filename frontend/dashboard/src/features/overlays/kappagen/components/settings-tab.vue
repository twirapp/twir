<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Plus, X } from 'lucide-vue-next'
import { useFieldArray } from 'vee-validate'
import CommandButton from '@/features/commands/ui/command-button.vue'
import { SliderRange, SliderRoot, SliderThumb, SliderTrack } from 'radix-vue'
import { Checkbox } from '@/components/ui/checkbox'

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
</script>

<template>
	<div class="bg-stone-700/40 rounded-md w-full">
		<div class="flex flex-col gap-4 p-2">
			<div class="flex justify-between items-center">
				<span>Command</span>

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

			<FormField name="size.rationSmall" v-slot="{ value, handleChange }">
				<FormItem class="flex flex-col">
					<FormLabel>Emote Ratio</FormLabel>
					<div class="flex flex-row items-center gap-2">
						<Input disabled :value="value" class="w-[80px]" />
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
							<SliderThumb
								class="block w-5 h-5 bg-white cursor-pointer rounded-[10px]"
								aria-label="Volume"
							/>
						</SliderRoot>
					</div>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField name="size.rationNormal" v-slot="{ value, handleChange }">
				<FormItem class="flex flex-col">
					<FormLabel>Small emote ratio</FormLabel>
					<div class="flex flex-row items-center gap-2">
						<Input disabled :value="value" class="w-[80px]" />
						<SliderRoot
							:model-value="[value]"
							@update:model-value="([v]) => handleChange(v)"
							class="relative flex items-center select-none touch-none w-full h-5"
							:min="0.02"
							:max="0.07"
							:step="1"
						>
							<SliderTrack class="bg-gray-500 relative grow rounded-full h-[3px]">
								<SliderRange class="absolute bg-indigo-400 rounded-full h-full" />
							</SliderTrack>
							<SliderThumb
								class="block w-5 h-5 bg-white cursor-pointer rounded-[10px]"
								aria-label="Volume"
							/>
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

			<Button type="button" variant="outline" size="sm" @click="() => addExcludedEmote('')">
				<Plus class="h-4 w-4 mr-2" />
				Add Excluded Emote
			</Button>
		</div>
	</div>
</template>
