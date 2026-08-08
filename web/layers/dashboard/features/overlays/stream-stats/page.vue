<script setup lang="ts">
import { VueMonacoEditor } from '@guolao/vue-monaco-editor'
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import { useProfile, useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth'
import { useCopyOverlayLink } from '~~/layers/dashboard/components/overlays/copyOverlayLink.js'

import { useStreamStatsOverlayApi } from './api'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import {
	ChannelRolePermissionEnum,
	StreamStatsOverlayDesign,
	StreamStatsOverlayViewersMode,
} from '~/gql/graphql.js'
import { StreamStatsOverlayUpdateInputSchema } from '~/gql/validation-schemas.js'

const { t } = useI18n()

const { data: profile } = useProfile()

const defaultSettings = {
	design: StreamStatsOverlayDesign.Bar,
	viewersEnabled: true,
	viewersMode: StreamStatsOverlayViewersMode.Cumulative,
	messagesEnabled: true,
	uptimeEnabled: true,
	subscribersEnabled: true,
	followersEnabled: true,
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
			viewersEnabled: overlay.viewersEnabled,
			viewersMode: overlay.viewersMode,
			messagesEnabled: overlay.messagesEnabled,
			uptimeEnabled: overlay.uptimeEnabled,
			subscribersEnabled: overlay.subscribersEnabled,
			followersEnabled: overlay.followersEnabled,
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
				viewersEnabled: values.viewersEnabled,
				viewersMode: values.viewersMode as StreamStatsOverlayViewersMode,
				messagesEnabled: values.messagesEnabled,
				uptimeEnabled: values.uptimeEnabled,
				subscribersEnabled: values.subscribersEnabled,
				followersEnabled: values.followersEnabled,
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
	{ value: StreamStatsOverlayDesign.Bar, labelKey: 'bar' },
	{ value: StreamStatsOverlayDesign.Cards, labelKey: 'cards' },
	{ value: StreamStatsOverlayDesign.Minimal, labelKey: 'minimal' },
] as const

const counters = [
	{ name: 'viewersEnabled', icon: 'lucide:users', labelKey: 'viewers' },
	{ name: 'messagesEnabled', icon: 'lucide:message-square', labelKey: 'messages' },
	{ name: 'uptimeEnabled', icon: 'lucide:clock', labelKey: 'uptime' },
	{ name: 'subscribersEnabled', icon: 'lucide:star', labelKey: 'subscribers' },
	{ name: 'followersEnabled', icon: 'lucide:user-plus', labelKey: 'followers' },
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
									<div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
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
												<!-- bar design mock -->
												<div
													v-if="option.value === StreamStatsOverlayDesign.Bar"
													class="bg-muted flex h-7 items-center gap-2 rounded-full border px-3"
												>
													<template
														v-for="i of 3"
														:key="i"
													>
														<span class="bg-primary/70 size-1.5 rounded-full" />
														<span class="bg-muted-foreground/40 h-1 w-5 rounded" />
													</template>
												</div>

												<!-- cards design mock -->
												<div
													v-else-if="option.value === StreamStatsOverlayDesign.Cards"
													class="flex gap-1.5"
												>
													<div
														v-for="i of 3"
														:key="i"
														class="bg-muted flex h-8 w-9 flex-col items-center justify-center gap-1 rounded-md border"
													>
														<span class="bg-primary/70 size-1.5 rounded-full" />
														<span class="bg-muted-foreground/40 h-1 w-5 rounded" />
													</div>
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
							{{ t('overlays.streamStats.settings.counters.label') }}
						</h3>

						<div class="space-y-4">
							<template
								v-for="counter of counters"
								:key="counter.name"
							>
								<FormField
									v-slot="{ value, handleChange }"
									:name="counter.name"
								>
									<FormItem>
										<div class="flex items-center justify-between gap-2">
											<FormLabel class="mt-0! flex items-center gap-2">
												<Icon
													:name="counter.icon"
													class="text-muted-foreground size-4"
												/>
												{{
													t(`overlays.streamStats.settings.counters.${counter.labelKey}`)
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
									v-if="counter.name === 'viewersEnabled' && form.values.viewersEnabled"
									v-slot="{ value, handleChange }"
									name="viewersMode"
								>
									<FormItem class="border-l pl-4 ml-2">
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
							</template>

							<p class="text-muted-foreground text-xs">
								{{ t('overlays.streamStats.settings.counters.twitchOnlyNote') }}
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
