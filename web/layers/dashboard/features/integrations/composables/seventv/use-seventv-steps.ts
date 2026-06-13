import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'

import type { VNode } from 'vue'

import { useSeventvData } from '@/features/integrations/composables/seventv/use-seventv-data.ts'
import EditorStep from '@/features/integrations/ui/seventv/steps/editor-step.vue'

interface Step {
	step: number
	title: string
	description: () => VNode | string
	completed: boolean
}

export const useSeventvSteps = createGlobalState(() => {
	const { isUserProfileExists, isEmoteSetSelected, botProfile, userProfile, isEditor } = useSeventvData()

	const steps = computed<Array<Step>>(() => {
		return [
			{
				step: 1,
				title: 'Register on 7tv.app',
				description: () => 'Go to 7tv.app and register your account',
				completed: isUserProfileExists.value,
			},
			{
				step: 2,
				title: 'Configure emote set',
				description: () => 'Create and select active emote set on 7tv.app',
				completed: isEmoteSetSelected.value,
			},
			{
				step: 3,
				title: 'Add bot as editor',
				description: () => h(EditorStep, { botProfile: botProfile.value, userProfile: userProfile.value }),
				completed: isEditor.value,
			},
		]
	})

	const currentStep = computed(() => {
		return steps.value.findIndex((step) => !step.completed)
	})

	return {
		steps,
		currentStep,
	}
})
