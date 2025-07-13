<script setup lang="ts">
import { FieldArray, useFieldArray } from 'vee-validate'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Switch } from '@/components/ui/switch'
import type { KappagenOverlayAnimationsSettings } from '@/gql/graphql'
import { PlayIcon, PlusIcon, SettingsIcon, XIcon } from 'lucide-vue-next'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'

import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useKappagenInstance } from '@/features/overlays/kappagen/composables/use-kappagen-instance.ts'

const { fields: animations, update: updateAnimations } =
	useFieldArray<KappagenOverlayAnimationsSettings>('animations')

function switchEnabled(index: number, value: boolean) {
	const animation = animations.value[index]
	if (animation) {
		updateAnimations(index, { ...animation.value, enabled: value })
	}
}

const kappagen = useKappagenInstance()
</script>

<template>
	<div class="grid grid-cols-1 gap-2 max-h-[40dvh] overflow-y-auto">
		<div
			v-for="(animation, index) of animations"
			:key="animation.key"
			class="flex flex-row gap-2 flex-wrap items-center justify-between bg-stone-700/40 p-2 rounded-md"
		>
			<div class="flex gap-2 items-center">
				<button
					@click="kappagen.playAnimation(animation.value)"
					class="p-2 rounded-md bg-indigo-400/70"
				>
					<PlayIcon class="size-4" />
				</button>
				{{ animation.value.style }}
			</div>

			<div class="flex gap-2 items-center">
				<Popover v-if="animation.value.prefs != null || animation.value.count != null">
					<PopoverTrigger as-child>
						<button
							class="p-1 border-border border rounded-md bg-zinc-600/50 hover:bg-zinc-600/30 transition-colors"
						>
							<SettingsIcon class="size-4" />
						</button>
					</PopoverTrigger>
					<PopoverContent class="w-80 bg-zinc-800/60 backdrop-blur-md">
						<div class="flex flex-col gap-2">
							<FormField
								v-if="typeof animation.value.count === 'number'"
								v-slot="{ componentField }"
								:name="`animations.${index}.count`"
							>
								<FormItem>
									<FormLabel>Emotes count</FormLabel>
									<FormControl>
										<Input
											v-bind="componentField"
											type="number"
											class="input"
											:min="1"
											:max="100"
										/>
									</FormControl>
									<FormDescription> Set count of emotes in animation. </FormDescription>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								v-if="typeof animation.value.prefs?.size === 'number'"
								v-slot="{ componentField }"
								:name="`animations.${index}.prefs.size`"
							>
								<FormItem>
									<FormLabel>Size</FormLabel>
									<FormControl>
										<Input
											v-bind="componentField"
											type="number"
											class="input"
											:min="1"
											:max="100"
										/>
									</FormControl>
									<FormDescription> Set the size of the animation. </FormDescription>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								v-if="typeof animation.value.prefs?.speed === 'number'"
								v-slot="{ componentField }"
								:name="`animations.${index}.prefs.speed`"
							>
								<FormItem>
									<FormLabel>Speed</FormLabel>
									<FormControl>
										<Input
											v-bind="componentField"
											type="number"
											class="input"
											:min="1"
											:max="100"
										/>
									</FormControl>
									<FormDescription> Set the speed of the animation. </FormDescription>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								v-if="typeof animation.value.prefs?.center === 'boolean'"
								v-slot="{ field }"
								:name="`animations.${index}.prefs.center`"
							>
								<FormItem class="flex flex-row items-center justify-between">
									<FormLabel class="text-base"> Center Animation </FormLabel>

									<FormControl>
										<Switch
											:checked="field.value"
											default-checked
											@update:checked="field['onUpdate:modelValue']"
										/>
									</FormControl>
								</FormItem>
							</FormField>

							<FormField
								v-if="typeof animation.value.prefs?.faces === 'boolean'"
								v-slot="{ field }"
								:name="`animations.${index}.prefs.faces`"
							>
								<FormItem class="flex flex-row items-center justify-between">
									<FormLabel class="text-base"> Use Faces </FormLabel>

									<FormControl>
										<Switch
											:checked="field.value"
											default-checked
											@update:checked="field['onUpdate:modelValue']"
										/>
									</FormControl>
								</FormItem>
							</FormField>

							<template
								v-if="
									animation.value.prefs?.message && Array.isArray(animation.value.prefs?.message)
								"
							>
								<FieldArray
									:name="`animations[${index}].prefs.message`"
									v-slot="{ fields, push, remove }"
								>
									<div class="flex items-cente justify-between">
										<h3>Messages</h3>
										<button
											class="p-1 border-border border rounded-md bg-green-600/50 hover:bg-green-600/30 transition-colors"
											@click="push('Twir')"
										>
											<PlusIcon class="size-4" />
										</button>
									</div>
									<ScrollArea class="h-[300px] max-h-[300px] pr-4">
										<div
											v-for="(message, messageIndex) of fields"
											:key="`responses-text-${message.key}`"
										>
											<FormField
												v-slot="{ componentField }"
												:name="`animations[${index}].prefs.message[${messageIndex}]`"
											>
												<FormItem class="flex items-center gap-2 space-y-0">
													<FormControl>
														<Input v-bind="componentField" />
													</FormControl>
													<button
														class="p-1 border-border border rounded-md bg-red-600/50 hover:bg-red-600/30 transition-colors"
														@click="remove(messageIndex)"
													>
														<XIcon />
													</button>
												</FormItem>
											</FormField>
										</div>
									</ScrollArea>
								</FieldArray>
							</template>
						</div>
					</PopoverContent>
				</Popover>
				<Switch :checked="animation.value.enabled" @update:checked="switchEnabled(index, $event)" />
			</div>
		</div>
	</div>
</template>
