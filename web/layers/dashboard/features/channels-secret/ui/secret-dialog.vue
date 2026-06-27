<script setup lang="ts">
import { useForm } from 'vee-validate'
import { z } from 'zod'

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

import type { Secret } from '../composables/use-secrets-api'

import { useSecretsApi } from '../composables/use-secrets-api'

const props = defineProps<{
	secret: Secret | null
}>()

const open = defineModel<boolean>('open', { default: false })

const secretsApi = useSecretsApi()
const createMutation = secretsApi.useMutationCreateSecret()
const updateMutation = secretsApi.useMutationUpdateSecret()

const formSchema = z.object({
	name: z.string().min(1, 'Name is required').max(100),
	description: z.string().max(500).optional(),
	value: z.string().min(1, 'Value is required').max(10000),
})

const { handleSubmit, resetForm, setValues } = useForm({
	validationSchema: formSchema,
	initialValues: {
		name: '',
		description: '',
		value: '',
	},
})

watch(open, (isOpen) => {
	if (isOpen) {
		if (props.secret) {
			setValues({
				name: props.secret.name,
				description: props.secret.description ?? '',
				value: '',
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
							<Input
								v-bind="componentField"
								type="password"
								placeholder="Enter secret value"
							/>
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
