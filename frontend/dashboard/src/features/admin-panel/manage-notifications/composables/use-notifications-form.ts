import { toTypedSchema } from '@vee-validate/zod'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import * as z from 'zod'

import { useAdminNotifications } from '@/api/admin/notifications'
import { useFormField } from '@/composables/use-form-field'

const formSchema = toTypedSchema(z.object({
	userId: z.string().nullable(),
	message: z.string()
}))

export const useNotificationsForm = defineStore('admin-panel/notifications-form', () => {
	const userIdField = useFormField<string | null>('userId', null)
	const messageField = useFormField<string>('message', '')

	const formValues = computed(() => {
		return {
			userId: userIdField.fieldModel.value,
			message: messageField.fieldModel.value
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
			if (!value)
				return

			if (editableMessageId.value) {
				await updateNotification({
					id: editableMessageId.value,
					opts: { text: value.message }
				})
			} else {
				await createNotification({
					text: value.message,
					userId: value.userId
				})
			}

			onReset()
		} catch (err) {
			console.error(err)
		}
	}

	function onReset(): void {
		messageField.reset()
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
		messageField,
		isEditableForm,
		editableMessageId,
		onSubmit,
		onReset,
		resetFieldUserId
	}
})
