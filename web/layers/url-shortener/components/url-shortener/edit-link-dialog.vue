<script setup lang="ts">
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import Button from '@/components/ui/button/Button.vue'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import Input from '@/components/ui/input/Input.vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'
import { useOapi } from '~/composables/use-oapi'
import { toast } from 'vue-sonner'

const props = defineProps<{
	open: boolean
	linkId: string
	currentShortId: string
	currentUrl: string
}>()

const emit = defineEmits<{
	(e: 'update:open', value: boolean): void
	(e: 'updated'): void
}>()

const api = useOapi()
const isSubmitting = ref(false)

const formSchema = z.object({
	shortId: z
		.string()
		.min(3, 'Short ID must be at least 3 characters')
		.max(50, 'Short ID must be at most 50 characters')
		.regex(/^[a-zA-Z0-9]+$/, 'Short ID can only contain letters and numbers'),
	url: z.string().url('Must be a valid URL').min(1, 'URL is required'),
})

const { handleSubmit, setFieldValue, resetForm } = useForm({
	validationSchema: formSchema,
})

// Initialize form with current values when dialog opens
watch(
	() => props.open,
	(isOpen) => {
		if (isOpen) {
			setFieldValue('shortId', props.currentShortId)
			setFieldValue('url', props.currentUrl)
		}
	}
)

onMounted(() => {
	setFieldValue('shortId', props.currentShortId)
	setFieldValue('url', props.currentUrl)
})

const onSubmit = handleSubmit(async (values) => {
	isSubmitting.value = true

	try {
		// Only send changed fields
		const body: { new_short_id?: string; url?: string } = {}

		if (values.shortId !== props.currentShortId) {
			body.new_short_id = values.shortId
		}

		if (values.url !== props.currentUrl) {
			body.url = values.url
		}

		// If nothing changed, just close
		if (Object.keys(body).length === 0) {
			toast.info('No changes detected')
			closeDialog()
			return
		}

		await api.v1.shortUrlUpdate(props.linkId, body)

		toast.success('Success', {
			description: 'Short link updated successfully',
		})

		emit('updated')
		closeDialog()
		resetForm()
	} catch (err: any) {
		console.error('Failed to update link:', err)
		const errorMessage = err?.data?.detail || 'Failed to update link'
		toast.error('Error', {
			description: errorMessage,
		})
	} finally {
		isSubmitting.value = false
	}
})

function closeDialog() {
	emit('update:open', false)
}
</script>

<template>
	<Dialog
		:open="open"
		@update:open="closeDialog"
	>
		<DialogContent class="max-w-md">
			<DialogHeader>
				<DialogTitle>Edit Short Link</DialogTitle>
				<DialogDescription>
					Update your short link's ID or target URL. Changes take effect immediately.
				</DialogDescription>
			</DialogHeader>

			<form
				@submit="onSubmit"
				class="space-y-4"
			>
				<FormField
					v-slot="{ componentField }"
					name="shortId"
				>
					<FormItem>
						<FormLabel>Short ID</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								placeholder="abc123"
								:disabled="isSubmitting"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField
					v-slot="{ componentField }"
					name="url"
				>
					<FormItem>
						<FormLabel>Target URL</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								type="url"
								placeholder="https://example.com"
								:disabled="isSubmitting"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<DialogFooter>
					<Button
						type="button"
						variant="outline"
						@click="closeDialog"
						:disabled="isSubmitting"
					>
						Cancel
					</Button>
					<Button
						type="submit"
						:disabled="isSubmitting"
					>
						<Icon
							v-if="isSubmitting"
							name="lucide:loader-2"
							class="h-4 w-4 mr-2 animate-spin"
						/>
						Save Changes
					</Button>
				</DialogFooter>
			</form>
		</DialogContent>
	</Dialog>
</template>
