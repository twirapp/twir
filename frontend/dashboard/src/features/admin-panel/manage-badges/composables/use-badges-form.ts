import { useMutation, useQueryClient } from '@tanstack/vue-query';
import type { CreateBadgeRequest } from '@twir/api/messages/admin_badges/admin_badges';
import { toTypedSchema } from '@vee-validate/zod';
import { defineStore } from 'pinia';
import { useForm } from 'vee-validate';
import { computed } from 'vue';
import * as z from 'zod';

import { adminApiClient } from '@/api/twirp';

const formSchema = toTypedSchema(z.object({
	name: z.string(),
	image: z.any(),
}));

export const useBadgesForm = defineStore('admin-panel/badges-form', () => {
	const form = useForm({
		validationSchema: formSchema,
	});

	const queryClient = useQueryClient();

	const uploader = useMutation({
		mutationFn: async (data: CreateBadgeRequest) => {
			return await adminApiClient.badgesCreate(data);
		},
		onSuccess() {
			queryClient.invalidateQueries(['admin/badges']);
		},
	});

	const image = computed(() => form.values.image);

	const onSubmit = form.handleSubmit(async (values) => {
		const file: File = values.image;
		const fileBuffer = await file.arrayBuffer();

		await uploader.mutateAsync({
			name: values.name,
			enabled: true,
			fileMimeType: file.type,
			fileBytes: new Uint8Array(fileBuffer),
		});
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
