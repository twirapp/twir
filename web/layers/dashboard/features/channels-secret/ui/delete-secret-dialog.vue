<script setup lang="ts">
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
const deleteMutation = secretsApi.useMutationRemoveSecret()

async function handleDelete() {
	await deleteMutation.executeMutation({
		id: props.secret.id,
	})

	open.value = false
}
</script>

<template>
	<Dialog v-model:open="open">
		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>Delete Secret</DialogTitle>
				<DialogDescription>
					Are you sure you want to delete <strong>{{ secret.name }}</strong>?
					This action cannot be undone.
				</DialogDescription>
			</DialogHeader>

			<DialogFooter>
				<Button type="button" variant="outline" @click="open = false">
					Cancel
				</Button>
				<Button
					type="button"
					variant="destructive"
					@click="handleDelete"
				>
					Delete
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>
