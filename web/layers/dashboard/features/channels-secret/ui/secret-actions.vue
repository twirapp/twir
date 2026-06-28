<script setup lang="ts">
import { ref } from 'vue'

import { Button } from '@/components/ui/button'

import type { Secret } from '../composables/use-secrets-api'

import DeleteSecretDialog from './delete-secret-dialog.vue'
import SecretDialog from './secret-dialog.vue'

const props = defineProps<{
	secret: Secret
}>()

const isEditOpen = ref(false)
const isDeleteOpen = ref(false)
</script>

<template>
	<div class="flex items-center gap-2">
		<Button
			variant="secondary"
			size="icon"
			@click="isEditOpen = true"
		>
			<Icon
				name="lucide:pencil"
				class="h-4 w-4"
			/>
		</Button>
		<Button
			variant="destructive"
			size="icon"
			@click="isDeleteOpen = true"
		>
			<Icon
				name="lucide:trash"
				class="h-4 w-4"
			/>
		</Button>

		<SecretDialog
			v-model:open="isEditOpen"
			:secret="secret"
		/>

		<DeleteSecretDialog
			v-model:open="isDeleteOpen"
			:secret="secret"
		/>
	</div>
</template>
