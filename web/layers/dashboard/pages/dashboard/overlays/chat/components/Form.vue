<script setup lang="ts">
import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useChatOverlayApi } from '#layers/dashboard/api/overlays/chat'
import { useCopyOverlayLink } from '#layers/dashboard/components/overlays/copyOverlayLink.ts'
import { type Font, FontSelector } from '#layers/dashboard/lib/fontsource'
import { ChatBox, type Settings as ChatBoxSettings, type Message } from '@twir/frontend-chat'
import { useIntervalFn } from '@vueuse/core'
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import { z } from 'zod'

import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { ColorPicker } from '@/components/ui/color-picker'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import InputWithIcon from '@/components/ui/InputWithIcon/InputWithIcon.vue'
import { Label } from '@/components/ui/label'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { Slider } from '@/components/ui/slider'
import { Switch } from '@/components/ui/switch'
import { ChannelRolePermissionEnum, ChatOverlayAnimation } from '@/gql/graphql'

import { globalBadges } from '../constants.ts'
import * as faker from '../faker.ts'
import { defaultChatSettings } from './default-settings'
import { useChatOverlayForm } from './form.ts'

const { t } = useI18n()
const { copyOverlayLink } = useCopyOverlayLink('chat')
const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)

const { data: formValue, reset: resetForm } = useChatOverlayForm()

// Preview state
const messagesMock = ref<Message[]>([])
const previewContainerRef = ref<HTMLDivElement>()

// Override scrollIntoView to prevent page scrolling
watch(
	previewContainerRef,
	(container) => {
		if (!container) return

		const observer = new MutationObserver(() => {
			const elements = container.querySelectorAll('*')
			elements.forEach((element) => {
				if (!element.hasAttribute('data-scroll-overridden')) {
					element.setAttribute('data-scroll-overridden', 'true')
					element.scrollIntoView = function () {
						return
					}
				}
			})
		})

		observer.observe(container, { childList: true, subtree: true })

		const elements = container.querySelectorAll('*')
		elements.forEach((element) => {
			element.setAttribute('data-scroll-overridden', 'true')
			element.scrollIntoView = function () {
				return
			}
		})
	},
	{ immediate: true }
)

// Generate mock messages
useIntervalFn(() => {
	if (!formValue.value) return

	const internalId = crypto.randomUUID()
	const words = [
		'Hello',
		'Hi',
		'Thanks',
		'Cool',
		'Nice',
		'Wow',
		'Amazing',
		'Great',
		'Awesome',
		'Perfect',
	]
	const randomWord = words[Math.floor(Math.random() * words.length)]

	const isKickMessage = messagesMock.value.length % 2 === 1

	messagesMock.value.push({
		sender: faker.firstName(),
		chunks: [{ type: 'text', value: randomWord! }],
		createdAt: new Date(),
		internalId,
		isAnnounce: faker.boolean(),
		isItalic: false,
		type: 'message',
		platform: isKickMessage ? 'kick' : 'twitch',
		senderColor: faker.rgb(),
		announceColor: '',
		badges: isKickMessage ? undefined : { [faker.randomArrayItem(globalBadges)!.set_id]: '1' },
		kickBadges: isKickMessage
			? [
					{
						type: faker.randomArrayItem(['subscriber', 'moderator', 'vip']) ?? 'subscriber',
						text: '',
					},
				]
			: undefined,
		id: crypto.randomUUID(),
		senderDisplayName: faker.firstName(),
	})

	if (formValue.value.messageHideTimeout !== 0) {
		setTimeout(() => {
			messagesMock.value = messagesMock.value.filter((m) => m.internalId !== internalId)
		}, formValue.value.messageHideTimeout * 1000)
	}

	if (messagesMock.value.length >= 20) {
		messagesMock.value = messagesMock.value.slice(1)
	}
}, 1000)

// Auto-scroll preview when messages change
watch(messagesMock, async () => {
	if (!formValue.value || !previewContainerRef.value) return

	await nextTick()
	const chatContainer = previewContainerRef.value.querySelector('.chat') as HTMLElement
	if (!chatContainer) return

	if (formValue.value.direction === 'bottom' || formValue.value.direction === 'top') {
		if (formValue.value.direction === 'bottom') {
			chatContainer.scrollTop = 0
		} else {
			chatContainer.scrollTop = chatContainer.scrollHeight
		}
	} else if (formValue.value.direction === 'left' || formValue.value.direction === 'right') {
		const maxScrollLeft = chatContainer.scrollWidth - chatContainer.clientWidth
		if (formValue.value.direction === 'right') {
			chatContainer.scrollLeft = 0
		} else {
			chatContainer.scrollLeft = maxScrollLeft
		}
	}
})

// Clear mocks on direction change
watch(
	() => formValue.value?.direction,
	() => {
		messagesMock.value = []
	}
)

const chatBoxSettings = computed<ChatBoxSettings>(() => {
	return {
		channelId: '',
		channelName: '',
		channelDisplayName: '',
		globalBadges,
		channelBadges: [],
		...formValue.value,
	}
})

const formSchema = z.object({
	preset: z.enum(['clean', 'boxed']),
	direction: z.enum(['top', 'right', 'bottom', 'left']),
	animation: z.enum(ChatOverlayAnimation),
	hideBots: z.boolean(),
	hideCommands: z.boolean(),
	showBadges: z.boolean(),
	showAnnounceBadge: z.boolean(),
	showPlatformIcon: z.boolean(),
	paddingContainer: z.number().min(0).max(256),
	fontFamily: z.string(),
	fontWeight: z.number(),
	fontStyle: z.string(),
	fontSize: z.number().min(12).max(80),
	chatBackgroundColor: z.string(),
	textShadowColor: z.string(),
	textShadowSize: z.number().min(0).max(30),
	messageHideTimeout: z.number().min(0).max(60),
	messageShowDelay: z.number().min(0).max(60),
})

const { handleSubmit, setValues, values } = useForm({
	validationSchema: formSchema,
})

const styleSelectOptions = [
	{ label: 'Clean', value: 'clean' },
	{ label: 'Boxed', value: 'boxed' },
]

const directionOptions = computed(() => {
	return ['top', 'right', 'bottom', 'left'].map((direction) => ({
		value: direction,
		label: t(`overlays.chat.directions.${direction}`),
	}))
})

const fontData = ref<Font | null>(null)

// Update form when font changes
watch(
	() => fontData.value,
	(font) => {
		if (!font) return
		// Update all font-related fields
		setValues({
			fontFamily: font.id,
			fontWeight: font.weights.includes(formValue.value?.fontWeight ?? 400)
				? formValue.value!.fontWeight
				: font.weights[0],
			fontStyle: font.styles.includes(formValue.value?.fontStyle ?? 'normal')
				? formValue.value!.fontStyle
				: font.styles[0],
		})
	}
)

const fontWeightOptions = computed(() => {
	if (!fontData.value) return []
	return fontData.value.weights.map((weight) => ({ label: `${weight}`, value: weight }))
})

const fontStyleOptions = computed(() => {
	if (!fontData.value) return []
	return fontData.value.styles.map((style) => ({ label: style, value: style }))
})

// Sync form values with formValue from composable
watch(
	formValue,
	(newValue) => {
		if (!newValue) return
		setValues({
			preset: newValue.preset as 'clean' | 'boxed',
			direction: newValue.direction as 'top' | 'right' | 'bottom' | 'left',
			animation: newValue.animation,
			hideBots: newValue.hideBots,
			hideCommands: newValue.hideCommands,
			showBadges: newValue.showBadges,
			showAnnounceBadge: newValue.showAnnounceBadge,
			showPlatformIcon: newValue.showPlatformIcon,
			paddingContainer: newValue.paddingContainer,
			fontFamily: newValue.fontFamily,
			fontWeight: newValue.fontWeight,
			fontStyle: newValue.fontStyle,
			fontSize: newValue.fontSize,
			chatBackgroundColor: newValue.chatBackgroundColor,
			textShadowColor: newValue.textShadowColor,
			textShadowSize: newValue.textShadowSize,
			messageHideTimeout: newValue.messageHideTimeout,
			messageShowDelay: newValue.messageShowDelay,
		})
	},
	{ deep: true, immediate: true }
)

// Sync form values back to formValue
watch(
	values,
	(newValues) => {
		if (!formValue.value || !newValues) return
		Object.assign(formValue.value, newValues)
	},
	{ deep: true }
)

const manager = useChatOverlayApi()
const updater = manager.useOverlayUpdate()

const onSubmit = handleSubmit(async (formValues) => {
	if (!formValue.value?.id) return

	const input = { ...formValues }
	const id = formValue.value.id

	await updater.executeMutation({
		id,
		input,
	})

	toast.success(t('sharedTexts.saved'))
})

function handleCopyLink() {
	if (!formValue.value?.id) return
	copyOverlayLink({ id: formValue.value.id })
}

function handleReset() {
	resetForm()
}
</script>

<template>
	<div
		v-if="!formValue"
		class="flex min-h-[60vh] flex-col items-center justify-center gap-6"
	>
		<div class="space-y-4 text-center">
			<div
				class="border-primary mx-auto h-8 w-8 animate-spin rounded-full border-4 border-t-transparent"
			></div>
			<p class="text-muted-foreground">Loading preset...</p>
		</div>
	</div>

	<div
		v-else
		class="flex w-full gap-8"
	>
		<!-- Form (Left side) -->
		<div class="max-w-2xl flex-1">
			<form
				class="animate-in fade-in slide-in-from-bottom-2 space-y-6 pb-6 duration-300"
				@submit="onSubmit"
			>
				<!-- Action Buttons -->
				<div
					class="bg-background/95 supports-backdrop-filter:bg-background/60 sticky top-0 z-10 -mt-1 flex flex-col gap-2 pb-4 backdrop-blur"
				>
					<Button
						type="submit"
						variant="default"
						size="sm"
						class="w-full shadow-sm transition-all hover:scale-[1.02]"
						:disabled="!userCanEditOverlays"
					>
						<Icon
							name="lucide:save"
							class="mr-2 h-4 w-4"
						/>
						{{ t('sharedButtons.save') }}
					</Button>
					<div class="grid grid-cols-2 gap-2">
						<Button
							type="button"
							variant="outline"
							size="sm"
							class="transition-all hover:scale-[1.02]"
							:disabled="!userCanEditOverlays"
							@click="handleReset"
						>
							<Icon
								name="lucide:rotate-ccw"
								class="mr-2 h-4 w-4"
							/>
							{{ t('sharedButtons.reset') || 'Reset' }}
						</Button>
						<Button
							type="button"
							variant="outline"
							size="sm"
							class="transition-all hover:scale-[1.02]"
							:disabled="!formValue.id || !userCanEditOverlays"
							@click="handleCopyLink"
						>
							<Icon
								name="lucide:link-2"
								class="mr-2 h-4 w-4"
							/>
							{{ t('overlays.copyLink') || 'Copy Link' }}
						</Button>
					</div>
				</div>

				<Separator />

				<!-- General Settings -->
				<div class="animate-in fade-in slide-in-from-bottom-3 space-y-4 delay-100 duration-500">
					<h3 class="text-foreground/90 text-sm font-semibold tracking-tight">
						{{ t('overlays.chat.generalSettings') || 'General' }}
					</h3>

					<FormField
						v-slot="{ componentField }"
						name="preset"
					>
						<FormItem>
							<FormLabel class="text-sm font-medium">{{ t('overlays.chat.style') }}</FormLabel>
							<Select v-bind="componentField">
								<FormControl>
									<SelectTrigger class="transition-all">
										<SelectValue />
									</SelectTrigger>
								</FormControl>
								<SelectContent>
									<SelectItem
										v-for="option in styleSelectOptions"
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

					<FormField
						v-slot="{ componentField }"
						name="direction"
					>
						<FormItem>
							<FormLabel class="text-sm font-medium">{{ t('overlays.chat.direction') }}</FormLabel>
							<Select v-bind="componentField">
								<FormControl>
									<SelectTrigger class="transition-all">
										<SelectValue />
									</SelectTrigger>
								</FormControl>
								<SelectContent>
									<SelectItem
										v-for="option in directionOptions"
										:key="option.value"
										:value="option.value"
									>
										{{ option.label }}
									</SelectItem>
								</SelectContent>
							</Select>
							<p class="text-muted-foreground mt-1 text-xs">
								{{ t('overlays.chat.directionWarning') }}
							</p>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField
						v-slot="{ componentField }"
						name="animation"
					>
						<FormItem>
							<FormLabel class="text-sm font-medium">{{
								t('overlays.chat.animation') || 'Animation'
							}}</FormLabel>
							<Select v-bind="componentField">
								<FormControl>
									<SelectTrigger class="transition-all">
										<SelectValue />
									</SelectTrigger>
								</FormControl>
								<SelectContent>
									<SelectItem :value="ChatOverlayAnimation.Default">Default</SelectItem>
									<SelectItem :value="ChatOverlayAnimation.Disabled">Disabled</SelectItem>
								</SelectContent>
							</Select>
							<FormMessage />
						</FormItem>
					</FormField>

					<div class="bg-muted/30 hover:bg-muted/50 space-y-3 rounded-lg border p-3 transition-all">
						<FormField
							v-slot="{ value, handleChange }"
							name="hideBots"
						>
							<FormItem class="flex items-center justify-between">
								<FormLabel
									class="hover:text-foreground cursor-pointer text-sm font-medium transition-colors"
								>
									{{ t('overlays.chat.hideBots') }}
								</FormLabel>
								<FormControl>
									<Switch
										:model-value="value"
										class="transition-all"
										@update:model-value="handleChange"
									/>
								</FormControl>
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ value, handleChange }"
							name="hideCommands"
						>
							<FormItem class="flex items-center justify-between">
								<FormLabel
									class="hover:text-foreground cursor-pointer text-sm font-medium transition-colors"
								>
									{{ t('overlays.chat.hideCommands') }}
								</FormLabel>
								<FormControl>
									<Switch
										:model-value="value"
										class="transition-all"
										@update:model-value="handleChange"
									/>
								</FormControl>
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ value, handleChange }"
							name="showBadges"
						>
							<FormItem class="flex items-center justify-between">
								<FormLabel
									class="hover:text-foreground cursor-pointer text-sm font-medium transition-colors"
								>
									{{ t('overlays.chat.showBadges') }}
								</FormLabel>
								<FormControl>
									<Switch
										:model-value="value"
										class="transition-all"
										@update:model-value="handleChange"
									/>
								</FormControl>
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ value, handleChange }"
							name="showPlatformIcon"
						>
							<FormItem class="flex items-center justify-between">
								<FormLabel
									class="hover:text-foreground cursor-pointer text-sm font-medium transition-colors"
								>
									{{ t('overlays.chat.showPlatformIcon') }}
								</FormLabel>
								<FormControl>
									<Switch
										:model-value="value"
										class="transition-all"
										@update:model-value="handleChange"
									/>
								</FormControl>
							</FormItem>
						</FormField>

						<FormField
							v-if="formValue.preset === 'boxed'"
							v-slot="{ value, handleChange }"
							name="showAnnounceBadge"
						>
							<FormItem class="flex items-center justify-between">
								<FormLabel
									class="hover:text-foreground cursor-pointer text-sm font-medium transition-colors"
								>
									{{ t('overlays.chat.showAnnounceBadge') }}
								</FormLabel>
								<FormControl>
									<Switch
										:model-value="value"
										:disabled="!formValue.showBadges"
										class="transition-all"
										@update:model-value="handleChange"
									/>
								</FormControl>
							</FormItem>
						</FormField>
					</div>

					<FormField
						v-slot="{ componentField }"
						name="paddingContainer"
					>
						<FormItem>
							<FormLabel class="text-sm font-medium"
								>{{ t('overlays.chat.paddingContainer') }} ({{
									formValue.paddingContainer
								}}px)</FormLabel
							>
							<FormControl>
								<Slider
									:model-value="[componentField.modelValue]"
									:min="0"
									:max="256"
									:step="1"
									@update:model-value="
										(value) => {
											if (value?.[0] !== undefined && componentField['onUpdate:modelValue'])
												componentField['onUpdate:modelValue'](value[0])
										}
									"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>

				<Separator />

				<!-- Font Settings -->
				<div class="animate-in fade-in slide-in-from-bottom-3 space-y-4 delay-200 duration-500">
					<h3 class="text-foreground/90 text-sm font-semibold tracking-tight">
						{{ t('overlays.chat.fontSettings') || 'Font' }}
					</h3>

					<div class="space-y-2 transition-all">
						<Label class="text-sm font-medium">{{ t('overlays.chat.fontFamily') }}</Label>
						<FontSelector
							:font-family="formValue.fontFamily"
							:font-weight="formValue.fontWeight"
							:font-style="formValue.fontStyle"
							@update:font="
								(font) => {
									if (font) fontData = font
								}
							"
						/>
					</div>

					<FormField
						v-slot="{ componentField }"
						name="fontWeight"
					>
						<FormItem>
							<FormLabel class="text-sm font-medium">{{ t('overlays.chat.fontWeight') }}</FormLabel>
							<Select v-bind="componentField">
								<FormControl>
									<SelectTrigger class="transition-all">
										<SelectValue />
									</SelectTrigger>
								</FormControl>
								<SelectContent>
									<SelectItem
										v-for="option in fontWeightOptions"
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

					<FormField
						v-slot="{ componentField }"
						name="fontStyle"
					>
						<FormItem>
							<FormLabel class="text-sm font-medium">{{ t('overlays.chat.fontStyle') }}</FormLabel>
							<Select v-bind="componentField">
								<FormControl>
									<SelectTrigger class="transition-all">
										<SelectValue />
									</SelectTrigger>
								</FormControl>
								<SelectContent>
									<SelectItem
										v-for="option in fontStyleOptions"
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

					<FormField
						v-slot="{ componentField }"
						name="fontSize"
					>
						<FormItem>
							<FormLabel class="text-sm font-medium"
								>{{ t('overlays.chat.fontSize') }} ({{ formValue.fontSize }}px)</FormLabel
							>
							<FormControl>
								<Slider
									:model-value="[componentField.modelValue]"
									:min="12"
									:max="80"
									:step="1"
									@update:model-value="
										(value) => {
											if (value?.[0] !== undefined && componentField['onUpdate:modelValue'])
												componentField['onUpdate:modelValue'](value[0])
										}
									"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>

				<Separator />

				<!-- Color Settings -->
				<div class="animate-in fade-in slide-in-from-bottom-3 space-y-4 delay-300 duration-500">
					<h3 class="text-foreground/90 text-sm font-semibold tracking-tight">
						{{ t('overlays.chat.colorSettings') || 'Colors' }}
					</h3>

					<FormField
						v-slot="{ componentField }"
						name="chatBackgroundColor"
					>
						<FormItem>
							<div class="flex items-center justify-between">
								<FormLabel class="text-sm font-medium">{{
									t('overlays.chat.backgroundColor')
								}}</FormLabel>
								<Button
									type="button"
									size="sm"
									variant="ghost"
									class="transition-all hover:scale-105"
									@click="formValue.chatBackgroundColor = defaultChatSettings.chatBackgroundColor"
								>
									<Icon
										name="lucide:rotate-ccw"
										class="mr-1 h-3 w-3"
									/>
									{{ t('overlays.chat.resetToDefault') }}
								</Button>
							</div>
							<FormControl>
								<InputWithIcon v-bind="componentField">
									<ColorPicker v-bind="componentField" />
								</InputWithIcon>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<div class="space-y-2 transition-all">
						<Label class="text-sm font-medium"
							>{{ t('overlays.chat.textShadow') }} ({{ formValue.textShadowSize }}px)</Label
						>
						<div class="flex items-center gap-2">
							<FormField
								v-slot="{ componentField }"
								name="textShadowColor"
							>
								<FormItem>
									<FormControl>
										<ColorPicker
											v-bind="componentField"
											class="h-10 w-10 shrink-0 transition-all hover:scale-105"
										/>
									</FormControl>
								</FormItem>
							</FormField>

							<FormField
								v-slot="{ componentField }"
								name="textShadowSize"
							>
								<FormItem class="flex-1">
									<FormControl>
										<Slider
											:model-value="[componentField.modelValue]"
											:min="0"
											:max="30"
											:step="1"
											class="flex-1"
											@update:model-value="
												(value) => {
													if (value?.[0] !== undefined && componentField['onUpdate:modelValue'])
														componentField['onUpdate:modelValue'](value[0])
												}
											"
										/>
									</FormControl>
								</FormItem>
							</FormField>
						</div>
					</div>
				</div>

				<Separator />

				<!-- Timing Settings -->
				<div class="animate-in fade-in slide-in-from-bottom-3 space-y-4 delay-400 duration-500">
					<h3 class="text-foreground/90 text-sm font-semibold tracking-tight">
						{{ t('overlays.chat.timingSettings') || 'Timing' }}
					</h3>

					<FormField
						v-slot="{ componentField }"
						name="messageHideTimeout"
					>
						<FormItem>
							<FormLabel class="text-sm font-medium"
								>{{ t('overlays.chat.hideTimeout') }} ({{
									formValue.messageHideTimeout
								}}s)</FormLabel
							>
							<FormControl>
								<Slider
									:model-value="[componentField.modelValue]"
									:min="0"
									:max="60"
									:step="1"
									@update:model-value="
										(value) => {
											if (value?.[0] !== undefined && componentField['onUpdate:modelValue'])
												componentField['onUpdate:modelValue'](value[0])
										}
									"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField
						v-slot="{ componentField }"
						name="messageShowDelay"
					>
						<FormItem>
							<FormLabel class="text-sm font-medium"
								>{{ t('overlays.chat.showDelay') }} ({{ formValue.messageShowDelay }}s)</FormLabel
							>
							<FormControl>
								<Slider
									:model-value="[componentField.modelValue]"
									:min="0"
									:max="60"
									:step="1"
									@update:model-value="
										(value) => {
											if (value?.[0] !== undefined && componentField['onUpdate:modelValue'])
												componentField['onUpdate:modelValue'](value[0])
										}
									"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
			</form>
		</div>

		<!-- Preview (Right side) -->
		<div class="sticky top-8 h-fit flex-1">
			<div class="mb-3 px-1">
				<h3 class="text-muted-foreground text-sm font-semibold tracking-wide uppercase">
					Live Preview
				</h3>
			</div>
			<div
				ref="previewContainerRef"
				class="preview-container"
			>
				<Card
					class="chatBox"
					:class="{
						'is-horizontal': formValue?.direction === 'left' || formValue?.direction === 'right',
					}"
				>
					<ChatBox
						class="chatBox-inner"
						:messages="messagesMock"
						:settings="chatBoxSettings"
					/>
				</Card>
			</div>
		</div>
	</div>
</template>

<style scoped>
.preview-container {
	position: relative;
	height: 600px;
	width: 100%;
	overflow: hidden;
	isolation: isolate;
	box-shadow:
		0 4px 6px -1px rgb(0 0 0 / 0.1),
		0 2px 4px -2px rgb(0 0 0 / 0.1);
	border-radius: 0.5rem;
	transition: box-shadow 0.2s ease-in-out;
}

.preview-container:hover {
	box-shadow:
		0 10px 15px -3px rgb(0 0 0 / 0.1),
		0 4px 6px -4px rgb(0 0 0 / 0.1);
}

.chatBox {
	height: 100%;
	width: 100%;
	overflow: hidden;
	display: flex;
	flex-direction: column;
}

.chatBox.is-horizontal {
	contain: layout size style;
}

.chatBox-inner {
	height: 100% !important;
	width: 100% !important;
	max-width: 100% !important;
	overflow: hidden;
	flex: 1;
}

/* Force chat component to respect container width */
:deep(.chat) {
	height: 100% !important;
	width: 100% !important;
	max-width: 100% !important;
	box-sizing: border-box !important;
}

/* For vertical direction, allow vertical overflow */
.chatBox:not(.is-horizontal) :deep(.chat) {
	overflow-x: hidden !important;
	overflow-y: auto !important;
	scrollbar-width: none !important;
	-ms-overflow-style: none !important;
}

/* For horizontal direction, allow horizontal overflow */
.chatBox.is-horizontal :deep(.chat) {
	overflow-x: auto !important;
	overflow-y: hidden !important;
	scrollbar-width: none !important;
	-ms-overflow-style: none !important;
}

/* Hide scrollbar for horizontal chat in Chrome, Safari and Opera */
.chatBox.is-horizontal :deep(.chat)::-webkit-scrollbar {
	display: none !important;
}

/* Hide scrollbar for vertical chat in Chrome, Safari and Opera */
.chatBox:not(.is-horizontal) :deep(.chat)::-webkit-scrollbar {
	display: none !important;
}

/* Force messages container to respect parent width and isolate scrolling */
:deep(.messages) {
	box-sizing: border-box !important;
	overscroll-behavior: contain !important;
	scrollbar-width: none !important;
	-ms-overflow-style: none !important;
}

/* Keep the preview stable when workspace SFC styles are not yet hydrated by Nuxt. */
.chatBox :deep(.profile) {
	display: inline-flex !important;
	align-items: center !important;
	margin-right: 4px !important;
}

.chatBox :deep(.badges) {
	display: inline-flex !important;
	flex: none !important;
	align-items: center !important;
	gap: 4px !important;
	margin-right: 4px !important;
}

.chatBox :deep(.badge),
.chatBox :deep(.kick-badge),
.chatBox :deep(.platform-icon) {
	display: inline-block !important;
	width: 1em !important;
	height: 1em !important;
	flex: none !important;
}

.chatBox :deep(.username) {
	font-weight: 700 !important;
}

/* Vertical direction messages container */
.chatBox:not(.is-horizontal) :deep(.messages) {
	max-width: 100% !important;
	width: 100% !important;
	height: max-content !important;
	min-height: 100% !important;
	overflow: visible !important;
}

/* Horizontal direction messages container - MUST override chat-box.vue overflow: hidden */
.chatBox.is-horizontal :deep(.messages) {
	display: flex !important;
	flex-direction: row !important;
	min-width: 100% !important;
	width: max-content !important;
	height: 100% !important;
	flex-wrap: nowrap !important;
	overflow: visible !important;
}

/* Force horizontal messages to stay in one line */
.chatBox.is-horizontal :deep(.messages > *) {
	flex-shrink: 0 !important;
}

/* Hide scrollbar for Chrome, Safari and Opera */
:deep(.messages)::-webkit-scrollbar {
	display: none !important;
}

/* Constrain messages based on direction */
:deep(.message) {
	box-sizing: border-box !important;
}

/* Vertical direction - wrap text */
.chatBox:not(.is-horizontal) :deep(.message) {
	max-width: 400px !important;
	white-space: normal !important;
	word-wrap: break-word !important;
	overflow-wrap: break-word !important;
}

/* Horizontal direction - single line */
.chatBox.is-horizontal :deep(.message) {
	white-space: nowrap !important;
	overflow: hidden !important;
}
</style>
