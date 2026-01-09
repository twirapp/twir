import { useAdminNotifications } from '#layers/dashboard/api/admin/notifications'
import { toTypedSchema } from '@vee-validate/zod'
import { createGlobalState } from '@vueuse/core'
import { computed, ref } from 'vue'
import * as z from 'zod'

import { useFormField } from '~/composables/use-form-field'

const formSchema = toTypedSchema(
	z.object({
		userId: z.string().nullable(),
		editorJsJson: z.string(),
	})
)

export const useNotificationsForm = createGlobalState(() => {
	const userIdField = useFormField<string | null>('userId', null)
	const editorJsJsonField = useFormField<string>('editorJsJson', '')

	const formValues = computed(() => {
		return {
			userId: userIdField.fieldModel.value,
			editorJsJson: editorJsJsonField.fieldModel.value,
		}
	})

	const editableMessageId = ref<string | null>(null)
	const isEditableForm = computed(() => Boolean(editableMessageId.value))

	const notificationsApi = useAdminNotifications()
	const { executeMutation: createNotification } = notificationsApi.useMutationCreateNotification()
	const { executeMutation: updateNotification } = notificationsApi.useMutationUpdateNotifications()

	async function onSubmit(event: Event) {
		event.preventDefault()

		try {
			const { value } = await formSchema.parse(formValues.value)
			if (!value) return

			if (editableMessageId.value) {
				await updateNotification({
					id: editableMessageId.value,
					opts: { editorJsJson: value.editorJsJson },
				})
			} else {
				await createNotification({
					editorJsJson: value.editorJsJson,
					userId: value.userId,
				})
			}

			onReset()
		} catch (err) {
			console.error(err)
		}
	}

	function onReset(): void {
		editorJsJsonField.reset()
		userIdField.reset()
		editableMessageId.value = null
	}

	function resetFieldUserId(event: Event): void {
		event.stopPropagation()
		userIdField.reset()
	}

	return {
		formValues,
		userIdField,
		editorJsJsonField,
		isEditableForm,
		editableMessageId,
		onSubmit,
		onReset,
		resetFieldUserId,
	}
})
