import { useMutation } from '@tanstack/vue-query'

import { protectedApiClient } from '@/api/twirp'

function b64EncodeUnicode(str: string) {
	return btoa(
		encodeURIComponent(str).replace(/%([0-9A-F]{2})/g, function toSolidBytes(_, p1) {
			return String.fromCharCode(Number.parseInt(`0x${p1}`))
		})
	)
}

function b64DecodeUnicode(str: string) {
	return decodeURIComponent(
		atob(str)
			.split('')
			.map(function (c) {
				return `%${`00${c.charCodeAt(0).toString(16)}`.slice(-2)}`
			})
			.join('')
	)
}

export const useOverlaysParseHtml = () =>
	useMutation({
		mutationFn: async (htmlString: string) => {
			if (!htmlString) {
				return ''
			}
			const req = await protectedApiClient.overlaysParseHtml({
				html: b64EncodeUnicode(htmlString),
			})

			return b64DecodeUnicode(req.response.html)
		},
	})
