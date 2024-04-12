import type { Cache, UpdateResolver } from '@urql/exchange-graphcache';
import { cacheExchange } from '@urql/exchange-graphcache';
import { Client, fetchExchange, subscriptionExchange } from '@urql/vue';
import { createClient as createWS, type SubscribePayload } from 'graphql-ws';

import type { GraphCacheConfig, Mutation } from '@/gql/graphcache';

const wsUrl = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api-new/query`;
const gqlApiUrl = `${window.location.protocol}//${window.location.host}/api-new/query`;

const gqlWs = createWS({
	url: wsUrl,
	lazy: true,
	shouldRetry: () => true,
});

const invalidateCache = (
	cache: Cache,
	name: string,
	args?: { input: { id: any } },
) => {
	args
		? cache.invalidate({ __typename: name, id: args.input.id })
		: cache
			.inspectFields('Query')
			.filter((field) => field.fieldName === name)
			.forEach((field) => {
				cache.invalidate('Query', field.fieldKey);
			});
};

type MutationQueryKeys = {
	[K in keyof Mutation]?: string[]
}

const notificationByAdmin = 'notificationsByAdmin';
const mutationQueryKeys: MutationQueryKeys = {
	notificationsCreate: [notificationByAdmin],
	notificationsDelete: [notificationByAdmin],
	notificationsUpdate: [notificationByAdmin],
};

const graphCacheConfig: GraphCacheConfig = {
	updates: {
		Mutation: {},
	},
};

for (const [mutationKey, queryKeys] of Object.entries(mutationQueryKeys)) {
	const updateResolver: UpdateResolver = (_parent, _args, cache, _info) => {
		queryKeys.forEach((queryKey) => invalidateCache(cache, queryKey));
	};

	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	// @ts-ignore
	graphCacheConfig.updates.Mutation[mutationKey] = updateResolver;
}

export const urqlClient = new Client({
	url: gqlApiUrl,
	exchanges: [
		cacheExchange(graphCacheConfig),
		fetchExchange,
		subscriptionExchange({
			enableAllOperations: true,
			forwardSubscription: (operation) => ({
				subscribe: (sink) => ({
					unsubscribe: gqlWs.subscribe(operation as SubscribePayload, sink),
				}),
			}),
		}),
	],
	// requestPolicy: 'cache-first',
	fetchOptions: {
		credentials: 'include',
	},
});
