<script setup lang="ts">
import { VueMonacoEditor } from '@guolao/vue-monaco-editor'
import { useForm } from 'vee-validate'
import { VueDraggable } from 'vue-draggable-plus'
import { toast } from 'vue-sonner'
import { useProfile, useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth'
import { useCopyOverlayLink } from '~~/layers/dashboard/components/overlays/copyOverlayLink.js'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { ColorPicker } from '@/components/ui/color-picker'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import InputWithIcon from '@/components/ui/InputWithIcon/InputWithIcon.vue'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import {
	ChannelRolePermissionEnum,
	StreamStatsOverlayCounter,
	StreamStatsOverlayDesign,
	StreamStatsOverlayVariant,
	StreamStatsOverlayViewersMode,
} from '~/gql/graphql.js'
import { StreamStatsOverlayUpdateInputSchema } from '~/gql/validation-schemas.js'

import { useStreamStatsOverlayApi } from './api'

const { t } = useI18n()

const { data: profile } = useProfile()

const defaultSettings = {
	design: StreamStatsOverlayDesign.Glass,
	variant: StreamStatsOverlayVariant.Horizontal,
	viewersEnabled: true,
	viewersMode: StreamStatsOverlayViewersMode.Cumulative,
	platformIconsEnabled: false,
	messagesEnabled: true,
	uptimeEnabled: true,
	subscribersEnabled: true,
	followersEnabled: true,
	counterOrder: [
		StreamStatsOverlayCounter.Viewers,
		StreamStatsOverlayCounter.Messages,
		StreamStatsOverlayCounter.Uptime,
		StreamStatsOverlayCounter.Subscribers,
		StreamStatsOverlayCounter.Followers,
	],
	viewersColor: '',
	messagesColor: '',
	uptimeColor: '',
	subscribersColor: '',
	followersColor: '',
	customHtmlEnabled: false,
	customHtml: '',
	customCss: '',
}

const form = useForm({
	validationSchema: StreamStatsOverlayUpdateInputSchema,
	initialValues: defaultSettings,
	keepValuesOnUnmount: true,
})

const api = useStreamStatsOverlayApi()
const {
	data: settings,
	fetching: isSettingsLoading,
	error: settingsError,
} = api.useQueryStreamStats()
const updateMutation = api.useMutationUpdateStreamStats()

const isSettingsError = computed(() => !!settingsError.value)

const selectedDashboard = computed(() => {
	return profile.value?.availableDashboards.find(
		(dashboard) => dashboard.id === profile.value?.selectedDashboardId
	)
})

const selectedDashboardApiKey = computed(() => {
	return selectedDashboard.value?.channelApiKey || profile.value?.channelApiKey || ''
})

watch(
	settings,
	(v) => {
		if (!v?.overlaysStreamStats) return
		const overlay = v.overlaysStreamStats
		form.setValues({
			design: overlay.design,
			variant: overlay.variant,
			viewersEnabled: overlay.viewersEnabled,
			viewersMode: overlay.viewersMode,
			platformIconsEnabled: overlay.platformIconsEnabled,
			messagesEnabled: overlay.messagesEnabled,
			uptimeEnabled: overlay.uptimeEnabled,
			subscribersEnabled: overlay.subscribersEnabled,
			followersEnabled: overlay.followersEnabled,
			counterOrder: overlay.counterOrder,
			viewersColor: overlay.viewersColor,
			messagesColor: overlay.messagesColor,
			uptimeColor: overlay.uptimeColor,
			subscribersColor: overlay.subscribersColor,
			followersColor: overlay.followersColor,
			customHtmlEnabled: overlay.customHtmlEnabled,
			customHtml: overlay.customHtml,
			customCss: overlay.customCss,
		})
	},
	{ immediate: true }
)

const iframeRef = ref<HTMLIFrameElement | null>(null)
const requestUrl = useRequestURL()
const iframeUrl = computed(() => {
	if (!selectedDashboardApiKey.value) return null

	return `${requestUrl.origin}/overlays/${selectedDashboardApiKey.value}/stream-stats?preview=1`
})

function sendIframeMessage(key: string, data?: any) {
	if (!iframeRef.value) return
	const win = iframeRef.value

	win.contentWindow?.postMessage(
		JSON.stringify({
			key,
			data: toRaw(data),
		})
	)
}

function sendSettings() {
	sendIframeMessage('settings', toRaw(form.values))
}

watch(iframeRef, (v) => {
	if (!v) return

	v.contentWindow?.addEventListener('message', (e) => {
		const parsed = JSON.parse(e.data)
		if (parsed.key !== 'getSettings') return

		sendSettings()
	})
})

watch(
	() => form.values,
	() => {
		if (!iframeRef.value) return

		sendSettings()
	},
	{ deep: true }
)

const { canCopyOverlayLink, copyOverlayLink } = useCopyOverlayLink('stream-stats')

const save = form.handleSubmit(async (values) => {
	try {
		await updateMutation.executeMutation({
			input: {
				design: values.design as StreamStatsOverlayDesign,
				variant: values.variant as StreamStatsOverlayVariant,
				viewersEnabled: values.viewersEnabled,
				viewersMode: values.viewersMode as StreamStatsOverlayViewersMode,
				platformIconsEnabled: values.platformIconsEnabled,
				messagesEnabled: values.messagesEnabled,
				uptimeEnabled: values.uptimeEnabled,
				subscribersEnabled: values.subscribersEnabled,
				followersEnabled: values.followersEnabled,
				counterOrder: values.counterOrder as StreamStatsOverlayCounter[],
				viewersColor: values.viewersColor,
				messagesColor: values.messagesColor,
				uptimeColor: values.uptimeColor,
				subscribersColor: values.subscribersColor,
				followersColor: values.followersColor,
				customHtmlEnabled: values.customHtmlEnabled,
				customHtml: values.customHtml,
				customCss: values.customCss,
			},
		})

		toast.success(t('sharedTexts.saved'))
	} catch (e) {
		toast.error('Error occurred while saving Stream Stats overlay')
		console.error(e)
	}
})

const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)

const canCopyLink = computed(() => {
	return canCopyOverlayLink.value && userCanEditOverlays.value
})

function setDefaultSettings() {
	form.setValues(structuredClone(defaultSettings))
}

const designOptions = [
	{ value: StreamStatsOverlayDesign.Glass, labelKey: 'glass' },
	{ value: StreamStatsOverlayDesign.Cards, labelKey: 'cards' },
	{ value: StreamStatsOverlayDesign.Neon, labelKey: 'neon' },
	{ value: StreamStatsOverlayDesign.Solid, labelKey: 'solid' },
	{ value: StreamStatsOverlayDesign.Minimal, labelKey: 'minimal' },
	{ value: StreamStatsOverlayDesign.Terminal, labelKey: 'terminal' },
	{ value: StreamStatsOverlayDesign.Outline, labelKey: 'outline' },
] as const

const cardsMockAccents = ['#8b5cf6', '#38bdf8', '#fbbf24'] as const

const variantOptions = [
	{ value: StreamStatsOverlayVariant.Horizontal, labelKey: 'horizontal' },
	{ value: StreamStatsOverlayVariant.HorizontalCompact, labelKey: 'horizontalCompact' },
	{ value: StreamStatsOverlayVariant.Vertical, labelKey: 'vertical' },
	{ value: StreamStatsOverlayVariant.VerticalCompact, labelKey: 'verticalCompact' },
	{ value: StreamStatsOverlayVariant.Large, labelKey: 'large' },
] as const

const countersMeta = {
	[StreamStatsOverlayCounter.Viewers]: {
		enabledField: 'viewersEnabled',
		icon: 'lucide:users',
		labelKey: 'viewers',
	},
	[StreamStatsOverlayCounter.Messages]: {
		enabledField: 'messagesEnabled',
		icon: 'lucide:message-square',
		labelKey: 'messages',
	},
	[StreamStatsOverlayCounter.Uptime]: {
		enabledField: 'uptimeEnabled',
		icon: 'lucide:clock',
		labelKey: 'uptime',
	},
	[StreamStatsOverlayCounter.Subscribers]: {
		enabledField: 'subscribersEnabled',
		icon: 'lucide:star',
		labelKey: 'subscribers',
	},
	[StreamStatsOverlayCounter.Followers]: {
		enabledField: 'followersEnabled',
		icon: 'lucide:user-plus',
		labelKey: 'followers',
	},
} as const

const counterOrder = computed<StreamStatsOverlayCounter[]>({
	get: () => (form.values.counterOrder ?? []) as StreamStatsOverlayCounter[],
	set: (value) => form.setFieldValue('counterOrder', value),
})

const colorCounters = [
	{ name: 'viewersColor', icon: 'lucide:users', labelKey: 'viewers' },
	{ name: 'messagesColor', icon: 'lucide:message-square', labelKey: 'messages' },
	{ name: 'uptimeColor', icon: 'lucide:clock', labelKey: 'uptime' },
	{ name: 'subscribersColor', icon: 'lucide:star', labelKey: 'subscribers' },
	{ name: 'followersColor', icon: 'lucide:user-plus', labelKey: 'followers' },
] as const

const viewersModeOptions = [
	{ value: StreamStatsOverlayViewersMode.Cumulative, labelKey: 'cumulative' },
	{ value: StreamStatsOverlayViewersMode.Separate, labelKey: 'separate' },
] as const

const templatePlaceholders = [
	'{{viewers}}',
	'{{messages}}',
	'{{uptime}}',
	'{{subscribers}}',
	'{{followers}}',
]

const monacoOptions = {
	automaticLayout: true,
	minimap: { enabled: false },
	fontSize: 14,
	lineNumbers: 'on' as const,
	scrollBeyondLastLine: false,
	wordWrap: 'on' as const,
	tabSize: 2,
}
</script>

<template>
	<div class="flex flex-col gap-4 p-4 lg:flex-row">
		<Card class="flex-1">
			<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-4">
				<h2 class="text-2xl font-bold">
					{{ t('overlays.streamStats.title', 'Stream Stats') }}
				</h2>
				<div class="flex gap-2">
					<Button
						variant="outline"
						@click="setDefaultSettings"
					>
						{{ t('sharedButtons.setDefaultSettings') }}
					</Button>
					<Button
						variant="outline"
						:disabled="isSettingsError || isSettingsLoading || !canCopyLink"
						@click="copyOverlayLink()"
					>
						{{ t('overlays.copyOverlayLink') }}
					</Button>
					<Button @click="save">
						{{ t('sharedButtons.save') }}
					</Button>
				</div>
			</CardHeader>

			<CardContent class="space-y-6">
				<form
					@submit="save"
					class="space-y-6"
				>
					<div class="space-y-4">
						<Separator />
						<h3 class="text-lg font-semibold">
							{{ t('overlays.streamStats.settings.design.label') }}
						</h3>

						<FormField
							v-slot="{ value, handleChange }"
							name="design"
						>
							<FormItem>
								<FormControl>
									<div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
										<button
											v-for="option of designOptions"
											:key="option.value"
											type="button"
											class="hover:border-primary/50 flex flex-col gap-2 rounded-lg border p-3 text-left transition-colors"
											:class="value === option.value ? 'ring-primary border-primary ring-2' : ''"
											@click="handleChange(option.value)"
										>
											<div
												class="bg-muted/30 flex h-16 items-center justify-center rounded-md border"
											>
												<!-- glass design mock -->
												<div
													v-if="option.value === StreamStatsOverlayDesign.Glass"
													class="bg-muted/60 flex h-7 items-center gap-2 rounded-full border px-3 backdrop-blur-sm"
												>
													<template
														v-for="i of 3"
														:key="i"
													>
														<span
															v-if="i > 1"
															class="bg-muted-foreground/30 h-3 w-px"
														/>
														<span class="bg-primary/70 size-1.5 rounded-full" />
														<span class="bg-muted-foreground/40 h-1 w-4 rounded" />
													</template>
												</div>

												<!-- cards design mock -->
												<div
													v-else-if="option.value === StreamStatsOverlayDesign.Cards"
													class="flex gap-1.5"
												>
													<div
														v-for="accent of cardsMockAccents"
														:key="accent"
														class="bg-muted flex h-8 w-9 flex-col items-center justify-center gap-1 rounded-md border border-t-2"
														:style="{ borderTopColor: accent }"
													>
														<span class="bg-primary/70 size-1.5 rounded-full" />
														<span class="bg-muted-foreground/40 h-1 w-5 rounded" />
													</div>
												</div>

												<!-- neon design mock -->
												<div
													v-else-if="option.value === StreamStatsOverlayDesign.Neon"
													class="bg-muted flex h-8 w-9 flex-col items-center justify-center gap-1 rounded-md border border-[#c084fc] shadow-[0_0_8px_rgba(192,132,252,0.5)]"
												>
													<span class="size-1.5 rounded-full bg-[#c084fc]" />
													<span class="bg-muted-foreground/40 h-1 w-5 rounded" />
												</div>

												<!-- solid design mock -->
												<div
													v-else-if="option.value === StreamStatsOverlayDesign.Solid"
													class="flex h-7 items-center gap-2 rounded-full bg-[linear-gradient(135deg,#7c3aed,#db2777)] px-3"
												>
													<template
														v-for="i of 3"
														:key="i"
													>
														<span class="size-1.5 rounded-full bg-white/80" />
														<span class="h-1 w-4 rounded bg-white/60" />
													</template>
												</div>

												<!-- terminal design mock -->
												<div
													v-else-if="option.value === StreamStatsOverlayDesign.Terminal"
													class="flex h-7 items-center rounded-md bg-[#0a0f0a] px-2"
												>
													<span class="font-mono text-[10px] text-[#22c55e]">> 1,342_</span>
												</div>

												<!-- outline design mock -->
												<div
													v-else-if="option.value === StreamStatsOverlayDesign.Outline"
													class="flex h-7 items-center gap-2 rounded-full border-2 border-[#8b5cf6] px-3"
												>
													<span class="size-1.5 rounded-full bg-[#8b5cf6]" />
													<span class="bg-muted-foreground/40 h-1 w-4 rounded" />
												</div>

												<!-- minimal design mock -->
												<div
													v-else
													class="flex flex-col items-start gap-1.5"
												>
													<span class="bg-muted-foreground/40 h-1.5 w-24 rounded" />
													<span class="bg-muted-foreground/40 h-1.5 w-16 rounded" />
												</div>
											</div>

											<div>
												<div class="text-sm font-medium">
													{{ t(`overlays.streamStats.settings.design.${option.labelKey}.name`) }}
												</div>
												<div class="text-muted-foreground text-xs">
													{{
														t(`overlays.streamStats.settings.design.${option.labelKey}.description`)
													}}
												</div>
											</div>
										</button>
									</div>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<Separator />
						<h3 class="text-lg font-semibold">
							{{ t('overlays.streamStats.settings.variant.label') }}
						</h3>

						<FormField
							v-slot="{ value, handleChange }"
							name="variant"
						>
							<FormItem>
								<FormControl>
									<div class="grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-5">
										<button
											v-for="option of variantOptions"
											:key="option.value"
											type="button"
											class="hover:border-primary/50 rounded-lg border p-3 text-left text-sm font-medium transition-colors"
											:class="value === option.value ? 'ring-primary border-primary ring-2' : ''"
											@click="handleChange(option.value)"
										>
											{{ t(`overlays.streamStats.settings.variant.${option.labelKey}`) }}
										</button>
									</div>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<Separator />
						<h3 class="text-lg font-semibold">
							{{ t('overlays.streamStats.settings.counters.label') }}
						</h3>

						<div class="space-y-4">
							<VueDraggable
								v-model="counterOrder"
								:animation="150"
								handle=".drag-handle"
								ghost-class="opacity-30"
								class="space-y-4"
							>
								<div
									v-for="counter of counterOrder"
									:key="counter"
									class="space-y-4"
								>
									<FormField
										v-slot="{ value, handleChange }"
										:name="countersMeta[counter].enabledField"
									>
										<FormItem>
											<div class="flex items-center justify-between gap-2">
												<FormLabel class="mt-0! flex items-center gap-2">
													<Icon
														name="lucide:grip-vertical"
														class="drag-handle text-muted-foreground size-4 cursor-grab"
													/>
													<Icon
														:name="countersMeta[counter].icon"
														class="text-muted-foreground size-4"
													/>
													{{
														t(
															`overlays.streamStats.settings.counters.${countersMeta[counter].labelKey}`
														)
													}}
												</FormLabel>
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

									<FormField
										v-if="
											counter === StreamStatsOverlayCounter.Viewers && form.values.viewersEnabled
										"
										v-slot="{ value, handleChange }"
										name="viewersMode"
									>
										<FormItem class="ml-2 border-l pl-4">
											<FormLabel>{{
												t('overlays.streamStats.settings.counters.viewersMode.label')
											}}</FormLabel>
											<FormControl>
												<div class="grid grid-cols-1 gap-2 sm:grid-cols-2">
													<button
														v-for="option of viewersModeOptions"
														:key="option.value"
														type="button"
														class="hover:border-primary/50 flex flex-col gap-0.5 rounded-lg border p-3 text-left transition-colors"
														:class="
															value === option.value ? 'ring-primary border-primary ring-2' : ''
														"
														@click="handleChange(option.value)"
													>
														<span class="text-sm font-medium">
															{{
																t(
																	`overlays.streamStats.settings.counters.viewersMode.${option.labelKey}.name`
																)
															}}
														</span>
														<span class="text-muted-foreground text-xs">
															{{
																t(
																	`overlays.streamStats.settings.counters.viewersMode.${option.labelKey}.description`
																)
															}}
														</span>
													</button>
												</div>
											</FormControl>
											<FormMessage />
										</FormItem>
									</FormField>

									<FormField
										v-if="
											counter === StreamStatsOverlayCounter.Viewers &&
											form.values.viewersEnabled &&
											form.values.viewersMode === StreamStatsOverlayViewersMode.Separate
										"
										v-slot="{ value, handleChange }"
										name="platformIconsEnabled"
									>
										<FormItem class="ml-2 border-l pl-4">
											<div class="flex items-center justify-between gap-2">
												<FormLabel class="mt-0!">
													{{ t('overlays.streamStats.settings.counters.platformIcons') }}
												</FormLabel>
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
							</VueDraggable>

							<p class="text-muted-foreground text-xs">
								{{ t('overlays.streamStats.settings.counters.dragHint') }}
							</p>

							<p class="text-muted-foreground text-xs">
								{{ t('overlays.streamStats.settings.counters.twitchOnlyNote') }}
							</p>
						</div>

						<Separator />
						<h3 class="text-lg font-semibold">
							{{ t('overlays.streamStats.settings.colors.label') }}
						</h3>

						<div class="space-y-4">
							<FormField
								v-for="colorCounter of colorCounters"
								:key="colorCounter.name"
								v-slot="{ componentField }"
								:name="colorCounter.name"
							>
								<FormItem>
									<FormLabel class="flex items-center gap-2">
										<Icon
											:name="colorCounter.icon"
											class="text-muted-foreground size-4"
										/>
										{{ t(`overlays.streamStats.settings.counters.${colorCounter.labelKey}`) }}
									</FormLabel>
									<FormControl>
										<InputWithIcon
											v-model="componentField.modelValue"
											@update:model-value="componentField['onUpdate:modelValue']"
										>
											<ColorPicker
												output-format="hex"
												v-model="componentField.modelValue"
												@update:model-value="componentField['onUpdate:modelValue']"
											/>
										</InputWithIcon>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<p class="text-muted-foreground text-xs">
								{{ t('overlays.streamStats.settings.colors.hint') }}
							</p>
						</div>

						<Separator />
						<h3 class="text-lg font-semibold">
							{{ t('overlays.streamStats.settings.customTemplate.label') }}
						</h3>

						<div class="space-y-4">
							<FormField
								v-slot="{ value, handleChange }"
								name="customHtmlEnabled"
							>
								<FormItem>
									<div class="flex items-center gap-2">
										<FormControl>
											<Switch
												:model-value="value"
												@update:model-value="handleChange"
											/>
										</FormControl>
										<FormLabel class="mt-0!">{{
											t('overlays.streamStats.settings.customTemplate.enabled')
										}}</FormLabel>
									</div>
									<FormMessage />
								</FormItem>
							</FormField>

							<template v-if="form.values.customHtmlEnabled">
								<p class="text-muted-foreground text-xs">
									{{ t('overlays.streamStats.settings.customTemplate.placeholders') }}
									<code
										v-for="placeholder of templatePlaceholders"
										:key="placeholder"
										class="bg-muted mx-0.5 rounded px-1 py-0.5 font-mono text-[11px]"
									>
										{{ placeholder }}
									</code>
									{{ t('overlays.streamStats.settings.customTemplate.replacesDesign') }}
								</p>

								<FormField
									v-slot="{ value, handleChange }"
									name="customHtml"
								>
									<FormItem>
										<FormLabel>{{
											t('overlays.streamStats.settings.customTemplate.html')
										}}</FormLabel>
										<FormControl>
											<div class="h-[200px] overflow-hidden rounded-md border">
												<VueMonacoEditor
													:value="value"
													language="html"
													theme="vs-dark"
													:options="monacoOptions"
													class="h-full"
													@update:value="handleChange"
												/>
											</div>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>

								<FormField
									v-slot="{ value, handleChange }"
									name="customCss"
								>
									<FormItem>
										<FormLabel>{{
											t('overlays.streamStats.settings.customTemplate.css')
										}}</FormLabel>
										<FormControl>
											<div class="h-[200px] overflow-hidden rounded-md border">
												<VueMonacoEditor
													:value="value"
													language="css"
													theme="vs-dark"
													:options="monacoOptions"
													class="h-full"
													@update:value="handleChange"
												/>
											</div>
										</FormControl>
										<FormMessage />
									</FormItem>
								</FormField>
							</template>
						</div>
					</div>
				</form>
			</CardContent>
		</Card>

		<div class="relative flex-1">
			<div class="sticky top-4">
				<iframe
					v-if="iframeUrl"
					ref="iframeRef"
					:src="iframeUrl"
					class="h-[600px] w-full rounded-lg border"
				/>
			</div>
		</div>
	</div>
</template>
