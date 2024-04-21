import type { RpcOptions, UnaryCall } from '@protobuf-ts/runtime-rpc';
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import type { MaybeRefOrGetter } from 'vue';
import { unref } from 'vue';

import { adminApiClient, protectedApiClient, unprotectedApiClient } from './twirp.js';

type CallFunc<
	Req extends Record<any, any>,
	Res extends Record<any, any>
> = (input: Req, options?: RpcOptions) => UnaryCall<Req, Res>;

export const createCrudManager = <
	GetAll extends CallFunc<any, any>,
	GetOne extends CallFunc<any, any>,
	Delete extends CallFunc<any, any>,
	Patch extends CallFunc<any, any>,
	Create extends CallFunc<any, any>,
	Update extends CallFunc<any, any>,
>(opts: {
	client: typeof protectedApiClient | typeof adminApiClient | typeof unprotectedApiClient,
	getAll: GetAll,
	getOne?: GetOne | null,
	deleteOne: Delete,
	patch?: Patch | null,
	create: Create,
	update: Update,

	queryKey: string,
	invalidateAdditionalQueries?: string[],
}) => {
	const queryClient = useQueryClient();

	for (const [key, value] of Object.entries(opts)) {
		if (typeof value === 'function') {
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore
			opts[key] = value.bind(opts.client);
		}
	}

	return {
		getAll: (req: MaybeRefOrGetter<Parameters<typeof opts.getAll>[0]>) => {
			return useQuery<Awaited<ReturnType<typeof opts.getAll>['response']>>({
				queryKey: [opts.queryKey, req],
				queryFn: async () => {
					const unrefedReq = unref(req);
					const call = await opts.getAll(unrefedReq);
					return call.response;
				},
			});
		},
		getOne: opts.getOne
			? (req: Parameters<typeof opts.getOne>[0] & {
				isQueryDisabled?: boolean
			}) => useQuery<Awaited<ReturnType<typeof opts.getOne>['response']>>({
				queryKey: [opts.queryKey],
				queryFn: async () => {
					const call = await opts.getOne!(req);
					return call.response;
				},
				enabled: !req.isQueryDisabled,
			})
			: null,
		deleteOne: useMutation({
			mutationFn: async (req: Parameters<typeof opts.deleteOne>[0]) => {
				await opts.deleteOne(req);
			},
			onSuccess: () => {
				queryClient.refetchQueries([opts.queryKey]);

				for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
					queryClient.refetchQueries([queryKey]);
				}
			},
		}),
		patch: opts.patch ? useMutation<Awaited<ReturnType<typeof opts.patch>['response']>, any, Parameters<typeof opts.patch>[0]>({
			mutationFn: async (req: Parameters<typeof opts.patch>[0]) => {
				const r = await opts.patch!(req);
				return r.response;
			},
			onSuccess: () => {
				queryClient.refetchQueries([opts.queryKey]);

				for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
					queryClient.refetchQueries([queryKey]);
				}
			},
		}) : null,
		create: useMutation<Awaited<ReturnType<typeof opts.create>['response']>, any, Parameters<typeof opts.create>[0]>({
			mutationFn: async (req: Parameters<typeof opts.create>[0]) => {
				const r = await opts.create(req);
				return r.response;
			},
			onSuccess: () => {
				queryClient.refetchQueries([opts.queryKey]);

				for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
					queryClient.refetchQueries([queryKey]);
				}
			},
		}),
		update: useMutation<Awaited<ReturnType<typeof opts.update>['response']>, any, Parameters<typeof opts.update>[0]>({
			mutationFn: async (req: Parameters<typeof opts.update>[0]) => {
				const r = await opts.update(req);
				return r.response;
			},
			onSuccess: () => {
				queryClient.refetchQueries([opts.queryKey]);

				for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
					queryClient.refetchQueries([queryKey]);
				}
			},
		}),
	};
};

export const useCommandsGroupsManager = () => createCrudManager({
	client: protectedApiClient,
	queryKey: 'commands/groups',
	getAll: protectedApiClient?.commandsGroupGetAll,
	update: protectedApiClient?.commandsGroupUpdate,
	create: protectedApiClient?.commandsGroupCreate,
	patch: null,
	deleteOne: protectedApiClient?.commandsGroupDelete,
	getOne: null,
	invalidateAdditionalQueries: ['commands'],
});

export const useGreetingsManager = () => createCrudManager({
	client: protectedApiClient,
	queryKey: 'greetings',
	getAll: protectedApiClient?.greetingsGetAll,
	update: protectedApiClient?.greetingsUpdate,
	create: protectedApiClient?.greetingsCreate,
	patch: protectedApiClient?.greetingsEnableOrDisable,
	deleteOne: protectedApiClient?.greetingsDelete,
	getOne: protectedApiClient?.greetingsGetById,
});

export const useKeywordsManager = () => createCrudManager({
	client: protectedApiClient,
	queryKey: 'keywords',
	getAll: protectedApiClient?.keywordsGetAll,
	update: protectedApiClient?.keywordsUpdate,
	create: protectedApiClient?.keywordsCreate,
	patch: protectedApiClient?.keywordsEnableOrDisable,
	deleteOne: protectedApiClient?.keywordsDelete,
	getOne: protectedApiClient?.keywordsGetById,
});

export const useTimersManager = () => createCrudManager({
	client: protectedApiClient,
	queryKey: 'timers',
	getAll: protectedApiClient?.timersGetAll,
	update: protectedApiClient?.timersUpdate,
	create: protectedApiClient?.timersCreate,
	patch: protectedApiClient?.timersEnableOrDisable,
	deleteOne: protectedApiClient?.timersDelete,
	getOne: null,
});

export const useVariablesManager = () => createCrudManager({
	client: protectedApiClient,
	queryKey: 'variables',
	getAll: protectedApiClient?.variablesGetAll,
	update: protectedApiClient?.variablesUpdate,
	create: protectedApiClient?.variablesCreate,
	patch: null,
	deleteOne: protectedApiClient?.variablesDelete,
	getOne: protectedApiClient?.variablesGetById,
});

export const useEventsManager = () => createCrudManager({
	client: protectedApiClient,
	queryKey: 'events',
	getAll: protectedApiClient?.eventsGetAll,
	update: protectedApiClient?.eventsUpdate,
	create: protectedApiClient?.eventsCreate,
	patch: protectedApiClient?.eventsEnableOrDisable,
	deleteOne: protectedApiClient?.eventsDelete,
	getOne: protectedApiClient?.eventsGetById,
});

export const useRolesManager = () => createCrudManager({
	client: protectedApiClient,
	queryKey: 'roles',
	getAll: protectedApiClient?.rolesGetAll,
	update: protectedApiClient?.rolesUpdate,
	create: protectedApiClient?.rolesCreate,
	patch: null,
	deleteOne: protectedApiClient?.rolesDelete,
	getOne: null,
	invalidateAdditionalQueries: ['commands'],
});

export const useAlertsManager = () => createCrudManager({
	client: protectedApiClient,
	queryKey: 'alerts',
	getAll: protectedApiClient?.alertsGetAll,
	update: protectedApiClient?.alertsUpdate,
	create: protectedApiClient?.alertsCreate,
	patch: null,
	deleteOne: protectedApiClient?.alertsDelete,
	getOne: null,
});

export const useOverlaysRegistry = () => createCrudManager({
	client: protectedApiClient,
	queryKey: 'registry/overlays',
	getAll: protectedApiClient?.overlaysGetAll,
	update: protectedApiClient?.overlaysUpdate,
	create: protectedApiClient?.overlaysCreate,
	patch: null,
	deleteOne: protectedApiClient?.overlaysDelete,
	getOne: protectedApiClient?.overlaysGetOne,
});

export const useModerationManager = () => createCrudManager({
	client: protectedApiClient,
	queryKey: 'moderation',
	getAll: protectedApiClient?.moderationGetAll,
	update: protectedApiClient?.moderationUpdate,
	create: protectedApiClient?.moderationCreate,
	deleteOne: protectedApiClient?.moderationDelete,
	patch: protectedApiClient?.moderationEnableOrDisable,
});
