<script setup lang="ts">
import { ExternalLink } from 'lucide-vue-next'

import { useDonatelloIntegration } from '@/api/index.js'
import DonatelloSVG from '@/assets/integrations/donatello.svg?use'
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue'
import WithSettings from '@/components/integrations/variants/withSettings.vue'
import { Alert, AlertDescription } from '@/components/ui/alert'
import CopyInput from '@/components/ui/copy-input/CopyInput.vue'

const { data: donatelloData } = useDonatelloIntegration()

const webhookUrl = `${window.location.origin}/api/webhooks/integrations/donatello`
</script>

<template>
	<WithSettings
		title="Donatello"
		:icon="DonatelloSVG"
		icon-width="48px"
	>
		<template #description>
			<DonateDescription />
		</template>
		<template #settings>
			<div class="flex flex-col gap-6">
				<!-- Step 1 -->
				<div class="flex flex-col gap-2">
					<div class="flex items-center gap-2">
						<div class="flex h-8 w-8 items-center justify-center rounded-full bg-muted text-sm font-medium">
							1
						</div>
						<h3 class="font-medium">
							Go to Donatello settings
						</h3>
					</div>
					<Alert>
						<AlertDescription class="flex items-center gap-2">
							<a
								href="https://donatello.to/panel/settings"
								target="_blank"
								rel="noopener noreferrer"
								class="inline-flex items-center gap-1 font-medium underline underline-offset-4"
							>
								https://donatello.to/panel/settings
								<ExternalLink class="h-4 w-4" />
							</a>
							and scroll to "Вихідний API" section
						</AlertDescription>
					</Alert>
				</div>

				<!-- Step 2 -->
				<div class="flex flex-col gap-2">
					<div class="flex items-center gap-2">
						<div class="flex h-8 w-8 items-center justify-center rounded-full bg-muted text-sm font-medium">
							2
						</div>
						<h3 class="font-medium">
							Copy API key
						</h3>
					</div>
					<div class="flex flex-col gap-2 pl-10">
						<span class="text-sm text-muted-foreground">Copy api key and paste into "Api Key" input</span>
						<CopyInput :text="donatelloData?.integrationId ?? ''" />
					</div>
				</div>

				<!-- Step 3 -->
				<div class="flex flex-col gap-2">
					<div class="flex items-center gap-2">
						<div class="flex h-8 w-8 items-center justify-center rounded-full bg-muted text-sm font-medium">
							3
						</div>
						<h3 class="font-medium">
							Set webhook URL
						</h3>
					</div>
					<div class="flex flex-col gap-2 pl-10">
						<span class="text-sm text-muted-foreground">Copy link and paste into "Link" field</span>
						<CopyInput :text="webhookUrl" />
					</div>
				</div>
			</div>
		</template>
	</WithSettings>
</template>
