import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import { useWebhookNotificationsApi } from '@/api/webhook-notifications'
import { toast } from 'vue-sonner'

export const formSchema = z.object({
	id: z.string().optional(),
	enabled: z.boolean().default(false),
	githubIssues: z.boolean().default(true),
	githubPullRequests: z.boolean().default(true),
	githubCommits: z.boolean().default(true),
})

export type FormSchema = z.infer<typeof formSchema>

export function useWebhookNotifications() {
	const { t } = useI18n()
	const isLoading = ref(false)

	const webhookNotificationsApi = useWebhookNotificationsApi()
	const { data, fetching, error } = webhookNotificationsApi.useQueryWebhookNotifications()
	const createMutation = webhookNotificationsApi.useMutationCreateWebhookNotifications()
	const updateMutation = webhookNotificationsApi.useMutationUpdateWebhookNotifications()

	const settings = computed(() => data.value?.webhookNotifications)
	const exists = computed(() => !!settings.value?.id)

	const form = useForm<FormSchema>({
		validationSchema: toTypedSchema(formSchema),
		initialValues: {
			enabled: false,
			githubIssues: true,
			githubPullRequests: true,
			githubCommits: true,
		},
		validateOnMount: false,
		keepValuesOnUnmount: true,
	})

	const setFormValues = () => {
		if (!settings.value) return

		form.setValues({
			id: settings.value.id,
			enabled: settings.value.enabled,
			githubIssues: settings.value.githubIssues,
			githubPullRequests: settings.value.githubPullRequests,
			githubCommits: settings.value.githubCommits,
		})
	}

	watch(data, (v) => {
		if (!v?.webhookNotifications) return

		setFormValues()
	})

	const handleSubmit = form.handleSubmit(async (values) => {
		isLoading.value = true
		try {
			if (settings.value?.id) {
				const result = await updateMutation.executeMutation({
					id: settings.value.id,
					input: {
						enabled: values.enabled,
						githubIssues: values.githubIssues,
						githubPullRequests: values.githubPullRequests,
						githubCommits: values.githubCommits,
					},
				})

				if (result.error) {
					toast.error(result.error.message || 'Error updating webhook notifications')
					return
				}
			} else {
				const result = await createMutation.executeMutation({
					input: {
						enabled: values.enabled,
						githubIssues: values.githubIssues,
						githubPullRequests: values.githubPullRequests,
						githubCommits: values.githubCommits,
					},
				})

				if (result.error) {
					toast.error(result.error.message || 'Error creating webhook notifications')
					return
				}
			}

			toast.success(t('sharedTexts.saved'))
		} catch (err) {
			console.error(err)
			toast.error('An error occurred')
		} finally {
			isLoading.value = false
		}
	})

	return {
		form,
		handleSubmit,
		isLoading,
		fetching,
		error,
		settings,
		exists,
	}
}
