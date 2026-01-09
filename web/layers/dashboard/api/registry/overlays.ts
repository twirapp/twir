import { useChannelOverlayParseHtml } from '@/api/overlays/custom.js'

export const useOverlaysParseHtml = () => {
	const parseHtml = useChannelOverlayParseHtml()

	const mutate = async (htmlString: string) => {
		if (!htmlString) {
			return ''
		}
		const result = await parseHtml.executeMutation({ html: htmlString })
		return result.data?.channelOverlayParseHtml ?? ''
	}

	return {
		mutate,
		mutateAsync: mutate,
	}
}
