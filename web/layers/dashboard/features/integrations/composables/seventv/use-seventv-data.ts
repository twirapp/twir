import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import { useCommandsApi } from '@/api/commands/commands.ts'
import { useSevenTvIntegration } from '@/api/integrations/seventv.ts'

export const useSeventvData = createGlobalState(() => {
	const { subscription } = useSevenTvIntegration()

	const isUserProfileExists = computed(() => {
		return subscription.data?.value?.sevenTvData?.userSeventvProfile !== null
	})

	const isEmoteSetSelected = computed(() => {
		return subscription.data?.value?.sevenTvData?.emoteSetId !== null
	})

	const botProfile = computed(() => {
		return subscription.data?.value?.sevenTvData?.botSeventvProfile
	})

	const userProfile = computed(() => {
		return subscription.data?.value?.sevenTvData?.userSeventvProfile
	})

	const isEditor = computed(() => {
		return !!subscription.data?.value?.sevenTvData?.isEditor
	})

	const commandsApi = useCommandsApi()
	const { data: commands } = commandsApi.useQueryCommands()

	const sevenTvCommands = computed(() => {
		return commands.value?.commands?.filter(c => c.module === '7tv')
	})

	return {
		isUserProfileExists,
		isEmoteSetSelected,
		botProfile,
		userProfile,
		isEditor,
		sevenTvCommands,
	}
})
