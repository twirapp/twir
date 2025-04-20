import { createGlobalState , useIntervalFn } from '@vueuse/core'
import { computed, ref } from 'vue'

import type { Giveaway } from '@/api/giveaways.js'

import { useGiveawaysApi } from '@/api/giveaways.js'

export const useGiveaways = createGlobalState(() => {
	const giveawaysApi = useGiveawaysApi()
	const { data: giveaways, fetching: giveawaysListFetching } = giveawaysApi.useGiveawaysList()

	const giveawaysList = computed<Giveaway[]>(() => {
		return giveaways.value?.giveaways as Giveaway[] ?? []
	})

	const participants = ref<number[]>(
		Array.from({ length: 10000 }, (_, i) => i),
	)

	useIntervalFn(() => {
		participants.value.push(participants.value.length + 1)
	}, 1000, { immediate: true })

	return {
		giveawaysList,
		giveawaysListFetching,
		participants,
	}
})
