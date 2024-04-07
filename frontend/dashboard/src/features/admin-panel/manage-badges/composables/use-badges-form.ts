import { toTypedSchema } from '@vee-validate/zod';
import { defineStore } from 'pinia';
import { useForm } from 'vee-validate';
import { computed, ref } from 'vue';
import * as z from 'zod';

import { useBadges } from './use-badges.js';

const formSchema = toTypedSchema(z.object({
	name: z.string(),
	image: z.any(),
}));

export const useBadgesForm = defineStore('admin-panel/badges-form', () => {
	const { badgesUpload, badgesUpdate } = useBadges();

	const editableBadgeId = ref<string | null>(null);

	const form = useForm({ validationSchema: formSchema });
	const image = computed(() => form.values.image);
	const isFormDirty = computed(() => form.isFieldDirty('name') || form.isFieldDirty('image'));

	const onSubmit = form.handleSubmit(async (values) => {
		const parsedImage = await parseImage(values.image);

		if (editableBadgeId.value) {
			await badgesUpdate.mutateAsync({
				id: editableBadgeId.value,
				name: values.name,
				enabled: true,
				...parsedImage,
			});
		} else {
			if (!parsedImage) return;
			await badgesUpload.mutateAsync({
				name: values.name,
				enabled: true,
				...parsedImage,
			});
		}

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

	async function parseImage(image: string | File) {
		if (!(image instanceof File)) return;
		const fileBuffer = await image.arrayBuffer();
		return {
			fileMimeType: image.type,
			fileBytes: new Uint8Array(fileBuffer),
		};
	}

	return {
		form,
		editableBadgeId,
		isFormDirty,
		image,
		onSubmit,
		onReset,
		setImageField,
	};
});
