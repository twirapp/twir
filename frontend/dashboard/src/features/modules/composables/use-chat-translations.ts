import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import { useChatTranslationApi } from '@/api/chat-translation'
import { useToast } from '@/components/ui/toast/use-toast'

export const formSchema = z.object({
	id: z.string().optional(),
	enabled: z.boolean().default(false),
	targetLanguage: z.string().min(2, 'Target language is required'),
	excludedLanguages: z.array(z.string()).default([]),
	useItalic: z.boolean().default(false),
	excludedUsersIDs: z.array(z.string()).default([]),
})

export type FormSchema = z.infer<typeof formSchema>

export function useChatTranslations() {
	const { t } = useI18n()
	const toast = useToast()
	const isLoading = ref(false)

	const chatTranslationApi = useChatTranslationApi()
	const { data, fetching, error } = chatTranslationApi.useQueryChatTranslation()
	const createMutation = chatTranslationApi.useMutationCreateChatTranslation()
	const updateMutation = chatTranslationApi.useMutationUpdateChatTranslation()

	const chatTranslation = computed(() => data.value?.chatTranslation)
	const exists = computed(() => !!chatTranslation.value?.id)

	const translationsForm = useForm<FormSchema>({
		validationSchema: toTypedSchema(formSchema),
		initialValues: {
			enabled: false,
			targetLanguage: 'en',
			excludedLanguages: [],
			useItalic: false,
			excludedUsersIDs: [],
		},
		validateOnMount: false,
		keepValuesOnUnmount: true,
	})

	// Set form values when data is loaded
	const setFormValues = () => {
		if (chatTranslation.value) {
			translationsForm.setValues({
				id: chatTranslation.value.id,
				enabled: chatTranslation.value.enabled,
				targetLanguage: chatTranslation.value.targetLanguage,
				excludedLanguages: chatTranslation.value.excludedLanguages,
				useItalic: chatTranslation.value.useItalic,
				excludedUsersIDs: chatTranslation.value.excludedUsersIDs,
			})
		}
	}

	// Watch for data changes and update form values

	watch(data, (v) => {
		if (!v?.chatTranslation) return

		setFormValues()
	})

	const handleSubmit = translationsForm.handleSubmit(async (values) => {
		isLoading.value = true
		try {
			if (chatTranslation.value?.id) {
				// Update existing translation
				const result = await updateMutation.executeMutation({
					id: chatTranslation.value.id,
					input: {
						enabled: values.enabled,
						targetLanguage: values.targetLanguage,
						excludedLanguages: values.excludedLanguages,
						useItalic: values.useItalic,
						excludedUsersIDs: values.excludedUsersIDs,
					},
				})

				if (result.error) {
					toast.toast({
						title: result.error.message || 'Error updating chat translation',
						variant: 'destructive',
					})
					return
				}
			} else {
				// Create new translation
				const result = await createMutation.executeMutation({
					input: {
						enabled: values.enabled,
						targetLanguage: values.targetLanguage,
						excludedLanguages: values.excludedLanguages,
						useItalic: values.useItalic,
						excludedUsersIDs: values.excludedUsersIDs,
					},
				})

				if (result.error) {
					toast.toast({
						title: result.error.message || 'Error creating chat translation',
						variant: 'destructive',
					})
					return
				}
			}

			toast.toast({
				title: t('sharedTexts.saved'),
				variant: 'success',
			})
		} catch (err) {
			console.error(err)
			toast.toast({
				title: 'An error occurred',
				variant: 'destructive',
			})
		} finally {
			isLoading.value = false
		}
	})

	return {
		translationsForm,
		handleSubmit,
		isLoading,
		fetching,
		error,
		chatTranslation,
		exists,
	}
}
