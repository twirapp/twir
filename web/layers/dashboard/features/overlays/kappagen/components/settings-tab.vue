<script setup lang="ts">
import { X } from 'lucide-vue-next'
import { SliderRange, SliderRoot, SliderThumb, SliderTrack } from 'reka-ui'
import { useFieldArray } from 'vee-validate'







import CommandButton from '~/features/commands/ui/command-button.vue'
import { KappagenEmojiStyle } from '~/gql/graphql.ts'

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

			<UiFormField v-slot="{ value, handleChange }" name="enableSpawn">
				<UiFormItem class="flex flex-row items-center justify-between">
					<div class="flex flex-col w-full">
						<UiFormLabel>Enable Spawn</UiFormLabel>
						<div class="text-[0.8rem] text-muted-foreground">Spawn emote on each chat message</div>
					</div>

					<UiSwitch :model-value="value" @update:model-value="handleChange" />
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField v-slot="{ value, handleChange }" name="enableRave">
				<UiFormItem class="flex flex-row items-center justify-between">
					<div class="flex flex-col w-full">
						<UiFormLabel>Enable Rave Mode</UiFormLabel>
						<div class="text-[0.8rem] text-muted-foreground">
							Enable special rave animations and effects
						</div>
					</div>
					<UiSwitch :model-value="value" @update:model-value="handleChange" />
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField v-slot="{ value, handleChange }" name="size.rationSmall">
				<UiFormItem class="flex flex-col">
					<UiFormLabel>Small Emote Ratio</UiFormLabel>
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
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField v-slot="{ value, handleChange }" name="size.rationNormal">
				<UiFormItem class="flex flex-col">
					<UiFormLabel>Normal emote ratio</UiFormLabel>
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
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<div class="space-y-2">
				<span>External emotes</span>
				<div class="grid grid-cols-2 gap-2">
					<UiFormField
						v-for="service of emotesCheckboxes"
						v-slot="{ value, handleChange }"
						:key="service.key"
						:name="`emotes.${service.key}Enabled`"
					>
						<UiFormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<UiFormControl class="flex items-center m-0!">
								<UiCheckbox :model-value="value" @update:model-value="handleChange" />
							</UiFormControl>
							<UiFormLabel class="select-none cursor-pointer m-0!">
								{{ service.label }}
							</UiFormLabel>
						</UiFormItem>
					</UiFormField>
				</div>
			</div>

			<UiFormField v-slot="{ componentField }" name="emotes.emojiStyle">
				<UiFormItem>
					<UiFormLabel>Emoji style</UiFormLabel>

					<UiSelect v-bind="componentField">
						<UiFormControl>
							<UiSelectTrigger>
								<UiSelectValue placeholder="Select emoji style" />
							</UiSelectTrigger>
						</UiFormControl>
						<UiSelectContent>
							<UiSelectGroup>
								<UiSelectItem
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
								</UiSelectItem>
							</UiSelectGroup>
						</UiSelectContent>
					</UiSelect>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<div class="flex flex-col gap-1">
				<UiFormField v-slot="{ value, handleChange }" name="emotes.time">
					<UiFormItem class="flex flex-col">
						<UiFormLabel>Max emote time</UiFormLabel>
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
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>

			<div class="flex flex-col gap-1">
				<UiFormField v-slot="{ value, handleChange }" name="emotes.max">
					<UiFormItem class="flex flex-col">
						<UiFormLabel>Max emotes count</UiFormLabel>
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
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>

			<div class="space-y-2">
				<div
					v-for="(field, index) in excludedEmotes"
					:key="field.key"
					class="flex items-center gap-2"
				>
					<UiFormField v-slot="{ componentField }" :name="`excludedEmotes[${index}]`">
						<UiFormItem class="flex-1">
							<UiInput v-bind="componentField" placeholder="Enter emote name" />
						</UiFormItem>
					</UiFormField>
					<UiButton type="button" variant="outline" size="sm" @click="removeExcludedEmote(index)">
						<X class="h-4 w-4" />
					</UiButton>
				</div>
			</div>

			<div class="space-y-2">
				<span>Appearing emotes animations</span>

				<div class="grid grid-cols-1 lg:grid-cols-2 gap-2">
					<UiFormField v-slot="{ value, handleChange }" name="animation.fadeIn">
						<UiFormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<UiFormControl class="flex items-center m-0!">
								<UiCheckbox :model-value="value" @update:model-value="handleChange" />
							</UiFormControl>
							<UiFormLabel class="select-none cursor-pointer m-0!"> Fade </UiFormLabel>
						</UiFormItem>
					</UiFormField>

					<UiFormField v-slot="{ value, handleChange }" name="animation.zoomIn">
						<UiFormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<UiFormControl class="flex items-center m-0!">
								<UiCheckbox :model-value="value" @update:model-value="handleChange" />
							</UiFormControl>
							<UiFormLabel class="select-none cursor-pointer m-0!"> Zoom </UiFormLabel>
						</UiFormItem>
					</UiFormField>
				</div>
			</div>

			<div class="space-y-2">
				<span>Disappearing emotes animations</span>

				<div class="grid grid-cols-1 lg:grid-cols-2 gap-2">
					<UiFormField v-slot="{ value, handleChange }" name="animation.fadeOut">
						<UiFormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<UiFormControl class="flex items-center m-0!">
								<UiCheckbox :model-value="value" @update:model-value="handleChange" />
							</UiFormControl>
							<UiFormLabel class="select-none cursor-pointer m-0!"> Fade </UiFormLabel>
						</UiFormItem>
					</UiFormField>

					<UiFormField v-slot="{ value, handleChange }" name="animation.zoomOut">
						<UiFormItem class="flex gap-2 items-center rounded-md bg-stone-700/40 p-2">
							<UiFormControl class="flex items-center m-0!">
								<UiCheckbox :model-value="value" @update:model-value="handleChange" />
							</UiFormControl>
							<UiFormLabel class="select-none cursor-pointer m-0!"> Zoom </UiFormLabel>
						</UiFormItem>
					</UiFormField>
				</div>
			</div>
		</div>
	</div>
</template>
