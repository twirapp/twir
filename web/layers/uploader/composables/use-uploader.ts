import type { ErrorModel, UploadedFileOutputDto } from '@twir/api/openapi'

import { useOapi } from '~/composables/use-oapi'

export type UploadedFile = UploadedFileOutputDto

// The upload response carries the full file DTO plus delete_link at runtime,
// but the generated client type only declares delete_link (unexported Go embed).
export type UploadedFileWithDeleteLink = UploadedFileOutputDto & {
	delete_link: string
}

export const useUploader = defineStore('uploader', () => {
	const api = useOapi()
	const latestUploads = ref<UploadedFile[]>([])
	const isUploading = ref(false)

	async function uploadFile(file: File) {
		isUploading.value = true
		try {
			const response = await api.v1.uploaderUploadFile({ file })
			const uploaded = response.data.data as UploadedFileWithDeleteLink

			latestUploads.value = [uploaded, ...latestUploads.value.slice(0, 2)]

			return {
				data: uploaded,
				error: null,
			}
		} catch (e) {
			return {
				data: null,
				error: await parseApiError(e),
			}
		} finally {
			isUploading.value = false
		}
	}

	async function refetchUploads(opts: { page?: number; perPage?: number } = {}) {
		try {
			const response = await api.v1.uploaderGetFiles({
				page: opts.page ?? 0,
				per_page: opts.perPage ?? 3,
			})

			latestUploads.value = response.data.data.files

			return {
				data: response.data.data,
				error: null,
			}
		} catch (e) {
			return {
				data: null,
				error: await parseApiError(e),
			}
		}
	}

	async function deleteUpload(publicId: string) {
		try {
			const response = await api.v1.uploaderDeleteFile(publicId)

			latestUploads.value = latestUploads.value.filter((file) => file.id !== publicId)

			return {
				data: response.data.data,
				error: null,
			}
		} catch (e) {
			return {
				data: null,
				error: await parseApiError(e),
			}
		}
	}

	return {
		uploadFile,
		refetchUploads,
		deleteUpload,
		latestUploads,
		isUploading,
	}
})

async function parseApiError(error: unknown) {
	if (error instanceof Response) {
		try {
			const errorData = (await error.json()) as ErrorModel
			return (
				errorData.detail ??
				errorData.errors?.map((err) => err.message)?.join(', ') ??
				errorData.title ??
				'Unknown error'
			)
		} catch {
			return `HTTP ${error.status}: ${error.statusText || 'Unknown error'}`
		}
	}

	if (error instanceof Error) {
		return error.message
	}

	return 'Unknown error'
}

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(useUploader, import.meta.hot))
}
