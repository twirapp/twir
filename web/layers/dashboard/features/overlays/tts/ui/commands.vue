<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { InfoIcon } from 'lucide-vue-next'

const { t } = useI18n()
</script>

<template>
	<div class="flex flex-col gap-4 p-4">
		<Card>
			<CardHeader>
				<CardTitle>{{ t('overlays.tts.tabs.commands') }}</CardTitle>
				<CardDescription>{{ t('overlays.tts.commandsDescription') }}</CardDescription>
			</CardHeader>

			<CardContent class="space-y-4">
				<Alert>
					<InfoIcon class="h-4 w-4" />
					<AlertTitle>{{ t('overlays.tts.commands.info') }}</AlertTitle>
					<AlertDescription>
						{{ t('overlays.tts.commands.infoDescription') }}
					</AlertDescription>
				</Alert>

				<!-- Main TTS Command -->
				<div class="space-y-2">
					<h3 class="text-lg font-semibold">{{ t('overlays.tts.commands.main') }}</h3>
					<div class="rounded-lg border p-4">
						<code class="text-sm">!tts [text]</code>
						<p class="mt-2 text-sm text-muted-foreground">
							{{ t('overlays.tts.commands.mainDescription') }}
						</p>
					</div>
				</div>

				<!-- TTS Settings Command -->
				<div class="space-y-2">
					<h3 class="text-lg font-semibold">{{ t('overlays.tts.commands.settings') }}</h3>
					<div class="rounded-lg border p-4 space-y-2">
						<div>
							<code class="text-sm">!tts voice [voice_name]</code>
							<p class="mt-1 text-sm text-muted-foreground">
								{{ t('overlays.tts.commands.voiceDescription') }}
							</p>
						</div>
						<div>
							<code class="text-sm">!tts rate [0-100]</code>
							<p class="mt-1 text-sm text-muted-foreground">
								{{ t('overlays.tts.commands.rateDescription') }}
							</p>
						</div>
						<div>
							<code class="text-sm">!tts pitch [0-100]</code>
							<p class="mt-1 text-sm text-muted-foreground">
								{{ t('overlays.tts.commands.pitchDescription') }}
							</p>
						</div>
						<div>
							<code class="text-sm">!tts volume [0-100]</code>
							<p class="mt-1 text-sm text-muted-foreground">
								{{ t('overlays.tts.commands.volumeDescription') }}
							</p>
						</div>
					</div>
				</div>

				<!-- Skip Command -->
				<div class="space-y-2">
					<h3 class="text-lg font-semibold">{{ t('overlays.tts.commands.skip') }}</h3>
					<div class="rounded-lg border p-4">
						<code class="text-sm">!skip</code>
						<p class="mt-2 text-sm text-muted-foreground">
							{{ t('overlays.tts.commands.skipDescription') }}
						</p>
					</div>
				</div>

				<!-- Example Usage -->
				<div class="space-y-2">
					<h3 class="text-lg font-semibold">{{ t('overlays.tts.commands.examples') }}</h3>
					<div class="rounded-lg border p-4 space-y-2">
						<div>
							<code class="text-sm bg-muted px-2 py-1 rounded">!tts Hello everyone!</code>
						</div>
						<div>
							<code class="text-sm bg-muted px-2 py-1 rounded">!tts voice elena</code>
						</div>
						<div>
							<code class="text-sm bg-muted px-2 py-1 rounded">!tts rate 75</code>
						</div>
					</div>
				</div>
			</CardContent>
		</Card>
	</div>
</template>

