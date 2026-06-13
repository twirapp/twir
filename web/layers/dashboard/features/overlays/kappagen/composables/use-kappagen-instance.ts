import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

import type { KappagenOverlayAnimationsSettings } from '@/gql/graphql.ts'
import type { Emote, KappagenAnimations, KappagenMethods } from '@twirapp/kappagen/types'

export const twirEmote: Emote = {
	url: 'https://cdn.7tv.app/emote/6548b7074789656a7be787e1/4x.webp',
	zwe: [
		{
			url: 'https://cdn.7tv.app/emote/6128ed55a50c52b1429e09dc/4x.webp',
		},
	],
}

export const useKappagenInstance = createGlobalState(() => {
	const kappagen = ref<KappagenMethods | null>(null)

	function setKappagenInstance(instance: KappagenMethods | null) {
		kappagen.value = instance
	}

	function showEmotes(emotes: Emote[]) {
		if (!kappagen.value) return
		kappagen.value.showEmotes(emotes)
	}

	function clear() {
		if (!kappagen.value) return
		kappagen.value.clear()
	}

	function playAnimation(animation: KappagenOverlayAnimationsSettings) {
		if (!kappagen.value) return Promise.resolve()

		const splittedAnimationStyle = animation.style.toLowerCase().split('_')

		const normalizedStyleName = splittedAnimationStyle
			.map((part) => part.charAt(0).toUpperCase() + part.slice(1))
			.join('')

		const normalizedAnimation = {
			...animation,
			style: normalizedStyleName,
		}

		return kappagen.value.playAnimation([twirEmote], normalizedAnimation as KappagenAnimations)
	}

	return {
		setKappagenInstance,
		showEmotes,
		clear,
		playAnimation,
	}
})
