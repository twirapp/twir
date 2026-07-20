import { describe, expect, test } from 'bun:test'
import { effectScope, isReadonly, nextTick, ref } from 'vue'

import { useArtworkPalette } from '../src/composables/use-artwork-palette.ts'
import {
	NEUTRAL_BASE,
	type Rgb,
	contrastRatio,
	derivePalette,
	selectArtworkColors,
} from '../src/utils/palette.ts'

const browserGlobals = ['window', 'Image', 'document'] as const
type BrowserGlobal = typeof browserGlobals[number]

function parseHexColor(color: string): Rgb {
	return {
		r: Number.parseInt(color.slice(1, 3), 16),
		g: Number.parseInt(color.slice(3, 5), 16),
		b: Number.parseInt(color.slice(5, 7), 16),
	}
}

function restoreGlobal(name: BrowserGlobal, descriptor: PropertyDescriptor | undefined) {
	if (descriptor) {
		Object.defineProperty(globalThis, name, descriptor)
	} else {
		Reflect.deleteProperty(globalThis, name)
	}
}

interface MockPixel extends Rgb {
	a: number
}

function installBrowserMocks(pixelsByUrl: Record<string, MockPixel[]>) {
	const descriptors = new Map(browserGlobals.map((name) => [
		name,
		Object.getOwnPropertyDescriptor(globalThis, name),
	]))
	const images: MockImage[] = []
	const canvases: Array<{ width: number; height: number }> = []
	const drawSizes: number[][] = []
	const readSizes: number[][] = []
	let drawnImage: MockImage | undefined
	let getImageDataError: Error | undefined

	class MockImage {
		crossOrigin: string | null = null
		onload: (() => void) | null = null
		onerror: (() => void) | null = null
		src = ''

		constructor() {
			images.push(this)
		}

		load() {
			this.onload?.()
		}

		fail() {
			this.onerror?.()
		}
	}

	Object.defineProperty(globalThis, 'window', {
		configurable: true,
		value: {},
	})
	Object.defineProperty(globalThis, 'Image', {
		configurable: true,
		value: MockImage,
	})
	Object.defineProperty(globalThis, 'document', {
		configurable: true,
		value: {
			createElement(tag: string) {
				expect(tag).toBe('canvas')
				const canvas = {
					width: 0,
					height: 0,
					getContext(kind: string) {
						expect(kind).toBe('2d')
						return {
							drawImage(image: MockImage, x: number, y: number, width: number, height: number) {
								drawnImage = image
								drawSizes.push([x, y, width, height])
							},
							getImageData(x: number, y: number, width: number, height: number) {
								readSizes.push([x, y, width, height])
								if (getImageDataError) throw new Error(getImageDataError.message)

								const data = new Uint8ClampedArray(width * height * 4)
								const pixels = pixelsByUrl[drawnImage?.src ?? ''] ?? []
								for (let index = 0; index < Math.min(pixels.length, width * height); index++) {
									const pixel = pixels[index]
									if (!pixel) continue
									const offset = index * 4
									data[offset] = pixel.r
									data[offset + 1] = pixel.g
									data[offset + 2] = pixel.b
									data[offset + 3] = pixel.a
								}
								return { data }
							},
						}
					},
				}
				canvases.push(canvas)
				return canvas
			},
		},
	})

	return {
		images,
		canvases,
		drawSizes,
		readSizes,
		failCanvas(error = new Error('Canvas is tainted')) {
			getImageDataError = error
		},
		restore() {
			for (const name of browserGlobals) {
				restoreGlobal(name, descriptors.get(name))
			}
		},
	}
}

describe('artwork palette', () => {
	test('selects a frequent blue dominant and a smaller saturated magenta accent', () => {
		const samples = [
			...Array.from({ length: 8 }, () => ({ r: 40, g: 70, b: 90 })),
			...Array.from({ length: 3 }, () => ({ r: 230, g: 40, b: 130 })),
		]

		const selected = selectArtworkColors(samples)

		expect(selected.dominant.b).toBeGreaterThan(selected.dominant.r)
		expect(selected.accent.r).toBeGreaterThan(selected.accent.g)
		expect(selected.accent.b).toBeGreaterThan(selected.accent.g)
	})

	test('uses the neutral base for empty, transparent, and invalid fallbacks', () => {
		const neutral = derivePalette([], `rgb(${NEUTRAL_BASE.r}, ${NEUTRAL_BASE.g}, ${NEUTRAL_BASE.b})`)

		expect(selectArtworkColors([])).toEqual({ dominant: NEUTRAL_BASE, accent: NEUTRAL_BASE })
		expect(derivePalette([], 'rgba(255, 0, 0, 0)')).toEqual(neutral)
		expect(derivePalette([], 'not-a-color')).toEqual(neutral)
	})

	test('supports opaque configured fallback formats and lets their color influence the palette', () => {
		const shortHex = derivePalette([], '#369')

		expect(derivePalette([], '#336699')).toEqual(shortHex)
		expect(derivePalette([], 'rgb(51, 102, 153)')).toEqual(shortHex)
		expect(derivePalette([], 'rgba(51, 102, 153, 1)')).toEqual(shortHex)
		expect(derivePalette([], '#cc3300')).not.toEqual(shortHex)
	})

	test('composites translucent fallback colors over the neutral base', () => {
		const translucentRed = derivePalette([], 'rgba(255, 0, 0, 0.5)')
		const opaqueRed = derivePalette([], 'rgba(255, 0, 0, 1)')

		expect(translucentRed).toEqual(derivePalette([], 'rgb(144, 19, 21)'))
		expect(translucentRed).not.toEqual(opaqueRed)
		expect(opaqueRed).toEqual(derivePalette([], '#ff0000'))
	})

	test('supports short and long hex alpha with CSS alpha compositing', () => {
		const cases: Array<[string, string]> = [
			['#f000', 'rgba(255, 0, 0, 0)'],
			['#ff000000', 'rgba(255, 0, 0, 0)'],
			['#ff000080', 'rgba(255, 0, 0, 0.5)'],
			['#f00f', 'rgba(255, 0, 0, 1)'],
			['#ff0000ff', 'rgba(255, 0, 0, 1)'],
			['#0f08', `rgba(0, 255, 0, ${8 / 15})`],
		]

		for (const [hex, rgba] of cases) {
			expect(derivePalette([], hex)).toEqual(derivePalette([], rgba))
		}
	})

	test('supports picker HSL output with wrapped hues and HSLA alpha forms', () => {
		const blue = derivePalette([], 'rgb(51, 102, 153)')

		expect(derivePalette([], 'hsl(210, 50%, 40%)')).toEqual(blue)
		expect(derivePalette([], 'hsl(570, 50%, 40%)')).toEqual(blue)
		expect(derivePalette([], 'hsl(-150, 50%, 40%)')).toEqual(blue)
		expect(derivePalette([], 'hsla(210, 50%, 40%, .5)')).toEqual(
			derivePalette([], 'rgba(51, 102, 153, 0.5)'),
		)
		expect(derivePalette([], 'hsla(210, 50%, 40%, 50%)')).toEqual(
			derivePalette([], 'rgba(51, 102, 153, 0.5)'),
		)
	})

	test('accepts CSS leading-dot alpha and picker numeric forms', () => {
		expect(derivePalette([], 'rgba(51, 102, 153, .5)')).toEqual(
			derivePalette([], 'rgba(51, 102, 153, 0.5)'),
		)
		expect(derivePalette([], 'hsla(-150, 50%, 40%, .5)')).toEqual(
			derivePalette([], 'hsla(210, 50%, 40%, 50%)'),
		)
		expect(derivePalette([], 'rgb(51, 102, 153)')).toEqual(
			derivePalette([], 'rgb(51.0, 102.0, 153.0)'),
		)
	})

	test('rejects malformed HSL and hex while clamping CSS alpha only', () => {
		const neutral = derivePalette([], `rgb(${NEUTRAL_BASE.r}, ${NEUTRAL_BASE.g}, ${NEUTRAL_BASE.b})`)
		const malformed = [
			'#12345',
			'#abcdz',
			'#11223344 trailing',
			'hsl(120, 50%, 50%) trailing',
			'hsl(120, 101%, 50%)',
			'hsl(120, -1%, 50%)',
			'hsl(120, 50%, 101%)',
			'hsl(120, 50, 50%)',
			'hsl(Infinity, 50%, 50%)',
			'hsla(120, 50%, 50%, NaN)',
			'hsla(120, 50%, 50%)',
			'hsl(120, 50%, 50%, .5)',
			'rgb(255., 0, 0)',
			'rgb(255, 0., 0)',
			'rgba(255, 0, 0, 0.)',
			'rgba(255, 0, 0, 50 %)',
			'hsl(120., 100%, 50%)',
			'hsl(120, 100.%, 50%)',
			'hsl(120, 100%, 50.%)',
			'hsl(120, 100 %, 50%)',
			'hsl(120, 100%, 50 %)',
			'hsla(120, 100%, 50%, 50 %)',
			'hsla(120, 100%, 50%, 50.0 %)',
			'hsla(120, 100%, 50%, 50. %)',
		]

		for (const color of malformed) {
			expect(derivePalette([], color)).toEqual(neutral)
		}
		expect(derivePalette([], 'hsla(120, 100%, 50%, 150%)')).toEqual(
			derivePalette([], 'rgb(0, 255, 0)'),
		)
		expect(derivePalette([], 'hsla(120, 100%, 50%, -25%)')).toEqual(neutral)
	})

	test('preserves an achromatic accent for grayscale artwork', () => {
		const palette = derivePalette([{ r: 128, g: 128, b: 128 }], '#fff')
		const accent = parseHexColor(palette.accent)

		expect(accent.r).toBe(accent.g)
		expect(accent.g).toBe(accent.b)
		expect(palette.glow).toBe(`rgba(${accent.r}, ${accent.g}, ${accent.b}, 0.35)`)
	})

	test('returns formatted colors with readable primary text', () => {
		const palette = derivePalette([{ r: 245, g: 245, b: 245 }], '#fff')
		const hexColor = /^#[0-9a-f]{6}$/

		for (const color of [
			palette.surface,
			palette.surfaceAlt,
			palette.accent,
			palette.text,
			palette.mutedText,
		]) {
			expect(color).toMatch(hexColor)
		}
		expect(palette.glow).toMatch(/^rgba\(\d+, \d+, \d+, 0\.35\)$/)
		expect(contrastRatio(palette.text, palette.surface)).toBeGreaterThanOrEqual(4.5)
	})

	test('keeps muted text readable against both surfaces for worst-case colors', () => {
		const samples = [
			{ r: 0, g: 0, b: 0 },
			{ r: 48, g: 48, b: 48 },
			{ r: 128, g: 128, b: 128 },
			{ r: 224, g: 224, b: 224 },
			{ r: 255, g: 255, b: 255 },
			{ r: 255, g: 0, b: 0 },
			{ r: 0, g: 255, b: 0 },
			{ r: 0, g: 0, b: 255 },
			{ r: 255, g: 255, b: 0 },
			{ r: 0, g: 255, b: 255 },
			{ r: 255, g: 0, b: 255 },
		]
		let minimumContrast = Number.POSITIVE_INFINITY

		for (const sample of samples) {
			const palette = derivePalette([sample], '#000')
			expect(contrastRatio(palette.text, palette.surface)).toBeGreaterThanOrEqual(4.5)
			for (const background of [palette.surface, palette.surfaceAlt]) {
				minimumContrast = Math.min(minimumContrast, contrastRatio(palette.mutedText, background))
			}
		}

		expect(minimumContrast).toBeGreaterThanOrEqual(4.5)
	})

	test('keeps accents distinguishable against both surfaces for worst-case colors', () => {
		const samples = [
			{ r: 0, g: 0, b: 0 },
			{ r: 48, g: 48, b: 48 },
			{ r: 128, g: 128, b: 128 },
			{ r: 224, g: 224, b: 224 },
			{ r: 255, g: 255, b: 255 },
			{ r: 255, g: 0, b: 0 },
			{ r: 0, g: 255, b: 0 },
			{ r: 0, g: 0, b: 255 },
			{ r: 255, g: 255, b: 0 },
			{ r: 0, g: 255, b: 255 },
			{ r: 255, g: 0, b: 255 },
		]
		let minimumContrast = Number.POSITIVE_INFINITY

		for (const sample of samples) {
			const palette = derivePalette([sample], '#000')
			for (const background of [palette.surface, palette.surfaceAlt]) {
				minimumContrast = Math.min(minimumContrast, contrastRatio(palette.accent, background))
			}
		}

		expect(minimumContrast).toBeGreaterThanOrEqual(3)
	})

	test('maintains contrast invariants across a representative RGB boundary grid', () => {
		const levels = [0, 64, 128, 192, 255]
		let minimumMutedContrast = Number.POSITIVE_INFINITY
		let minimumAccentContrast = Number.POSITIVE_INFINITY

		for (const r of levels) {
			for (const g of levels) {
				for (const b of levels) {
					const palette = derivePalette([{ r, g, b }], '#20262a')
					for (const background of [palette.surface, palette.surfaceAlt]) {
						minimumMutedContrast = Math.min(
							minimumMutedContrast,
							contrastRatio(palette.mutedText, background),
						)
						minimumAccentContrast = Math.min(
							minimumAccentContrast,
							contrastRatio(palette.accent, background),
						)
					}
				}
			}
		}

		expect(minimumMutedContrast).toBeGreaterThanOrEqual(4.5)
		expect(minimumAccentContrast).toBeGreaterThanOrEqual(3)
	})

	test('clamps edge channels and non-finite samples without exposing NaN', () => {
		const selected = selectArtworkColors([
			{ r: -100, g: Number.NaN, b: Number.POSITIVE_INFINITY },
			{ r: 999, g: Number.NEGATIVE_INFINITY, b: 255.8 },
		])
		const palette = derivePalette([
			{ r: -100, g: Number.NaN, b: Number.POSITIVE_INFINITY },
			{ r: 999, g: Number.NEGATIVE_INFINITY, b: 255.8 },
		], 'rgb(Infinity, 0, 0)')

		for (const color of [selected.dominant, selected.accent]) {
			for (const channel of [color.r, color.g, color.b]) {
				expect(Number.isFinite(channel)).toBe(true)
				expect(channel).toBeGreaterThanOrEqual(0)
				expect(channel).toBeLessThanOrEqual(255)
			}
		}
		expect(Object.values(palette).join('')).not.toContain('NaN')
		expect(contrastRatio('#000000', '#ffffff')).toBeCloseTo(21, 8)
		expect(Number.isFinite(contrastRatio('invalid', '#ffffff'))).toBe(true)
	})
})

test('useArtworkPalette exposes a synchronous readonly fallback style without browser globals', () => {
	const descriptors = new Map(browserGlobals.map((name) => [
		name,
		Object.getOwnPropertyDescriptor(globalThis, name),
	]))
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		for (const name of browserGlobals) Reflect.deleteProperty(globalThis, name)
		scope = effectScope()
		const result = scope.run(() => useArtworkPalette(ref(undefined), ref('#336699')))

		expect(result).toBeDefined()
		expect(isReadonly(result?.palette)).toBe(true)
		expect(result?.palette.value).toEqual(derivePalette([], '#336699'))
		expect(result?.paletteStyle.value).toEqual({
			'--np-surface': result?.palette.value.surface,
			'--np-surface-alt': result?.palette.value.surfaceAlt,
			'--np-accent': result?.palette.value.accent,
			'--np-text': result?.palette.value.text,
			'--np-muted': result?.palette.value.mutedText,
			'--np-glow': result?.palette.value.glow,
		})
	} finally {
		scope?.stop()
		for (const name of browserGlobals) restoreGlobal(name, descriptors.get(name))
	}
})

test('useArtworkPalette loads anonymous artwork through a 32x32 canvas and skips transparent pixels', () => {
	const browser = installBrowserMocks({
		'artwork-blue': [
			{ r: 30, g: 80, b: 220, a: 128 },
			{ r: 255, g: 0, b: 0, a: 127 },
		],
	})
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		scope = effectScope()
		const result = scope.run(() => useArtworkPalette(ref('artwork-blue'), ref('#111111')))
		const image = browser.images[0]
		if (!image) throw new Error('Image was not created')

		expect(image.crossOrigin).toBe('anonymous')
		image.load()

		expect(browser.canvases[0]).toMatchObject({ width: 32, height: 32 })
		expect(browser.drawSizes).toEqual([[0, 0, 32, 32]])
		expect(browser.readSizes).toEqual([[0, 0, 32, 32]])
		expect(result?.palette.value).toEqual(derivePalette([{ r: 30, g: 80, b: 220 }], '#111111'))
	} finally {
		scope?.stop()
		browser.restore()
	}
})

test('useArtworkPalette falls back without throwing on image and canvas failures', async () => {
	const browser = installBrowserMocks({ artwork: [{ r: 30, g: 80, b: 220, a: 255 }] })
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		scope = effectScope()
		const imageUrl = ref('image-error')
		const fallback = ref('#663399')
		const result = scope.run(() => useArtworkPalette(imageUrl, fallback))
		const failedImage = browser.images[0]
		if (!failedImage) throw new Error('Image was not created')

		expect(() => failedImage.fail()).not.toThrow()
		expect(result?.palette.value).toEqual(derivePalette([], fallback.value))

		imageUrl.value = 'artwork'
		await nextTick()
		const taintedImage = browser.images[1]
		if (!taintedImage) throw new Error('Replacement image was not created')
		browser.failCanvas()

		expect(() => taintedImage.load()).not.toThrow()
		expect(result?.palette.value).toEqual(derivePalette([], fallback.value))
	} finally {
		scope?.stop()
		browser.restore()
	}
})

test('useArtworkPalette ignores stale image completion and invalidates loads on cleanup', async () => {
	const browser = installBrowserMocks({
		first: [{ r: 30, g: 80, b: 220, a: 255 }],
		second: [{ r: 230, g: 40, b: 130, a: 255 }],
	})
	let scope: ReturnType<typeof effectScope> | undefined

	try {
		scope = effectScope()
		const imageUrl = ref('first')
		const fallback = ref('#111111')
		const result = scope.run(() => useArtworkPalette(imageUrl, fallback))
		const firstImage = browser.images[0]
		if (!firstImage) throw new Error('First image was not created')

		imageUrl.value = 'second'
		await nextTick()
		const secondImage = browser.images[1]
		if (!secondImage) throw new Error('Second image was not created')
		secondImage.load()
		const secondPalette = derivePalette([{ r: 230, g: 40, b: 130 }], fallback.value)
		expect(result?.palette.value).toEqual(secondPalette)
		expect(browser.canvases).toHaveLength(1)
		expect(browser.drawSizes).toHaveLength(1)
		expect(browser.readSizes).toHaveLength(1)
		expect(firstImage.onload).toBeNull()
		expect(firstImage.onerror).toBeNull()

		firstImage.load()
		expect(result?.palette.value).toEqual(secondPalette)
		expect(browser.canvases).toHaveLength(1)
		expect(browser.drawSizes).toHaveLength(1)
		expect(browser.readSizes).toHaveLength(1)

		imageUrl.value = 'first'
		await nextTick()
		const cleanupImage = browser.images[2]
		if (!cleanupImage) throw new Error('Cleanup image was not created')
		scope.stop()
		expect(cleanupImage.onload).toBeNull()
		expect(cleanupImage.onerror).toBeNull()
		cleanupImage.load()
		expect(result?.palette.value).toEqual(secondPalette)
		expect(browser.canvases).toHaveLength(1)
		expect(browser.drawSizes).toHaveLength(1)
		expect(browser.readSizes).toHaveLength(1)
	} finally {
		scope?.stop()
		browser.restore()
	}
})
