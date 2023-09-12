import type { RpcOptions, UnaryCall } from '@protobuf-ts/runtime-rpc';
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';

import { protectedApiClient } from './twirp.js';

type CallFunc<
	Req extends Record<any, any>,
	Res extends Record<any, any>
> = (input: Req, options?: RpcOptions) => UnaryCall<Req, Res>;

const createCrudManager = <
	GetAll extends CallFunc<any, any>,
	GetOne extends CallFunc<any, any>,
	Delete extends CallFunc<any, any>,
	Patch extends CallFunc<any, any>,
	Create extends CallFunc<any, any>,
	Update extends CallFunc<any, any>,
>(opts: {
	getAll: GetAll,
	getOne: GetOne | null,
	deleteOne: Delete,
	patch: Patch | null,
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
			opts[key] = value.bind(protectedApiClient);
		}
	}

	return {
		getAll: (req: Parameters<typeof opts.getAll>[0]) => useQuery<Awaited<ReturnType<typeof opts.getAll>['response']>>({
			queryKey: [opts.queryKey],
			queryFn: async () => {
				const call = await opts.getAll(req);
				return call.response;
			},
		}),
		getOne: opts.getOne
			? (req: Parameters<typeof opts.getOne>[0] & { isQueryDisabled?: boolean }) => useQuery<Awaited<ReturnType<typeof opts.getOne>['response']>>({
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
				queryClient.invalidateQueries([opts.queryKey]);

				for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
					queryClient.invalidateQueries([queryKey]);
				}
			},
		}),
		patch: opts.patch ? useMutation({
			mutationFn: async (req: Parameters<typeof opts.patch>[0]) => {
				await opts.patch!(req);
			},
			onSuccess: () => {
				queryClient.invalidateQueries([opts.queryKey]);

				for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
					queryClient.invalidateQueries([queryKey]);
				}
			},
		}) : null,
		create: useMutation({
			mutationFn: async (req: Parameters<typeof opts.create>[0]) => {
				await opts.create(req);
			},
			onSuccess: () => {
				queryClient.invalidateQueries([opts.queryKey]);

				for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
					queryClient.invalidateQueries([queryKey]);
				}
			},
		}),
		update: useMutation({
			mutationFn: async (req: Parameters<typeof opts.update>[0]) => {
				await opts.update(req);

				for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
					queryClient.invalidateQueries([queryKey]);
				}
			},
			onSuccess: () => {
				queryClient.invalidateQueries([opts.queryKey]);
			},
		}),
	};
};

export const useCommandsManager = () => createCrudManager({
	queryKey: 'commands',
	getAll: protectedApiClient?.commandsGetAll,
	update: protectedApiClient?.commandsUpdate,
	create: protectedApiClient?.commandsCreate,
	patch: protectedApiClient?.commandsEnableOrDisable,
	deleteOne: protectedApiClient?.commandsDelete,
	getOne: protectedApiClient?.commandsGetById,
	invalidateAdditionalQueries: ['commands/groups', 'alerts'],
});

export const useCommandsGroupsManager = () => createCrudManager({
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
	queryKey: 'greetings',
	getAll: protectedApiClient?.greetingsGetAll,
	update: protectedApiClient?.greetingsUpdate,
	create: protectedApiClient?.greetingsCreate,
	patch: protectedApiClient?.greetingsEnableOrDisable,
	deleteOne: protectedApiClient?.greetingsDelete,
	getOne: protectedApiClient?.greetingsGetById,
});

export const useKeywordsManager = () => createCrudManager({
	queryKey: 'keywords',
	getAll: protectedApiClient?.keywordsGetAll,
	update: protectedApiClient?.keywordsUpdate,
	create: protectedApiClient?.keywordsCreate,
	patch: protectedApiClient?.keywordsEnableOrDisable,
	deleteOne: protectedApiClient?.keywordsDelete,
	getOne: protectedApiClient?.keywordsGetById,
});

export const useTimersManager = () => createCrudManager({
	queryKey: 'timers',
	getAll: protectedApiClient?.timersGetAll,
	update: protectedApiClient?.timersUpdate,
	create: protectedApiClient?.timersCreate,
	patch: protectedApiClient?.timersEnableOrDisable,
	deleteOne: protectedApiClient?.timersDelete,
	getOne: null,
});

export const useVariablesManager = () => createCrudManager({
	queryKey: 'variables',
	getAll: protectedApiClient?.variablesGetAll,
	update: protectedApiClient?.variablesUpdate,
	create: protectedApiClient?.variablesCreate,
	patch: null,
	deleteOne: protectedApiClient?.variablesDelete,
	getOne: protectedApiClient?.variablesGetById,
});

export const useEventsManager = () => createCrudManager({
	queryKey: 'events',
	getAll: protectedApiClient?.eventsGetAll,
	update: protectedApiClient?.eventsUpdate,
	create: protectedApiClient?.eventsCreate,
	patch: protectedApiClient?.eventsEnableOrDisable,
	deleteOne: protectedApiClient?.eventsDelete,
	getOne: protectedApiClient?.eventsGetById,
});

export const useRolesManager = () => createCrudManager({
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
	queryKey: 'alerts',
	getAll: protectedApiClient?.alertsGetAll,
	update: protectedApiClient?.alertsUpdate,
	create: protectedApiClient?.alertsCreate,
	patch: null,
	deleteOne: protectedApiClient?.alertsDelete,
	getOne: null,
});

export const useOverlaysRegistry = () => createCrudManager({
	queryKey: 'registry/overlays',
	getAll: protectedApiClient?.overlaysGetAll,
	update: protectedApiClient?.overlaysUpdate,
	create: protectedApiClient?.overlaysCreate,
	patch: null,
	deleteOne: protectedApiClient?.overlaysDelete,
	getOne: protectedApiClient?.overlaysGetOne,
});
