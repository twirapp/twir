<script setup lang="ts">
import { RadioGroupItem, RadioGroupRoot } from 'reka-ui'
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import { z } from 'zod'
import { useSongRequestOverlaySettingsApi } from '~~/layers/dashboard/api/song-request-overlay-settings.js'

import type {
	SongRequestOverlaySettingsDashboardQuery,
	SongRequestOverlaySettingsUpdateInput,
} from '~/gql/graphql.js'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { ColorPicker } from '@/components/ui/color-picker'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { cn } from '@/lib/utils'
import SongRequestOverlayRenderer from '~/components/song-request-overlay/SongRequestOverlayRenderer.vue'
import {
	SONG_REQUEST_OVERLAY_DEFAULTS,
	SONG_REQUEST_OVERLAY_STYLES,
	type SongRequestOverlayStyle,
} from '~/components/song-request-overlay/types'

type OverlaySettings = SongRequestOverlaySettingsDashboardQuery['songRequestOverlaySettings']

const open = defineModel<boolean>({ default: false })
const { t } = useI18n()
const api = useSongRequestOverlaySettingsApi()
const settingsQuery = api.useSettingsQuery()
const updateMutation = api.useUpdateMutation()
const initializedForOpen = ref(false)
const hasLoadError = ref(false)
const isInitializing = ref(false)
const isSaving = ref(false)
let initializationRequest = 0

const colorPattern = /^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$/
const formSchema = computed(() =>
	z.object({
		style: z.enum(SONG_REQUEST_OVERLAY_STYLES),
		accentColor: z.string().regex(colorPattern, t('songRequests.overlaySettings.validation.color')),
		tickerBackgroundColor: z
			.string()
			.regex(colorPattern, t('songRequests.overlaySettings.validation.color')),
		tickerTextColor: z
			.string()
			.regex(colorPattern, t('songRequests.overlaySettings.validation.color')),
		tickerSpeed: z.coerce
			.number()
			.min(10, t('songRequests.overlaySettings.validation.speed'))
			.max(100, t('songRequests.overlaySettings.validation.speed')),
		hideOnPause: z.boolean(),
	})
)

const defaultValues: SongRequestOverlaySettingsUpdateInput = {
	style: SONG_REQUEST_OVERLAY_DEFAULTS.style,
	accentColor: SONG_REQUEST_OVERLAY_DEFAULTS.accentColor,
	tickerBackgroundColor: SONG_REQUEST_OVERLAY_DEFAULTS.tickerBackgroundColor,
	tickerTextColor: SONG_REQUEST_OVERLAY_DEFAULTS.tickerTextColor,
	tickerSpeed: SONG_REQUEST_OVERLAY_DEFAULTS.tickerSpeed,
	hideOnPause: true,
}

const { handleSubmit, resetForm, values } = useForm({
	validationSchema: formSchema,
	initialValues: defaultValues,
	keepValuesOnUnmount: true,
})

const previewStyle = computed(() => values.style as SongRequestOverlayStyle)
const styleOptions = computed(() =>
	SONG_REQUEST_OVERLAY_STYLES.map((style) => ({
		value: style,
		name: t(`songRequests.overlaySettings.styles.${style.toLowerCase()}.name`),
		description: t(`songRequests.overlaySettings.styles.${style.toLowerCase()}.description`),
	}))
)

const sample = {
	title: 'Midnight City - Live at The Greek Theatre',
	requester: 'twir_viewer',
	videoId: 'dX3k_QDnzHE',
	position: 84,
	duration: 257,
}

function toFormValues(value: OverlaySettings): SongRequestOverlaySettingsUpdateInput {
	return {
		style: value.style,
		accentColor: value.accentColor,
		tickerBackgroundColor: value.tickerBackgroundColor,
		tickerTextColor: value.tickerTextColor,
		tickerSpeed: value.tickerSpeed,
		hideOnPause: value.hideOnPause,
	}
}

function initializeForm(value: OverlaySettings) {
	resetForm({ values: toFormValues(value) })
	initializedForOpen.value = true
}

async function loadSettings() {
	const request = ++initializationRequest
	initializedForOpen.value = false
	hasLoadError.value = false
	isInitializing.value = true

	try {
		const response = await settingsQuery.executeQuery({ requestPolicy: 'network-only' })
		if (!open.value || request !== initializationRequest) return

		if (response.error.value) {
			hasLoadError.value = true
			return
		}

		const loadedSettings = response.data.value?.songRequestOverlaySettings
		if (!loadedSettings) {
			hasLoadError.value = true
			return
		}

		initializeForm(loadedSettings)
	} catch {
		if (open.value && request === initializationRequest) hasLoadError.value = true
	} finally {
		if (request === initializationRequest) isInitializing.value = false
	}
}

watch(
	open,
	(isOpen) => {
		if (!isOpen) {
			initializationRequest++
			initializedForOpen.value = false
			hasLoadError.value = false
			isInitializing.value = false
			return
		}

		void loadSettings()
	},
	{ immediate: true }
)

function setOpen(value: boolean) {
	if (isSaving.value) return
	open.value = value
}

const save = handleSubmit(async (formValues) => {
	isSaving.value = true

	try {
		const result = await updateMutation.executeMutation({
			opts: formValues as SongRequestOverlaySettingsUpdateInput,
		})

		if (result.error) {
			toast.error(t('songRequests.overlaySettings.feedback.error'), {
				description: result.error.message,
			})
			return
		}

		await settingsQuery.executeQuery({ requestPolicy: 'network-only' })
		toast.success(t('songRequests.overlaySettings.feedback.success'))
		open.value = false
	} catch (error) {
		toast.error(t('songRequests.overlaySettings.feedback.error'), {
			description: error instanceof Error ? error.message : undefined,
		})
	} finally {
		isSaving.value = false
	}
})
</script>

<template>
	<Dialog
		:open="open"
		@update:open="setOpen"
	>
		<DialogContent
			class="flex max-h-[90vh] w-[calc(100%-2rem)] max-w-6xl flex-col gap-0 overflow-hidden p-0 sm:max-w-6xl"
		>
			<DialogHeader class="shrink-0 p-6 pb-4">
				<DialogTitle>{{ t('songRequests.overlaySettings.title') }}</DialogTitle>
				<DialogDescription>
					{{ t('songRequests.overlaySettings.description') }}
				</DialogDescription>
			</DialogHeader>

			<div
				v-if="isInitializing"
				class="flex min-h-64 items-center justify-center"
				role="status"
			>
				<Icon
					name="lucide:loader-circle"
					class="text-muted-foreground size-8 animate-spin"
				/>
				<span class="sr-only">{{ t('songRequests.overlaySettings.loading') }}</span>
			</div>

			<div
				v-else-if="hasLoadError"
				class="p-6 pt-2"
			>
				<Alert variant="destructive">
					<Icon name="lucide:circle-alert" />
					<AlertTitle>{{ t('songRequests.overlaySettings.loadError.title') }}</AlertTitle>
					<AlertDescription class="flex flex-col items-start gap-4">
						<p>{{ t('songRequests.overlaySettings.loadError.description') }}</p>
						<Button
							type="button"
							variant="outline"
							:disabled="isInitializing"
							@click="loadSettings"
						>
							<Icon
								name="lucide:refresh-cw"
								data-icon="inline-start"
							/>
							{{ t('songRequests.overlaySettings.loadError.retry') }}
						</Button>
					</AlertDescription>
				</Alert>
			</div>

			<form
				v-else-if="initializedForOpen"
				class="flex min-h-0 flex-1 flex-col overflow-hidden"
				@submit="save"
			>
				<div class="min-h-0 flex-1 overflow-y-auto overscroll-contain">
					<div class="flex flex-col gap-6 px-6 pb-6">
						<FormField
							v-slot="{ value, handleChange }"
							name="style"
						>
							<FormItem>
								<FormLabel>{{ t('songRequests.overlaySettings.styleLabel') }}</FormLabel>
								<FormDescription>
									{{ t('songRequests.overlaySettings.styleDescription') }}
								</FormDescription>
								<FormControl>
									<RadioGroupRoot
										:model-value="value"
										class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3"
										:aria-label="t('songRequests.overlaySettings.styleLabel')"
										@update:model-value="handleChange"
									>
										<div
											v-for="option in styleOptions"
											:key="option.value"
											class="relative"
										>
											<RadioGroupItem
												:id="`song-request-overlay-style-${option.value.toLowerCase()}`"
												:value="option.value"
												class="peer sr-only"
											/>
											<label
												:for="`song-request-overlay-style-${option.value.toLowerCase()}`"
												:class="
													cn(
														'group bg-card hover:bg-accent/50 peer-focus-visible:ring-ring flex h-full cursor-pointer flex-col gap-3 rounded-lg border p-3 text-left transition-colors peer-focus-visible:ring-2 peer-focus-visible:outline-none',
														value === option.value && 'border-primary ring-primary/30 ring-2'
													)
												"
											>
												<div
													class="bg-muted relative aspect-video w-full overflow-hidden rounded-md"
												>
													<SongRequestOverlayRenderer
														:style="option.value"
														:title="sample.title"
														:requester="sample.requester"
														:video-id="sample.videoId"
														:position="sample.position"
														:duration="sample.duration"
														:accent-color="values.accentColor"
														:ticker-background-color="values.tickerBackgroundColor"
														:ticker-text-color="values.tickerTextColor"
														:ticker-speed="Number(values.tickerSpeed)"
														is-playing
													/>
												</div>
												<div class="flex items-start justify-between gap-3">
													<div class="min-w-0">
														<p class="font-medium">{{ option.name }}</p>
														<p class="text-muted-foreground mt-1 text-xs">
															{{ option.description }}
														</p>
													</div>
													<Icon
														v-if="value === option.value"
														name="lucide:circle-check"
														class="text-primary size-5 shrink-0"
													/>
												</div>
											</label>
										</div>
									</RadioGroupRoot>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<section
							class="flex flex-col gap-3"
							aria-labelledby="song-request-preview-title"
						>
							<div>
								<h3
									id="song-request-preview-title"
									class="font-medium"
								>
									{{ t('songRequests.overlaySettings.preview.title') }}
								</h3>
								<p class="text-muted-foreground text-sm">
									{{ t('songRequests.overlaySettings.preview.description') }}
								</p>
							</div>
							<div class="bg-muted relative h-64 overflow-hidden rounded-lg border sm:h-80 lg:h-96">
								<SongRequestOverlayRenderer
									:style="previewStyle"
									:title="sample.title"
									:requester="sample.requester"
									:video-id="sample.videoId"
									:position="sample.position"
									:duration="sample.duration"
									:accent-color="values.accentColor"
									:ticker-background-color="values.tickerBackgroundColor"
									:ticker-text-color="values.tickerTextColor"
									:ticker-speed="Number(values.tickerSpeed)"
									is-playing
								/>
							</div>
						</section>

						<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
							<FormField
								v-slot="{ componentField }"
								name="accentColor"
							>
								<FormItem>
									<FormLabel>{{ t('songRequests.overlaySettings.accentColor') }}</FormLabel>
									<div class="flex items-center gap-3">
										<FormControl>
											<ColorPicker
												v-bind="componentField"
												:aria-label="t('songRequests.overlaySettings.accentColor')"
												class="size-10"
												output-format="hex"
											/>
										</FormControl>
										<span class="text-muted-foreground font-mono text-sm">{{
											values.accentColor
										}}</span>
									</div>
									<FormMessage />
								</FormItem>
							</FormField>

							<template v-if="values.style === 'TICKER'">
								<FormField
									v-slot="{ componentField }"
									name="tickerBackgroundColor"
								>
									<FormItem>
										<FormLabel>{{
											t('songRequests.overlaySettings.ticker.backgroundColor')
										}}</FormLabel>
										<div class="flex items-center gap-3">
											<FormControl>
												<ColorPicker
													v-bind="componentField"
													:aria-label="t('songRequests.overlaySettings.ticker.backgroundColor')"
													class="size-10"
													output-format="hex"
												/>
											</FormControl>
											<span class="text-muted-foreground font-mono text-sm">{{
												values.tickerBackgroundColor
											}}</span>
										</div>
										<FormMessage />
									</FormItem>
								</FormField>

								<FormField
									v-slot="{ componentField }"
									name="tickerTextColor"
								>
									<FormItem>
										<FormLabel>{{ t('songRequests.overlaySettings.ticker.textColor') }}</FormLabel>
										<div class="flex items-center gap-3">
											<FormControl>
												<ColorPicker
													v-bind="componentField"
													:aria-label="t('songRequests.overlaySettings.ticker.textColor')"
													class="size-10"
													output-format="hex"
												/>
											</FormControl>
											<span class="text-muted-foreground font-mono text-sm">{{
												values.tickerTextColor
											}}</span>
										</div>
										<FormMessage />
									</FormItem>
								</FormField>

								<FormField
									v-slot="{ componentField }"
									name="tickerSpeed"
								>
									<FormItem>
										<FormLabel>{{ t('songRequests.overlaySettings.ticker.speed') }}</FormLabel>
										<FormDescription>
											{{ t('songRequests.overlaySettings.ticker.speedDescription') }}
										</FormDescription>
										<FormControl>
											<Input
												v-bind="componentField"
												type="number"
												min="10"
												max="100"
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</template>
						</div>

						<FormField
							v-slot="{ value, handleChange }"
							name="hideOnPause"
						>
							<FormItem class="flex flex-col gap-2 rounded-lg border p-4">
								<div class="flex items-start justify-between gap-4">
									<div class="flex flex-col gap-1">
										<FormLabel>{{ t('songRequests.overlaySettings.hideOnPause') }}</FormLabel>
										<FormDescription>
											{{ t('songRequests.overlaySettings.hideOnPauseDescription') }}
										</FormDescription>
									</div>
									<FormControl>
										<Switch
											:model-value="value"
											@update:model-value="handleChange"
										/>
									</FormControl>
								</div>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>
				</div>

				<DialogFooter class="bg-background relative z-10 shrink-0 border-t p-6">
					<Button
						type="button"
						variant="outline"
						:disabled="isSaving"
						@click="setOpen(false)"
					>
						{{ t('songRequests.overlaySettings.cancel') }}
					</Button>
					<Button
						type="submit"
						:disabled="isSaving || !initializedForOpen"
					>
						<Icon
							v-if="isSaving"
							name="lucide:loader-circle"
							class="animate-spin"
							data-icon="inline-start"
						/>
						{{ t('songRequests.overlaySettings.save') }}
					</Button>
				</DialogFooter>
			</form>

			<div
				v-else
				class="min-h-64"
			/>
		</DialogContent>
	</Dialog>
</template>
