<script setup lang="ts">
import { ExternalLink, Eye, EyeOff } from 'lucide-vue-next'
import { ref, watch } from 'vue'

import { useDonatepayIntegration } from '@/api/integrations/donatepay'
import { useIntegrationsPageData } from '@/api/integrations/integrations-page.ts'
import DonatePaySVG from '@/assets/integrations/donatepay.svg?use'
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue'
import WithSettings from '@/components/integrations/variants/withSettings.vue'
import { Button } from '@/components/ui/button'
import { FormItem } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'

function redirectToGetApiKey() {
	window.open('https://donatepay.ru/page/api', '_blank')
}

const integrationsPage = useIntegrationsPageData()
const manager = useDonatepayIntegration()
const { executeMutation } = manager.useUpdate()

const apiKey = ref<string>('')
const enabled = ref(true)
const showPassword = ref(false)

watch(
	() => integrationsPage.donatePayData.value,
	(value) => {
		if (value?.apiKey) {
			apiKey.value = value.apiKey
			enabled.value = value.enabled
		}
	},
	{ immediate: true }
)

async function save() {
	await executeMutation({
		apiKey: apiKey.value,
		enabled: enabled.value,
	})
}
</script>

<template>
	<WithSettings
		title="Donatepay"
		:save="save"
		:icon="DonatePaySVG"
		icon-width="80px"
		dialog-content-class="w-[600px]"
	>
		<template #description>
			<DonateDescription />
		</template>
		<template #settings>
			<div class="flex flex-col gap-2">
				<FormItem label="Api key">
					<div class="flex gap-2">
						<Button variant="outline" class="shrink-0" @click="redirectToGetApiKey">
							Get api key
							<ExternalLink class="ml-2 h-4 w-4" />
						</Button>
						<div class="relative flex-1">
							<Input
								v-model="apiKey"
								:type="showPassword ? 'text' : 'password'"
								placeholder="Api key"
								class="pr-10"
							/>
							<Button
								type="button"
								variant="ghost"
								size="icon"
								class="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
								@click="showPassword = !showPassword"
							>
								<Eye v-if="showPassword" class="h-4 w-4" />
								<EyeOff v-else class="h-4 w-4" />
								<span class="sr-only">
									{{ showPassword ? 'Hide password' : 'Show password' }}
								</span>
							</Button>
						</div>
					</div>
				</FormItem>
				<FormItem class="flex flex-row items-center gap-4">
					<span>Enabled</span>
					<Switch v-model:model-value="enabled" />
				</FormItem>
			</div>
		</template>
	</WithSettings>
</template>
