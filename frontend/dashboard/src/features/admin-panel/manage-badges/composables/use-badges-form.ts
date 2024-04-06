import { toTypedSchema } from '@vee-validate/zod';
import { defineStore } from 'pinia';
import { useForm } from 'vee-validate';
import { computed } from 'vue';
import * as z from 'zod';

const formSchema = toTypedSchema(z.object({
	name: z.string(),
	image: z.any(),
}));

export const useBadgesForm = defineStore('admin-panel/badges-form', () => {
	const form = useForm({
		validationSchema: formSchema,
	});

	const image = computed(() => form.values.image);

	const onSubmit = form.handleSubmit(async (values) => {
		console.log(values);
	});

	function setImageField(event: Event): void {
		const files = (event.target as HTMLInputElement).files;
		if (!files) return;
		form.setFieldValue('image', files[0]);
	}

	return {
		form,
		image,
		onSubmit,
		setImageField,
	};
});
