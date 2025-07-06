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
import { Slider } from '@/components/ui/slider'
import { Switch } from '@/components/ui/switch'
import { Textarea } from '@/components/ui/textarea'
import { KappagenEmojiStyle } from '@/gql/graphql'
import { Plus, Trash2 } from 'lucide-vue-next'

// Custom animations
const {
	fields: animations,
	push: addAnimation,
	remove: removeAnimation,
} = useFieldArray('animations')

const emojiStyleOptions = computed(() => {
	return Object.values(KappagenEmojiStyle).map((style) => ({
		value: style,
		label: style.charAt(0) + style.slice(1).toLowerCase(),
	}))
})

function createNewAnimation() {
	addAnimation({
		style: '',
		prefs: {
			size: 1,
			center: false,
			speed: 50,
			faces: false,
			message: [],
			time: 5000,
		},
		count: 1,
		enabled: true,
	})
}

function updateAnimationMessage(handleChange: (value: string[]) => void, value: string) {
	const messages = value
		.split(',')
		.map((s) => s.trim())
		.filter(Boolean)
	handleChange(messages)
}

function getAnimationMessageString(messages: string[]): string {
	return messages?.join(', ') || ''
}
</script>

<template>
	<div class="space-y-6">
		<!-- Basic Animation Settings -->
		<Card>
			<CardHeader>
				<CardTitle>Basic Animation Settings</CardTitle>
				<CardDescription>Configure basic animation behaviors for emotes</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="grid grid-cols-2 gap-4">
					<FormField name="animation.fadeIn" v-slot="{ value, handleChange }">
						<FormItem
							class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm"
						>
							<div class="space-y-0.5">
								<FormLabel>Fade In</FormLabel>
								<div class="text-[0.8rem] text-muted-foreground">Enable fade in animation</div>
							</div>
							<Switch :checked="value" @update:checked="handleChange" />
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="animation.fadeOut" v-slot="{ value, handleChange }">
						<FormItem
							class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm"
						>
							<div class="space-y-0.5">
								<FormLabel>Fade Out</FormLabel>
								<div class="text-[0.8rem] text-muted-foreground">Enable fade out animation</div>
							</div>
							<Switch :checked="value" @update:checked="handleChange" />
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="animation.zoomIn" v-slot="{ value, handleChange }">
						<FormItem
							class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm"
						>
							<div class="space-y-0.5">
								<FormLabel>Zoom In</FormLabel>
								<div class="text-[0.8rem] text-muted-foreground">Enable zoom in animation</div>
							</div>
							<Switch :checked="value" @update:checked="handleChange" />
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="animation.zoomOut" v-slot="{ value, handleChange }">
						<FormItem
							class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm"
						>
							<div class="space-y-0.5">
								<FormLabel>Zoom Out</FormLabel>
								<div class="text-[0.8rem] text-muted-foreground">Enable zoom out animation</div>
							</div>
							<Switch :checked="value" @update:checked="handleChange" />
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
			</CardContent>
		</Card>

		<!-- Emote Settings -->
		<Card>
			<CardHeader>
				<CardTitle>Emote Settings</CardTitle>
				<CardDescription>Configure emote behavior and appearance</CardDescription>
			</CardHeader>
			<CardContent class="space-y-6">
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<FormField name="emotes.time" v-slot="{ componentField }">
						<FormItem>
							<FormLabel>Display Time (ms)</FormLabel>
							<Input
								v-bind="componentField"
								type="number"
								min="1000"
								max="60000"
								placeholder="5000"
							/>
							<div class="text-[0.8rem] text-muted-foreground">
								How long emotes stay on screen (1000-60000ms)
							</div>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="emotes.max" v-slot="{ componentField }">
						<FormItem>
							<FormLabel>Maximum Emotes</FormLabel>
							<Input v-bind="componentField" type="number" min="1" max="500" placeholder="100" />
							<div class="text-[0.8rem] text-muted-foreground">
								Maximum number of emotes on screen (1-500)
							</div>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="emotes.queue" v-slot="{ componentField }">
						<FormItem>
							<FormLabel>Queue Size</FormLabel>
							<Input v-bind="componentField" type="number" min="1" max="1000" placeholder="100" />
							<div class="text-[0.8rem] text-muted-foreground">
								Maximum queue size for emotes (1-1000)
							</div>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="emotes.emojiStyle" v-slot="{ value, handleChange }">
						<FormItem>
							<FormLabel>Emoji Style</FormLabel>
							<Select :model-value="value" @update:model-value="handleChange">
								<SelectTrigger>
									<SelectValue placeholder="Select emoji style" />
								</SelectTrigger>
								<SelectContent>
									<SelectItem
										v-for="option in emojiStyleOptions"
										:key="option.value"
										:value="option.value"
									>
										{{ option.label }}
									</SelectItem>
								</SelectContent>
							</Select>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>

				<div class="grid grid-cols-3 gap-4">
					<FormField name="emotes.ffzEnabled" v-slot="{ value, handleChange }">
						<FormItem
							class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm"
						>
							<div class="space-y-0.5">
								<FormLabel>FFZ Emotes</FormLabel>
								<div class="text-[0.8rem] text-muted-foreground">Enable FrankerFaceZ emotes</div>
							</div>
							<Switch :checked="value" @update:checked="handleChange" />
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="emotes.bttvEnabled" v-slot="{ value, handleChange }">
						<FormItem
							class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm"
						>
							<div class="space-y-0.5">
								<FormLabel>BTTV Emotes</FormLabel>
								<div class="text-[0.8rem] text-muted-foreground">Enable BetterTTV emotes</div>
							</div>
							<Switch :checked="value" @update:checked="handleChange" />
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="emotes.sevenTvEnabled" v-slot="{ value, handleChange }">
						<FormItem
							class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm"
						>
							<div class="space-y-0.5">
								<FormLabel>7TV Emotes</FormLabel>
								<div class="text-[0.8rem] text-muted-foreground">Enable 7TV emotes</div>
							</div>
							<Switch :checked="value" @update:checked="handleChange" />
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
			</CardContent>
		</Card>

		<!-- Size Settings -->
		<Card>
			<CardHeader>
				<CardTitle>Size Settings</CardTitle>
				<CardDescription>Configure emote size scaling and limits</CardDescription>
			</CardHeader>
			<CardContent class="space-y-6">
				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					<FormField name="size.rationNormal" v-slot="{ value, handleChange }">
						<FormItem>
							<FormLabel>Normal Ratio ({{ value }})</FormLabel>
							<Slider
								:model-value="[value]"
								@update:model-value="(val) => handleChange(val[0])"
								:min="0.1"
								:max="5"
								:step="0.1"
								class="w-full"
							/>
							<div class="text-[0.8rem] text-muted-foreground">
								Size ratio for normal emotes (0.1-5)
							</div>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="size.rationSmall" v-slot="{ value, handleChange }">
						<FormItem>
							<FormLabel>Small Ratio ({{ value }})</FormLabel>
							<Slider
								:model-value="[value]"
								@update:model-value="(val) => handleChange(val[0])"
								:min="0.1"
								:max="5"
								:step="0.1"
								class="w-full"
							/>
							<div class="text-[0.8rem] text-muted-foreground">
								Size ratio for small emotes (0.1-5)
							</div>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="size.min" v-slot="{ componentField }">
						<FormItem>
							<FormLabel>Minimum Size (px)</FormLabel>
							<Input v-bind="componentField" type="number" min="10" max="200" placeholder="20" />
							<div class="text-[0.8rem] text-muted-foreground">
								Minimum emote size in pixels (10-200)
							</div>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="size.max" v-slot="{ componentField }">
						<FormItem>
							<FormLabel>Maximum Size (px)</FormLabel>
							<Input v-bind="componentField" type="number" min="50" max="500" placeholder="150" />
							<div class="text-[0.8rem] text-muted-foreground">
								Maximum emote size in pixels (50-500)
							</div>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
			</CardContent>
		</Card>

		<!-- Custom Animations -->
		<Card>
			<CardHeader>
				<CardTitle>Custom Animations</CardTitle>
				<CardDescription>Configure custom animation styles and their properties</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="space-y-6">
					<div
						v-for="(animation, index) in animations"
						:key="animation.key"
						class="border rounded-lg p-4 space-y-4"
					>
						<div class="flex items-center justify-between">
							<h4 class="font-medium">Animation {{ index + 1 }}</h4>
							<div class="flex items-center gap-2">
								<FormField :name="`animations[${index}].enabled`" v-slot="{ value, handleChange }">
									<FormItem class="flex items-center space-x-2">
										<Switch :checked="value" @update:checked="handleChange" />
										<FormLabel>Enabled</FormLabel>
									</FormItem>
								</FormField>
								<Button
									type="button"
									variant="destructive"
									size="sm"
									@click="removeAnimation(index)"
								>
									<Trash2 class="h-4 w-4" />
								</Button>
							</div>
						</div>

						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
							<FormField :name="`animations[${index}].style`" v-slot="{ componentField }">
								<FormItem>
									<FormLabel>Animation Style</FormLabel>
									<Input v-bind="componentField" placeholder="Enter animation style" />
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField :name="`animations[${index}].count`" v-slot="{ componentField }">
								<FormItem>
									<FormLabel>Count</FormLabel>
									<Input v-bind="componentField" type="number" min="1" max="100" placeholder="1" />
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField :name="`animations[${index}].prefs.time`" v-slot="{ componentField }">
								<FormItem>
									<FormLabel>Duration (ms)</FormLabel>
									<Input
										v-bind="componentField"
										type="number"
										min="100"
										max="30000"
										placeholder="5000"
									/>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField :name="`animations[${index}].prefs.size`" v-slot="{ value, handleChange }">
								<FormItem>
									<FormLabel>Size ({{ value }})</FormLabel>
									<Slider
										:model-value="[value]"
										@update:model-value="(val) => handleChange(val[0])"
										:min="0.1"
										:max="10"
										:step="0.1"
										class="w-full"
									/>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								:name="`animations[${index}].prefs.speed`"
								v-slot="{ value, handleChange }"
							>
								<FormItem>
									<FormLabel>Speed ({{ value }})</FormLabel>
									<Slider
										:model-value="[value]"
										@update:model-value="(val) => handleChange(val[0])"
										:min="1"
										:max="100"
										:step="1"
										class="w-full"
									/>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								:name="`animations[${index}].prefs.center`"
								v-slot="{ value, handleChange }"
							>
								<FormItem
									class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm"
								>
									<div class="space-y-0.5">
										<FormLabel>Center</FormLabel>
										<div class="text-[0.8rem] text-muted-foreground">Center animation</div>
									</div>
									<Switch :checked="value" @update:checked="handleChange" />
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								:name="`animations[${index}].prefs.faces`"
								v-slot="{ value, handleChange }"
							>
								<FormItem
									class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm"
								>
									<div class="space-y-0.5">
										<FormLabel>Faces</FormLabel>
										<div class="text-[0.8rem] text-muted-foreground">Enable face detection</div>
									</div>
									<Switch :checked="value" @update:checked="handleChange" />
									<FormMessage />
								</FormItem>
							</FormField>
						</div>

						<FormField
							:name="`animations[${index}].prefs.message`"
							v-slot="{ value, handleChange }"
						>
							<FormItem>
								<FormLabel>Messages</FormLabel>
								<Textarea
									:model-value="getAnimationMessageString(value)"
									placeholder="Enter messages separated by commas"
									@update:model-value="(val) => updateAnimationMessage(handleChange, val)"
								/>
								<div class="text-[0.8rem] text-muted-foreground">
									Enter messages separated by commas
								</div>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>
				</div>

				<Button type="button" variant="outline" @click="createNewAnimation">
					<Plus class="h-4 w-4 mr-2" />
					Add Animation
				</Button>
			</CardContent>
		</Card>
	</div>
</template>
