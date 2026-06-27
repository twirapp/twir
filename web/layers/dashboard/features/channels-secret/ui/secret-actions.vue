<script setup lang="ts">
import { ref } from 'vue'

import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
// import { Edit, Eye, MoreHorizontal, Trash } from 'lucide-vue-next'

import type { Secret } from '../composables/use-secrets-api'

import { useSecretsApi } from '../composables/use-secrets-api'
import DeleteSecretDialog from './delete-secret-dialog.vue'
import SecretDialog from './secret-dialog.vue'
import SecretRevealDialog from './secret-reveal-dialog.vue'

const props = defineProps<{
	secret: Secret
}>()

const isEditOpen = ref(false)
const isRevealOpen = ref(false)
const isDeleteOpen = ref(false)
</script>

<template>
	<div class="flex items-center gap-2">
		<Button
			variant="ghost"
			size="icon"
			@click="isRevealOpen = true"
		>
			<Icon
				name="lucide:eye"
				class="h-4 w-4"
			/>
		</Button>

		<DropdownMenu>
			<DropdownMenuTrigger as-child>
				<Button
					variant="ghost"
					size="icon"
				>
					<Icon
						name="lucide:more-horizontal"
						class="h-4 w-4"
					/>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="end">
				<DropdownMenuItem @click="isEditOpen = true">
					<Icon
						name="lucide:edit"
						class="mr-2 h-4 w-4"
					/>
					Edit
				</DropdownMenuItem>
				<DropdownMenuItem
					class="text-destructive"
					@click="isDeleteOpen = true"
				>
					<Icon
						name="lucide:trash"
						class="mr-2 h-4 w-4"
					/>
					Delete
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>

		<SecretDialog
			v-model:open="isEditOpen"
			:secret="secret"
		/>

		<SecretRevealDialog
			v-model:open="isRevealOpen"
			:secret="secret"
		/>

		<DeleteSecretDialog
			v-model:open="isDeleteOpen"
			:secret="secret"
		/>
	</div>
</template>
