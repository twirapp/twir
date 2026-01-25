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
import { storeToRefs } from 'pinia'
import { useUrlShortener } from '../../composables/use-url-shortener'

const props = defineProps<{
	open: boolean
	linkId: string
	currentShortId: string
	currentUrl: string
	currentShortUrl: string
}>()

const emit = defineEmits<{
	(e: 'update:open', value: boolean): void
	(e: 'updated'): void
}>()

const api = useOapi()
const isSubmitting = ref(false)
const urlShortener = useUrlShortener()
const { customDomain } = storeToRefs(urlShortener)
const hasInitializedDomain = ref(false)
const initialUseCustomDomain = ref(false)

const formSchema = z.object({
	shortId: z
		.string()
		.min(3, 'Short ID must be at least 3 characters')
		.max(50, 'Short ID must be at most 50 characters')
		.regex(/^[a-zA-Z0-9]+$/, 'Short ID can only contain letters and numbers'),
	url: z.string().url('Must be a valid URL').min(1, 'URL is required'),
	useCustomDomain: z.boolean().default(false),
})

const { handleSubmit, setFieldValue, resetForm } = useForm({
	validationSchema: toTypedSchema(formSchema),
	initialValues: {
		useCustomDomain: false,
	},
})

const hasVerifiedCustomDomain = computed(
	() => Boolean(customDomain.value?.domain && customDomain.value?.verified)
)
const customDomainLabel = computed(() => customDomain.value?.domain ?? '')
const currentShortHost = computed(() => {
	if (!props.currentShortUrl) return ''
	try {
		return new URL(props.currentShortUrl).hostname.toLowerCase()
	} catch {
		return ''
	}
})

async function syncCustomDomainField() {
	const usingCustomDomain =
		hasVerifiedCustomDomain.value &&
		customDomain.value?.domain?.toLowerCase() === currentShortHost.value

	initialUseCustomDomain.value = usingCustomDomain
	setFieldValue('useCustomDomain', usingCustomDomain)
	hasInitializedDomain.value = true
}

// Initialize form with current values when dialog opens
watch(
	() => props.open,
	async (isOpen) => {
		if (isOpen) {
			hasInitializedDomain.value = false
			setFieldValue('shortId', props.currentShortId)
			setFieldValue('url', props.currentUrl)

			if (!customDomain.value && import.meta.client) {
				await urlShortener.fetchCustomDomain()
			}

			await syncCustomDomainField()
		}
	}
)

onMounted(() => {
	setFieldValue('shortId', props.currentShortId)
	setFieldValue('url', props.currentUrl)
})

watch(
	() => customDomain.value,
	() => {
		if (props.open && !hasInitializedDomain.value) {
			syncCustomDomainField()
		}
	}
)

const onSubmit = handleSubmit(async (values) => {
	isSubmitting.value = true

	try {
		// Only send changed fields
		const body: { new_short_id?: string; url?: string; use_custom_domain?: boolean } = {}

		if (values.shortId !== props.currentShortId) {
			body.new_short_id = values.shortId
		}

		if (values.url !== props.currentUrl) {
			body.url = values.url
		}

		if (values.useCustomDomain !== initialUseCustomDomain.value) {
			body.use_custom_domain = values.useCustomDomain
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

				<FormField
					v-if="customDomain?.domain"
					v-slot="{ value, handleChange }"
					name="useCustomDomain"
				>
					<FormItem class="flex items-center justify-between gap-3">
						<div class="space-y-1">
							<FormLabel>Use custom domain</FormLabel>
							<p class="text-xs text-[hsl(240,11%,60%)]">
								{{ customDomainLabel }}
							</p>
						</div>
						<FormControl>
							<UiSwitch
								:model-value="value"
								:disabled="!hasVerifiedCustomDomain || isSubmitting"
								@update:model-value="handleChange"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<p
					v-if="customDomain?.domain && !customDomain?.verified"
					class="text-xs text-yellow-300/80"
				>
					Verify your custom domain to enable switching.
				</p>

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
