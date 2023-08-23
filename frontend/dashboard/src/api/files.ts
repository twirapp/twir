import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';

import { protectedApiClient } from './twirp';


export const useFileUpload = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['filesUpload'],
		mutationFn: async (f: File) => {
			const content = new Uint8Array(await f.arrayBuffer());

			await protectedApiClient.filesUpload({
				name: f.name,
				mimetype: f.type,
				content,
			});
		},
		async onSuccess() {
			await queryClient.invalidateQueries({ queryKey: ['files'] });
		},
	});
};

export const useFiles = () => useQuery({
	queryKey: ['files'],
	queryFn: async () => {
		const req = await protectedApiClient.filesGetAll({});
		return req.response;
	},
});

export const userFileDelete = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['filesDelete'],
		mutationFn: async (id: string) => {
			await protectedApiClient.filesDelete({ id });
		},
		async onSuccess() {
			await queryClient.invalidateQueries({ queryKey: ['files'] });
		},
	});
};
