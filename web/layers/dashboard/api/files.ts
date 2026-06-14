import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'
import { graphql } from '~/gql/gql.js'

export const useFilesApi = createGlobalState(() => {
	const query = () => useQuery({
		query: graphql(`
			query ChannelFiles {
				files {
					id
					channelId
					mimetype
					name
					size
				}
			}
		`),
		variables: {},
		context: {
			additionalTypenames: ['files'],
		},
	})

	const useDelete = () => useMutation(
		graphql(`
			mutation DeleteFile($id: UUID!) {
				filesRemove(id: $id)
			}
		`),
		['files'],
	)

	const useUpload = () => useMutation(
		graphql(`
			mutation UploadFile($file: Upload!) {
				filesUpload(file: $file) {
					id
				}
			}
		`),
		['files'],
	)

	function computeFileUrl(channelId: string, fileId: string) {
		return `${window.location.origin}/api/v1/channels/${channelId}/files/content/${fileId}`
	}

	return {
		useQuery: query,
		useDelete,
		useUpload,
		computeFileUrl,
	}
})
