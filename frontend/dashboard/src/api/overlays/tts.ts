import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { unref } from 'vue'

import type {
	GetInfoResponse,
	GetResponse,
	GetUsersSettingsResponse,
	PostRequest,
	SayRequest,
} from '@twir/api/messages/modules_tts/modules_tts'
import type { Ref } from 'vue'

import { protectedApiClient, unprotectedApiClient } from '@/api/twirp.js'

declare global {
	interface Window {
		webkitAudioContext: typeof AudioContext
	}
}

export function useTtsOverlayManager() {
	const queryClient = useQueryClient()
	const queryKey = ['ttsSettings']
	const usersQueryKey = ['ttsUsersSettings']

	return {
		getSettings: () => useQuery({
			queryKey,
			queryFn: async (): Promise<GetResponse> => {
				const call = await protectedApiClient.modulesTTSGet({})
				return call.response
			},
		}),
		updateSettings: () => useMutation({
			mutationKey: ['ttsUpdate'],
			mutationFn: async (opts: PostRequest) => {
				await protectedApiClient.modulesTTSUpdate(opts)
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(queryKey)
			},
		}),
		getInfo: () => useQuery({
			queryKey: ['ttsInfo'],
			queryFn: async (): Promise<GetInfoResponse> => {
				const call = await protectedApiClient.modulesTTSGetInfo({})
				return call.response
			},
		}),
		getUsersSettings: () => useQuery({
			queryKey: usersQueryKey,
			queryFn: async (): Promise<GetUsersSettingsResponse> => {
				const call = await protectedApiClient.modulesTTSGetUsersSettings({})
				return call.response
			},
		}),
		deleteUsersSettings: () => useMutation({
			mutationKey: ['ttsUsersSettingsDelete'],
			mutationFn: async (ids: string[] | Ref<string[]>) => {
				const usersIds = unref(ids)

				await protectedApiClient.modulesTTSUsersDelete({ usersIds })
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(usersQueryKey)
			},
		}),
		useSay: () => useMutation({
			mutationKey: ['ttsSay'],
			mutationFn: async (opts: SayRequest) => {
				const audioContext = new (window.AudioContext || window!.webkitAudioContext)()
				const gainNode = audioContext.createGain()

				const req = await unprotectedApiClient.modulesTTSSay(opts)

				const source = audioContext.createBufferSource()

				source.buffer = await audioContext.decodeAudioData(req.response.file.buffer as ArrayBuffer)

				gainNode.gain.value = opts.volume / 100
				source.connect(gainNode)
				gainNode.connect(audioContext.destination)

				return new Promise((resolve) => {
					source.onended = () => {
						resolve(null)
					}

					source.start(0)
				})
			},
		}),
	}
}
