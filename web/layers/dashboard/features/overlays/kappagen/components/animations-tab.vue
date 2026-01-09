<script setup lang="ts">
import { PlayIcon, PlusIcon, SettingsIcon, XIcon } from 'lucide-vue-next'
import { FieldArray, useFieldArray } from 'vee-validate'

import type { KappagenOverlayAnimationsSettings } from '~/gql/graphql'






import { useKappagenInstance } from '~/features/overlays/kappagen/composables/use-kappagen-instance.ts'

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
					class="p-2 rounded-md bg-indigo-400/70"
					@click="kappagen.playAnimation(animation.value)"
					type="button"
				>
					<PlayIcon class="size-4" />
				</button>
				{{ animation.value.style }}
			</div>

			<div class="flex gap-2 items-center">
				<UiPopover v-if="animation.value.prefs != null || animation.value.count != null">
					<UiPopoverTrigger as-child>
						<button
							class="p-1 border-border border rounded-md bg-zinc-600/50 hover:bg-zinc-600/30 transition-colors"
							type="button"
						>
							<SettingsIcon class="size-4" />
						</button>
					</UiPopoverTrigger>
					<UiPopoverContent class="w-80 bg-zinc-800/60 backdrop-blur-md">
						<div class="flex flex-col gap-2">
							<UiFormField
								v-if="typeof animation.value.count === 'number'"
								v-slot="{ componentField }"
								:name="`animations.${index}.count`"
							>
								<UiFormItem>
									<UiFormLabel>Emotes count</UiFormLabel>
									<UiFormControl>
										<UiInput
											v-bind="componentField"
											type="number"
											class="input"
											:min="1"
											:max="100"
										/>
									</UiFormControl>
									<UiFormDescription> Set count of emotes in animation. </UiFormDescription>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField
								v-if="typeof animation.value.prefs?.size === 'number'"
								v-slot="{ componentField }"
								:name="`animations.${index}.prefs.size`"
							>
								<UiFormItem>
									<UiFormLabel>Size</UiFormLabel>
									<UiFormControl>
										<UiInput
											v-bind="componentField"
											type="number"
											class="input"
											:min="1"
											:max="100"
										/>
									</UiFormControl>
									<UiFormDescription> Set the size of the animation. </UiFormDescription>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField
								v-if="typeof animation.value.prefs?.speed === 'number'"
								v-slot="{ componentField }"
								:name="`animations.${index}.prefs.speed`"
							>
								<UiFormItem>
									<UiFormLabel>Speed</UiFormLabel>
									<UiFormControl>
										<UiInput
											v-bind="componentField"
											type="number"
											class="input"
											:min="1"
											:max="100"
										/>
									</UiFormControl>
									<UiFormDescription> Set the speed of the animation. </UiFormDescription>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField
								v-if="typeof animation.value.prefs?.center === 'boolean'"
								v-slot="{ field }"
								:name="`animations.${index}.prefs.center`"
							>
								<UiFormItem class="flex flex-row items-center justify-between">
									<UiFormLabel class="text-base"> Center Animation </UiFormLabel>

									<UiFormControl>
										<UiSwitch
											:model-value="field.value"
											default-checked
											@update:model-value="field['onUpdate:modelValue']"
										/>
									</UiFormControl>
								</UiFormItem>
							</UiFormField>

							<UiFormField
								v-if="typeof animation.value.prefs?.faces === 'boolean'"
								v-slot="{ field }"
								:name="`animations.${index}.prefs.faces`"
							>
								<UiFormItem class="flex flex-row items-center justify-between">
									<UiFormLabel class="text-base"> Use Faces </UiFormLabel>

									<UiFormControl>
										<UiSwitch
											:model-value="field.value"
											default-checked
											@update:model-value="field['onUpdate:modelValue']"
										/>
									</UiFormControl>
								</UiFormItem>
							</UiFormField>

							<template
								v-if="
									animation.value.prefs?.message && Array.isArray(animation.value.prefs?.message)
								"
							>
								<FieldArray
									v-slot="{ fields, push, remove }"
									:name="`animations[${index}].prefs.message`"
								>
									<div class="flex items-cente justify-between">
										<h3>Messages</h3>
										<button
											class="p-1 border-border border rounded-md bg-green-600/50 hover:bg-green-600/30 transition-colors"
											@click="push('Twir')"
											type="button"
										>
											<PlusIcon class="size-4" />
										</button>
									</div>
									<UiScrollArea class="h-[300px] max-h-[300px] pr-4">
										<div
											v-for="(message, messageIndex) of fields"
											:key="`responses-text-${message.key}`"
										>
											<UiFormField
												v-slot="{ componentField }"
												:name="`animations[${index}].prefs.message[${messageIndex}]`"
											>
												<UiFormItem class="flex items-center gap-2 space-y-0">
													<UiFormControl>
														<UiInput v-bind="componentField" />
													</UiFormControl>
													<button
														class="p-1 border-border border rounded-md bg-red-600/50 hover:bg-red-600/30 transition-colors"
														@click="remove(messageIndex)"
														type="button"
													>
														<XIcon />
													</button>
												</UiFormItem>
											</UiFormField>
										</div>
									</UiScrollArea>
								</FieldArray>
							</template>
						</div>
					</UiPopoverContent>
				</UiPopover>
				<UiSwitch
					:model-value="animation.value.enabled"
					@update:model-value="switchEnabled(index, $event)"
				/>
			</div>
		</div>
	</div>
</template>
