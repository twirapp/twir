<script setup lang="ts">
import { ChatBox, type Settings as ChatBoxSettings, type Message } from '@twir/frontend-chat'
import { type Font, FontSelector } from '@twir/fontsource'
import { useIntervalFn } from '@vueuse/core'
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Link2, RotateCcw, Save } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'

import { defaultChatSettings } from './default-settings'
import { useChatOverlayForm } from './form.js'
import { globalBadges } from '../constants.js'
import * as faker from '../faker.js'

import { useChatOverlayApi, useUserAccessFlagChecker } from '@/api'
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.js'
import { ChannelRolePermissionEnum, ChatOverlayAnimation } from '@/gql/graphql'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Slider } from '@/components/ui/slider'
import { Separator } from '@/components/ui/separator'
import { ColorPicker } from '@/components/ui/color-picker'
import { toast } from 'vue-sonner'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import InputWithIcon from '@/components/ui/InputWithIcon.vue'
import { Card } from '@/components/ui/card'

const { t } = useI18n()
const { copyOverlayLink } = useCopyOverlayLink('chat')
const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)

const { data: formValue, reset: resetForm } = useChatOverlayForm()

// Preview state
const messagesMock = ref<Message[]>([])
const previewContainerRef = ref<HTMLDivElement>()

// Override scrollIntoView to prevent page scrolling
watch(previewContainerRef, (container) => {
	if (!container) return

	const observer = new MutationObserver(() => {
		const elements = container.querySelectorAll('*')
		elements.forEach((element) => {
			if (!element.hasAttribute('data-scroll-overridden')) {
				element.setAttribute('data-scroll-overridden', 'true')
				element.scrollIntoView = function() {
					return
				}
			}
		})
	})

	observer.observe(container, { childList: true, subtree: true })

	const elements = container.querySelectorAll('*')
	elements.forEach((element) => {
		element.setAttribute('data-scroll-overridden', 'true')
		element.scrollIntoView = function() {
			return
		}
	})
}, { immediate: true })

// Generate mock messages
useIntervalFn(() => {
	if (!formValue.value) return

	const internalId = crypto.randomUUID()
	const words = ['Hello', 'Hi', 'Thanks', 'Cool', 'Nice', 'Wow', 'Amazing', 'Great', 'Awesome', 'Perfect']
	const randomWord = words[Math.floor(Math.random() * words.length)]

	messagesMock.value.push({
		sender: faker.firstName(),
		chunks: [{ type: 'text', value: randomWord }],
		createdAt: new Date(),
		internalId,
		isAnnounce: faker.boolean(),
		isItalic: false,
		type: 'message',
		senderColor: faker.rgb(),
		announceColor: '',
		badges: { [faker.randomArrayItem(globalBadges).set_id]: '1' },
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
watch(() => formValue.value?.direction, () => {
	messagesMock.value = []
})

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

const formSchema = toTypedSchema(
	z.object({
		preset: z.enum(['clean', 'boxed']),
		direction: z.enum(['top', 'right', 'bottom', 'left']),
		animation: z.nativeEnum(ChatOverlayAnimation),
		hideBots: z.boolean(),
		hideCommands: z.boolean(),
		showBadges: z.boolean(),
		showAnnounceBadge: z.boolean(),
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
)

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
	<div v-if="!formValue" class="flex flex-col items-center justify-center min-h-[60vh] gap-6">
		<div class="text-center space-y-4">
			<div class="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full mx-auto"></div>
			<p class="text-muted-foreground">Loading preset...</p>
		</div>
	</div>

	<div v-else class="flex gap-8 w-full">
		<!-- Form (Left side) -->
		<div class="flex-1 max-w-2xl">
			<form class="space-y-6 pb-6 animate-in fade-in slide-in-from-bottom-2 duration-300" @submit="onSubmit">
				<!-- Action Buttons -->
		<div class="flex flex-col gap-2 sticky top-0 bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60 z-10 pb-4 -mt-1">
			<Button
				type="submit"
				variant="default"
				size="sm"
				class="w-full transition-all hover:scale-[1.02] shadow-sm"
				:disabled="!userCanEditOverlays"
			>
				<Save class="h-4 w-4 mr-2" />
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
					<RotateCcw class="h-4 w-4 mr-2" />
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
					<Link2 class="h-4 w-4 mr-2" />
					{{ t('overlays.copyLink') || 'Copy Link' }}
				</Button>
			</div>
		</div>

		<Separator />

		<!-- General Settings -->
		<div class="space-y-4 animate-in fade-in slide-in-from-bottom-3 duration-500 delay-100">
			<h3 class="text-sm font-semibold tracking-tight text-foreground/90">{{ t('overlays.chat.generalSettings') || 'General' }}</h3>

			<FormField v-slot="{ componentField }" name="preset">
				<FormItem>
					<FormLabel class="text-sm font-medium">{{ t('overlays.chat.style') }}</FormLabel>
					<Select v-bind="componentField">
						<FormControl>
							<SelectTrigger class="transition-all">
								<SelectValue />
							</SelectTrigger>
						</FormControl>
						<SelectContent>
							<SelectItem v-for="option in styleSelectOptions" :key="option.value" :value="option.value">
								{{ option.label }}
							</SelectItem>
						</SelectContent>
					</Select>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="direction">
				<FormItem>
					<FormLabel class="text-sm font-medium">{{ t('overlays.chat.direction') }}</FormLabel>
					<Select v-bind="componentField">
						<FormControl>
							<SelectTrigger class="transition-all">
								<SelectValue />
							</SelectTrigger>
						</FormControl>
						<SelectContent>
							<SelectItem v-for="option in directionOptions" :key="option.value" :value="option.value">
								{{ option.label }}
							</SelectItem>
						</SelectContent>
					</Select>
					<p class="text-xs text-muted-foreground mt-1">
						{{ t('overlays.chat.directionWarning') }}
					</p>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="animation">
				<FormItem>
					<FormLabel class="text-sm font-medium">{{ t('overlays.chat.animation') || 'Animation' }}</FormLabel>
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

			<div class="space-y-3 p-3 rounded-lg border bg-muted/30 transition-all hover:bg-muted/50">
				<FormField v-slot="{ value, handleChange }" name="hideBots">
					<FormItem class="flex items-center justify-between">
						<FormLabel class="cursor-pointer text-sm font-medium transition-colors hover:text-foreground">
							{{ t('overlays.chat.hideBots') }}
						</FormLabel>
						<FormControl>
							<Switch :model-value="value" class="transition-all" @update:model-value="handleChange" />
						</FormControl>
					</FormItem>
				</FormField>

				<FormField v-slot="{ value, handleChange }" name="hideCommands">
					<FormItem class="flex items-center justify-between">
						<FormLabel class="cursor-pointer text-sm font-medium transition-colors hover:text-foreground">
							{{ t('overlays.chat.hideCommands') }}
						</FormLabel>
						<FormControl>
							<Switch :model-value="value" class="transition-all" @update:model-value="handleChange" />
						</FormControl>
					</FormItem>
				</FormField>

				<FormField v-slot="{ value, handleChange }" name="showBadges">
					<FormItem class="flex items-center justify-between">
						<FormLabel class="cursor-pointer text-sm font-medium transition-colors hover:text-foreground">
							{{ t('overlays.chat.showBadges') }}
						</FormLabel>
						<FormControl>
							<Switch :model-value="value" class="transition-all" @update:model-value="handleChange" />
						</FormControl>
					</FormItem>
				</FormField>

				<FormField v-if="formValue.preset === 'boxed'" v-slot="{ value, handleChange }" name="showAnnounceBadge">
					<FormItem class="flex items-center justify-between">
						<FormLabel class="cursor-pointer text-sm font-medium transition-colors hover:text-foreground">
							{{ t('overlays.chat.showAnnounceBadge') }}
						</FormLabel>
						<FormControl>
							<Switch :model-value="value" :disabled="!formValue.showBadges" class="transition-all" @update:model-value="handleChange" />
						</FormControl>
					</FormItem>
				</FormField>
			</div>

			<FormField v-slot="{ componentField }" name="paddingContainer">
				<FormItem>
					<FormLabel class="text-sm font-medium">{{ t('overlays.chat.paddingContainer') }} ({{ formValue.paddingContainer }}px)</FormLabel>
					<FormControl>
						<Slider
							:model-value="[componentField.modelValue]"
							:min="0"
							:max="256"
							:step="1"
							@update:model-value="(value) => { if (value?.[0] !== undefined && componentField['onUpdate:modelValue']) componentField['onUpdate:modelValue'](value[0]) }"
						/>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>
		</div>

		<Separator />

		<!-- Font Settings -->
		<div class="space-y-4 animate-in fade-in slide-in-from-bottom-3 duration-500 delay-200">
			<h3 class="text-sm font-semibold tracking-tight text-foreground/90">{{ t('overlays.chat.fontSettings') || 'Font' }}</h3>

			<div class="space-y-2 transition-all">
				<Label class="text-sm font-medium">{{ t('overlays.chat.fontFamily') }}</Label>
				<FontSelector
					:font-family="formValue.fontFamily"
					:font-weight="formValue.fontWeight"
					:font-style="formValue.fontStyle"
					@update:font="(font) => { if (font) fontData = font }"
				/>
			</div>

			<FormField v-slot="{ componentField }" name="fontWeight">
				<FormItem>
					<FormLabel class="text-sm font-medium">{{ t('overlays.chat.fontWeight') }}</FormLabel>
					<Select v-bind="componentField">
						<FormControl>
							<SelectTrigger class="transition-all">
								<SelectValue />
							</SelectTrigger>
						</FormControl>
						<SelectContent>
							<SelectItem v-for="option in fontWeightOptions" :key="option.value" :value="option.value">
								{{ option.label }}
							</SelectItem>
						</SelectContent>
					</Select>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="fontStyle">
				<FormItem>
					<FormLabel class="text-sm font-medium">{{ t('overlays.chat.fontStyle') }}</FormLabel>
					<Select v-bind="componentField">
						<FormControl>
							<SelectTrigger class="transition-all">
								<SelectValue />
							</SelectTrigger>
						</FormControl>
						<SelectContent>
							<SelectItem v-for="option in fontStyleOptions" :key="option.value" :value="option.value">
								{{ option.label }}
							</SelectItem>
						</SelectContent>
					</Select>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="fontSize">
				<FormItem>
					<FormLabel class="text-sm font-medium">{{ t('overlays.chat.fontSize') }} ({{ formValue.fontSize }}px)</FormLabel>
					<FormControl>
						<Slider
							:model-value="[componentField.modelValue]"
							:min="12"
							:max="80"
							:step="1"
							@update:model-value="(value) => { if (value?.[0] !== undefined && componentField['onUpdate:modelValue']) componentField['onUpdate:modelValue'](value[0]) }"
						/>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>
		</div>

		<Separator />

		<!-- Color Settings -->
		<div class="space-y-4 animate-in fade-in slide-in-from-bottom-3 duration-500 delay-300">
			<h3 class="text-sm font-semibold tracking-tight text-foreground/90">{{ t('overlays.chat.colorSettings') || 'Colors' }}</h3>

			<FormField v-slot="{ componentField }" name="chatBackgroundColor">
				<FormItem>
					<div class="flex items-center justify-between">
						<FormLabel class="text-sm font-medium">{{ t('overlays.chat.backgroundColor') }}</FormLabel>
						<Button
							type="button"
							size="sm"
							variant="ghost"
							class="transition-all hover:scale-105"
							@click="formValue.chatBackgroundColor = defaultChatSettings.chatBackgroundColor"
						>
							<RotateCcw class="h-3 w-3 mr-1" />
							{{ t('overlays.chat.resetToDefault') }}
						</Button>
					</div>
					<FormControl>
						<InputWithIcon
							v-bind="componentField"
						>
							<ColorPicker
								v-bind="componentField"
							/>
						</InputWithIcon>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<div class="space-y-2 transition-all">
				<Label class="text-sm font-medium">{{ t('overlays.chat.textShadow') }} ({{ formValue.textShadowSize }}px)</Label>
				<div class="flex items-center gap-2">
					<FormField v-slot="{ componentField }" name="textShadowColor">
						<FormItem>
							<FormControl>
								<ColorPicker
									v-bind="componentField"
									class="w-10 h-10 shrink-0 transition-all hover:scale-105"
								/>
							</FormControl>
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="textShadowSize">
						<FormItem class="flex-1">
							<FormControl>
								<Slider
									:model-value="[componentField.modelValue]"
									:min="0"
									:max="30"
									:step="1"
									class="flex-1"
									@update:model-value="(value) => { if (value?.[0] !== undefined && componentField['onUpdate:modelValue']) componentField['onUpdate:modelValue'](value[0]) }"
								/>
							</FormControl>
						</FormItem>
					</FormField>
				</div>
			</div>
		</div>

		<Separator />

		<!-- Timing Settings -->
		<div class="space-y-4 animate-in fade-in slide-in-from-bottom-3 duration-500 delay-400">
			<h3 class="text-sm font-semibold tracking-tight text-foreground/90">{{ t('overlays.chat.timingSettings') || 'Timing' }}</h3>

			<FormField v-slot="{ componentField }" name="messageHideTimeout">
				<FormItem>
					<FormLabel class="text-sm font-medium">{{ t('overlays.chat.hideTimeout') }} ({{ formValue.messageHideTimeout }}s)</FormLabel>
					<FormControl>
						<Slider
							:model-value="[componentField.modelValue]"
							:min="0"
							:max="60"
							:step="1"
							@update:model-value="(value) => { if (value?.[0] !== undefined && componentField['onUpdate:modelValue']) componentField['onUpdate:modelValue'](value[0]) }"
						/>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="messageShowDelay">
				<FormItem>
					<FormLabel class="text-sm font-medium">{{ t('overlays.chat.showDelay') }} ({{ formValue.messageShowDelay }}s)</FormLabel>
					<FormControl>
						<Slider
							:model-value="[componentField.modelValue]"
							:min="0"
							:max="60"
							:step="1"
							@update:model-value="(value) => { if (value?.[0] !== undefined && componentField['onUpdate:modelValue']) componentField['onUpdate:modelValue'](value[0]) }"
						/>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>
		</div>
			</form>
		</div>

		<!-- Preview (Right side) -->
		<div class="flex-1 sticky top-8 h-fit">
			<div class="mb-3 px-1">
				<h3 class="text-sm font-semibold text-muted-foreground uppercase tracking-wide">Live Preview</h3>
			</div>
			<div ref="previewContainerRef" class="preview-container">
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
	box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
	border-radius: 0.5rem;
	transition: box-shadow 0.2s ease-in-out;
}

.preview-container:hover {
	box-shadow: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1);
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
