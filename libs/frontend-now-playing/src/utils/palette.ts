export interface Rgb {
	r: number
	g: number
	b: number
}

export interface ArtworkPalette {
	surface: string
	surfaceAlt: string
	accent: string
	text: string
	mutedText: string
	glow: string
}

export const NEUTRAL_BASE: Rgb = { r: 32, g: 38, b: 42 }

const DARK_BASE: Rgb = { r: 16, g: 20, b: 25 }
const PRIMARY_TEXT: Rgb = { r: 248, g: 249, b: 250 }
const LIGHT_ENDPOINT: Rgb = { r: 255, g: 255, b: 255 }
const DARK_ENDPOINT: Rgb = { r: 0, g: 0, b: 0 }

interface ColorBucket {
	count: number
	r: number
	g: number
	b: number
}

interface Hsl {
	h: number
	s: number
	l: number
}

interface ParsedColor extends Rgb {
	alpha: number
}

function clamp(value: number, minimum: number, maximum: number): number {
	return Math.min(Math.max(Number.isFinite(value) ? value : minimum, minimum), maximum)
}

function clampChannel(value: number): number {
	return Math.round(clamp(value, 0, 255))
}

function normalizeRgb(color: Rgb): Rgb {
	return {
		r: clampChannel(color.r),
		g: clampChannel(color.g),
		b: clampChannel(color.b),
	}
}

function bucketColor(bucket: ColorBucket): Rgb {
	return normalizeRgb({
		r: bucket.r / bucket.count,
		g: bucket.g / bucket.count,
		b: bucket.b / bucket.count,
	})
}

function rgbToHsl(color: Rgb): Hsl {
	const { r, g, b } = normalizeRgb(color)
	const red = r / 255
	const green = g / 255
	const blue = b / 255
	const maximum = Math.max(red, green, blue)
	const minimum = Math.min(red, green, blue)
	const delta = maximum - minimum
	const lightness = (maximum + minimum) / 2

	if (delta === 0) {
		return { h: 0, s: 0, l: lightness }
	}

	let hue: number
	if (maximum === red) {
		hue = ((green - blue) / delta) % 6
	} else if (maximum === green) {
		hue = (blue - red) / delta + 2
	} else {
		hue = (red - green) / delta + 4
	}

	return {
		h: (hue * 60 + 360) % 360,
		s: delta / (1 - Math.abs(2 * lightness - 1)),
		l: lightness,
	}
}

function hslToRgb({ h, s, l }: Hsl): Rgb {
	const hue = ((Number.isFinite(h) ? h : 0) % 360 + 360) % 360
	const saturation = clamp(s, 0, 1)
	const lightness = clamp(l, 0, 1)
	const chroma = (1 - Math.abs(2 * lightness - 1)) * saturation
	const segment = hue / 60
	const secondary = chroma * (1 - Math.abs(segment % 2 - 1))
	let red = 0
	let green = 0
	let blue = 0

	if (segment < 1) {
		red = chroma
		green = secondary
	} else if (segment < 2) {
		red = secondary
		green = chroma
	} else if (segment < 3) {
		green = chroma
		blue = secondary
	} else if (segment < 4) {
		green = secondary
		blue = chroma
	} else if (segment < 5) {
		red = secondary
		blue = chroma
	} else {
		red = chroma
		blue = secondary
	}

	const offset = lightness - chroma / 2
	return normalizeRgb({
		r: (red + offset) * 255,
		g: (green + offset) * 255,
		b: (blue + offset) * 255,
	})
}

function mix(first: Rgb, second: Rgb, secondWeight: number): Rgb {
	const weight = clamp(secondWeight, 0, 1)
	return normalizeRgb({
		r: first.r * (1 - weight) + second.r * weight,
		g: first.g * (1 - weight) + second.g * weight,
		b: first.b * (1 - weight) + second.b * weight,
	})
}

function toHex(color: Rgb): string {
	const { r, g, b } = normalizeRgb(color)
	return `#${[r, g, b].map((channel) => channel.toString(16).padStart(2, '0')).join('')}`
}

function parseAlpha(value: string | undefined, percentage: boolean): number | null {
	if (value === undefined) return 1
	const alpha = Number(value)
	if (!Number.isFinite(alpha)) return null
	return clamp(percentage ? alpha / 100 : alpha, 0, 1)
}

function parseColor(value: string): ParsedColor | null {
	const color = value.trim().toLowerCase()
	const hex = /^#([0-9a-f]{3}|[0-9a-f]{4}|[0-9a-f]{6}|[0-9a-f]{8})$/.exec(color)
	if (hex?.[1]) {
		const expanded = hex[1].length <= 4
			? [...hex[1]].map((digit) => `${digit}${digit}`).join('')
			: hex[1]
		return {
			r: Number.parseInt(expanded.slice(0, 2), 16),
			g: Number.parseInt(expanded.slice(2, 4), 16),
			b: Number.parseInt(expanded.slice(4, 6), 16),
			alpha: expanded.length === 8 ? Number.parseInt(expanded.slice(6, 8), 16) / 255 : 1,
		}
	}

	const rgb = /^rgba?\(\s*([+-]?(?:\d+(?:\.\d+)?|\.\d+))\s*,\s*([+-]?(?:\d+(?:\.\d+)?|\.\d+))\s*,\s*([+-]?(?:\d+(?:\.\d+)?|\.\d+))(?:\s*,\s*([+-]?(?:\d+(?:\.\d+)?|\.\d+))(%)?)?\s*\)$/.exec(color)
	if (rgb) {
		if (color.startsWith('rgba') !== (rgb[4] !== undefined)) return null
		const channels = rgb.slice(1, 4).map(Number)
		const alpha = parseAlpha(rgb[4], rgb[5] === '%')
		if (channels.some((channel) => !Number.isFinite(channel)) || alpha === null) return null

		return {
			...normalizeRgb({
				r: channels[0] ?? 0,
				g: channels[1] ?? 0,
				b: channels[2] ?? 0,
			}),
			alpha,
		}
	}

	const hsl = /^hsla?\(\s*([+-]?(?:\d+(?:\.\d+)?|\.\d+))\s*,\s*([+-]?(?:\d+(?:\.\d+)?|\.\d+))%\s*,\s*([+-]?(?:\d+(?:\.\d+)?|\.\d+))%(?:\s*,\s*([+-]?(?:\d+(?:\.\d+)?|\.\d+))(%)?)?\s*\)$/.exec(color)
	if (!hsl) return null
	if (color.startsWith('hsla') !== (hsl[4] !== undefined)) return null

	const hue = Number(hsl[1])
	const saturation = Number(hsl[2])
	const lightness = Number(hsl[3])
	const alpha = parseAlpha(hsl[4], hsl[5] === '%')
	if (
		![hue, saturation, lightness].every(Number.isFinite)
		|| saturation < 0
		|| saturation > 100
		|| lightness < 0
		|| lightness > 100
		|| alpha === null
	) return null

	return {
		...hslToRgb({
			h: hue,
			s: saturation / 100,
			l: lightness / 100,
		}),
		alpha,
	}
}

function resolveColor(value: string): Rgb | null {
	const parsed = parseColor(value)
	return parsed ? mix(NEUTRAL_BASE, parsed, parsed.alpha) : null
}

function relativeLuminance(color: Rgb): number {
	const channels = [color.r, color.g, color.b].map((channel) => {
		const normalized = clampChannel(channel) / 255
		return normalized <= 0.04045
			? normalized / 12.92
			: ((normalized + 0.055) / 1.055) ** 2.4
	})

	return (channels[0] ?? 0) * 0.2126
		+ (channels[1] ?? 0) * 0.7152
		+ (channels[2] ?? 0) * 0.0722
}

function rgbContrastRatio(foreground: Rgb, background: Rgb): number {
	const foregroundLuminance = relativeLuminance(foreground)
	const backgroundLuminance = relativeLuminance(background)
	const lighter = Math.max(foregroundLuminance, backgroundLuminance)
	const darker = Math.min(foregroundLuminance, backgroundLuminance)
	return (lighter + 0.05) / (darker + 0.05)
}

function minimumContrast(color: Rgb, backgrounds: Rgb[]): number {
	return Math.min(...backgrounds.map((background) => rgbContrastRatio(color, background)))
}

function adjustContrast(color: Rgb, backgrounds: Rgb[], threshold: number): Rgb {
	const normalized = normalizeRgb(color)
	if (minimumContrast(normalized, backgrounds) >= threshold) return normalized

	const endpoint = minimumContrast(LIGHT_ENDPOINT, backgrounds)
		>= minimumContrast(DARK_ENDPOINT, backgrounds)
		? LIGHT_ENDPOINT
		: DARK_ENDPOINT
	for (let step = 1; step <= 20; step++) {
		const candidate = mix(normalized, endpoint, step / 20)
		if (minimumContrast(candidate, backgrounds) >= threshold) return candidate
	}

	return { ...endpoint }
}

export function selectArtworkColors(samples: Rgb[]): { dominant: Rgb; accent: Rgb } {
	if (samples.length === 0) {
		return { dominant: { ...NEUTRAL_BASE }, accent: { ...NEUTRAL_BASE } }
	}

	const buckets = new Map<string, ColorBucket>()
	for (const sample of samples) {
		const color = normalizeRgb(sample)
		const key = [color.r, color.g, color.b]
			.map((channel) => Math.floor(channel / 32))
			.join(':')
		const bucket = buckets.get(key) ?? { count: 0, r: 0, g: 0, b: 0 }
		bucket.count++
		bucket.r += color.r
		bucket.g += color.g
		bucket.b += color.b
		buckets.set(key, bucket)
	}

	const rankedBuckets = [...buckets.values()]
	const dominantBucket = rankedBuckets.reduce((best, candidate) => (
		candidate.count > best.count ? candidate : best
	))
	const accentBucket = rankedBuckets.reduce((best, candidate) => {
		const bestScore = rgbToHsl(bucketColor(best)).s * Math.sqrt(best.count)
		const candidateScore = rgbToHsl(bucketColor(candidate)).s * Math.sqrt(candidate.count)
		return candidateScore > bestScore ? candidate : best
	})

	return {
		dominant: bucketColor(dominantBucket),
		accent: bucketColor(accentBucket),
	}
}

export function derivePalette(samples: Rgb[], fallbackColor: string): ArtworkPalette {
	const fallback = resolveColor(fallbackColor) ?? NEUTRAL_BASE
	const selected = selectArtworkColors(samples.length > 0 ? samples : [fallback])
	let surface = mix(DARK_BASE, selected.dominant, 0.28)
	const surfaceAlt = mix(DARK_BASE, selected.dominant, 0.38)
	const accentHsl = rgbToHsl(selected.accent)
	const accentBase = accentHsl.s <= 0.01
		? normalizeRgb(selected.accent)
		: hslToRgb({
			h: accentHsl.h,
			s: Math.max(accentHsl.s, 0.65),
			l: clamp(accentHsl.l, 0.42, 0.64),
		})
	const text = { ...PRIMARY_TEXT }

	for (let adjustment = 0; adjustment < 16 && rgbContrastRatio(text, surface) < 4.5; adjustment++) {
		surface = mix(surface, DARK_BASE, 0.2)
	}
	if (rgbContrastRatio(text, surface) < 4.5) surface = { ...DARK_BASE }

	const backgrounds = [surface, surfaceAlt]
	const mutedText = adjustContrast(mix(surface, text, 0.65), backgrounds, 4.5)
	const accent = adjustContrast(accentBase, backgrounds, 3)
	return {
		surface: toHex(surface),
		surfaceAlt: toHex(surfaceAlt),
		accent: toHex(accent),
		text: toHex(text),
		mutedText: toHex(mutedText),
		glow: `rgba(${accent.r}, ${accent.g}, ${accent.b}, 0.35)`,
	}
}

export function contrastRatio(foreground: string, background: string): number {
	return rgbContrastRatio(
		resolveColor(foreground) ?? NEUTRAL_BASE,
		resolveColor(background) ?? NEUTRAL_BASE,
	)
}
