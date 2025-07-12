<script setup lang="ts">
import { useFieldArray } from 'vee-validate'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Switch } from '@/components/ui/switch'
import type { KappagenOverlayAnimationsSettings } from '@/gql/graphql'
import { PlayIcon, SettingsIcon } from 'lucide-vue-next'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'

const { fields: animations, update } =
	useFieldArray<KappagenOverlayAnimationsSettings>('animations')

function switchEnabled(index: number, value: boolean) {
	const animation = animations.value[index]
	if (animation) {
		update(index, { ...animation.value, enabled: value })
	}
}
</script>

<template>
	<div class="space-y-6">
		<Card>
			<CardHeader>
				<CardTitle>Animations</CardTitle>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="grid grid-cols-1 md:grid-cols-2 gap-2">
					<div
						v-for="(animation, index) of animations"
						:key="animation.key"
						class="flex flex-row items-center justify-between bg-background/60 p-2 rounded-md"
					>
						<div class="flex gap-2 items-center">
							<button class="p-2 rounded-md bg-indigo-400/70">
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
											v-if="animation.value.prefs?.size !== null"
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
											v-if="animation.value.prefs?.speed !== null"
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
											v-if="animation.value.prefs?.center !== null"
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
											v-if="animation.value.prefs?.faces !== null"
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
									</div>

									{{ animation.value }}
								</PopoverContent>
							</Popover>
							<Switch
								:checked="animation.value.enabled"
								@update:checked="switchEnabled(index, $event)"
							/>
						</div>
					</div>
				</div>
			</CardContent>
		</Card>
	</div>
</template>
