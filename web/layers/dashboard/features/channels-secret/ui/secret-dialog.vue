<script setup lang="ts">
import { useForm } from 'vee-validate'

import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { SecretCreateInputSchema, SecretUpdateInputSchema } from '~/gql/validation-schemas.js'

import type { Secret } from '../composables/use-secrets-api'

import { useSecretsApi } from '../composables/use-secrets-api'

const props = defineProps<{
	secret: Secret | null
}>()

const open = defineModel<boolean>('open', { default: false })

const secretsApi = useSecretsApi()
const createMutation = secretsApi.useMutationCreateSecret()
const updateMutation = secretsApi.useMutationUpdateSecret()

const { handleSubmit, resetForm, setValues } = useForm({
	validationSchema: props.secret ? SecretUpdateInputSchema : SecretCreateInputSchema,
	initialValues: {
		name: '',
		description: '',
		value: '',
	},
})

const isRevealed = ref(false)

watch(open, (isOpen) => {
	if (isOpen) {
		isRevealed.value = false

		if (props.secret) {
			setValues({
				name: props.secret.name,
				description: props.secret.description ?? '',
				value: props.secret.value ?? '',
			})
		} else {
			resetForm()
		}
	}
})

const onSubmit = handleSubmit(async (values) => {
	if (props.secret) {
		await updateMutation.executeMutation({
			id: props.secret.id,
			opts: {
				name: values.name,
				description: values.description || null,
				value: values.value,
			},
		})
	} else {
		await createMutation.executeMutation({
			opts: {
				name: values.name,
				description: values.description || null,
				value: values.value,
			},
		})
	}

	open.value = false
})
</script>

<template>
	<Dialog v-model:open="open">
		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>
					{{ secret ? 'Edit Secret' : 'Create Secret' }}
				</DialogTitle>
				<DialogDescription>
					{{ secret ? 'Update the secret value.' : 'Add a new secret to your channel.' }}
				</DialogDescription>
			</DialogHeader>

			<form
				@submit="onSubmit"
				class="space-y-4"
			>
				<FormField
					v-slot="{ componentField }"
					name="name"
				>
					<FormItem>
						<FormLabel>Name</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								placeholder="MY_API_TOKEN"
								:disabled="!!secret"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField
					v-slot="{ componentField }"
					name="description"
				>
					<FormItem>
						<FormLabel>Description (optional)</FormLabel>
						<FormControl>
							<Textarea
								v-bind="componentField"
								placeholder="What this secret is used for"
								rows="2"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField
					v-slot="{ componentField }"
					name="value"
				>
					<FormItem>
						<FormLabel>Value</FormLabel>
						<FormControl>
							<div class="relative">
								<Input
									v-bind="componentField"
									:type="isRevealed ? 'text' : 'password'"
									placeholder="Enter secret value"
									class="pr-10"
								/>
								<Button
									type="button"
									variant="ghost"
									size="icon"
									class="absolute top-0 right-0 h-full px-3"
									@click="isRevealed = !isRevealed"
								>
									<Icon
										:name="isRevealed ? 'lucide:eye-off' : 'lucide:eye'"
										class="h-4 w-4"
									/>
								</Button>
							</div>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<DialogFooter>
					<Button
						type="button"
						variant="outline"
						@click="open = false"
					>
						Cancel
					</Button>
					<Button type="submit">
						{{ secret ? 'Update' : 'Create' }}
					</Button>
				</DialogFooter>
			</form>
		</DialogContent>
	</Dialog>
</template>
