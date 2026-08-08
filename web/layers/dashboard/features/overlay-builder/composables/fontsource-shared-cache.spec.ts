import { useFontSource } from '@twir/fontsource'
import { beforeEach, describe, expect, it, vi } from 'vitest'

class MockFontFace {
	family: string
	status: FontFaceStatus = 'unloaded'

	constructor(family: string, _source: string) {
		this.family = family
	}

	async load(): Promise<FontFace> {
		this.status = 'loaded'
		return this as unknown as FontFace
	}
}

function createFontFaceSetMock() {
	const faces = new Set<MockFontFace>()
	return {
		add(face: MockFontFace) {
			faces.add(face)
		},
		entries() {
			return faces.entries()
		},
		[Symbol.iterator]() {
			return faces[Symbol.iterator]()
		},
	}
}

const fontMetadata = {
	id: 'inter',
	family: 'Inter',
	subsets: ['latin'],
	weights: [400, 700],
	styles: ['normal', 'italic'],
	variants: {
		400: {
			normal: {
				latin: { url: { woff2: 'https://cdn.example.com/inter-400-latin.woff2' } },
			},
			italic: {
				latin: { url: { woff2: 'https://cdn.example.com/inter-400-italic-latin.woff2' } },
			},
		},
		700: {
			normal: {
				latin: { url: { woff2: 'https://cdn.example.com/inter-700-latin.woff2' } },
			},
		},
	},
}

describe('useFontSource shared font cache', () => {
	beforeEach(() => {
		vi.stubGlobal('FontFace', MockFontFace)
		Object.defineProperty(document, 'fonts', {
			value: createFontFaceSetMock(),
			configurable: true,
		})
		vi.stubGlobal(
			'fetch',
			vi.fn(async () => ({
				ok: true,
				json: async () => fontMetadata,
			}))
		)
	})

	it('resolves a font loaded by another composable instance', async () => {
		const first = useFontSource(false)
		const firstFont = await first.loadFont('inter', 400, 'normal')
		expect(firstFont?.id).toBe('inter')

		// a second instance (e.g. builder canvas preview vs FontSelector) must see
		// the already-loaded font instead of falling back to undefined
		const second = useFontSource(false)
		const secondFont = await second.loadFont('inter', 400, 'normal')
		expect(secondFont?.id).toBe('inter')
	})

	it('loads different weights of the same font independently', async () => {
		const instance = useFontSource(false)
		await instance.loadFont('inter', 400, 'normal')

		const bold = await instance.loadFont('inter', 700, 'normal')
		expect(bold?.id).toBe('inter')
		expect(instance.getFont('inter')?.weights).toContain(700)
	})
})
