import { toTypedSchema } from '@vee-validate/zod';
import { defineStore } from 'pinia';
import { useForm } from 'vee-validate';
import { computed, ref } from 'vue';
import * as z from 'zod';

import { useBadges } from './use-badges';

const formSchema = toTypedSchema(z.object({
	name: z.string(),
	image: z.any(),
}));

export const useBadgesForm = defineStore('admin-panel/badges-form', () => {
	const { badgesUpload } = useBadges();

	const editableBadgeId = ref<string | null>(null);
	const isEditableForm = computed(() => Boolean(editableBadgeId.value));
	const form = useForm({ validationSchema: formSchema });
	const image = computed(() => form.values.image);

	const onSubmit = form.handleSubmit(async (values) => {
		const file: File = values.image;
		const fileBuffer = await file.arrayBuffer();

		await badgesUpload.mutateAsync({
			name: values.name,
			enabled: true,
			fileMimeType: file.type,
			fileBytes: new Uint8Array(fileBuffer),
		});

		onReset();
	});

	function onReset(): void {
		form.resetForm();
		editableBadgeId.value = null;
	}

	function setImageField(event: Event): void {
		const files = (event.target as HTMLInputElement).files;
		if (!files) return;
		form.setFieldValue('image', files[0]);
	}

	return {
		form,
		editableBadgeId,
		isEditableForm,
		image,
		onSubmit,
		onReset,
		setImageField,
	};
});
