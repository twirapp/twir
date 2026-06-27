<script setup lang="ts">
import { ref, watch } from 'vue'

import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'

import type { Secret } from '../composables/use-secrets-api'

import { useSecretsApi } from '../composables/use-secrets-api'

const props = defineProps<{
	secret: Secret
}>()

const open = defineModel<boolean>('open', { default: false })

const secretsApi = useSecretsApi()
const { data, executeQuery } = secretsApi.useQuerySecretValue(props.secret.id)

const isRevealed = ref(false)
const secretValue = ref('')

watch(open, async (isOpen) => {
	if (isOpen) {
		isRevealed.value = false
		secretValue.value = ''
		await executeQuery()
	}
})

watch(data, (newData) => {
	if (newData?.secretValue) {
		secretValue.value = newData.secretValue
	}
})

function toggleReveal() {
	isRevealed.value = !isRevealed.value
}

function copyToClipboard() {
	navigator.clipboard.writeText(secretValue.value)
}
</script>

<template>
	<Dialog v-model:open="open">
		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>Secret Value</DialogTitle>
				<DialogDescription>
					Viewing value for <strong>{{ secret.name }}</strong>
				</DialogDescription>
			</DialogHeader>

			<div class="space-y-4">
				<div class="relative">
					<div class="bg-muted flex items-center gap-2 rounded-md p-3 font-mono text-sm">
						<span class="flex-1 overflow-hidden text-ellipsis">
							{{ isRevealed ? secretValue : '••••••••' }}
						</span>
						<Button
							variant="ghost"
							size="icon"
							@click="toggleReveal"
						>
							<Icon
								:name="isRevealed ? 'lucide:eye-off' : 'lucide:eye'"
								class="h-4 w-4"
							/>
						</Button>
						<Button
							variant="ghost"
							size="icon"
							@click="copyToClipboard"
						>
							<Icon
								name="copy"
								class="h-4 w-4"
							/>
						</Button>
					</div>
				</div>
			</div>

			<DialogFooter>
				<Button
					type="button"
					variant="outline"
					@click="open = false"
				>
					Close
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>
