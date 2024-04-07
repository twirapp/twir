import { toTypedSchema } from '@vee-validate/zod';
import { defineStore } from 'pinia';
import { useForm } from 'vee-validate';
import { computed, ref } from 'vue';
import * as z from 'zod';

import { useAdminNotifications } from '@/api/notifications';

const formSchema = toTypedSchema(z.object({
	userId: z.string().optional(),
	message: z.string(),
}));

export const useNotificationsForm = defineStore('admin-panel/notifications-form', () => {
	const editableMessageId = ref<string | null>(null);
	const isEditableForm = computed(() => Boolean(editableMessageId.value));
	const notifications = useAdminNotifications();

	const form = useForm({ validationSchema: formSchema });
	const message = computed(() => form.values.message?.trim());

	const onSubmit = form.handleSubmit(async (values) => {
		if (editableMessageId.value) {
			await notifications.update.mutateAsync({
				id: editableMessageId.value,
				message: values.message,
			});
		} else {
			await notifications.create.mutateAsync(values);
		}

		onReset();
	});

	function onReset(): void {
		form.resetForm();
		editableMessageId.value = null;
	}

	function resetFieldUserId(event: Event): void {
		event.stopPropagation();
		form.resetField('userId');
	}

	return {
		form,
		message,
		isEditableForm,
		editableMessageId,
		onSubmit,
		onReset,
		resetFieldUserId,
	};
});
