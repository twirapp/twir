import { ref } from 'vue'

interface MetaData {
	url: string
	title?: string
	description?: string
	ogImage?: string
	favicon?: string
	domain?: string
}

export function useMetaExtractor() {
	const loading = ref(false)
	const error = ref<string | null>(null)

	const extractMetaFromUrl = async (url: string): Promise<MetaData | null> => {
		loading.value = true
		error.value = null

		try {
			const urlObj = new URL(url)
			const domain = urlObj.hostname.replace('www.', '')

			const favicon = `https://www.google.com/s2/favicons?domain=${domain}&sz=64`

			let metadata: MetaData = {
				url,
				favicon,
				domain,
			}

			try {
				const response = await $fetch(`https://api.microlink.io`, {
					params: { url },
					timeout: 3000,
				})

				metadata = {
					...metadata,
					title: response.data?.title,
					description: response.data?.description,
					ogImage: response.data?.image?.url,
				}
			} catch {
				console.log('Microlink failed, using basic metadata')
			}

			return metadata
		} catch (err) {
			error.value = 'Failed to extract metadata'
			console.error('Meta extraction error:', err)
			return null
		} finally {
			loading.value = false
		}
	}

	return {
		extractMetaFromUrl,
		loading,
		error,
	}
}
