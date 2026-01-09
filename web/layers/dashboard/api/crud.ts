import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { unref } from 'vue'

import { protectedApiClient } from './twirp.js'

import type { RpcOptions, UnaryCall } from '@protobuf-ts/runtime-rpc'
import type { MaybeRefOrGetter } from 'vue'

type CallFunc<Req extends Record<any, any>, Res extends Record<any, any>> = (
	input: Req,
	options?: RpcOptions
) => UnaryCall<Req, Res>

export function createCrudManager<
	GetAll extends CallFunc<any, any>,
	GetOne extends CallFunc<any, any>,
	Delete extends CallFunc<any, any>,
	Patch extends CallFunc<any, any>,
	Create extends CallFunc<any, any>,
	Update extends CallFunc<any, any>,
>(opts: {
	client: typeof protectedApiClient
	getAll: GetAll
	getOne?: GetOne | null
	deleteOne: Delete
	patch?: Patch | null
	create: Create
	update: Update

	queryKey: string
	invalidateAdditionalQueries?: string[]
}) {
	const queryClient = useQueryClient()

	for (const [key, value] of Object.entries(opts)) {
		if (typeof value === 'function') {
			//
			// @ts-expect-error
			opts[key] = value.bind(opts.client)
		}
	}

	return {
		getAll: (req: MaybeRefOrGetter<Parameters<typeof opts.getAll>[0]>) => {
			return useQuery<Awaited<ReturnType<typeof opts.getAll>['response']>>({
				queryKey: [opts.queryKey, req],
				queryFn: async () => {
					const unrefedReq = unref(req)
					const call = await opts.getAll(unrefedReq)
					return call.response
				},
			})
		},
		getOne: opts.getOne
			? (
					req: Parameters<typeof opts.getOne>[0] & {
						isQueryDisabled?: boolean
					}
				) =>
					useQuery<Awaited<ReturnType<typeof opts.getOne>['response']>>({
						queryKey: [opts.queryKey],
						queryFn: async () => {
							const call = await opts.getOne!(req)
							return call.response
						},
						enabled: !req.isQueryDisabled,
					})
			: null,
		deleteOne: useMutation({
			mutationFn: async (req: Parameters<typeof opts.deleteOne>[0]) => {
				await opts.deleteOne(req)
			},
			onSuccess: () => {
				queryClient.refetchQueries([opts.queryKey])

				for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
					queryClient.refetchQueries([queryKey])
				}
			},
		}),
		patch: opts.patch
			? useMutation<
					Awaited<ReturnType<typeof opts.patch>['response']>,
					any,
					Parameters<typeof opts.patch>[0]
				>({
					mutationFn: async (req: Parameters<typeof opts.patch>[0]) => {
						const r = await opts.patch!(req)
						return r.response
					},
					onSuccess: () => {
						queryClient.refetchQueries([opts.queryKey])

						for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
							queryClient.refetchQueries([queryKey])
						}
					},
				})
			: null,
		create: useMutation<
			Awaited<ReturnType<typeof opts.create>['response']>,
			any,
			Parameters<typeof opts.create>[0]
		>({
			mutationFn: async (req: Parameters<typeof opts.create>[0]) => {
				const r = await opts.create(req)
				return r.response
			},
			onSuccess: () => {
				queryClient.refetchQueries([opts.queryKey])

				for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
					queryClient.refetchQueries([queryKey])
				}
			},
		}),
		update: useMutation<
			Awaited<ReturnType<typeof opts.update>['response']>,
			any,
			Parameters<typeof opts.update>[0]
		>({
			mutationFn: async (req: Parameters<typeof opts.update>[0]) => {
				const r = await opts.update(req)
				return r.response
			},
			onSuccess: () => {
				queryClient.refetchQueries([opts.queryKey])

				for (const queryKey of opts.invalidateAdditionalQueries ?? []) {
					queryClient.refetchQueries([queryKey])
				}
			},
		}),
	}
}
