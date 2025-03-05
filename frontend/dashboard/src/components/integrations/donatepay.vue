<script setup lang="ts">
import { ExternalLink, Eye, EyeOff } from 'lucide-vue-next'
import { ref, watch } from 'vue'

import { useDonatepayIntegration } from '@/api/index.js'
import DonatePaySVG from '@/assets/integrations/donatepay.svg?use'
import DonateDescription from '@/components/integrations/helpers/donateDescription.vue'
import WithSettings from '@/components/integrations/variants/withSettings.vue'
import { Button } from '@/components/ui/button'
import { FormItem } from '@/components/ui/form'
import { Input } from '@/components/ui/input'

function redirectToGetApiKey() {
	window.open('https://donatepay.ru/page/api', '_blank')
}

const manager = useDonatepayIntegration()
const { data } = manager.useGetData()
const { mutateAsync } = manager.usePost()

const apiKey = ref<string>('')
const showPassword = ref(false)

watch(data, (value) => {
	if (value?.apiKey) {
		apiKey.value = value.apiKey
	}
})

async function save() {
	await mutateAsync(apiKey.value)
}
</script>

<template>
	<WithSettings
		title="Donatepay"
		:save="save"
		:icon="DonatePaySVG"
		icon-width="80px"
	>
		<template #description>
			<DonateDescription />
		</template>
		<template #settings>
			<FormItem label="Api key">
				<div class="flex gap-2">
					<Button
						variant="outline"
						class="shrink-0"
						@click="redirectToGetApiKey"
					>
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
							<Eye
								v-if="showPassword"
								class="h-4 w-4"
							/>
							<EyeOff
								v-else
								class="h-4 w-4"
							/>
							<span class="sr-only">
								{{ showPassword ? 'Hide password' : 'Show password' }}
							</span>
						</Button>
					</div>
				</div>
			</FormItem>
		</template>
	</WithSettings>
</template>
