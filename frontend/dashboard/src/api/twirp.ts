import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport'
import { ProtectedClient, UnProtectedClient } from '@twir/api/api.client'

import type { MethodInfo, NextUnaryFn, RpcOptions, UnaryCall } from '@protobuf-ts/runtime-rpc'

const transport = new TwirpFetchTransport({
	baseUrl: `${window.location.origin}/api-old/v1`,
	sendJson: import.meta.env.DEV,
	interceptors: [
		{
			interceptUnary(next: NextUnaryFn, method: MethodInfo, input: object, options: RpcOptions): UnaryCall {
				const locationQuery = new URLSearchParams(window.location.search)
				const apiKey = locationQuery.get('apiKey')

				if (apiKey) {
					options.meta = {
						...options.meta,
						'Api-Key': apiKey,
					}
				}

				return next(method, input, options)
			},
		},
	],
})
export const protectedApiClient = new ProtectedClient(transport)
export const unprotectedApiClient = new UnProtectedClient(transport)
