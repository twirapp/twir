import {
	type ComputedRef,
	type Ref,
	computed,
	onScopeDispose,
	readonly,
	ref,
	watch,
} from 'vue'

import { type ArtworkPalette, type Rgb, derivePalette } from '../utils/palette.js'

interface ArtworkPaletteState {
	palette: Readonly<Ref<ArtworkPalette>>
	paletteStyle: ComputedRef<Record<string, string>>
}

export function useArtworkPalette(
	imageUrl: Ref<string | null | undefined>,
	fallbackColor: Ref<string>,
): ArtworkPaletteState {
	const palette = ref(derivePalette([], fallbackColor.value))
	let requestToken = 0

	watch([imageUrl, fallbackColor], ([url, fallback], _previous, onCleanup) => {
		const token = ++requestToken
		let cancelled = false
		let image: HTMLImageElement | undefined
		const isCurrent = () => !cancelled && token === requestToken
		const applyPalette = (samples: Rgb[]) => {
			if (isCurrent()) {
				palette.value = derivePalette(samples, fallback)
			}
		}
		onCleanup(() => {
			cancelled = true
			if (!image) return

			image.onload = null
			image.onerror = null
			try {
				image.src = ''
			} catch {
				// Some image implementations reject source changes during teardown.
			}
		})

		if (
			!url
			|| typeof window === 'undefined'
			|| typeof Image === 'undefined'
			|| typeof document === 'undefined'
		) {
			applyPalette([])
			return
		}

		try {
			const loadedImage = new Image()
			image = loadedImage
			loadedImage.crossOrigin = 'anonymous'
			loadedImage.onload = () => {
				if (!isCurrent()) return

				try {
					const canvas = document.createElement('canvas')
					canvas.width = 32
					canvas.height = 32
					const context = canvas.getContext('2d')
					if (!context) throw new Error('Canvas 2D context is unavailable')

					context.drawImage(loadedImage, 0, 0, 32, 32)
					const pixels = context.getImageData(0, 0, 32, 32).data
					const samples: Rgb[] = []
					for (let index = 0; index < pixels.length; index += 4) {
						if ((pixels[index + 3] ?? 0) < 128) continue
						samples.push({
							r: pixels[index] ?? 0,
							g: pixels[index + 1] ?? 0,
							b: pixels[index + 2] ?? 0,
						})
					}
					applyPalette(samples)
				} catch {
					applyPalette([])
				}
			}
			loadedImage.onerror = () => applyPalette([])
			loadedImage.src = url
		} catch {
			applyPalette([])
		}
	}, { immediate: true })

	onScopeDispose(() => {
		requestToken++
	})

	const paletteStyle = computed(() => ({
		'--np-surface': palette.value.surface,
		'--np-surface-alt': palette.value.surfaceAlt,
		'--np-accent': palette.value.accent,
		'--np-text': palette.value.text,
		'--np-muted': palette.value.mutedText,
		'--np-glow': palette.value.glow,
	}))

	return {
		palette: readonly(palette),
		paletteStyle,
	}
}
