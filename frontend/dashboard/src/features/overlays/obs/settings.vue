<script setup lang="ts">
import { CheckCircle2, InfoIcon, XCircle } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { useObsForm } from '@/features/overlays/obs/composables/use-obs-form'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useToast } from '@/components/ui/toast'
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink'

const { t } = useI18n()
const { toast } = useToast()
const { onSubmit, isLoading, isSaving, settings, isConnected } = useObsForm()
const { copyOverlayLink } = useCopyOverlayLink('obs')

async function handleSave() {
	await onSubmit()
	toast({
		title: t('sharedTexts.saved'),
		description: t('overlays.obs.settingsSavedHint'),
		variant: 'default',
	})
}
</script>

<template>
	<div class="flex flex-col gap-4 p-4">
		<Card>
			<CardHeader class="flex flex-row flex-wrap items-center justify-between space-y-0 pb-4">
				<div class="flex items-center flex-wrap gap-3">
					<h2 class="text-2xl font-bold">
						{{ t('overlays.obs.title') }}
					</h2>
					<Badge v-if="isConnected" variant="default" class="bg-green-600 hover:bg-green-600">
						<CheckCircle2 class="h-3 w-3 mr-1" />
						{{ t('overlays.obs.connected') }}
					</Badge>
					<Badge v-else variant="destructive">
						<XCircle class="h-3 w-3 mr-1" />
						{{ t('overlays.obs.notConnected') }}
					</Badge>
				</div>
				<div class="flex gap-2">
					<Button variant="outline" :disabled="isLoading" @click="copyOverlayLink()">
						{{ t('overlays.copyOverlayLink') }}
					</Button>
					<Button :disabled="isLoading || isSaving" @click="handleSave">
						{{ t('sharedButtons.save') }}
					</Button>
				</div>
			</CardHeader>

			<CardContent class="space-y-6">
				<Alert>
					<InfoIcon class="h-4 w-4" />
					<AlertDescription>
						{{ t('overlays.obs.description') }}
					</AlertDescription>
				</Alert>

				<form class="space-y-4" @submit.prevent="handleSave">
					<FormField v-slot="{ componentField }" name="serverAddress">
						<FormItem>
							<FormLabel>{{ t('overlays.obs.address') }}</FormLabel>
							<FormControl>
								<Input type="text" placeholder="localhost" v-bind="componentField" />
							</FormControl>
							<FormDescription>
								{{ t('overlays.obs.addressDescription') }}
							</FormDescription>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="serverPort">
						<FormItem>
							<FormLabel>{{ t('overlays.obs.port') }}</FormLabel>
							<FormControl>
								<Input type="number" placeholder="4455" v-bind="componentField" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="serverPassword">
						<FormItem>
							<FormLabel>{{ t('overlays.obs.password') }}</FormLabel>
							<FormControl>
								<Input type="password" placeholder="Password" v-bind="componentField" />
							</FormControl>
							<FormDescription>
								{{ t('overlays.obs.passwordDescription') }}
							</FormDescription>
							<FormMessage />
						</FormItem>
					</FormField>
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
							<Card class="bg-muted/50">
								<CardContent class="pt-4">
									<div class="text-2xl font-bold">{{ settings.scenes?.length ?? 0 }}</div>
									<p class="text-sm text-muted-foreground">{{ t('overlays.obs.scenesCount') }}</p>
								</CardContent>
							</Card>
							<Card class="bg-muted/50">
								<CardContent class="pt-4">
									<div class="text-2xl font-bold">{{ settings.sources?.length ?? 0 }}</div>
									<p class="text-sm text-muted-foreground">{{ t('overlays.obs.sourcesCount') }}</p>
								</CardContent>
							</Card>
							<Card class="bg-muted/50">
								<CardContent class="pt-4">
									<div class="text-2xl font-bold">{{ settings.audioSources?.length ?? 0 }}</div>
									<p class="text-sm text-muted-foreground">
										{{ t('overlays.obs.audioSourcesCount') }}
									</p>
								</CardContent>
							</Card>
						</div>
					</div>
				</div>
			</CardContent>
		</Card>
	</div>
</template>
