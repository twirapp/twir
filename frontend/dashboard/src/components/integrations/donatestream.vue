<script setup lang="ts">
import { ExternalLink, Info } from 'lucide-vue-next'
import { computed, ref } from 'vue'

import { useIntegrationsPageData } from '@/api/integrations/integrations-page.ts'
import { useIntegrations } from '@/api/integrations/integrations.ts'
import DonateStreamSVG from '@/assets/integrations/donatestream.svg?use'
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue'
import WithSettings from '@/components/integrations/variants/withSettings.vue'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import CopyInput from '@/components/ui/copy-input/CopyInput.vue'
import { Input } from '@/components/ui/input'

const integrationsPage = useIntegrationsPageData()
const manager = useIntegrations()
const { executeMutation } = manager.donateStreamPostCode()

const currentPageUrl = `${window.location.origin}/api/webhooks/integrations/donatestream`
const webhookUrl = computed(() => {
	return `${currentPageUrl}/${integrationsPage.donateStreamData.value?.integrationId}`
})

const secret = ref('')

async function saveSecret() {
	if (!secret.value) return
	await executeMutation({ secret: secret.value })
}
</script>

<template>
	<WithSettings
		title="Donate.stream"
		:icon="DonateStreamSVG"
		icon-width="9rem"
		dialog-content-class="w-[700px]"
	>
		<template #description>
			<DonateDescription />
		</template>

		<template #settings>
			<div class="space-y-8">
				<!-- Step 1 -->
				<div
					class="relative pl-8 before:absolute before:left-3 before:top-[5px] before:h-full before:w-[1px] before:bg-border"
				>
					<div
						class="absolute left-0 top-1 flex h-6 w-6 items-center justify-center rounded-full border bg-background text-sm font-medium"
					>
						1
					</div>
					<div class="space-y-2">
						<div class="flex items-center gap-2">
							<h4 class="font-medium leading-none">Copy webhook URL</h4>
							<a
								href="https://lk.donate.stream/settings/api"
								target="_blank"
								class="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
							>
								Open donate.stream
								<ExternalLink class="h-3 w-3" />
							</a>
						</div>
						<p class="text-sm text-muted-foreground">
							Paste this link into the input field on donate.stream settings page
						</p>
						<CopyInput
							:text="webhookUrl"
							:disabled="!integrationsPage.donateStreamData.value?.integrationId"
							class="relative"
						/>
					</div>
				</div>

				<!-- Step 2 -->
				<div
					class="relative pl-8 before:absolute before:left-3 before:top-[5px] before:h-full before:w-[1px] before:bg-border"
				>
					<div
						class="absolute left-0 top-1 flex h-6 w-6 items-center justify-center rounded-full border bg-background text-sm font-medium"
					>
						2
					</div>
					<div class="space-y-2">
						<div class="flex items-center gap-2">
							<h4 class="font-medium leading-none">Enter secret key</h4>
							<a
								href="https://i.imgur.com/OtW97pV.png"
								target="_blank"
								class="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
							>
								View example
								<ExternalLink class="h-3 w-3" />
							</a>
						</div>
						<p class="text-sm text-muted-foreground">
							Copy the secret key from donate.stream and paste it below
						</p>
						<div class="flex gap-2">
							<Input
								v-model="secret"
								placeholder="Secret key from donate.stream"
								class="max-w-md"
							/>
							<Button variant="default" :disabled="!secret" @click="saveSecret"> Save </Button>
						</div>
					</div>
				</div>

				<!-- Step 3 -->
				<div class="relative pl-8">
					<div
						class="absolute left-0 top-1 flex h-6 w-6 items-center justify-center rounded-full border bg-background text-sm font-medium"
					>
						3
					</div>
					<div class="space-y-2">
						<h4 class="font-medium leading-none">Confirm integration</h4>
						<p class="text-sm text-muted-foreground">
							Return to donate.stream and click the "Confirm" button to complete the setup
						</p>
					</div>
				</div>

				<Alert>
					<Info class="h-4 w-4" />
					<AlertDescription>
						After completing these steps, your Donate.stream integration will be active
					</AlertDescription>
				</Alert>
			</div>
		</template>
	</WithSettings>
</template>
