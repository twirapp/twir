import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport'

export function createTwirpTransport(baseUrl: string) {
	return new TwirpFetchTransport({
		baseUrl,
		fetchInit: {
			credentials: 'include',
		},
	})
}
