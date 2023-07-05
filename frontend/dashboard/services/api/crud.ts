import type { RpcOptions, UnaryCall } from '@protobuf-ts/runtime-rpc';
import { useMutation, useQuery } from '@tanstack/react-query';

import { protectedApiClient, queryClient } from '@/services/api';

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
}) => {

	return {
		getAll: (req: Parameters<typeof opts.getAll>[0]) => useQuery<ReturnType<typeof opts.getAll>['response']>({
			queryKey: [opts.queryKey],
			queryFn: async () => {
				const call = await opts.getAll(req);
				return call.response;
			},
		}),
		getOne: opts.getOne ? (req: Parameters<typeof opts.getOne>[0]) => useQuery<ReturnType<typeof opts.getOne>['response']>({
			queryKey: [opts.queryKey],
			queryFn: async () => {
				const call = await opts.getOne!(req);
				return call.response;
			},
		}) : null,
		deleteOne: useMutation({
			mutationFn: async (req: Parameters<typeof opts.deleteOne>[0]) => {
				await opts.deleteOne(req);
			},
			onSuccess: () => {
				queryClient.invalidateQueries([opts.queryKey]);
			},
		}),
		patch: opts.patch ? useMutation({
			mutationFn: async (req: Parameters<typeof opts.patch>[0]) => {
				await opts.patch!(req);
			},
			onSuccess: () => {
				queryClient.invalidateQueries([opts.queryKey]);
			},
		}) : null,
		create: useMutation({
			mutationFn: async (req: Parameters<typeof opts.create>[0]) => {
				await opts.create(req);
			},
			onSuccess: () => {
				queryClient.invalidateQueries([opts.queryKey]);
			},
		}),
		update: useMutation({
			mutationFn: async (req: Parameters<typeof opts.update>[0]) => {
				await opts.update(req);
			},
			onSuccess: () => {
				queryClient.invalidateQueries([opts.queryKey]);
			},
		}),
	};
};

export const commandsManager = () => createCrudManager({
	queryKey: 'commands',
	getAll: protectedApiClient.commandsGetAll,
	update: protectedApiClient.commandsUpdate,
	create: protectedApiClient.commandsCreate,
	patch: protectedApiClient.commandsEnableOrDisable,
	deleteOne: protectedApiClient.commandsDelete,
	getOne: protectedApiClient.commandsGetById,
});

export const greetingsManager = () => createCrudManager({
	queryKey: 'greetings',
	getAll: protectedApiClient.greetingsGetAll,
	update: protectedApiClient.greetingsUpdate,
	create: protectedApiClient.greetingsCreate,
	patch: protectedApiClient.greetingsEnableOrDisable,
	deleteOne: protectedApiClient.greetingsDelete,
	getOne: protectedApiClient.greetingsGetById,
});

export const keywordsManager = () => createCrudManager({
	queryKey: 'keywords',
	getAll: protectedApiClient.keywordsGetAll,
	update: protectedApiClient.keywordsUpdate,
	create: protectedApiClient.keywordsCreate,
	patch: protectedApiClient.keywordsEnableOrDisable,
	deleteOne: protectedApiClient.keywordsDelete,
	getOne: protectedApiClient.keywordsGetById,
});

export const timersManager = () => createCrudManager({
	queryKey: 'timers',
	getAll: protectedApiClient.timersGet,
	update: protectedApiClient.timersUpdate,
	create: protectedApiClient.timersCreate,
	patch: protectedApiClient.timersEnableOrDisable,
	deleteOne: protectedApiClient.timersDelete,
	getOne: null,
});

export const variablesManager = () => createCrudManager({
	queryKey: 'variables',
	getAll: protectedApiClient.variablesGetAll,
	update: protectedApiClient.variablesUpdate,
	create: protectedApiClient.variablesCreate,
	patch: protectedApiClient.variablesEnableOrDisable,
	deleteOne: protectedApiClient.variablesDelete,
	getOne: protectedApiClient.variablesGetById,
});

export const eventsManager = () => createCrudManager({
	queryKey: 'events',
	getAll: protectedApiClient.eventsGetAll,
	update: protectedApiClient.eventsUpdate,
	create: protectedApiClient.eventsCreate,
	patch: protectedApiClient.eventsEnableOrDisable,
	deleteOne: protectedApiClient.eventsDelete,
	getOne: protectedApiClient.eventsGetById,
});

export const commandsGroupsManager = () => createCrudManager({
	queryKey: 'commands/groups',
	getAll: protectedApiClient.commandsGroupsGetAll,
	update: protectedApiClient.commandsGroupsUpdate,
	create: protectedApiClient.commandsGroupsCreate,
	patch: protectedApiClient.commandsGroupsEnableOrDisable,
	deleteOne: protectedApiClient.commandsGroupsDelete,
	getOne: protectedApiClient.commandsGroupsGetById,
});

export const rolesManager = () => createCrudManager({
	queryKey: 'roles',
	getAll: protectedApiClient.rolesGetAll,
	update: protectedApiClient.rolesUpdate,
	create: protectedApiClient.rolesCreate,
	patch: null,
	deleteOne: protectedApiClient.rolesDelete,
	getOne: null,
});
