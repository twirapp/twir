<script setup lang="ts">
import { CheckCircle2, InfoIcon, XCircle } from 'lucide-vue-next'


import { useObsForm } from '~/features/overlays/obs/composables/use-obs-form'






import { toast } from 'vue-sonner'
import { useCopyOverlayLink } from '#layers/dashboard/components/overlays/copyOverlayLink'

const { t } = useI18n()
const { onSubmit, isLoading, isSaving, settings, isConnected } = useObsForm()
const { copyOverlayLink } = useCopyOverlayLink('obs')

async function handleSave() {
	await onSubmit()
	toast.success(t('sharedTexts.saved'), {
		description: t('overlays.obs.settingsSavedHint'),
	})
}
</script>

<template>
	<div class="flex flex-col gap-4 p-4">
		<UiCard>
			<UiCardHeader class="flex flex-row flex-wrap items-center justify-between space-y-0 pb-4">
				<div class="flex items-center flex-wrap gap-3">
					<h2 class="text-2xl font-bold">
						{{ t('overlays.obs.title') }}
					</h2>
					<UiBadge v-if="isConnected" variant="default" class="bg-green-600 hover:bg-green-600">
						<CheckCircle2 class="h-3 w-3 mr-1" />
						{{ t('overlays.obs.connected') }}
					</UiBadge>
					<UiBadge v-else variant="destructive">
						<XCircle class="h-3 w-3 mr-1" />
						{{ t('overlays.obs.notConnected') }}
					</UiBadge>
				</div>
				<div class="flex gap-2">
					<UiButton variant="outline" :disabled="isLoading" @click="copyOverlayLink()">
						{{ t('overlays.copyOverlayLink') }}
					</UiButton>
					<UiButton :disabled="isLoading || isSaving" @click="handleSave">
						{{ t('sharedButtons.save') }}
					</UiButton>
				</div>
			</UiCardHeader>

			<UiCardContent class="space-y-6">
				<UiAlert>
					<InfoIcon class="h-4 w-4" />
					<UiAlertDescription>
						{{ t('overlays.obs.description') }}
					</UiAlertDescription>
				</UiAlert>

				<form class="space-y-4" @submit.prevent="handleSave">
					<UiFormField v-slot="{ componentField }" name="serverAddress">
						<UiFormItem>
							<UiFormLabel>{{ t('overlays.obs.address') }}</UiFormLabel>
							<UiFormControl>
								<UiInput type="text" placeholder="localhost" v-bind="componentField" />
							</UiFormControl>
							<UiFormDescription>
								{{ t('overlays.obs.addressDescription') }}
							</UiFormDescription>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>

					<UiFormField v-slot="{ componentField }" name="serverPort">
						<UiFormItem>
							<UiFormLabel>{{ t('overlays.obs.port') }}</UiFormLabel>
							<UiFormControl>
								<UiInput type="number" placeholder="4455" v-bind="componentField" />
							</UiFormControl>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>

					<UiFormField v-slot="{ componentField }" name="serverPassword">
						<UiFormItem>
							<UiFormLabel>{{ t('overlays.obs.password') }}</UiFormLabel>
							<UiFormControl>
								<UiInput type="password" placeholder="Password" v-bind="componentField" />
							</UiFormControl>
							<UiFormDescription>
								{{ t('overlays.obs.passwordDescription') }}
							</UiFormDescription>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>
				</form>

				<!-- Connection info -->
				<div
					v-if="
						settings &&
						(settings.scenes.length || settings.sources.length || settings.audioSources.length)
					"
					class="space-y-4"
				>
					<div class="border-t pt-4">
						<h3 class="text-lg font-semibold mb-3">{{ t('overlays.obs.connectionInfo') }}</h3>
						<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
							<UiCard class="bg-muted/50">
								<UiCardContent class="pt-4">
									<div class="text-2xl font-bold">{{ settings.scenes?.length ?? 0 }}</div>
									<p class="text-sm text-muted-foreground">{{ t('overlays.obs.scenesCount') }}</p>
								</UiCardContent>
							</UiCard>
							<UiCard class="bg-muted/50">
								<UiCardContent class="pt-4">
									<div class="text-2xl font-bold">{{ settings.sources?.length ?? 0 }}</div>
									<p class="text-sm text-muted-foreground">{{ t('overlays.obs.sourcesCount') }}</p>
								</UiCardContent>
							</UiCard>
							<UiCard class="bg-muted/50">
								<UiCardContent class="pt-4">
									<div class="text-2xl font-bold">{{ settings.audioSources?.length ?? 0 }}</div>
									<p class="text-sm text-muted-foreground">
										{{ t('overlays.obs.audioSourcesCount') }}
									</p>
								</UiCardContent>
							</UiCard>
						</div>
					</div>
				</div>
			</UiCardContent>
		</UiCard>
	</div>
</template>
