import { toTypedSchema } from '@vee-validate/zod';
import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import * as z from 'zod';

import { useBadges } from './use-badges.js';

import { useFormField } from '@/composables/use-form-field.js';

const formSchema = toTypedSchema(z.object({
	name: z.string(),
	image: z.any(),
}));

export const useBadgesForm = defineStore('admin-panel/badges-form', () => {
	const { badgesUpload, badgesUpdate } = useBadges();

	const editableBadgeId = ref<string | null>(null);
	const isEditableForm = computed(() => Boolean(editableBadgeId.value));

	const nameField = useFormField<string>('name', '');
	const fileField = useFormField<string | File | null>('image', '');
	const fileInputRef = computed({
		get() {
			return fileField.fieldRef.value;
		},
		set(el: any) {
			fileField.fieldRef.value = el?.$el;
		},
	});

	const formValues = computed(() => {
		return {
			name: nameField.fieldModel.value,
			file: fileField.fieldModel.value,
		};
	});
	const isFormDirty = computed(() => Boolean(formValues.value.name || formValues.value.file));

	async function onSubmit(event: Event) {
		event.preventDefault();

		const parsedImage = await parseImage(formValues.value.file);

		try {
			const { value } = await formSchema.parse(formValues.value);
			if (!value) return;

			if (editableBadgeId.value) {
				await badgesUpdate.mutateAsync({
					id: editableBadgeId.value,
					name: value.name,
					enabled: true,
					...parsedImage,
				});
			} else {
				if (!parsedImage) return;
				await badgesUpload.mutateAsync({
					name: value.name,
					enabled: true,
					...parsedImage,
				});
			}

			onReset();
		} catch (err) {
			console.error(err);
		}

		onReset();
	}

	function onReset(): void {
		nameField.reset();
		fileField.reset();
		editableBadgeId.value = null;
	}

	function setImageField(event: Event): void {
		const files = (event.target as HTMLInputElement).files;
		if (!files) return;
		fileField.fieldModel.value = files[0];
	}

	async function parseImage(image: string | File | null) {
		if (!(image instanceof File)) return;
		const fileBuffer = await image.arrayBuffer();
		return {
			fileMimeType: image.type,
			fileBytes: new Uint8Array(fileBuffer),
		};
	}

	return {
		isFormDirty,
		formValues,
		nameField,
		fileField,
		fileInputRef,
		isEditableForm,
		editableBadgeId,
		onSubmit,
		onReset,
		setImageField,
	};
});
