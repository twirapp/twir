import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

import type { Emote, KappagenMethods } from '@twirapp/kappagen/types'
import { useFormValues } from 'vee-validate'
import type { KappagenFormSchema } from '@/features/overlays/kappagen/kappagen-form-schema.ts'

export const twirEmote: Emote = {
	url: 'https://cdn.7tv.app/emote/6548b7074789656a7be787e1/4x.webp',
	zwe: [
		{
			url: 'https://cdn.7tv.app/emote/6128ed55a50c52b1429e09dc/4x.webp',
		},
	],
}

export const useKappagenInstance = createGlobalState(() => {
	const { value: formSettings } = useFormValues<KappagenFormSchema>()

	const kappagen = ref<KappagenMethods>()

	function setKappagenInstance(instance: KappagenMethods) {
		kappagen.value = instance
	}

	function showEmotes(emotes: Emote[]) {
		console.log(kappagen.value.showEmotes, emotes)
		if (!kappagen.value) return
		kappagen.value.showEmotes(emotes)
	}

	function clear() {
		if (!kappagen.value) return
		kappagen.value.clear()
	}

	function playAnimation(animation: KappagenAnimations) {
		if (!kappagen.value) return Promise.resolve()

		const splittedAnimationStyle = animation.style.toLowerCase().split('_')

		const normalizedStyleName = splittedAnimationStyle
			.map((part) => part.charAt(0).toUpperCase() + part.slice(1))
			.join('')

		return kappagen.value.playAnimation([twirEmote], {
			...animation,
			style: normalizedStyleName,
		})
	}

	return {
		setKappagenInstance,
		showEmotes,
		clear,
		playAnimation,
	}
})
